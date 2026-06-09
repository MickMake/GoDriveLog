package sensors

import (
	"context"
	"testing"
	"time"
)

func TestRaceDemoScenarioMilestones(t *testing.T) {
	scenario := NewRaceDemoScenario()

	idle := scenario.SampleAt(3 * time.Second)
	if !idle.EngineOn || idle.EngineFailed || idle.SpeedKPH != 0 || idle.RPM < 800 || idle.RPM > 900 {
		t.Fatalf("idle sample = %#v, want engine on, stopped, about 850 rpm", idle)
	}

	rev := scenario.SampleAt(9 * time.Second)
	if rev.SpeedKPH != 0 || rev.RPM < 4000 || rev.ThrottlePercent < 80 || rev.WarningLevel != RaceDemoWarningNone {
		t.Fatalf("stationary rev sample = %#v, want high rpm/throttle at zero speed and no warning", rev)
	}

	burnout := scenario.SampleAt(18 * time.Second)
	if burnout.Gear != "1" || burnout.SpeedKPH > 8 || burnout.RPM < 4500 || burnout.ThrottlePercent != 100 {
		t.Fatalf("burnout sample = %#v, want high rpm/throttle with low speed", burnout)
	}

	oilWarning := scenario.SampleAt(48 * time.Second)
	if oilWarning.Gear != "4" || oilWarning.WarningLevel != RaceDemoWarningWarning || oilWarning.WarningCode != RaceDemoCodeOilTempHigh || oilWarning.OilTempC < 115 {
		t.Fatalf("oil warning sample = %#v, want oil warning just after 4th gear", oilWarning)
	}

	cruise := scenario.SampleAt(75 * time.Second)
	if cruise.Gear != "6" || cruise.SpeedKPH < 145 || cruise.SpeedKPH > 152 || cruise.WarningLevel != RaceDemoWarningWarning {
		t.Fatalf("cruise sample = %#v, want 6th gear around 150 km/h with warning", cruise)
	}

	critical := scenario.SampleAt(95 * time.Second)
	if critical.WarningLevel != RaceDemoWarningCritical || critical.WarningCode != RaceDemoCodeOilTempCritical || critical.OilTempC < 130 {
		t.Fatalf("critical sample = %#v, want critical oil temperature", critical)
	}

	downshift := scenario.SampleAt(111 * time.Second)
	if downshift.Gear != "3" || downshift.SpeedKPH > 60 || downshift.RPM < 2500 || downshift.WarningLevel != RaceDemoWarningCritical {
		t.Fatalf("downshift sample = %#v, want 4th to 3rd downshift under critical warning", downshift)
	}

	failure := scenario.SampleAt(118 * time.Second)
	if failure.EngineOn || !failure.EngineFailed || failure.FailureCode != RaceDemoFailureThrownRod || failure.RPM != 0 || failure.OilPressureKPA != 0 || failure.SpeedKPH <= 0 {
		t.Fatalf("failure sample = %#v, want thrown rod with rpm/oil pressure zero while moving", failure)
	}
	if failure.WarningLevel != RaceDemoWarningCritical || failure.WarningCode != RaceDemoCodeCatastrophicEngineFailure || !failure.RequiresReset {
		t.Fatalf("failure alert = %#v, want latched catastrophic critical failure", failure)
	}

	coasting := scenario.SampleAt(125 * time.Second)
	if coasting.RPM != 0 || coasting.SpeedKPH <= 0 || !coasting.EngineFailed || !coasting.RequiresReset {
		t.Fatalf("coasting sample = %#v, want engine failed, rpm zero, speed still positive", coasting)
	}

	stopped := scenario.SampleAt(132 * time.Second)
	if stopped.SpeedKPH != 0 || stopped.RPM != 0 || stopped.EngineOn || !stopped.EngineFailed || stopped.WarningCode != RaceDemoCodeCatastrophicEngineFailure || !stopped.RequiresReset {
		t.Fatalf("stopped sample = %#v, want stopped with critical failure latched", stopped)
	}
}

func TestRaceDemoReaderReturnsPIDValues(t *testing.T) {
	reader := NewRaceDemoReaderAt(time.Now())

	rpm, unit, err := reader.ReadAt("010C", 118*time.Second)
	if err != nil {
		t.Fatalf("ReadAt rpm returned error: %v", err)
	}
	if rpm != 0 || unit != "rpm" {
		t.Fatalf("ReadAt rpm = %v %q, want 0 rpm after thrown rod", rpm, unit)
	}

	speed, unit, err := reader.ReadAt("010D", 125*time.Second)
	if err != nil {
		t.Fatalf("ReadAt speed returned error: %v", err)
	}
	if speed <= 0 || unit != "km/h" {
		t.Fatalf("ReadAt speed = %v %q, want positive km/h while coasting", speed, unit)
	}

	pressure, unit, err := reader.ReadAt("DEMO_OIL_PRESSURE", 118*time.Second)
	if err != nil {
		t.Fatalf("ReadAt oil pressure returned error: %v", err)
	}
	if pressure != 0 || unit != "kPa" {
		t.Fatalf("ReadAt oil pressure = %v %q, want 0 kPa", pressure, unit)
	}

	failed, unit, err := reader.ReadAt("DEMO_ENGINE_FAILED", 118*time.Second)
	if err != nil {
		t.Fatalf("ReadAt engine failed returned error: %v", err)
	}
	if failed != 1 || unit != "bool" {
		t.Fatalf("ReadAt engine failed = %v %q, want 1 bool", failed, unit)
	}

	if _, _, err := reader.Read(context.Background(), "DEMO_UNKNOWN"); err == nil {
		t.Fatal("Read with unknown PID succeeded, want error")
	}
}
