package telegrafsensor

import (
	"testing"
	"time"

	"go.viam.com/test"
)

func TestNetSamplerFirstReadStoresBaselineOnly(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	_, _, ok := ns.sample("eth0", 1_000, 2_000, 1000, t0)
	test.That(t, ok, test.ShouldBeFalse)
}

func TestNetSamplerSecondReadComputesUtilization(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	ns.sample("eth0", 0, 0, 1000, t0)

	// 125 MB sent, 12.5 MB recv over 10s on a 1 Gbps link.
	// sent: 125e6 * 8 / (1e9 * 10) * 100 = 10%
	// recv: 12.5e6 * 8 / (1e9 * 10) * 100 = 1%
	sent, recv, ok := ns.sample("eth0", 125_000_000, 12_500_000, 1000, t0.Add(10*time.Second))
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, sent, test.ShouldAlmostEqual, 10.0, 1e-9)
	test.That(t, recv, test.ShouldAlmostEqual, 1.0, 1e-9)
}

func TestNetSamplerSkipsWhenSpeedUnknown(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	ns.sample("wlan0", 0, 0, -1, t0)
	_, _, ok := ns.sample("wlan0", 1_000_000, 1_000_000, -1, t0.Add(5*time.Second))
	test.That(t, ok, test.ShouldBeFalse)
}

func TestNetSamplerSkipsWhenCounterDecreases(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	ns.sample("eth0", 1_000_000, 1_000_000, 1000, t0)
	// Interface flap: counters reset. Rejecting avoids underflow.
	_, _, ok := ns.sample("eth0", 500, 500, 1000, t0.Add(5*time.Second))
	test.That(t, ok, test.ShouldBeFalse)
}

func TestNetSamplerSkipsWhenGapTooLarge(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	ns.sample("eth0", 0, 0, 1000, t0)
	_, _, ok := ns.sample("eth0", 1_000_000, 1_000_000, 1000, t0.Add(maxSampleGap+time.Second))
	test.That(t, ok, test.ShouldBeFalse)
}

func TestNetSamplerSkipsWhenGapNonPositive(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	ns.sample("eth0", 0, 0, 1000, t0)
	_, _, ok := ns.sample("eth0", 1_000_000, 1_000_000, 1000, t0)
	test.That(t, ok, test.ShouldBeFalse)
}

func TestNetSamplerTracksInterfacesIndependently(t *testing.T) {
	ns := newNetSampler()
	t0 := time.Unix(1_000_000, 0)
	ns.sample("eth0", 0, 0, 1000, t0)
	// A new interface's first sample should still return ok=false even though
	// eth0 already has a baseline.
	_, _, ok := ns.sample("eth1", 1_000_000, 1_000_000, 1000, t0.Add(5*time.Second))
	test.That(t, ok, test.ShouldBeFalse)
}
