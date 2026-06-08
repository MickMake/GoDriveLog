package decoders

import (
	"fmt"

	"github.com/MickMake/GoDriveLog/internal/config"
)

type DecodeFunc func(config.DashboardDecoderConfig, Inputs) (Value, error)

type Registry struct {
	decoders map[string]DecodeFunc
}

func NewRegistry() *Registry {
	r := &Registry{decoders: map[string]DecodeFunc{}}
	r.Register(config.DashboardDecoderNormalize, decodeNormalize)
	r.Register(config.DashboardDecoderThreshold, decodeThreshold)
	r.Register(config.DashboardDecoderFrameIndex, decodeFrameIndex)
	r.Register(config.DashboardDecoderFormatNumber, decodeFormatNumber)
	r.Register(config.DashboardDecoderDigits, decodeDigits)
	r.Register(config.DashboardDecoderBoolean, decodeBoolean)
	return r
}

func (r *Registry) Register(decoderType string, fn DecodeFunc) {
	if decoderType == "" || fn == nil {
		return
	}
	r.decoders[decoderType] = fn
}

func (r *Registry) Decode(decoder config.DashboardDecoderConfig, inputs Inputs) (Value, error) {
	fn, ok := r.decoders[decoder.Type]
	if !ok {
		return Value{}, fmt.Errorf("decoder %q type %q is not registered", decoder.ID, decoder.Type)
	}
	return fn(decoder, inputs)
}

func Execute(decoderConfigs []config.DashboardDecoderConfig, inputs Inputs) (map[string]Value, error) {
	registry := NewRegistry()
	values := map[string]Value{}
	for k, v := range inputs.Values {
		values[k] = v
	}

	for _, decoder := range decoderConfigs {
		decoderInputs := Inputs{Sensors: inputs.Sensors, Values: values}
		value, err := registry.Decode(decoder, decoderInputs)
		if err != nil {
			return nil, err
		}
		values[decoder.ID] = value
	}

	return values, nil
}
