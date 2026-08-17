package telegrafsensor

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readPackageThrottleStats returns the CPU package thermal throttle counters
// from sysfs. Returns nil on systems where /sys/devices/system/cpu/cpu0/thermal_throttle
// is not exposed (macOS, some VMs, some kernel configs), which lets callers
// omit the metric silently rather than error.
//
// Reads cpu0's counters, which are per-package: on a single-socket machine
// every CPU reports the same value. Multi-socket systems only see the first package.
func readPackageThrottleStats() map[string]interface{} {
	const cpu0Dir = "/sys/devices/system/cpu/cpu0/thermal_throttle"

	count, err := readUintFile(filepath.Join(cpu0Dir, "package_throttle_count"))
	if err != nil {
		return nil
	}

	result := map[string]interface{}{
		"package_count": count,
	}

	// Kernel changed this file's name from package_throttle_total_time to
	// package_throttle_total_time_ms. Try the modern name first.
	for _, name := range []string{"package_throttle_total_time_ms", "package_throttle_total_time"} {
		if t, err := readUintFile(filepath.Join(cpu0Dir, name)); err == nil {
			result["package_total_time_ms"] = t
			break
		}
	}

	return result
}

func readUintFile(path string) (uint64, error) {
	data, err := os.ReadFile(path) //nolint:gosec // reading kernel sysfs files by design
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}
