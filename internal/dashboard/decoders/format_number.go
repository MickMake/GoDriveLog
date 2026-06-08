package decoders

import (
	"fmt"
	"math"
	"strconv"
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
		format = "0"
	}

	text, err := formatNumber(decoder.ID, number, format)
	if err != nil {
		return Value{}, err
	}
	return Value{Type: ValueTypeText, Text: text, Status: input.Status, Error: input.Error}, nil
}

func formatNumber(decoderID string, number float64, format string) (string, error) {
	if isZeroMask(format) {
		return formatZeroMask(number, len(format)), nil
	}

	if strings.Contains(format, "%") {
		text := fmt.Sprintf(format, number)
		if strings.Contains(text, "%!") {
			return "", fmt.Errorf("decoder %q format_number format %q is invalid for numeric input", decoderID, format)
		}
		return text, nil
	}

	return "", fmt.Errorf("decoder %q format_number format %q must be a zero mask like 0000 or a valid Go numeric format", decoderID, format)
}

func isZeroMask(format string) bool {
	if format == "" {
		return false
	}
	for _, r := range format {
		if r != '0' {
			return false
		}
	}
	return true
}

func formatZeroMask(number float64, width int) string {
	rounded := int64(math.Round(number))
	negative := rounded < 0
	if negative {
		rounded = -rounded
	}

	text := strconv.FormatInt(rounded, 10)
	if len(text) < width {
		text = strings.Repeat("0", width-len(text)) + text
	}
	if negative {
		text = "-" + text
	}
	return text
}
