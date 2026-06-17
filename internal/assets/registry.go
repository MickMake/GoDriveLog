package assets

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

const (
	IndicatorStateOff     = "off"
	IndicatorStateOn      = "on"
	IndicatorStateUnknown = "unknown"
)

var requiredDigitCharacters = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

// Registry contains decoded v3 assets keyed by global asset family and ID.
// It is built once from resolved asset search paths so widgets can reuse
// decoded image data instead of loading or decoding in the hot render path.
type Registry struct {
	searchPaths   []string
	Images        map[string]ImageSet
	DigitSets     map[string]DigitSet
	BarSets       map[string]BarSet
	FrameSets     map[string]FrameSet
	IndicatorSets map[string]IndicatorSet
}

type ImageSet struct {
	ID         string
	Image      *ImageAsset
	Background *ImageAsset
	Foreground *ImageAsset
}

type DigitSet struct {
	ID           string
	Background   *ImageAsset
	Characters   map[string]ImageAsset
	DecimalPoint *ImageAsset
	Foreground   *ImageAsset
	Spacing      int
}

type BarSet struct {
	ID         string
	Background *ImageAsset
	Cells      map[string]ImageAsset
	Foreground *ImageAsset
	Spacing    int
}

type FrameSet struct {
	ID         string
	Background *ImageAsset
	Frames     map[int]ImageAsset
	First      int
	Last       int
	Foreground *ImageAsset
}

type IndicatorSet struct {
	ID         string
	Background *ImageAsset
	States     map[string]ImageAsset
	Foreground *ImageAsset
}

type ImageAsset struct {
	Path   string
	Data   []byte
	Image  image.Image
	Bounds image.Rectangle
}

