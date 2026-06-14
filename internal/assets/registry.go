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

type Registry struct {
	repoRoot      string
	Images        map[string]ImageSet
	DigitSets     map[string]DigitSet
	BarSets       map[string]BarSet
	FrameSets     map[string]FrameSet
	IndicatorSets map[string]IndicatorSet
}

type ImageSet struct{ ID string; Image, Background, Foreground *ImageAsset }
type DigitSet struct{ ID string; Background *ImageAsset; Characters map[string]ImageAsset; DecimalPoint, Foreground *ImageAsset; Spacing int }
type BarSet struct{ ID string; Background *ImageAsset; Cells map[string]ImageAsset; Foreground *ImageAsset; Spacing int }
type FrameSet struct{ ID string; Background *ImageAsset; Frames map[int]ImageAsset; First, Last int; Foreground *ImageAsset }
type IndicatorSet struct{ ID string; Background *ImageAsset; States map[string]ImageAsset; Foreground *ImageAsset }
type ImageAsset struct{ Path string; Data []byte; Image image.Image; Bounds image.Rectangle }

func Load(cfg v3config.AssetConfig, repoRoot string) (*Registry, error) {
	root, err := cleanRepoRoot(repoRoot); if err != nil { return nil, err }
	r := &Registry{repoRoot: root, Images: map[string]ImageSet{}, DigitSets: map[string]DigitSet{}, BarSets: map[string]BarSet{}, FrameSets: map[string]FrameSet{}, IndicatorSets: map[string]IndicatorSet{}}
	for _, id := range sortedKeys(cfg.ImageSets) { set, err := loadImageSet(root, id, cfg.ImageSets[id]); if err != nil { return nil, err }; r.Images[id] = set }
	for _, id := range sortedKeys(cfg.DigitSets) { set, err := loadDigitSet(root, id, cfg.DigitSets[id]); if err != nil { return nil, err }; r.DigitSets[id] = set }
	for _, id := range sortedKeys(cfg.BarSets) { set, err := loadBarSet(root, id, cfg.BarSets[id]); if err != nil { return nil, err }; r.BarSets[id] = set }
	for _, id := range sortedKeys(cfg.FrameSets) { set, err := loadFrameSet(root, id, cfg.FrameSets[id]); if err != nil { return nil, err }; r.FrameSets[id] = set }
	for _, id := range sortedKeys(cfg.IndicatorSets) { set, err := loadIndicatorSet(root, id, cfg.IndicatorSets[id]); if err != nil { return nil, err }; r.IndicatorSets[id] = set }
	return r, nil
}

func (r *Registry) RepoRoot() string { if r == nil { return "" }; return r.repoRoot }
func (r *Registry) ImageSet(id string) (ImageSet, bool) { if r == nil { return ImageSet{}, false }; set, ok := r.Images[id]; return set, ok }
func (r *Registry) DigitSet(id string) (DigitSet, bool) { if r == nil { return DigitSet{}, false }; set, ok := r.DigitSets[id]; return set, ok }
func (r *Registry) BarSet(id string) (BarSet, bool) { if r == nil { return BarSet{}, false }; set, ok := r.BarSets[id]; return set, ok }
func (r *Registry) FrameSet(id string) (FrameSet, bool) { if r == nil { return FrameSet{}, false }; set, ok := r.FrameSets[id]; return set, ok }
func (r *Registry) IndicatorSet(id string) (IndicatorSet, bool) { if r == nil { return IndicatorSet{}, false }; set, ok := r.IndicatorSets[id]; return set, ok }

func loadImageSet(root, id string, cfg v3config.ImageSetConfig) (ImageSet, error) {
	set := ImageSet{ID: id}; var err error
	if cfg.Image != "" { set.Image, err = loadOptionalImage(root, cfg.Image, fmt.Sprintf("assets.image_sets.%s.image", id)); if err != nil { return ImageSet{}, err } }
	if cfg.Background != "" { set.Background, err = loadOptionalImage(root, cfg.Background, fmt.Sprintf("assets.image_sets.%s.background", id)); if err != nil { return ImageSet{}, err } }
	if cfg.Foreground != "" { set.Foreground, err = loadOptionalImage(root, cfg.Foreground, fmt.Sprintf("assets.image_sets.%s.foreground", id)); if err != nil { return ImageSet{}, err } }
	if set.Image == nil && set.Background == nil && set.Foreground == nil { return ImageSet{}, fmt.Errorf("assets.image_sets.%s must define image, background, or foreground", id) }
	return set, nil
}

