// Package main is the module entrypoint for the Telegraf sensor.
package main

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"

	"github.com/viam-modules/viam-telegraf-sensor/telegrafsensor"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("telegraf-sensor"))
}

func mainWithArgs(ctx context.Context, _ []string, _ logging.Logger) error {
	sensorModule, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	if err := sensorModule.AddModelFromRegistry(ctx, sensor.API, telegrafsensor.Model); err != nil {
		return err
	}

	if err := sensorModule.Start(ctx); err != nil {
		return err
	}
	defer sensorModule.Close(ctx)

	<-ctx.Done()
	return nil
}
