package sensors

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	RaceDemoScenarioName = "RaceDemoScenario"

	RaceDemoWarningNone     = "none"
	RaceDemoWarningWarning  = "warning"
	RaceDemoWarningCritical = "critical"

	RaceDemoCodeOilTempHigh                = "OIL_TEMP_HIGH"
	RaceDemoCodeOilTempCritical            = "OIL_TEMP_CRITICAL"
	RaceDemoCodeCatastrophicEngineFailure  = "CATASTROPHIC_ENGINE_FAILURE"
	RaceDemoFailureThrownRod               = "THROWN_ROD"
	RaceDemoStatusEngineFailureThrownRod   = "ENGINE FAILURE - THROWN ROD"
	RaceDemoStatusOilTempWarningIgnored    = "OIL TEMP WARNING - DRIVER IGNORED"
	RaceDemoStatusOilTempCriticalIgnored   = "OIL TEMP CRITICAL - DRIVER STILL MOVING"
	RaceDemoStatusCatastrophicFailureLatch = "ENGINE FAILED - CRITICAL ALERT LATCHED"
)

type RaceDemoScenario struct{}

type RaceDemoSample struct {
	TimestampMS       int64   `json:"timestamp_ms"`
	ScenarioName      string  `json:"scenario_name"`
	ScenarioPhase     string  `json:"scenario_phase"`
	EngineOn          bool    `json:"engine_on"`
	EngineFailed      bool    `json:"engine_failed"`
	FailureCode       string  `json:"failure_code,omitempty"`
	SpeedKPH          float64 `json:"speed_kph"`
	RPM               float64 `json:"rpm"`
	Gear              string  `json:"gear"`
	ThrottlePercent   float64 `json:"throttle_percent"`
	BrakePercent      float64 `json:"brake_percent"`
	OilTempC          float64 `json:"oil_temp_c"`
	CoolantTempC      float64 `json:"coolant_temp_c"`
	OilPressureKPA    float64 `json:"oil_pressure_kpa"`
	EngineLoadPercent float64 `json:"engine_load_percent"`
	BatteryV          float64 `json:"battery_v"`
	WarningLevel      string  `json:"warning_level"`
	WarningCode       string  `json:"warning_code,omitempty"`
	StatusMessage     string  `json:"status_message,omitempty"`
	RequiresReset     bool    `json:"requires_reset"`
}

type RaceDemoReader struct {
	mu       sync.Mutex
	start    time.Time
	now      func() time.Time
	scenario RaceDemoScenario
}

func NewRaceDemoScenario() RaceDemoScenario {
	return RaceDemoScenario{}
}

func NewRaceDemoReader() *RaceDemoReader {
	return &RaceDemoReader{
		start:    time.Now(),
		now:      time.Now,
		scenario: NewRaceDemoScenario(),
	}
}

func NewRaceDemoReaderAt(start time.Time) *RaceDemoReader {
	return &RaceDemoReader{
		start:    start,
		now:      time.Now,
		scenario: NewRaceDemoScenario(),
	}
}

