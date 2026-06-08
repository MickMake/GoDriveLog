package decoders

import (
	"fmt"
	"sort"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func decodeThreshold(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	input, err := resolveInput(decoder, inputs)
	if err != nil {
		return Value{}, err
	}
	number, err := input.NumberValue()
	if err != nil {
		return Value{}, fmt.Errorf("decoder %q threshold input: %w", decoder.ID, err)
	}
	if len(decoder.Thresholds) == 0 {
		return Value{}, fmt.Errorf("decoder %q threshold requires at least one threshold", decoder.ID)
	}

	thresholds := append([]config.ThresholdConfig(nil), decoder.Thresholds...)
	sort.Slice(thresholds, func(i, j int) bool { return thresholds[i].At < thresholds[j].At })

	selected := thresholds[0].Value
	for _, threshold := range thresholds {
		if number < threshold.At {
			break
		}
		selected = threshold.Value
	}

	return Value{Type: ValueTypeText, Text: selected, Status: input.Status, Error: input.Error}, nil
}
