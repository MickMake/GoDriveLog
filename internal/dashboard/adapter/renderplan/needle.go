package renderplan

import "github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"

type NeedleLikePart struct {
	X      float64
	Y      float64
	Scale  float64
	Needle bool
	Shadow bool
	Angle  float64
	Alpha  float64
	PivotX float64
	PivotY float64
}

func BuildNeedleLikePart(part v3dashboard.Part, assetWidth, assetHeight int, baseX, baseY, gaugeWidth, gaugeHeight, widgetScale float64) (NeedleLikePart, bool) {
	if part.Kind != v3dashboard.PartKindNeedle &&
		part.Kind != v3dashboard.PartKindNeedleShadow &&
		part.Kind != v3dashboard.PartKindNeedleMin &&
		part.Kind != v3dashboard.PartKindNeedleMax {
		return NeedleLikePart{}, false
	}

	faceX := baseX + part.FacePivot.X*gaugeWidth*widgetScale
	faceY := baseY + part.FacePivot.Y*gaugeHeight*widgetScale
	if len(part.Position) >= 2 {
		faceX += float64(part.Position[0]) * widgetScale
		faceY += float64(part.Position[1]) * widgetScale
	}

	return NeedleLikePart{
		X:      faceX,
		Y:      faceY,
		Scale:  widgetScale,
		Needle: true,
		Shadow: part.Kind == v3dashboard.PartKindNeedleShadow,
		Angle:  part.Angle,
		Alpha:  part.Alpha,
		PivotX: part.NeedlePivot.X * float64(assetWidth),
		PivotY: part.NeedlePivot.Y * float64(assetHeight),
	}, true
}
