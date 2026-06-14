package v3dashboard

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const (
	PartKindBackground   = "background"
	PartKindImage        = "image"
	PartKindCharacter    = "character"
	PartKindDecimalPoint = "decimal_point"
	PartKindState        = "state"
	PartKindForeground   = "foreground"
)

// Runtime owns selected v3 dashboard scene state. It consumes sensor state/events
// produced by the sensor runtime; it never reads endpoint or OBD code directly.
type Runtime struct {
	dashboards []Dashboard
	states     map[string]sensors.SensorState
	signatures map[string]string
}

type Dashboard struct {
	ID      string
	Config  v3config.DashboardConfig
	Assets  *v3assets.Registry
}

type Scene struct {
	DashboardID string
	Display     string
	Size        v3config.SizeConfig
	Widgets     []Widget
}

type Widget struct {
	ID       string
	Type     string
	SensorID string
	AssetID  string
	Position []int
	Status   string
	Text     string
	Parts    []Part
	Error    string
}

type Part struct {
	Kind      string
	AssetPath string
	Slot      int
	Character string
	State     string
}

// NewRuntime builds the selected-dashboard runtime from an already resolved
// RuntimePlan. Unselected dashboards are not present in the plan and stay inert.
func NewRuntime(plan v3config.RuntimePlan, registry *v3assets.Registry) (*Runtime, error) {
	if registry == nil {
		return nil, fmt.Errorf("v3 dashboard runtime requires an asset registry")
	}
	dashboards := make([]Dashboard, 0, len(plan.Dashboards))
	for _, selected := range plan.Dashboards {
		dashboard := Dashboard{ID: selected.ID, Config: selected.Config, Assets: registry}
		if err := dashboard.Validate(); err != nil {
			return nil, err
		}
		dashboards = append(dashboards, dashboard)
	}
	return &Runtime{
		dashboards: dashboards,
		states:     map[string]sensors.SensorState{},
		signatures: map[string]string{},
	}, nil
}

func (r *Runtime) DashboardCount() int {
	if r == nil {
		return 0
	}
	return len(r.dashboards)
}

func (r *Runtime) SetState(state sensors.SensorState) {
	if r == nil || state.ID == "" {
		return
	}
	r.states[state.ID] = state
}

func (r *Runtime) ApplyEvent(event sensors.SensorEvent) ([]Scene, bool, error) {
	if r == nil {
		return nil, false, fmt.Errorf("v3 dashboard runtime is nil")
	}
	if event.SensorID != "" {
		state := event.State
		if state.ID == "" {
			state.ID = event.SensorID
		}
		r.states[event.SensorID] = state
	}

	scenes, err := r.Snapshot()
	if err != nil {
		return nil, false, err
	}
	changed := false
	for _, scene := range scenes {
		signature := sceneSignature(scene)
		if r.signatures[scene.DashboardID] != signature {
			r.signatures[scene.DashboardID] = signature
			changed = true
		}
	}
	if !changed {
		return nil, false, nil
	}
	return scenes, true, nil
}

func (r *Runtime) Snapshot() ([]Scene, error) {
	if r == nil {
		return nil, fmt.Errorf("v3 dashboard runtime is nil")
	}
	scenes := make([]Scene, 0, len(r.dashboards))
	for _, dashboard := range r.dashboards {
		scene, err := dashboard.Render(r.states)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, scene)
	}
	return scenes, nil
}

func (d Dashboard) Validate() error {
	if d.Assets == nil {
		return fmt.Errorf("dashboard %q requires an asset registry", d.ID)
	}
	for _, widget := range d.Config.Widgets {
		if strings.TrimSpace(widget.ID) == "" {
			return fmt.Errorf("dashboard %q has widget with empty id", d.ID)
		}
		if _, err := d.renderWidget(widget, map[string]sensors.SensorState{}); err != nil && !isMissingSensorOnly(err) {
			return err
		}
	}
	return nil
}

