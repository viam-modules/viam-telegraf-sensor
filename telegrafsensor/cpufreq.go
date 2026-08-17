package telegrafsensor

import (
	"path/filepath"
)

// readCPUFreqs returns current per-CPU frequencies in kHz from sysfs. Returns
// nil where cpufreq isn't exposed (macOS, some VMs, some kernel configs), so
// callers can omit the metric silently rather than error.
func readCPUFreqs() map[string]interface{} {
	matches, err := filepath.Glob("/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq")
	if err != nil || len(matches) == 0 {
		return nil
	}

	result := map[string]interface{}{}
	for _, path := range matches {
		khz, err := readUintFile(path)
		if err != nil {
			continue
		}
		// path looks like /sys/devices/system/cpu/cpuN/cpufreq/scaling_cur_freq
		cpuName := filepath.Base(filepath.Dir(filepath.Dir(path)))
		result[cpuName] = map[string]interface{}{
			"current_khz": khz,
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}