func (r *RaceDemoReader) Read(ctx context.Context, pid string) (float64, string, error) {
	select {
	case <-ctx.Done():
		return 0, "", ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.ReadAt(pid, r.now().Sub(r.start))
}

func (r *RaceDemoReader) ReadAt(pid string, elapsed time.Duration) (float64, string, error) {
	return valueForRaceDemoPID(r.scenario.SampleAt(elapsed), pid)
}

func (RaceDemoScenario) SampleAt(elapsed time.Duration) RaceDemoSample {
	if elapsed < 0 {
		elapsed = 0
	}
	return raceDemoSampleAtSeconds(elapsed.Seconds())
}

func valueForRaceDemoPID(sample RaceDemoSample, pid string) (float64, string, error) {
	scenarioSeconds := float64(sample.TimestampMS) / 1000

	switch pid {
	case "0104", "DEMO_ENGINE_LOAD":
		return sample.EngineLoadPercent, "%", nil
	case "0105", "DEMO_COOLANT_TEMP":
		return sample.CoolantTempC, "C", nil
	case "0106", "DEMO_SHORT_FUEL_TRIM_BANK1":
		return round1(2 * math.Sin(scenarioSeconds/6)), "%", nil
	case "0107", "DEMO_LONG_FUEL_TRIM_BANK1":
		return 1.5, "%", nil
	case "010B", "DEMO_INTAKE_MANIFOLD_PRESSURE":
		return round1(30 + sample.EngineLoadPercent*0.75), "kPa", nil
	case "010C", "DEMO_RPM":
		return sample.RPM, "rpm", nil
	case "010D", "DEMO_SPEED":
		return sample.SpeedKPH, "km/h", nil
	case "010E", "DEMO_TIMING_ADVANCE":
		return round1(8 + (sample.RPM/5200)*18 - (sample.EngineLoadPercent/100)*8), "deg", nil
	case "010F", "DEMO_INTAKE_AIR_TEMP":
		return round1(28 + (sample.CoolantTempC-70)*0.35), "C", nil
	case "0110", "DEMO_MASS_AIR_FLOW":
		return round1((sample.RPM / 1000) * (10 + sample.EngineLoadPercent/4)), "g/s", nil
	case "0111", "DEMO_THROTTLE":
		return sample.ThrottlePercent, "%", nil
	case "011F", "DEMO_RUN_TIME_SINCE_ENGINE_START":
		return math.Round(scenarioSeconds), "s", nil
	case "0142", "DEMO_BATTERY":
		return sample.BatteryV, "V", nil
	case "015C", "DEMO_OIL_TEMP":
		return sample.OilTempC, "C", nil
	case "DEMO_BRAKE":
		return sample.BrakePercent, "%", nil
	case "DEMO_GEAR":
		return gearValue(sample.Gear), "gear", nil
	case "DEMO_ENGINE_ON":
		return boolValue(sample.EngineOn), "bool", nil
	case "DEMO_ENGINE_FAILED":
		return boolValue(sample.EngineFailed), "bool", nil
	case "DEMO_OIL_PRESSURE":
		return sample.OilPressureKPA, "kPa", nil
	case "DEMO_WARNING_LEVEL":
		return warningLevelValue(sample.WarningLevel), "level", nil
	case "DEMO_REQUIRES_RESET":
		return boolValue(sample.RequiresReset), "bool", nil
	default:
		return 0, "", fmt.Errorf("unsupported race demo PID %s", pid)
	}
}

func raceDemoSampleAtSeconds(t float64) RaceDemoSample {
	s := RaceDemoSample{
		TimestampMS:       int64(math.Round(t * 1000)),
		ScenarioName:      RaceDemoScenarioName,
		ScenarioPhase:     "boot",
		EngineOn:          true,
		EngineFailed:      false,
		Gear:              "neutral",
		SpeedKPH:          0,
		RPM:               0,
		ThrottlePercent:   0,
		BrakePercent:      100,
		OilTempC:          82,
		CoolantTempC:      70,
		OilPressureKPA:    280,
		EngineLoadPercent: 5,
		BatteryV:          14.1,
		WarningLevel:      RaceDemoWarningNone,
	}

	switch {
	case t < 3:
		s.ScenarioPhase = "engine_start"
		s.RPM = lerp(0, 850, t/3)
		s.ThrottlePercent = 6
		s.BrakePercent = 100
		s.EngineLoadPercent = 10
	case t < 5:
		s.ScenarioPhase = "idle"
		s.RPM = 850
		s.ThrottlePercent = 4
		s.BrakePercent = 100
		s.EngineLoadPercent = 8
		s.OilTempC = lerp(84, 88, (t-3)/2)
		s.CoolantTempC = lerp(72, 78, (t-3)/2)
	case t < 15:
		s = stationaryRevSample(s, t)
	case t < 20:
		s.ScenarioPhase = "burnout_launch"
		u := (t - 15) / 5
		s.Gear = "1"
		s.SpeedKPH = lerp(0, 8, u)
		s.RPM = lerp(4700, 5000, u) + 160*math.Sin(u*math.Pi*4)
		s.ThrottlePercent = 100
		s.BrakePercent = lerp(100, 0, u)
		s.EngineLoadPercent = 98
		s.OilTempC = lerp(90, 96, u)
		s.CoolantTempC = lerp(82, 88, u)
		s.OilPressureKPA = lerp(420, 360, u)
	case t < 28:
		s.ScenarioPhase = "hard_accel_1st"
		u := (t - 20) / 8
		s.Gear = "1"
		s.SpeedKPH = lerp(8, 45, u)
		s.RPM = lerp(3500, 5200, u)
		s.ThrottlePercent = 100
		s.EngineLoadPercent = 96
		s.OilTempC = lerp(96, 100, u)
		s.CoolantTempC = lerp(88, 91, u)
		s.OilPressureKPA = lerp(360, 330, u)
	case t < 36:
		s.ScenarioPhase = "hard_accel_2nd"
		u := (t - 28) / 8
		s.Gear = "2"
		s.SpeedKPH = lerp(45, 82, u)
		s.RPM = lerp(3200, 5300, u)
		s.ThrottlePercent = 100
		s.EngineLoadPercent = 97
		s.OilTempC = lerp(100, 108, u)
		s.CoolantTempC = lerp(91, 95, u)
		s.OilPressureKPA = lerp(330, 290, u)
	case t < 46:
		s.ScenarioPhase = "hard_accel_3rd"
		u := (t - 36) / 10
		s.Gear = "3"
		s.SpeedKPH = lerp(82, 120, u)
		s.RPM = lerp(3400, 5000, u)
		s.ThrottlePercent = 96
		s.EngineLoadPercent = 93
		s.OilTempC = lerp(108, 114, u)
		s.CoolantTempC = lerp(95, 99, u)
		s.OilPressureKPA = lerp(290, 240, u)
	case t < 48:
		s.ScenarioPhase = "shift_3_to_4"
		u := (t - 46) / 2
		s.Gear = "4"
		s.SpeedKPH = lerp(120, 126, u)
		s.RPM = lerp(3600, 3900, u)
		s.ThrottlePercent = 92
		s.EngineLoadPercent = 88
		s.OilTempC = lerp(114, 118, u)
		s.CoolantTempC = lerp(99, 101, u)
		s.OilPressureKPA = lerp(240, 210, u)
	case t < 60:
		s.ScenarioPhase = "ignored_oil_warning_4th"
		u := (t - 48) / 12
		s.Gear = "4"
		s.SpeedKPH = lerp(126, 145, u)
		s.RPM = lerp(3900, 4800, u)
		s.ThrottlePercent = 90
		s.EngineLoadPercent = 86
		s.OilTempC = lerp(118, 123, u)
		s.CoolantTempC = lerp(101, 104, u)
		s.OilPressureKPA = lerp(210, 185, u)
		s.WarningLevel = RaceDemoWarningWarning
		s.WarningCode = RaceDemoCodeOilTempHigh
		s.StatusMessage = RaceDemoStatusOilTempWarningIgnored
	case t < 70:
		s.ScenarioPhase = "shift_4_to_5"
		u := (t - 60) / 10
		s.Gear = "5"
		s.SpeedKPH = lerp(145, 150, u)
		s.RPM = lerp(3500, 4200, u)
		s.ThrottlePercent = 78
		s.EngineLoadPercent = 74
		s.OilTempC = lerp(123, 127, u)
		s.CoolantTempC = lerp(104, 106, u)
		s.OilPressureKPA = lerp(185, 160, u)
		s.WarningLevel = RaceDemoWarningWarning
		s.WarningCode = RaceDemoCodeOilTempHigh
		s.StatusMessage = RaceDemoStatusOilTempWarningIgnored
	case t < 75:
		s.ScenarioPhase = "shift_5_to_6"
		u := (t - 70) / 5
		s.Gear = "6"
		s.SpeedKPH = 150
		s.RPM = lerp(3600, 3100, u)
		s.ThrottlePercent = 62
		s.EngineLoadPercent = 65
		s.OilTempC = lerp(127, 128, u)
		s.CoolantTempC = lerp(106, 107, u)
		s.OilPressureKPA = lerp(160, 145, u)
		s.WarningLevel = RaceDemoWarningWarning
		s.WarningCode = RaceDemoCodeOilTempHigh
		s.StatusMessage = RaceDemoStatusOilTempWarningIgnored
	case t < 95:
		s.ScenarioPhase = "casual_150_cruise"
		u := (t - 75) / 20
		s.Gear = "6"
		s.SpeedKPH = 150 + 2*math.Sin((t-75)*math.Pi/5)
		s.RPM = 3100 + 90*math.Sin((t-75)*math.Pi/5)
		s.ThrottlePercent = 58
		s.EngineLoadPercent = 64
		s.OilTempC = lerp(128, 132, u)
		s.CoolantTempC = lerp(107, 109, u)
		s.OilPressureKPA = lerp(145, 115, u)
		s.WarningLevel = RaceDemoWarningWarning
		s.WarningCode = RaceDemoCodeOilTempHigh
		s.StatusMessage = RaceDemoStatusOilTempWarningIgnored
	case t < 100:
		s.ScenarioPhase = "oil_temp_critical_cruise"
		u := (t - 95) / 5
		s.Gear = "6"
		s.SpeedKPH = lerp(150, 150, u)
		s.RPM = lerp(3100, 3000, u)
		s.ThrottlePercent = 52
		s.EngineLoadPercent = 60
		s.OilTempC = lerp(132, 134, u)
		s.CoolantTempC = lerp(109, 110, u)
		s.OilPressureKPA = lerp(115, 100, u)
		s.WarningLevel = RaceDemoWarningCritical
		s.WarningCode = RaceDemoCodeOilTempCritical
		s.StatusMessage = RaceDemoStatusOilTempCriticalIgnored
	case t < 118:
		s = emergencyDecelSample(s, t)
	default:
		s = thrownRodSample(s, t)
	}

	if t >= 132 {
		s.SpeedKPH = 0
		s.Gear = "neutral"
		s.ScenarioPhase = "stopped_failed_latched"
		s.StatusMessage = RaceDemoStatusCatastrophicFailureLatch
	}

	return roundSample(s)
}

func stationaryRevSample(s RaceDemoSample, t float64) RaceDemoSample {
	s.ScenarioPhase = "stationary_revving"
	s.Gear = "neutral"
	s.SpeedKPH = 0
	s.BrakePercent = 100
	s.EngineLoadPercent = 28
	s.OilTempC = lerp(88, 90, (t-5)/10)
	s.CoolantTempC = lerp(78, 82, (t-5)/10)
	s.OilPressureKPA = 360

	switch {
	case t < 7:
		s.RPM = triangularPulse(t, 5, 6, 7, 850, 3200, 900)
		s.ThrottlePercent = triangularPulse(t, 5, 6, 7, 5, 72, 5)
	case t < 10:
		s.RPM = triangularPulse(t, 7, 9, 10, 900, 4100, 1000)
		s.ThrottlePercent = triangularPulse(t, 7, 9, 10, 5, 88, 5)
	case t < 13:
		s.RPM = triangularPulse(t, 10, 12, 13, 1000, 4700, 1200)
		s.ThrottlePercent = triangularPulse(t, 10, 12, 13, 5, 95, 5)
	default:
		s.RPM = lerp(1200, 1500, (t-13)/2)
		s.ThrottlePercent = lerp(8, 100, (t-13)/2)
	}

	return s
}

func emergencyDecelSample(s RaceDemoSample, t float64) RaceDemoSample {
	s.EngineOn = true
	s.EngineFailed = false
	s.ThrottlePercent = 0
	s.BrakePercent = 85
	s.WarningLevel = RaceDemoWarningCritical
	s.WarningCode = RaceDemoCodeOilTempCritical
	s.StatusMessage = RaceDemoStatusOilTempCriticalIgnored
	s.OilTempC = lerp(134, 136, (t-100)/18)
	s.CoolantTempC = lerp(110, 112, (t-100)/18)
	s.OilPressureKPA = lerp(100, 80, (t-100)/18)
	s.BatteryV = 14.0

	switch {
	case t < 102:
		u := (t - 100) / 2
		s.ScenarioPhase = "emergency_decel_6th"
		s.Gear = "6"
		s.SpeedKPH = lerp(150, 110, u)
		s.RPM = lerp(3000, 2200, u)
		s.EngineLoadPercent = lerp(35, 20, u)
	case t < 106:
		u := (t - 102) / 4
		s.ScenarioPhase = "downshift_6_to_5"
		s.Gear = "5"
		s.SpeedKPH = lerp(110, 80, u)
		s.RPM = downshiftRPM(u, 4100, 3000)
		s.EngineLoadPercent = 25
	case t < 111:
		u := (t - 106) / 5
		s.ScenarioPhase = "downshift_5_to_4"
		s.Gear = "4"
		s.SpeedKPH = lerp(80, 55, u)
		s.RPM = downshiftRPM(u, 3900, 2800)
		s.EngineLoadPercent = 22
	case t < 117:
		u := (t - 111) / 6
		s.ScenarioPhase = "downshift_4_to_3"
		s.Gear = "3"
		s.SpeedKPH = lerp(55, 45, u)
		s.RPM = downshiftRPM(u, 3600, 3200)
		s.EngineLoadPercent = 20
	default:
		u := (t - 117)
		s.ScenarioPhase = "downshift_3_to_2_begins"
		s.Gear = "2"
		s.SpeedKPH = lerp(45, 43, u)
		s.RPM = lerp(4800, 5200, u)
		s.EngineLoadPercent = 18
	}

	return s
}

func thrownRodSample(s RaceDemoSample, t float64) RaceDemoSample {
	u := clamp01((t - 118) / 12)
	s.ScenarioPhase = "catastrophic_engine_failure"
	s.EngineOn = false
	s.EngineFailed = true
	s.FailureCode = RaceDemoFailureThrownRod
	s.Gear = "2"
	s.SpeedKPH = lerp(43, 0, u)
	s.RPM = 0
	s.ThrottlePercent = 0
	s.BrakePercent = lerp(85, 100, u)
	s.OilTempC = lerp(136, 130, u)
	s.CoolantTempC = lerp(112, 108, u)
	s.OilPressureKPA = 0
	s.EngineLoadPercent = 0
	s.BatteryV = lerp(12.3, 12.0, u)
	s.WarningLevel = RaceDemoWarningCritical
	s.WarningCode = RaceDemoCodeCatastrophicEngineFailure
	s.StatusMessage = RaceDemoStatusEngineFailureThrownRod
	s.RequiresReset = true
	return s
}

func downshiftRPM(u, peak, end float64) float64 {
	if u < 0.25 {
		return lerp(2600, peak, u/0.25)
	}
	return lerp(peak, end, (u-0.25)/0.75)
}

func triangularPulse(t, start, peak, end, lowStart, high, lowEnd float64) float64 {
	if t <= peak {
		return lerp(lowStart, high, (t-start)/(peak-start))
	}
	return lerp(high, lowEnd, (t-peak)/(end-peak))
}

func lerp(a, b, u float64) float64 {
	return a + (b-a)*clamp01(u)
}

func clamp01(u float64) float64 {
	if u < 0 {
		return 0
	}
	if u > 1 {
		return 1
	}
	return u
}

func roundSample(s RaceDemoSample) RaceDemoSample {
	s.SpeedKPH = round1(s.SpeedKPH)
	s.RPM = round0(s.RPM)
	s.ThrottlePercent = round1(s.ThrottlePercent)
	s.BrakePercent = round1(s.BrakePercent)
	s.OilTempC = round1(s.OilTempC)
	s.CoolantTempC = round1(s.CoolantTempC)
	s.OilPressureKPA = round1(s.OilPressureKPA)
	s.EngineLoadPercent = round1(s.EngineLoadPercent)
	s.BatteryV = round1(s.BatteryV)
	return s
}

func round0(v float64) float64 { return math.Round(v) }
func round1(v float64) float64 { return math.Round(v*10) / 10 }

func boolValue(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func warningLevelValue(level string) float64 {
	switch level {
	case RaceDemoWarningWarning:
		return 1
	case RaceDemoWarningCritical:
		return 2
	default:
		return 0
	}
}

func gearValue(gear string) float64 {
	switch gear {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	default:
		return 0
	}
}