// DefaultSearchPaths builds the standard v3 asset search path order.
// Relative asset paths are checked in this order:
//  1. config directory / vehicle ID
//  2. current working directory / vehicle ID
//  3. config directory
//  4. current working directory
func DefaultSearchPaths(configPath, vehicleID string) ([]string, error) {
	configDir := "."
	if strings.TrimSpace(configPath) != "" {
		configDir = filepath.Dir(configPath)
	}
	configDir, err := filepath.Abs(configDir)
	if err != nil {
		return nil, fmt.Errorf("config directory for %q could not be resolved: %w", configPath, err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("current working directory could not be resolved: %w", err)
	}

	vehicleID = strings.TrimSpace(vehicleID)
	candidates := []string{}
	if vehicleID != "" {
		candidates = append(candidates, filepath.Join(configDir, vehicleID), filepath.Join(pwd, vehicleID))
	}
	candidates = append(candidates, configDir, pwd)
	return cleanSearchPaths(candidates)
}

// Load builds the v3 asset registry for image-backed dashboard asset families.
// It keeps the original one-root call shape for tests and simple callers.
func Load(cfg v3config.AssetConfig, root string) (*Registry, error) {
	return LoadWithSearchPaths(cfg, []string{root})
}

// LoadWithSearchPaths builds the v3 asset registry using ordered asset search
// paths. The first path containing a requested relative asset wins.
func LoadWithSearchPaths(cfg v3config.AssetConfig, searchPaths []string) (*Registry, error) {
	paths, err := cleanSearchPaths(searchPaths)
	if err != nil {
		return nil, err
	}

	registry := &Registry{
		searchPaths:   paths,
		Images:        make(map[string]ImageSet, len(cfg.ImageSets)),
		DigitSets:     make(map[string]DigitSet, len(cfg.DigitSets)),
		BarSets:       make(map[string]BarSet, len(cfg.BarSets)),
		FrameSets:     make(map[string]FrameSet, len(cfg.FrameSets)),
		IndicatorSets: make(map[string]IndicatorSet, len(cfg.IndicatorSets)),
	}

	for _, id := range sortedKeys(cfg.ImageSets) {
		set, err := loadImageSet(paths, id, cfg.ImageSets[id])
		if err != nil {
			return nil, err
		}
		registry.Images[id] = set
	}
	for _, id := range sortedKeys(cfg.DigitSets) {
		set, err := loadDigitSet(paths, id, cfg.DigitSets[id])
		if err != nil {
			return nil, err
		}
		registry.DigitSets[id] = set
	}
	for _, id := range sortedKeys(cfg.BarSets) {
		set, err := loadBarSet(paths, id, cfg.BarSets[id])
		if err != nil {
			return nil, err
		}
		registry.BarSets[id] = set
	}
	for _, id := range sortedKeys(cfg.FrameSets) {
		set, err := loadFrameSet(paths, id, cfg.FrameSets[id])
		if err != nil {
			return nil, err
		}
		registry.FrameSets[id] = set
	}
	for _, id := range sortedKeys(cfg.IndicatorSets) {
		set, err := loadIndicatorSet(paths, id, cfg.IndicatorSets[id])
		if err != nil {
			return nil, err
		}
		registry.IndicatorSets[id] = set
	}

	return registry, nil
}

func (r *Registry) RepoRoot() string {
	if r == nil || len(r.searchPaths) == 0 {
		return ""
	}
	return r.searchPaths[0]
}

func (r *Registry) SearchPaths() []string {
	if r == nil {
		return nil
	}
	return append([]string(nil), r.searchPaths...)
}

func (r *Registry) ImageSet(id string) (ImageSet, bool) {
	if r == nil {
		return ImageSet{}, false
	}
	set, ok := r.Images[id]
	return set, ok
}

func (r *Registry) DigitSet(id string) (DigitSet, bool) {
	if r == nil {
		return DigitSet{}, false
	}
	set, ok := r.DigitSets[id]
	return set, ok
}

func (r *Registry) BarSet(id string) (BarSet, bool) {
	if r == nil {
		return BarSet{}, false
	}
	set, ok := r.BarSets[id]
	return set, ok
}

func (r *Registry) FrameSet(id string) (FrameSet, bool) {
	if r == nil {
		return FrameSet{}, false
	}
	set, ok := r.FrameSets[id]
	return set, ok
}

func (r *Registry) IndicatorSet(id string) (IndicatorSet, bool) {
	if r == nil {
		return IndicatorSet{}, false
	}
	set, ok := r.IndicatorSets[id]
	return set, ok
}

func loadImageSet(searchPaths []string, id string, cfg v3config.ImageSetConfig) (ImageSet, error) {
	set := ImageSet{ID: id}
	var err error
	if cfg.Image != "" {
		set.Image, err = loadOptionalImage(searchPaths, cfg.Image, fmt.Sprintf("assets.image_sets.%s.image", id))
		if err != nil {
			return ImageSet{}, err
		}
	}
	if cfg.Background != "" {
		set.Background, err = loadOptionalImage(searchPaths, cfg.Background, fmt.Sprintf("assets.image_sets.%s.background", id))
		if err != nil {
			return ImageSet{}, err
		}
	}
	if cfg.Foreground != "" {
		set.Foreground, err = loadOptionalImage(searchPaths, cfg.Foreground, fmt.Sprintf("assets.image_sets.%s.foreground", id))
		if err != nil {
			return ImageSet{}, err
		}
	}
	if set.Image == nil && set.Background == nil && set.Foreground == nil {
		return ImageSet{}, fmt.Errorf("assets.image_sets.%s must define image, background, or foreground", id)
	}
	return set, nil
}

func loadDigitSet(searchPaths []string, id string, cfg v3config.DigitSetConfig) (DigitSet, error) {
	set := DigitSet{ID: id, Characters: make(map[string]ImageAsset, len(cfg.Characters)), Spacing: cfg.Spacing}
	var err error
	if cfg.Background != "" {
		set.Background, err = loadOptionalImage(searchPaths, cfg.Background, fmt.Sprintf("assets.digit_sets.%s.background", id))
		if err != nil {
			return DigitSet{}, err
		}
	}
	if cfg.DecimalPoint != "" {
		set.DecimalPoint, err = loadOptionalImage(searchPaths, cfg.DecimalPoint, fmt.Sprintf("assets.digit_sets.%s.decimal_point", id))
		if err != nil {
			return DigitSet{}, err
		}
	}
	if cfg.Foreground != "" {
		set.Foreground, err = loadOptionalImage(searchPaths, cfg.Foreground, fmt.Sprintf("assets.digit_sets.%s.foreground", id))
		if err != nil {
			return DigitSet{}, err
		}
	}
	for _, ch := range requiredDigitCharacters {
		asset, err := loadRequiredImage(searchPaths, cfg.Characters[ch], fmt.Sprintf("assets.digit_sets.%s.characters.%s", id, ch))
		if err != nil {
			return DigitSet{}, err
		}
		set.Characters[ch] = asset
	}
	for _, ch := range sortedKeys(cfg.Characters) {
		if _, exists := set.Characters[ch]; exists {
			continue
		}
		asset, err := loadRequiredImage(searchPaths, cfg.Characters[ch], fmt.Sprintf("assets.digit_sets.%s.characters.%s", id, ch))
		if err != nil {
			return DigitSet{}, err
		}
		set.Characters[ch] = asset
	}
	return set, nil
}

func loadBarSet(searchPaths []string, id string, cfg v3config.BarSetConfig) (BarSet, error) {
	set := BarSet{ID: id, Cells: make(map[string]ImageAsset, len(cfg.Cells)), Spacing: cfg.Spacing}
	var err error
	if cfg.Background != "" {
		set.Background, err = loadOptionalImage(searchPaths, cfg.Background, fmt.Sprintf("assets.bar_sets.%s.background", id))
		if err != nil {
			return BarSet{}, err
		}
	}
	if cfg.Foreground != "" {
		set.Foreground, err = loadOptionalImage(searchPaths, cfg.Foreground, fmt.Sprintf("assets.bar_sets.%s.foreground", id))
		if err != nil {
			return BarSet{}, err
		}
	}
	if strings.TrimSpace(cfg.Cells["off"]) == "" {
		return BarSet{}, fmt.Errorf("assets.bar_sets.%s.cells.off path must not be empty", id)
	}
	var cellBounds image.Rectangle
	for _, cell := range sortedKeys(cfg.Cells) {
		asset, err := loadRequiredImage(searchPaths, cfg.Cells[cell], fmt.Sprintf("assets.bar_sets.%s.cells.%s", id, cell))
		if err != nil {
			return BarSet{}, err
		}
		if cellBounds.Empty() {
			cellBounds = asset.Bounds
		} else if asset.Bounds != cellBounds {
			return BarSet{}, fmt.Errorf("assets.bar_sets.%s.cells.%s dimensions %s must match cell dimensions %s", id, cell, asset.Bounds, cellBounds)
		}
		set.Cells[cell] = asset
	}
	return set, nil
}

func loadFrameSet(searchPaths []string, id string, cfg v3config.FrameSetConfig) (FrameSet, error) {
	if cfg.Frames.First > cfg.Frames.Last {
		return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.frames.first must be less than or equal to last", id)
	}
	if cfg.Frames.First < cfg.Frames.Last && !hasPrintfPlaceholder(cfg.Frames.Path) {
		return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.frames.path must contain a printf placeholder for a multi-frame range", id)
	}
	set := FrameSet{ID: id, Frames: make(map[int]ImageAsset, cfg.Frames.Last-cfg.Frames.First+1), First: cfg.Frames.First, Last: cfg.Frames.Last}
	var err error
	if cfg.Background != "" {
		set.Background, err = loadOptionalImage(searchPaths, cfg.Background, fmt.Sprintf("assets.frame_sets.%s.background", id))
		if err != nil {
			return FrameSet{}, err
		}
	}
	if cfg.Foreground != "" {
		set.Foreground, err = loadOptionalImage(searchPaths, cfg.Foreground, fmt.Sprintf("assets.frame_sets.%s.foreground", id))
		if err != nil {
			return FrameSet{}, err
		}
	}
	var frameBounds image.Rectangle
	for frame := cfg.Frames.First; frame <= cfg.Frames.Last; frame++ {
		assetPath := formatFramePath(cfg.Frames.Path, frame)
		asset, err := loadRequiredImage(searchPaths, assetPath, fmt.Sprintf("assets.frame_sets.%s.frames.%d", id, frame))
		if err != nil {
			return FrameSet{}, err
		}
		if frameBounds.Empty() {
			frameBounds = asset.Bounds
		} else if asset.Bounds != frameBounds {
			return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.frames.%d dimensions %s must match frame dimensions %s", id, frame, asset.Bounds, frameBounds)
		}
		set.Frames[frame] = asset
	}
	if set.Background != nil && set.Background.Bounds != frameBounds {
		return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.background dimensions %s must match frame dimensions %s", id, set.Background.Bounds, frameBounds)
	}
	if set.Foreground != nil && set.Foreground.Bounds != frameBounds {
		return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.foreground dimensions %s must match frame dimensions %s", id, set.Foreground.Bounds, frameBounds)
	}
	return set, nil
}

func formatFramePath(pattern string, frame int) string {
	if !strings.Contains(pattern, "%") {
		return pattern
	}
	return fmt.Sprintf(pattern, frame)
}

func hasPrintfPlaceholder(pattern string) bool {
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != '%' {
			continue
		}
		i++
		if i >= len(pattern) {
			return false
		}
		if pattern[i] == '%' {
			continue
		}
		return true
	}
	return false
}

