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
	default:
		return 0, "", fmt.Errorf("unsupported OBD PID %s", pid)
	}
}
