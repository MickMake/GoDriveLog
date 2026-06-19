package fyne

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	fyneui "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

const sceneGap = 12

// Adapter renders v3 dashboard scene output with Fyne. It deliberately consumes
// only dashboard scene data and resolved asset paths; it does not read sensors,
// poll OBD endpoints, or own dashboard state.
type Adapter struct {
	repoRoot string
	root     *fyneui.Container
	assets   map[string]cachedAsset
	images   []*canvas.Image
}

type cachedAsset struct {
	resource fyneui.Resource
	size     fyneui.Size
}

type renderedPart struct {
	asset cachedAsset
	size  fyneui.Size
	x     float32
	y     float32
}

// New creates a Fyne adapter for v3 dashboard scene output.
func New(repoRoot string) (*Adapter, error) {
	root, err := cleanRepoRoot(repoRoot)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		repoRoot: root,
		root:     container.NewWithoutLayout(),
		assets:   map[string]cachedAsset{},
	}, nil
}

// CanvasObject returns the visible Fyne object managed by the adapter.
func (a *Adapter) CanvasObject() fyneui.CanvasObject {
	return a.root
}

// Update renders the latest v3 dashboard scenes. It reuses existing Fyne image
// objects when the rendered part count is unchanged so fast dashboard updates do
// not build native canvas/image object churn.
func (a *Adapter) Update(scenes []v3dashboard.Scene) error {
	if a == nil {
		return fmt.Errorf("v3 Fyne adapter is nil")
	}

	parts, size, err := a.renderParts(scenes)
	if err != nil {
		return err
	}

	if len(a.images) != len(parts) {
		a.rebuildImages(parts)
	} else {
		a.updateImages(parts)
	}
	if size.Width <= 0 {
		size.Width = 1
	}
	if size.Height <= 0 {
		size.Height = 1
	}
	a.root.Resize(size)
	if fyneui.CurrentApp() != nil {
		a.root.Refresh()
	}
	return nil
}

func (a *Adapter) renderParts(scenes []v3dashboard.Scene) ([]renderedPart, fyneui.Size, error) {
	parts := []renderedPart{}
	yOffset := float32(0)
	maxWidth := float32(0)
	for _, scene := range scenes {
		sceneParts, sceneSize, err := a.renderSceneParts(scene, yOffset)
		if err != nil {
			return nil, fyneui.Size{}, err
		}
		parts = append(parts, sceneParts...)
		yOffset += sceneSize.Height + sceneGap
		if sceneSize.Width > maxWidth {
			maxWidth = sceneSize.Width
		}
	}
	if len(scenes) == 0 {
		return parts, fyneui.NewSize(1, 1), nil
	}
	return parts, fyneui.NewSize(maxWidth, yOffset-sceneGap), nil
}

func (a *Adapter) renderSceneParts(scene v3dashboard.Scene, yOffset float32) ([]renderedPart, fyneui.Size, error) {
	parts := []renderedPart{}
	for _, widget := range scene.Widgets {
		widgetParts, err := a.renderWidgetParts(widget, yOffset)
		if err != nil {
			return nil, fyneui.Size{}, fmt.Errorf("dashboard %q widget %q: %w", scene.DashboardID, widget.ID, err)
		}
		parts = append(parts, widgetParts...)
	}
	size := fyneui.NewSize(float32(scene.Size.Width), float32(scene.Size.Height))
	if size.Width <= 0 {
		size.Width = partsWidth(parts)
	}
	if size.Height <= 0 {
		size.Height = partsHeight(parts, yOffset)
	}
	return parts, size, nil
}

func (a *Adapter) renderWidgetParts(widget v3dashboard.Widget, yOffset float32) ([]renderedPart, error) {
	parts := make([]renderedPart, 0, len(widget.Parts))
	baseX, baseY := widgetPosition(widget)
	baseY += yOffset
	widgetScale := widget.Scale
	if widgetScale <= 0 {
		widgetScale = 1
	}

	for index, part := range widget.Parts {
		asset, err := a.loadAsset(part.AssetPath)
		if err != nil {
			return nil, fmt.Errorf("part %d %q asset %q: %w", index, part.Kind, part.AssetPath, err)
		}
		size := scaledSize(asset.size, widgetScale)
		x, y := partPosition(baseX, baseY, size, widgetScale, part)
		parts = append(parts, renderedPart{asset: asset, size: size, x: x, y: y})
	}
	return parts, nil
}

