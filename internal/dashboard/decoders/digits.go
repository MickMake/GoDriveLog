package decoders

import (
	"fmt"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func decodeDigits(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	input, err := resolveInput(decoder, inputs)
	if err != nil {
		return Value{}, err
	}

	text := input.Text
	status := input.Status
	errorText := input.Error
	if input.Type != ValueTypeText {
		formatted, err := decodeFormatNumber(decoder, inputs)
		if err != nil {
			return Value{}, err
		}
		text = formatted.Text
		status = formatted.Status
		errorText = formatted.Error
	}

	digits := make([]string, 0, len(text))
	for _, r := range text {
		if r < '0' || r > '9' {
			return Value{}, fmt.Errorf("decoder %q digits output contains non-digit character %q", decoder.ID, r)
		}
		digits = append(digits, string(r))
	}
	if len(digits) == 0 {
		return Value{}, fmt.Errorf("decoder %q digits output must not be empty", decoder.ID)
	}

	return Value{Type: ValueTypeDigits, Text: strings.Join(digits, ""), Digits: digits, Status: status, Error: errorText}, nil
}
