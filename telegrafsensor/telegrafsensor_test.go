package telegrafsensor

import (
	"testing"

	"go.viam.com/rdk/logging"
	"go.viam.com/test"
)

func TestModelName(t *testing.T) {
	test.That(t, Model.Name, test.ShouldEqual, "telegrafsensor")
}

func TestToMapKeyableByTag(t *testing.T) {
	logger := logging.NewTestLogger(t)
	input := map[string][]Metric{
		"temp": {
			{Name: "temp", Tags: map[string]interface{}{"sensor": "PMU tdie1"}, Fields: map[string]interface{}{"temp": 45.0}},
			{Name: "temp", Tags: map[string]interface{}{"sensor": "NAND CH0 temp"}, Fields: map[string]interface{}{"temp": 37.0}},
		},
	}

	result := toMap(input, logger)

	tempMap, ok := result["temp"].(map[string]interface{})
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, tempMap, test.ShouldContainKey, "PMU_tdie1")
	test.That(t, tempMap, test.ShouldContainKey, "NAND_CH0_temp")
}