func loadDigitSet(root, id string, cfg v3config.DigitSetConfig) (DigitSet, error) {
	set := DigitSet{ID: id, Characters: map[string]ImageAsset{}, Spacing: cfg.Spacing}; var err error
	if cfg.Background != "" { set.Background, err = loadOptionalImage(root, cfg.Background, fmt.Sprintf("assets.digit_sets.%s.background", id)); if err != nil { return DigitSet{}, err } }
	if cfg.DecimalPoint != "" { set.DecimalPoint, err = loadOptionalImage(root, cfg.DecimalPoint, fmt.Sprintf("assets.digit_sets.%s.decimal_point", id)); if err != nil { return DigitSet{}, err } }
	if cfg.Foreground != "" { set.Foreground, err = loadOptionalImage(root, cfg.Foreground, fmt.Sprintf("assets.digit_sets.%s.foreground", id)); if err != nil { return DigitSet{}, err } }
	for _, ch := range requiredDigitCharacters { asset, err := loadRequiredImage(root, cfg.Characters[ch], fmt.Sprintf("assets.digit_sets.%s.characters.%s", id, ch)); if err != nil { return DigitSet{}, err }; set.Characters[ch] = asset }
	for _, ch := range sortedKeys(cfg.Characters) { if _, ok := set.Characters[ch]; ok { continue }; asset, err := loadRequiredImage(root, cfg.Characters[ch], fmt.Sprintf("assets.digit_sets.%s.characters.%s", id, ch)); if err != nil { return DigitSet{}, err }; set.Characters[ch] = asset }
	return set, nil
}

func loadBarSet(root, id string, cfg v3config.BarSetConfig) (BarSet, error) {
	set := BarSet{ID: id, Cells: map[string]ImageAsset{}, Spacing: cfg.Spacing}; var err error
	if cfg.Background != "" { set.Background, err = loadOptionalImage(root, cfg.Background, fmt.Sprintf("assets.bar_sets.%s.background", id)); if err != nil { return BarSet{}, err } }
	if cfg.Foreground != "" { set.Foreground, err = loadOptionalImage(root, cfg.Foreground, fmt.Sprintf("assets.bar_sets.%s.foreground", id)); if err != nil { return BarSet{}, err } }
	if strings.TrimSpace(cfg.Cells["off"]) == "" { return BarSet{}, fmt.Errorf("assets.bar_sets.%s.cells.off path must not be empty", id) }
	var bounds image.Rectangle
	for _, cell := range sortedKeys(cfg.Cells) { asset, err := loadRequiredImage(root, cfg.Cells[cell], fmt.Sprintf("assets.bar_sets.%s.cells.%s", id, cell)); if err != nil { return BarSet{}, err }; if bounds.Empty() { bounds = asset.Bounds } else if asset.Bounds != bounds { return BarSet{}, fmt.Errorf("assets.bar_sets.%s.cells.%s dimensions %s must match cell dimensions %s", id, cell, asset.Bounds, bounds) }; set.Cells[cell] = asset }
	return set, nil
}

