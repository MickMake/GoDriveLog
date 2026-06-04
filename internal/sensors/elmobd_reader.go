package sensors

import (
	"context"
	"fmt"
	"sync"

	"github.com/rzetterberg/elmobd"
)

type ELMOBDReader struct {
	mu  sync.Mutex
	dev *elmobd.Device
}

func NewELMOBDReader(addr string, debug bool) (*ELMOBDReader, error) {
	dev, err := elmobd.NewDevice(addr, debug)
	if err != nil {
		return nil, err
	}
	return &ELMOBDReader{dev: dev}, nil
}

func (r *ELMOBDReader) Read(ctx context.Context, pid string) (float64, string, error) {
	select {
	case <-ctx.Done():
		return 0, "", ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	switch pid {
	case "0104":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewEngineLoad())
		if err != nil {
			return 0, "%", err
		}
		value, ok := cmd.(*elmobd.EngineLoad)
		if !ok {
			return 0, "%", fmt.Errorf("unexpected engine load command type %T", cmd)
		}
		return float64(value.Value), "%", nil
	case "0105":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewCoolantTemperature())
		if err != nil {
			return 0, "C", err
		}
		value, ok := cmd.(*elmobd.CoolantTemperature)
		if !ok {
			return 0, "C", fmt.Errorf("unexpected coolant command type %T", cmd)
		}
		return float64(value.Value), "C", nil
	case "0106":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewShortFuelTrim1())
		if err != nil {
			return 0, "%", err
		}
		value, ok := cmd.(*elmobd.ShortFuelTrim1)
		if !ok {
			return 0, "%", fmt.Errorf("unexpected short fuel trim bank 1 command type %T", cmd)
		}
		return float64(value.Value), "%", nil
	case "0107":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewLongFuelTrim1())
		if err != nil {
			return 0, "%", err
		}
		value, ok := cmd.(*elmobd.LongFuelTrim1)
		if !ok {
			return 0, "%", fmt.Errorf("unexpected long fuel trim bank 1 command type %T", cmd)
		}
		return float64(value.Value), "%", nil
	case "010B":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewIntakeManifoldPressure())
		if err != nil {
			return 0, "kPa", err
		}
		value, ok := cmd.(*elmobd.IntakeManifoldPressure)
		if !ok {
			return 0, "kPa", fmt.Errorf("unexpected intake manifold pressure command type %T", cmd)
		}
		return float64(value.Value), "kPa", nil
	case "010C":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewEngineRPM())
		if err != nil {
			return 0, "rpm", err
		}
		value, ok := cmd.(*elmobd.EngineRPM)
		if !ok {
			return 0, "rpm", fmt.Errorf("unexpected rpm command type %T", cmd)
		}
		return float64(value.Value), "rpm", nil
	case "010D":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewVehicleSpeed())
		if err != nil {
			return 0, "km/h", err
		}
		value, ok := cmd.(*elmobd.VehicleSpeed)
		if !ok {
			return 0, "km/h", fmt.Errorf("unexpected speed command type %T", cmd)
		}
		return float64(value.Value), "km/h", nil
	case "010F":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewIntakeAirTemperature())
		if err != nil {
			return 0, "C", err
		}
		value, ok := cmd.(*elmobd.IntakeAirTemperature)
		if !ok {
			return 0, "C", fmt.Errorf("unexpected intake air temperature command type %T", cmd)
		}
		return float64(value.Value), "C", nil
	case "0111":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewThrottlePosition())
		if err != nil {
			return 0, "%", err
		}
		value, ok := cmd.(*elmobd.ThrottlePosition)
		if !ok {
			return 0, "%", fmt.Errorf("unexpected throttle position command type %T", cmd)
		}
		return float64(value.Value), "%", nil
	case "012F":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewFuel())
		if err != nil {
			return 0, "%", err
		}
		value, ok := cmd.(*elmobd.Fuel)
		if !ok {
			return 0, "%", fmt.Errorf("unexpected fuel level command type %T", cmd)
		}
		return float64(value.Value), "%", nil
	case "0142":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewControlModuleVoltage())
		if err != nil {
			return 0, "V", err
		}
		value, ok := cmd.(*elmobd.ControlModuleVoltage)
		if !ok {
			return 0, "V", fmt.Errorf("unexpected control module voltage command type %T", cmd)
		}
		return float64(value.Value), "V", nil
	case "015C":
		cmd, err := r.dev.RunOBDCommand(elmobd.NewEngineOilTemperature())
		if err != nil {
			return 0, "C", err
		}
		value, ok := cmd.(*elmobd.EngineOilTemperature)
		if !ok {
			return 0, "C", fmt.Errorf("unexpected engine oil temperature command type %T", cmd)
		}
		return float64(value.Value), "C", nil
	default:
		return 0, "", fmt.Errorf("unsupported OBD PID %s", pid)
	}
}