func (d Dashboard) Render(states map[string]sensors.SensorState) (Scene, error) {
	scene := Scene{DashboardID: d.ID, Display: d.Config.Display, Size: d.Config.Size}
	for _, configWidget := range d.Config.Widgets {
		widget, err := d.renderWidget(configWidget, states)
		if err != nil {
			return Scene{}, err
		}
		scene.Widgets = append(scene.Widgets, widget)
	}
	return scene, nil
}

func (d Dashboard) renderWidget(configWidget v3config.WidgetConfig, states map[string]sensors.SensorState) (Widget, error) {
	widget := Widget{
		ID:       configWidget.ID,
		Type:     configWidget.Type,
		SensorID: configWidget.Sensor,
		AssetID:  configWidget.Asset,
		Position: append([]int(nil), configWidget.Position...),
	}

	switch configWidget.Type {
	case v3config.WidgetTypeImage:
		set, ok := d.Assets.ImageSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.image_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		widget.Status = sensors.StatusOK
		widget.Parts = appendImageSetParts(widget.Parts, set)
	case v3config.WidgetTypeDigitDisplay:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.DigitSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.digit_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		if state.Status != sensors.StatusOK {
			widget.Parts = appendDigitLayerParts(widget.Parts, set, configWidget.Digits)
			return widget, nil
		}
		text := formatValue(configWidget.Format, state.Value)
		parts, err := digitParts(set, text, configWidget.Digits, d.ID, configWidget.ID, configWidget.Asset)
		if err != nil {
			return Widget{}, err
		}
		widget.Text = text
		widget.Parts = parts
	case v3config.WidgetTypeIndicator:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.IndicatorSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.indicator_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		indicatorState := indicatorStateFor(state)
		widget.Parts = appendIndicatorParts(widget.Parts, set, indicatorState)
	case v3config.WidgetTypeBarDisplay, v3config.WidgetTypeFrameGauge:
		return Widget{}, fmt.Errorf("dashboard %q widget %q type %q belongs to a later v3 slice", d.ID, configWidget.ID, configWidget.Type)
	default:
		return Widget{}, fmt.Errorf("dashboard %q widget %q type %q is not supported", d.ID, configWidget.ID, configWidget.Type)
	}

	return widget, nil
}

func stateForWidget(widget v3config.WidgetConfig, states map[string]sensors.SensorState) sensors.SensorState {
	state, ok := states[widget.Sensor]
	if !ok {
		return sensors.SensorState{ID: widget.Sensor, Status: sensors.StatusMissingUnsupported}
	}
	if state.ID == "" {
		state.ID = widget.Sensor
	}
	if state.Status == "" {
		state.Status = sensors.StatusUnknown
	}
	return state
}

func formatValue(format string, value float64) string {
	if strings.TrimSpace(format) == "" {
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
	return fmt.Sprintf(format, value)
}

func digitParts(set v3assets.DigitSet, text string, slots int, dashboardID, widgetID, assetID string) ([]Part, error) {
	if slots <= 0 {
		return nil, fmt.Errorf("dashboard %q widget %q digits must be greater than zero", dashboardID, widgetID)
	}
	characters, decimalSlots, err := splitDigitText(text, slots)
	if err != nil {
		return nil, fmt.Errorf("dashboard %q widget %q: %w", dashboardID, widgetID, err)
	}

	parts := make([]Part, 0, len(characters)*3+len(decimalSlots))
	for slot, ch := range characters {
		if set.Background != nil {
			parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path, Slot: slot})
		}
		if ch != " " {
			asset, ok := set.Characters[ch]
			if !ok {
				return nil, fmt.Errorf("dashboard %q widget %q digit set %q has no character asset for %q", dashboardID, widgetID, assetID, ch)
			}
			parts = append(parts, Part{Kind: PartKindCharacter, AssetPath: asset.Path, Slot: slot, Character: ch})
		}
		if set.Foreground != nil {
			parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path, Slot: slot})
		}
	}
	for _, slot := range decimalSlots {
		if set.DecimalPoint == nil {
			return nil, fmt.Errorf("dashboard %q widget %q digit set %q requires decimal_point", dashboardID, widgetID, assetID)
		}
		parts = append(parts, Part{Kind: PartKindDecimalPoint, AssetPath: set.DecimalPoint.Path, Slot: slot})
	}
	return parts, nil
}

