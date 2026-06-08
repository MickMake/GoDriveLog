package decoders

import (
	"fmt"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func decodeNormalize(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	input, err := resolveInput(decoder, inputs)
	if err != nil {
		return Value{}, err
	}
	if !input.HasRange {
		return Value{}, fmt.Errorf("decoder %q normalize input must have a non-zero range", decoder.ID)
	}
	normalized := (input.Number - input.Min) / (input.Max - input.Min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}
	return Value{Type: ValueTypeNumber, Number: normalized, Status: input.Status, Error: input.Error, Min: 0, Max: 1, HasRange: true}, nil
}
