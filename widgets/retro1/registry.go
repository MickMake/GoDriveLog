package retro1

import "github.com/MickMake/GoDriveLog/widgets/model"

const (
	StyleRamp1 = "retro1_ramp1"
	StyleRamp2 = "retro1_ramp2"
	StyleRamp3 = "retro1_ramp3"

	Style7Seg1 = "retro1_7seg1"
	Style7Seg2 = "retro1_7seg2"
	Style7Seg3 = "retro1_7seg3"

	// StyleSeg7Legacy keeps the original handoff ID working while the cleaner
	// retro1_7seg1 name becomes the primary ID.
	StyleSeg7Legacy = "retro1_seg7_1"
)

func NewRamp1(cfg model.GaugeConfig) model.Widget { return newRamp(cfg, StyleRamp1, 1) }
func NewRamp2(cfg model.GaugeConfig) model.Widget { return newRamp(cfg, StyleRamp2, 2) }
func NewRamp3(cfg model.GaugeConfig) model.Widget { return newRamp(cfg, StyleRamp3, 3) }

func New7Seg1(cfg model.GaugeConfig) model.Widget      { return newSevenSeg(cfg, Style7Seg1, 1) }
func New7Seg2(cfg model.GaugeConfig) model.Widget      { return newSevenSeg(cfg, Style7Seg2, 2) }
func New7Seg3(cfg model.GaugeConfig) model.Widget      { return newSevenSeg(cfg, Style7Seg3, 3) }
func New7SegLegacy(cfg model.GaugeConfig) model.Widget { return newSevenSeg(cfg, StyleSeg7Legacy, 3) }