func loadFrameSet(root, id string, cfg v3config.FrameSetConfig) (FrameSet, error) {
	if cfg.Frames.First > cfg.Frames.Last { return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.frames.first must be less than or equal to last", id) }
	if cfg.Frames.First < cfg.Frames.Last && !hasPrintfPlaceholder(cfg.Frames.Path) { return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.frames.path must contain a printf placeholder for a multi-frame range", id) }
	set := FrameSet{ID: id, Frames: map[int]ImageAsset{}, First: cfg.Frames.First, Last: cfg.Frames.Last}; var err error
	if cfg.Background != "" { set.Background, err = loadOptionalImage(root, cfg.Background, fmt.Sprintf("assets.frame_sets.%s.background", id)); if err != nil { return FrameSet{}, err } }
	if cfg.Foreground != "" { set.Foreground, err = loadOptionalImage(root, cfg.Foreground, fmt.Sprintf("assets.frame_sets.%s.foreground", id)); if err != nil { return FrameSet{}, err } }
	var bounds image.Rectangle
	for frame := cfg.Frames.First; frame <= cfg.Frames.Last; frame++ { assetPath := formatFramePath(cfg.Frames.Path, frame); asset, err := loadRequiredImage(root, assetPath, fmt.Sprintf("assets.frame_sets.%s.frames.%d", id, frame)); if err != nil { return FrameSet{}, err }; if bounds.Empty() { bounds = asset.Bounds } else if asset.Bounds != bounds { return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.frames.%d dimensions %s must match frame dimensions %s", id, frame, asset.Bounds, bounds) }; set.Frames[frame] = asset }
	if set.Background != nil && set.Background.Bounds != bounds { return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.background dimensions %s must match frame dimensions %s", id, set.Background.Bounds, bounds) }
	if set.Foreground != nil && set.Foreground.Bounds != bounds { return FrameSet{}, fmt.Errorf("assets.frame_sets.%s.foreground dimensions %s must match frame dimensions %s", id, set.Foreground.Bounds, bounds) }
	return set, nil
}

func formatFramePath(pattern string, frame int) string { if !strings.Contains(pattern, "%") { return pattern }; return fmt.Sprintf(pattern, frame) }
func hasPrintfPlaceholder(pattern string) bool { for i := 0; i < len(pattern); i++ { if pattern[i] != '%' { continue }; i++; if i >= len(pattern) { return false }; if pattern[i] == '%' { continue }; return true }; return false }

func loadIndicatorSet(root, id string, cfg v3config.IndicatorSetConfig) (IndicatorSet, error) {
	set := IndicatorSet{ID: id, States: map[string]ImageAsset{}}; var err error
	if cfg.Background != "" { set.Background, err = loadOptionalImage(root, cfg.Background, fmt.Sprintf("assets.indicator_sets.%s.background", id)); if err != nil { return IndicatorSet{}, err } }
	if cfg.Foreground != "" { set.Foreground, err = loadOptionalImage(root, cfg.Foreground, fmt.Sprintf("assets.indicator_sets.%s.foreground", id)); if err != nil { return IndicatorSet{}, err } }
	for _, state := range []string{IndicatorStateOff, IndicatorStateOn, IndicatorStateUnknown} { asset, err := loadRequiredImage(root, cfg.States[state], fmt.Sprintf("assets.indicator_sets.%s.states.%s", id, state)); if err != nil { return IndicatorSet{}, err }; set.States[state] = asset }
	for _, state := range sortedKeys(cfg.States) { if _, ok := set.States[state]; ok { continue }; asset, err := loadRequiredImage(root, cfg.States[state], fmt.Sprintf("assets.indicator_sets.%s.states.%s", id, state)); if err != nil { return IndicatorSet{}, err }; set.States[state] = asset }
	return set, nil
}

func loadOptionalImage(root, assetPath, label string) (*ImageAsset, error) { if strings.TrimSpace(assetPath) == "" { return nil, nil }; asset, err := loadRequiredImage(root, assetPath, label); if err != nil { return nil, err }; return &asset, nil }

func loadRequiredImage(root, assetPath, label string) (ImageAsset, error) {
	if strings.TrimSpace(assetPath) == "" { return ImageAsset{}, fmt.Errorf("%s path must not be empty", label) }
	resolved, err := resolveAssetPath(root, assetPath); if err != nil { return ImageAsset{}, fmt.Errorf("%s path %q is invalid: %w", label, assetPath, err) }
	data, err := os.ReadFile(resolved); if err != nil { return ImageAsset{}, fmt.Errorf("%s path %q could not be loaded: %w", label, resolved, err) }
	decoded, _, err := image.Decode(strings.NewReader(string(data))); if err != nil { return ImageAsset{}, fmt.Errorf("%s path %q could not be decoded as an image: %w", label, resolved, err) }
	return ImageAsset{Path: resolved, Data: data, Image: decoded, Bounds: decoded.Bounds()}, nil
}

func cleanRepoRoot(repoRoot string) (string, error) { if strings.TrimSpace(repoRoot) == "" { repoRoot = "." }; cleaned, err := filepath.Abs(repoRoot); if err != nil { return "", fmt.Errorf("repository root %q could not be resolved: %w", repoRoot, err) }; return filepath.Clean(cleaned), nil }

func resolveAssetPath(root, assetPath string) (string, error) {
	trimmed := strings.TrimSpace(assetPath)
	if strings.Contains(trimmed, "://") { return "", fmt.Errorf("must be repository-root relative, not remote or URL-like") }
	if filepath.IsAbs(trimmed) || path.IsAbs(filepath.ToSlash(trimmed)) { return "", fmt.Errorf("must be repository-root relative") }
	slashPath := filepath.ToSlash(trimmed); cleaned := path.Clean(slashPath)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || hasUpwardEscapeSegment(slashPath) { return "", fmt.Errorf("must be repository-root relative") }
	resolved := filepath.Clean(filepath.Join(root, filepath.FromSlash(cleaned))); if !isWithinRoot(root, resolved) { return "", fmt.Errorf("must stay within repository root") }
	return resolved, nil
}

func isWithinRoot(root, candidate string) bool { rel, err := filepath.Rel(root, candidate); if err != nil { return false }; return rel == "." || (!strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel)) }
func hasUpwardEscapeSegment(slashPath string) bool { for _, segment := range strings.Split(slashPath, "/") { if segment == ".." { return true } }; return false }
func sortedKeys[T any](items map[string]T) []string { keys := make([]string, 0, len(items)); for key := range items { keys = append(keys, key) }; sort.Strings(keys); return keys }
