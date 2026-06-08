package decoders

import (
	"fmt"
	"math"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func decodeFrameIndex(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	input, err := resolveInput(decoder, inputs)
	if err != nil {
		return Value{}, err
	}
	number, err := input.NumberValue()
	if err != nil {
		return Value{}, fmt.Errorf("decoder %q frame_index input: %w", decoder.ID, err)
	}
	if decoder.FrameCount <= 0 {
		return Value{}, fmt.Errorf("decoder %q frame_index requires positive frame_count", decoder.ID)
	}

	var normalized float64
	if input.HasRange {
		normalized = (number - input.Min) / (input.Max - input.Min)
	} else {
		normalized = number
	}
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	index := int(math.Round(normalized * float64(decoder.FrameCount-1)))
	return Value{Type: ValueTypeFrameIndex, Number: float64(index), FrameIndex: index, Status: input.Status, Error: input.Error, Min: 0, Max: float64(decoder.FrameCount - 1), HasRange: true}, nil
}
