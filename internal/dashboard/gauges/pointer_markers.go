package gauges

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type PointerMarkersConfig struct {
	Max     bool           `yaml:"max,omitempty"`
	Min     bool           `yaml:"min,omitempty"`
	Average bool           `yaml:"average,omitempty"`
	Window  *time.Duration `yaml:"window,omitempty"`
}

func (c *PointerMarkersConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("realism pointer_markers must be a mapping")
	}

	var decoded PointerMarkersConfig
	for index := 0; index+1 < len(node.Content); index += 2 {
		keyNode := node.Content[index]
		valueNode := node.Content[index+1]
		switch keyNode.Value {
		case "max":
			value, err := decodePointerMarkerBool("max", valueNode)
			if err != nil {
				return err
			}
			decoded.Max = value
		case "min":
			value, err := decodePointerMarkerBool("min", valueNode)
			if err != nil {
				return err
			}
			decoded.Min = value
		case "average":
			value, err := decodePointerMarkerBool("average", valueNode)
			if err != nil {
				return err
			}
			decoded.Average = value
		case "window":
			value, err := decodePointerMarkerWindow(valueNode)
			if err != nil {
				return err
			}
			decoded.Window = &value
		default:
			return fmt.Errorf("realism pointer_markers field %q is not supported", keyNode.Value)
		}
	}

	*c = decoded
	return nil
}

func (c *PointerMarkersConfig) Enabled() bool {
	return c != nil && (c.Max || c.Min || c.Average)
}

type PointerMarkerValueState struct {
	Set                bool
	NormalizedPosition float64
	RecordedAt         time.Time
}

type PointerMarkerSample struct {
	NormalizedPosition float64
	RecordedAt         time.Time
}

type PointerMarkerState struct {
	LocalDayKey string
	Min         PointerMarkerValueState
	Max         PointerMarkerValueState
	Average     PointerMarkerValueState
	Samples     []PointerMarkerSample
}

func decodePointerMarkerBool(name string, node *yaml.Node) (bool, error) {
	if node.Kind != yaml.ScalarNode {
		return false, fmt.Errorf("realism pointer_markers %s must be a boolean", name)
	}

	var enabled bool
	if err := node.Decode(&enabled); err != nil {
		return false, fmt.Errorf("realism pointer_markers %s must be a boolean", name)
	}
	return enabled, nil
}

func decodePointerMarkerWindow(node *yaml.Node) (time.Duration, error) {
	if node.Kind != yaml.ScalarNode {
		return 0, fmt.Errorf("realism pointer_markers window must be a duration")
	}

	raw := strings.TrimSpace(node.Value)
	if raw == "" {
		return 0, fmt.Errorf("realism pointer_markers window must be a positive duration")
	}

	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("realism pointer_markers window %q is not a valid duration", raw)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("realism pointer_markers window must be greater than zero")
	}
	return duration, nil
}
