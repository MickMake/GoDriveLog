package config

import "sort"

type RuntimePID struct {
	Key     string
	RawPID  string
	Unit    string
	Refresh int
	Log     bool
	Display DisplayConfig
	Min     float64
	Max     float64
}

func ActivePIDs(cfg Config) []RuntimePID {
	keys := make([]string, 0, len(cfg.Vehicle.PIDs))
	for key := range cfg.Vehicle.PIDs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	active := make([]RuntimePID, 0, len(keys))
	for _, key := range keys {
		pid := cfg.Vehicle.PIDs[key]
		if pid.Type != "obd" {
			continue
		}
		if !pid.Log && !pid.Display.Enabled {
			continue
		}

		active = append(active, RuntimePID{
			Key:     key,
			RawPID:  pid.PID,
			Unit:    pid.Unit,
			Refresh: pid.Refresh,
			Log:     pid.Log,
			Display: pid.Display,
			Min:     pid.Min,
			Max:     pid.Max,
		})
	}

	return active
}