func loadIndicatorSet(searchPaths []string, id string, cfg v3config.IndicatorSetConfig) (IndicatorSet, error) {
	set := IndicatorSet{ID: id, States: make(map[string]ImageAsset, len(cfg.States))}
	var err error
	if cfg.Background != "" {
		set.Background, err = loadOptionalImage(searchPaths, cfg.Background, fmt.Sprintf("assets.indicator_sets.%s.background", id))
		if err != nil {
			return IndicatorSet{}, err
		}
	}
	if cfg.Foreground != "" {
		set.Foreground, err = loadOptionalImage(searchPaths, cfg.Foreground, fmt.Sprintf("assets.indicator_sets.%s.foreground", id))
		if err != nil {
			return IndicatorSet{}, err
		}
	}
	for _, state := range []string{IndicatorStateOff, IndicatorStateOn, IndicatorStateUnknown} {
		assetPath := cfg.States[state]
		asset, err := loadRequiredImage(searchPaths, assetPath, fmt.Sprintf("assets.indicator_sets.%s.states.%s", id, state))
		if err != nil {
			return IndicatorSet{}, err
		}
		set.States[state] = asset
	}
	for _, state := range sortedKeys(cfg.States) {
		if _, exists := set.States[state]; exists {
			continue
		}
		asset, err := loadRequiredImage(searchPaths, cfg.States[state], fmt.Sprintf("assets.indicator_sets.%s.states.%s", id, state))
		if err != nil {
			return IndicatorSet{}, err
		}
		set.States[state] = asset
	}
	return set, nil
}

