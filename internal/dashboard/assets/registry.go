package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/config"
)

const (
	TypeImage    = config.DashboardAssetImage
	TypeFrameSet = config.DashboardAssetFrameSet
	TypeCharset  = config.DashboardAssetCharset
)

type Registry struct {
	root   string
	assets map[string]Asset
}

type Asset struct {
	ID     string
	Type   string
	Path   string
	Data   []byte
	Frames []Frame
	Glyphs map[string]Glyph
}

type Frame struct {
	Index int
	Path  string
	Data  []byte
}

type Glyph struct {
	Char string
	Path string
	Data []byte
}

func Load(dashboard config.DashboardConfig, configPath string) (*Registry, error) {
	root := resolveRoot(dashboard.AssetRoot, configPath)
	registry := &Registry{root: root, assets: map[string]Asset{}}
	for _, assetConfig := range dashboard.Assets {
		asset, err := loadAsset(root, assetConfig)
		if err != nil {
			return nil, err
		}
		registry.assets[asset.ID] = asset
	}
	return registry, nil
}

func (r *Registry) Root() string {
	return r.root
}

func (r *Registry) Get(id string) (Asset, bool) {
	asset, ok := r.assets[id]
	return asset, ok
}

func (r *Registry) MustGet(id string) (Asset, error) {
	asset, ok := r.Get(id)
	if !ok {
		return Asset{}, fmt.Errorf("dashboard asset %q is not loaded", id)
	}
	return asset, nil
}

func (r *Registry) IDs() []string {
	ids := make([]string, 0, len(r.assets))
	for id := range r.assets {
		ids = append(ids, id)
	}
	return ids
}

func resolveRoot(assetRoot string, configPath string) string {
	base := "."
	if configPath != "" {
		base = filepath.Dir(configPath)
	}
	if assetRoot == "" {
		return filepath.Clean(base)
	}
	if filepath.IsAbs(assetRoot) {
		return filepath.Clean(assetRoot)
	}
	return filepath.Clean(filepath.Join(base, assetRoot))
}

func loadAsset(root string, assetConfig config.DashboardAssetConfig) (Asset, error) {
	switch assetConfig.Type {
	case TypeImage:
		path, data, err := readAssetFile(root, assetConfig.Path, fmt.Sprintf("asset %q image", assetConfig.ID))
		if err != nil {
			return Asset{}, err
		}
		return Asset{ID: assetConfig.ID, Type: assetConfig.Type, Path: path, Data: data}, nil
	case TypeFrameSet:
		frames, err := loadFrames(root, assetConfig)
		if err != nil {
			return Asset{}, err
		}
		return Asset{ID: assetConfig.ID, Type: assetConfig.Type, Frames: frames}, nil
	case TypeCharset:
		glyphs, err := loadGlyphs(root, assetConfig)
		if err != nil {
			return Asset{}, err
		}
		return Asset{ID: assetConfig.ID, Type: assetConfig.Type, Glyphs: glyphs}, nil
	default:
		return Asset{}, fmt.Errorf("asset %q type %q is not supported", assetConfig.ID, assetConfig.Type)
	}
}

func loadFrames(root string, assetConfig config.DashboardAssetConfig) ([]Frame, error) {
	framePaths := append([]string(nil), assetConfig.Frames...)
	if assetConfig.Pattern != "" {
		generated, err := expandFramePattern(assetConfig.Pattern, assetConfig.FrameCount)
		if err != nil {
			return nil, fmt.Errorf("asset %q frame_set: %w", assetConfig.ID, err)
		}
		framePaths = generated
	}
	if assetConfig.FrameCount > 0 && len(framePaths) != assetConfig.FrameCount {
		return nil, fmt.Errorf("asset %q frame_set expected %d frames, got %d", assetConfig.ID, assetConfig.FrameCount, len(framePaths))
	}
	frames := make([]Frame, 0, len(framePaths))
	for i, framePath := range framePaths {
		path, data, err := readAssetFile(root, framePath, fmt.Sprintf("asset %q frame %d", assetConfig.ID, i))
		if err != nil {
			return nil, err
		}
		frames = append(frames, Frame{Index: i, Path: path, Data: data})
	}
	return frames, nil
}

func loadGlyphs(root string, assetConfig config.DashboardAssetConfig) (map[string]Glyph, error) {
	glyphs := map[string]Glyph{}
	for char, glyphPath := range assetConfig.Glyphs {
		if char == "" {
			return nil, fmt.Errorf("asset %q charset glyph key must not be empty", assetConfig.ID)
		}
		path, data, err := readAssetFile(root, glyphPath, fmt.Sprintf("asset %q glyph %q", assetConfig.ID, char))
		if err != nil {
			return nil, err
		}
		glyphs[char] = Glyph{Char: char, Path: path, Data: data}
	}
	return glyphs, nil
}

func readAssetFile(root string, assetPath string, label string) (string, []byte, error) {
	if assetPath == "" {
		return "", nil, fmt.Errorf("%s path must not be empty", label)
	}
	if isRemotePath(assetPath) {
		return "", nil, fmt.Errorf("%s path %q must be local", label, assetPath)
	}
	resolved := assetPath
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(root, assetPath)
	}
	resolved = filepath.Clean(resolved)
	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", nil, fmt.Errorf("%s path %q could not be loaded: %w", label, resolved, err)
	}
	return resolved, data, nil
}

func isRemotePath(path string) bool {
	return strings.Contains(path, "://")
}

func expandFramePattern(pattern string, count int) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("frame_count must be positive")
	}
	paths := make([]string, 0, count)
	for i := 0; i < count; i++ {
		path, err := expandFramePath(pattern, i)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func expandFramePath(pattern string, index int) (string, error) {
	if strings.Contains(pattern, "{index}") {
		return strings.ReplaceAll(pattern, "{index}", strconv.Itoa(index)), nil
	}
	marker := "{index:"
	start := strings.Index(pattern, marker)
	if start < 0 {
		return "", fmt.Errorf("pattern %q must contain {index} or {index:0N}", pattern)
	}
	end := strings.Index(pattern[start:], "}")
	if end < 0 {
		return "", fmt.Errorf("pattern %q has an unterminated index marker", pattern)
	}
	end += start
	spec := pattern[start+len(marker) : end]
	if len(spec) < 2 || spec[0] != '0' {
		return "", fmt.Errorf("pattern %q index marker must use zero padding like {index:03}", pattern)
	}
	width, err := strconv.Atoi(spec[1:])
	if err != nil || width <= 0 {
		return "", fmt.Errorf("pattern %q index marker has invalid width", pattern)
	}
	value := fmt.Sprintf("%0*d", width, index)
	return pattern[:start] + value + pattern[end+1:], nil
}
