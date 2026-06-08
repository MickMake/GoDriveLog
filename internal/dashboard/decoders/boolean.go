package decoders

import (
	"fmt"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func decodeBoolean(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	input, err := resolveInput(decoder, inputs)
	if err != nil {
		return Value{}, err
	}

	switch input.Type {
	case ValueTypeBoolean:
		return Value{Type: ValueTypeBoolean, Bool: input.Bool, Status: input.Status, Error: input.Error}, nil
	case ValueTypeText:
		switch strings.ToLower(input.Text) {
		case "true", "yes", "on", "ok", "1":
			return Value{Type: ValueTypeBoolean, Bool: true, Status: input.Status, Error: input.Error}, nil
		case "false", "no", "off", "error", "0", "":
			return Value{Type: ValueTypeBoolean, Bool: false, Status: input.Status, Error: input.Error}, nil
		default:
			return Value{}, fmt.Errorf("decoder %q boolean text input %q is not recognised", decoder.ID, input.Text)
		}
	default:
		number, err := input.NumberValue()
		if err != nil {
			return Value{}, fmt.Errorf("decoder %q boolean input: %w", decoder.ID, err)
		}
		return Value{Type: ValueTypeBoolean, Bool: number != 0, Status: input.Status, Error: input.Error}, nil
	}
}
