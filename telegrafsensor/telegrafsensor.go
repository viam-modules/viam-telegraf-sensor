package telegrafsensor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"reflect"
	"strings"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	Model = resource.NewModel("viam", "viam-sensor", "telegrafsensor")
)

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newSensor,
		})
}

func newSensor(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (sensor.Sensor, error) {
	ts := TelegrafSensor{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	err := newTelegrafConf(conf, logger)

	if err != nil {
		return nil, err
	}
	return &ts, nil
}

type TelegrafSensor struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

type Metric struct {
	Name      string                 `json:"name"`
	Fields    map[string]interface{} `json:"fields"`
	Tags      map[string]interface{} `json:"tags"`
	Timestamp uint64                 `json:"timestamp"`
}

func (ts *TelegrafSensor) Readings(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
	metrics := map[string][]Metric{}

	telegrafOut, err := getTelegrafMetrics()
	if err != nil {
		ts.logger.Errorw("Error executing Telegraf", "error", err)
		return nil, err
	}

	for _, mline := range strings.Split(telegrafOut, "\n") {
		if mline == "" {
			continue
		}

		var metric Metric
		err := json.Unmarshal([]byte(mline), &metric)
		if err != nil {
			ts.logger.Errorw("Error parsing reading", "input", mline, "error", err)
		}

		metrics[metric.Name] = append(metrics[metric.Name], metric)
	}

	return toMap(metrics, ts.logger), nil
}

func toMap(metricsMap map[string][]Metric, logger logging.Logger) map[string]interface{} {
	results := map[string]interface{}{}

	metricsMap = mergeMetrics(metricsMap, logger)

	for name, metrics := range metricsMap {
		_, ok := results["host"]
		if !ok {
			results["host"] = metrics[0].Tags["host"]
		}

		if len(metrics) == 1 {
			results[name] = metricToMap(metrics[0])
			continue
		}

		metricsArray := []interface{}{}

		for _, metric := range metrics {
			metricsArray = append(metricsArray, metricToMap(metric))
		}

		results[name] = metricsArray
	}

	logger.Debugf("readings %v", results)

	return results
}

func metricToMap(m Metric) map[string]interface{} {
	mapM := m.Fields

	for _, tag := range metricsExtraFields[m.Name] {
		mapM[tag] = m.Tags[tag]
	}

	mapM["timestamp"] = m.Timestamp

	return mapM
}

var metricsExtraFields = map[string][]string{
	"disk":     {"device", "fstype", "path"},
	"temp":     {"sensor"},
	"diskio":   {"name"},
	"wireless": {"interface"},
	"net":      {"interface"},
}

// A given Telegraf metric may come in multiple json readings. If tags are the same, merge fields
// values to report only one Metric per set of tags.
func mergeMetrics(metricsMap map[string][]Metric, logger logging.Logger) map[string][]Metric {
	for name, metrics := range metricsMap {
		metric := metrics[0]
		merge := []Metric{metric}

		for i := 1; i < len(metrics); i++ {
			m := metrics[i]

			if reflect.DeepEqual(metric.Tags, m.Tags) {
				fields, err := appendFields(metric, m.Fields)
				if err != nil {
					logger.Errorw("Error appendig fields", "metric", metric.Name, "fields", metric.Fields, "new_fields", m.Fields)
					continue
				}
				metric.Fields = fields
			} else {
				merge = append(merge, m)
			}
		}

		metricsMap[name] = merge
	}
	return metricsMap
}

func appendFields(m Metric, newFields map[string]interface{}) (map[string]interface{}, error) {
	fields := m.Fields
	for key, val := range newFields {
		if _, ok := m.Fields[key]; ok {
			return nil, errors.New("duplicate field key")
		}

		fields[key] = val
	}

	return fields, nil
}

func getTelegrafMetrics() (string, error) {
	// telegraf must be configure to output in json format
	cmd := exec.Command("telegraf", "--config", telegrafConfPath, "--once")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
