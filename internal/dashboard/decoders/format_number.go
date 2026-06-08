package decoders

import (
	"fmt"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func decodeFormatNumber(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	input, err := resolveInput(decoder, inputs)
	if err != nil {
		return Value{}, err
	}
	number, err := input.NumberValue()
	if err != nil {
		return Value{}, fmt.Errorf("decoder %q format_number input: %w", decoder.ID, err)
	}

	format := decoder.Format
	if format == "" {
		format = "%.0f"
	}
	if !strings.Contains(format, "%") {
		return Value{}, fmt.Errorf("decoder %q format_number format must contain a format marker", decoder.ID)
	}

	return Value{Type: ValueTypeText, Text: fmt.Sprintf(format, number)}, nil
}