func (a *Adapter) rebuildImages(parts []renderedPart) {
	objects := make([]fyneui.CanvasObject, 0, len(parts))
	a.images = make([]*canvas.Image, 0, len(parts))
	for _, part := range parts {
		object := canvas.NewImageFromResource(part.asset.resource)
		object.FillMode = canvas.ImageFillStretch
		object.Move(fyneui.NewPos(part.x, part.y))
		object.Resize(part.size)
		a.images = append(a.images, object)
		objects = append(objects, object)
	}
	a.root.Objects = objects
}

func (a *Adapter) updateImages(parts []renderedPart) {
	for index, part := range parts {
		object := a.images[index]
		if object.Resource != part.asset.resource {
			object.Resource = part.asset.resource
			object.Refresh()
		}
		object.Move(fyneui.NewPos(part.x, part.y))
		object.Resize(part.size)
	}
}

func (a *Adapter) loadAsset(assetPath string) (cachedAsset, error) {
	fullPath, cacheKey, err := a.resolveAssetPath(assetPath)
	if err != nil {
		return cachedAsset{}, err
	}
	if cached, ok := a.assets[cacheKey]; ok {
		return cached, nil
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return cachedAsset{}, err
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return cachedAsset{}, err
	}
	cached := cachedAsset{
		resource: fyneui.NewStaticResource(cacheKey, data),
		size:     fyneui.NewSize(float32(config.Width), float32(config.Height)),
	}
	a.assets[cacheKey] = cached
	return cached, nil
}

func (a *Adapter) resolveAssetPath(assetPath string) (string, string, error) {
	trimmed := strings.TrimSpace(assetPath)
	if trimmed == "" {
		return "", "", fmt.Errorf("asset path must not be empty")
	}
	if strings.Contains(trimmed, "://") {
		return "", "", fmt.Errorf("asset path %q must be a filesystem path, not a URL", assetPath)
	}
	if filepath.IsAbs(trimmed) {
		cleaned := filepath.Clean(trimmed)
		return cleaned, cleaned, nil
	}

	cleaned, err := cleanRelativeAssetPath(trimmed)
	if err != nil {
		return "", "", err
	}
	fullPath := filepath.Join(a.repoRoot, filepath.FromSlash(cleaned))
	return fullPath, cleaned, nil
}

func cleanRepoRoot(repoRoot string) (string, error) {
	trimmed := strings.TrimSpace(repoRoot)
	if trimmed == "" {
		trimmed = "."
	}
	abs, err := filepath.Abs(trimmed)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func cleanRelativeAssetPath(assetPath string) (string, error) {
	if path.IsAbs(assetPath) || filepath.IsAbs(assetPath) {
		return "", fmt.Errorf("asset path %q must be relative", assetPath)
	}
	cleaned := path.Clean(strings.ReplaceAll(assetPath, "\\", "/"))
	if cleaned == "." || strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return "", fmt.Errorf("asset path %q must not escape the repository root", assetPath)
	}
	return cleaned, nil
}

func widgetPosition(widget v3dashboard.Widget) (float32, float32) {
	if len(widget.Position) < 2 {
		return 0, 0
	}
	return float32(widget.Position[0]), float32(widget.Position[1])
}

func scaledSize(size fyneui.Size, scale float64) fyneui.Size {
	if scale <= 0 || scale == 1 {
		return size
	}
	return fyneui.NewSize(size.Width*float32(scale), size.Height*float32(scale))
}

func partPosition(baseX, baseY float32, size fyneui.Size, scale float64, part v3dashboard.Part) (float32, float32) {
	x := baseX
	y := baseY
	if len(part.Position) >= 2 {
		x += float32(part.Position[0]) * float32(scale)
		y += float32(part.Position[1]) * float32(scale)
		return x, y
	}
	if part.Slot > 0 {
		x += float32(part.Slot) * size.Width
	}
	return x, y
}

func partsWidth(parts []renderedPart) float32 {
	max := float32(0)
	for _, part := range parts {
		width := part.x + part.size.Width
		if width > max {
			max = width
		}
	}
	return max
}

func partsHeight(parts []renderedPart, yOffset float32) float32 {
	max := yOffset
	for _, part := range parts {
		height := part.y + part.size.Height
		if height > max {
			max = height
		}
	}
	return max
}
