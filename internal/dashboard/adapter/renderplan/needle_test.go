package renderplan

import (
	"testing"

	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestBuildNeedleLikePartTreatsNeedleMinAsRotatingNeedlePart(t *testing.T) {
	rendered := buildSingleNeedleLikePart(v3dashboard.PartKindNeedleMin)

	if !rendered.Needle {
		t.Fatal("expected needle_min to use rotating needle path")
	}
	if rendered.Shadow {
		t.Fatal("expected needle_min not to be treated as a shadow part")
	}
	assertNeedleLikeGeometry(t, rendered)
}

func TestBuildNeedleLikePartTreatsNeedleMaxAsRotatingNeedlePart(t *testing.T) {
	rendered := buildSingleNeedleLikePart(v3dashboard.PartKindNeedleMax)

	if !rendered.Needle {
		t.Fatal("expected needle_max to use rotating needle path")
	}
	if rendered.Shadow {
		t.Fatal("expected needle_max not to be treated as a shadow part")
	}
	assertNeedleLikeGeometry(t, rendered)
}

func TestBuildNeedleLikePartKeepsLiveNeedleAndShadowBehavior(t *testing.T) {
	needle := buildSingleNeedleLikePart(v3dashboard.PartKindNeedle)
	if !needle.Needle {
		t.Fatal("expected live needle to use rotating needle path")
	}
	if needle.Shadow {
		t.Fatal("expected live needle not to be treated as a shadow part")
	}
	assertNeedleLikeGeometry(t, needle)

	shadow := buildSingleNeedleLikePart(v3dashboard.PartKindNeedleShadow)
	if !shadow.Needle {
		t.Fatal("expected needle shadow to use rotating needle path")
	}
	if !shadow.Shadow {
		t.Fatal("expected needle shadow to stay marked as a shadow part")
	}
	assertNeedleLikeGeometry(t, shadow)
}

func buildSingleNeedleLikePart(kind string) NeedleLikePart {
	part := v3dashboard.Part{
		Kind:        kind,
		Layer:       kind,
		Angle:       27.5,
		FacePivot:   v3gauges.Point{X: 0.25, Y: 0.75},
		NeedlePivot: v3gauges.Point{X: 0.4, Y: 0.6},
		Position:    []int{3, 4},
		Alpha:       0.35,
	}
	rendered, ok := BuildNeedleLikePart(part, 20, 40, 10, 20, 100, 80, 2)
	if !ok {
		panic("expected needle-like part to build")
	}
	return rendered
}

func assertNeedleLikeGeometry(t *testing.T, rendered NeedleLikePart) {
	t.Helper()

	if rendered.Angle != 27.5 {
		t.Fatalf("angle = %v, want 27.5", rendered.Angle)
	}
	if rendered.X != 66 {
		t.Fatalf("x = %v, want 66", rendered.X)
	}
	if rendered.Y != 148 {
		t.Fatalf("y = %v, want 148", rendered.Y)
	}
	if rendered.PivotX != 8 {
		t.Fatalf("pivotX = %v, want 8", rendered.PivotX)
	}
	if rendered.PivotY != 24 {
		t.Fatalf("pivotY = %v, want 24", rendered.PivotY)
	}
	if rendered.Scale != 2 {
		t.Fatalf("scale = %v, want 2", rendered.Scale)
	}
}