func loadOptionalImage(searchPaths []string, assetPath, label string) (*ImageAsset, error) {
	if strings.TrimSpace(assetPath) == "" {
		return nil, nil
	}
	asset, err := loadRequiredImage(searchPaths, assetPath, label)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func loadRequiredImage(searchPaths []string, assetPath, label string) (ImageAsset, error) {
	if strings.TrimSpace(assetPath) == "" {
		return ImageAsset{}, fmt.Errorf("%s path must not be empty", label)
	}
	resolved, searched, err := resolveAssetPath(searchPaths, assetPath)
	if err != nil {
		return ImageAsset{}, fmt.Errorf("%s path %q is invalid: %w", label, assetPath, err)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return ImageAsset{}, fmt.Errorf("%s path %q could not be loaded after searching %s: %w", label, assetPath, strings.Join(searched, ", "), err)
	}
	decoded, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return ImageAsset{}, fmt.Errorf("%s path %q could not be decoded as an image: %w", label, resolved, err)
	}
	return ImageAsset{Path: resolved, Data: data, Image: decoded, Bounds: decoded.Bounds()}, nil
}

func cleanSearchPaths(paths []string) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	cleaned := make([]string, 0, len(paths))
	seen := map[string]bool{}
	for _, item := range paths {
		if strings.TrimSpace(item) == "" {
			continue
		}
		absolute, err := filepath.Abs(item)
		if err != nil {
			return nil, fmt.Errorf("asset search path %q could not be resolved: %w", item, err)
		}
		absolute = filepath.Clean(absolute)
		if seen[absolute] {
			continue
		}
		seen[absolute] = true
		cleaned = append(cleaned, absolute)
	}
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("asset search paths must not be empty")
	}
	return cleaned, nil
}

func resolveAssetPath(searchPaths []string, assetPath string) (string, []string, error) {
	trimmed := strings.TrimSpace(assetPath)
	if strings.Contains(trimmed, "://") {
		return "", nil, fmt.Errorf("must be search-path relative, not remote or URL-like")
	}
	if filepath.IsAbs(trimmed) || path.IsAbs(filepath.ToSlash(trimmed)) {
		return "", nil, fmt.Errorf("must be search-path relative")
	}
	slashPath := filepath.ToSlash(trimmed)
	cleaned := path.Clean(slashPath)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || hasUpwardEscapeSegment(slashPath) {
		return "", nil, fmt.Errorf("must be search-path relative")
	}
	paths, err := cleanSearchPaths(searchPaths)
	if err != nil {
		return "", nil, err
	}
	searched := make([]string, 0, len(paths))
	for _, base := range paths {
		resolved := filepath.Clean(filepath.Join(base, filepath.FromSlash(cleaned)))
		if !isWithinRoot(base, resolved) {
			continue
		}
		searched = append(searched, resolved)
		if _, err := os.Stat(resolved); err == nil {
			return resolved, searched, nil
		} else if !os.IsNotExist(err) {
			return resolved, searched, err
		}
	}
	return "", searched, fmt.Errorf("could not be found in asset search paths: %s", strings.Join(searched, ", "))
}

func isWithinRoot(root, candidate string) bool {
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel))
}

func hasUpwardEscapeSegment(slashPath string) bool {
	for _, segment := range strings.Split(slashPath, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}

func sortedKeys[T any](items map[string]T) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
