package telegrafsensor

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// maxSampleGap drops stale baselines to avoid misleading averages across idle windows.
const maxSampleGap = 60 * time.Second

type netSample struct {
	bytesSent uint64
	bytesRecv uint64
	ts        time.Time
}

type netSampler struct {
	mu   sync.Mutex
	prev map[string]netSample
}

func newNetSampler() *netSampler {
	return &netSampler{prev: map[string]netSample{}}
}

func (ns *netSampler) sample(iface string, bytesSent, bytesRecv uint64, speedMbps int64, now time.Time) (sentPct, recvPct float64, ok bool) {
	ns.mu.Lock()
	prev, hadPrev := ns.prev[iface]
	ns.prev[iface] = netSample{bytesSent: bytesSent, bytesRecv: bytesRecv, ts: now}
	ns.mu.Unlock()

	if !hadPrev || speedMbps <= 0 {
		return 0, 0, false
	}
	elapsed := now.Sub(prev.ts)
	if elapsed <= 0 || elapsed > maxSampleGap {
		return 0, 0, false
	}
	// Counters reset to 0 on interface flap or module reload; skip rather than underflow.
	if bytesSent < prev.bytesSent || bytesRecv < prev.bytesRecv {
		return 0, 0, false
	}
	capBps := float64(speedMbps) * 1e6
	elapsedSec := elapsed.Seconds()
	sentPct = roundPct(float64(bytesSent-prev.bytesSent) * 8 / (capBps * elapsedSec) * 100)
	recvPct = roundPct(float64(bytesRecv-prev.bytesRecv) * 8 / (capBps * elapsedSec) * 100)
	return sentPct, recvPct, true
}

func roundPct(v float64) float64 {
	return math.Round(v*1000) / 1000
}

func readLinkSpeedMbps(iface string) (int64, error) {
	data, err := os.ReadFile(filepath.Join("/sys/class/net", iface, "speed")) //nolint:gosec // reading kernel sysfs files by design
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

func annotateNetUtil(result map[string]interface{}, ns *netSampler) {
	raw, ok := result["net"]
	if !ok {
		return
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return
	}
	now := time.Now()
	// toMap emits a flat shape when there's one interface, keyed-by-interface otherwise.
	if _, flat := m["bytes_sent"]; flat {
		if iface, _ := m["interface"].(string); iface != "" {
			injectNetUtil(m, iface, ns, now)
		}
		return
	}
	for iface, entry := range m {
		if sub, ok := entry.(map[string]interface{}); ok {
			injectNetUtil(sub, iface, ns, now)
		}
	}
}

func injectNetUtil(m map[string]interface{}, iface string, ns *netSampler, now time.Time) {
	bytesSent, ok1 := uint64FromFloat(m["bytes_sent"])
	bytesRecv, ok2 := uint64FromFloat(m["bytes_recv"])
	if !ok1 || !ok2 {
		return
	}
	speedMbps, err := readLinkSpeedMbps(iface)
	if err != nil {
		speedMbps = -1
	}
	if sent, recv, ok := ns.sample(iface, bytesSent, bytesRecv, speedMbps, now); ok {
		m["sent_utilization_pct"] = sent
		m["received_utilization_pct"] = recv
	}
}

func uint64FromFloat(v interface{}) (uint64, bool) {
	f, ok := v.(float64)
	if !ok || f < 0 {
		return 0, false
	}
	return uint64(f), true
}
