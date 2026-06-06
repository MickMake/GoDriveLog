package retro1

import "github.com/MickMake/GoDriveLog/widgets/model"

func newSevenSeg(cfg model.GaugeConfig, style string, level int) model.Widget {
	return newRamp(cfg, style, level)
}
