package decoders

import (
	"fmt"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

type Inputs struct {
	Sensors map[string]sensors.SensorState
	Values  map[string]Value
}

func resolveInput(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	if decoder.Sensor != "" {
		state, ok := inputs.Sensors[decoder.Sensor]
		if !ok {
			return Value{}, fmt.Errorf("decoder %q sensor %q is not available", decoder.ID, decoder.Sensor)
		}
		return Value{
			Type:     ValueTypeNumber,
			Number:   state.Value,
			Unit:     state.Unit,
			Status:   state.Status,
			Error:    state.Error,
			Min:      state.Min,
			Max:      state.Max,
			HasRange: state.Max != state.Min,
		}, nil
	}

	if decoder.Input != "" {
		value, ok := inputs.Values[decoder.Input]
		if !ok {
			return Value{}, fmt.Errorf("decoder %q input %q is not available", decoder.ID, decoder.Input)
		}
		return value, nil
	}

	return Value{}, fmt.Errorf("decoder %q must define sensor or input", decoder.ID)
}
