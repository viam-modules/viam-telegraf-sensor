package main

import (
	"context"

	"github.com/viam-modules/viam-telegraf-sensor/telegrafsensor"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("telegraf-sensor"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	sensorModule, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	sensorModule.AddModelFromRegistry(ctx, sensor.API, telegrafsensor.Model)

	err = sensorModule.Start(ctx)
	if err != nil {
		return err
	}
	defer sensorModule.Close(ctx)

	<-ctx.Done()
	return nil
}
