package telegrafsensor

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

const (
	telegrafConfPath = "/tmp/viam-telegraf.conf"
	baseTelegrafConf = "conf/base.conf"
)

//go:embed conf
var embeddedConfFS embed.FS

func newTelegrafConf(conf resource.Config, logger logging.Logger) error {
	logger.Debugf("configuring sensor with %+v", conf)

	baseConfigData, err := embeddedConfFS.ReadFile(baseTelegrafConf)
	if err != nil {
		return fmt.Errorf("error reading base.conf: %v", err)
	}

	if err := os.WriteFile(telegrafConfPath, baseConfigData, 0o600); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	emptyInputs := true

	for confName, disableField := range metricMap {
		switch disableField {
		// Wireless and Temp disabled by default as not all systems have the necessary hardware.
		case "disable_temp", "disable_wireless":
			if conf.Attributes.Bool(disableField, true) {
				logger.Debugf("Skipping config section for %s metric", confName)
				continue
			}
		default:
			if conf.Attributes.Bool(disableField, false) {
				logger.Debugf("Skipping config section for %s metric", confName)
				continue
			}
		}

		configFileName := fmt.Sprintf("conf/inputs/%s.conf", confName)
		templateData, err := embeddedConfFS.ReadFile(configFileName)
		if err != nil {
			logger.Errorf("Skipping %s. Error reading %v", configFileName, err)
			continue
		}

		destFile, err := os.OpenFile(telegrafConfPath, os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			logger.Errorf("Error opening file %s: %v", telegrafConfPath, err)
			continue
		}

		if _, err := destFile.Write(templateData); err != nil {
			logger.Errorf("Error writing config file %s: %v", telegrafConfPath, err)
		}
		if err := destFile.Close(); err != nil {
			logger.Errorf("Error closing config file %s: %v", telegrafConfPath, err)
		}
		emptyInputs = false
		logger.Debugf("Added config section for %s metric", confName)
	}

	if emptyInputs {
		return errors.New("all Telegraf input metrics disabled. At least one metric must be enabled")
	}

	return nil
}

// Config is the module configuration schema. Each Disable* field maps to a
// Telegraf input plugin that can be turned off via the machine config.
type Config struct {
	resource.TriviallyValidateConfig
	DisableCPU       bool `json:"disable_cpu"`
	DisableDisk      bool `json:"disable_disk"`
	DisableDiskIo    bool `json:"disable_disk_io"`
	DisableKernel    bool `json:"disable_kernel"`
	DisableMem       bool `json:"disable_mem"`
	DisableNet       bool `json:"disable_net"`
	DisableNetstat   bool `json:"disable_netstat"`
	DisableProcesses bool `json:"disable_processes"`
	DisableSwap      bool `json:"disable_swap"`
	DisableSystem    bool `json:"disable_system"`
	DisableTemp      bool `json:"disable_temp,omitempty"`
	DisableWireless  bool `json:"disable_wireless,omitempty"`
}

var metricMap = map[string]string{
	"cpu":       "disable_cpu",
	"disk":      "disable_disk",
	"diskio":    "disable_disk_io",
	"kernel":    "disable_kernel",
	"mem":       "disable_mem",
	"net":       "disable_net",
	"netstat":   "disable_netstat",
	"processes": "disable_processes",
	"swap":      "disable_swap",
	"system":    "disable_system",
	"temp":      "disable_temp",
	"wireless":  "disable_wireless",
}