func splitDigitText(text string, slots int) ([]string, []int, error) {
	characters := make([]string, 0, len(text))
	decimalSlots := []int{}
	lastSlot := -1
	for _, r := range text {
		ch := string(r)
		if ch == "." {
			if lastSlot < 0 {
				return nil, nil, fmt.Errorf("decimal separator has no preceding digit slot")
			}
			decimalSlots = append(decimalSlots, lastSlot)
			continue
		}
		characters = append(characters, ch)
		lastSlot = len(characters) - 1
	}
	if len(characters) > slots {
		return nil, nil, fmt.Errorf("formatted output %q needs %d character slots, dashboard config allows %d", text, len(characters), slots)
	}
	padded := make([]string, slots)
	padding := slots - len(characters)
	for i := 0; i < padding; i++ {
		padded[i] = " "
	}
	copy(padded[padding:], characters)
	if padding > 0 {
		for i := range decimalSlots {
			decimalSlots[i] += padding
		}
	}
	return padded, decimalSlots, nil
}

func appendDigitLayerParts(parts []Part, set v3assets.DigitSet, slots int) []Part {
	if slots <= 0 {
		return parts
	}
	for slot := 0; slot < slots; slot++ {
		if set.Background != nil {
			parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path, Slot: slot})
		}
		if set.Foreground != nil {
			parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path, Slot: slot})
		}
	}
	return parts
}

func indicatorStateFor(state sensors.SensorState) string {
	if state.Status != sensors.StatusOK {
		return v3assets.IndicatorStateUnknown
	}
	if state.Value != 0 {
		return v3assets.IndicatorStateOn
	}
	return v3assets.IndicatorStateOff
}

func appendImageSetParts(parts []Part, set v3assets.ImageSet) []Part {
	if set.Background != nil {
		parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path})
	}
	if set.Image != nil {
		parts = append(parts, Part{Kind: PartKindImage, AssetPath: set.Image.Path})
	}
	if set.Foreground != nil {
		parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
	}
	return parts
}

func appendIndicatorParts(parts []Part, set v3assets.IndicatorSet, state string) []Part {
	if set.Background != nil {
		parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path})
	}
	asset := set.States[state]
	parts = append(parts, Part{Kind: PartKindState, AssetPath: asset.Path, State: state})
	if set.Foreground != nil {
		parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
	}
	return parts
}

func isMissingSensorOnly(err error) bool {
	return err != nil && strings.Contains(err.Error(), string(sensors.StatusMissingUnsupported))
}

func sceneSignature(scene Scene) string {
	var b strings.Builder
	b.WriteString(scene.DashboardID)
	b.WriteString("|")
	for _, widget := range scene.Widgets {
		b.WriteString(widget.ID)
		b.WriteString(":")
		b.WriteString(widget.Type)
		b.WriteString(":")
		b.WriteString(widget.Status)
		b.WriteString(":")
		b.WriteString(widget.Text)
		b.WriteString(":")
		for _, part := range widget.Parts {
			b.WriteString(part.Kind)
			b.WriteString("@")
			b.WriteString(strconv.Itoa(part.Slot))
			b.WriteString("=")
			b.WriteString(part.AssetPath)
			b.WriteString("#")
			b.WriteString(part.Character)
			b.WriteString("#")
			b.WriteString(part.State)
			b.WriteString(";")
		}
		b.WriteString("|")
	}
	return b.String()
}

func sortedWidgetIDs(widgets []Widget) []string {
	ids := make([]string, 0, len(widgets))
	for _, widget := range widgets {
		ids = append(ids, widget.ID)
	}
	sort.Strings(ids)
	return ids
}
