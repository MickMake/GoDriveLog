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
}

type cachedAsset struct {
	resource fyneui.Resource
	size     fyneui.Size
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

// Update renders the latest v3 dashboard scenes. Multiple selected scenes are
// stacked vertically; most vehicles are expected to select one scene initially.
func (a *Adapter) Update(scenes []v3dashboard.Scene) error {
	if a == nil {
		return fmt.Errorf("v3 Fyne adapter is nil")
	}

	objects := make([]fyneui.CanvasObject, 0, len(scenes))
	yOffset := float32(0)
	maxWidth := float32(0)
	for _, scene := range scenes {
		sceneObject, size, err := a.renderScene(scene)
		if err != nil {
			return err
		}
		sceneObject.Move(fyneui.NewPos(0, yOffset))
		sceneObject.Resize(size)
		objects = append(objects, sceneObject)
		yOffset += size.Height + sceneGap
		if size.Width > maxWidth {
			maxWidth = size.Width
		}
	}

	a.root.Objects = objects
	if len(objects) == 0 {
		a.root.Resize(fyneui.NewSize(1, 1))
	} else {
		a.root.Resize(fyneui.NewSize(maxWidth, yOffset-sceneGap))
	}
	if fyneui.CurrentApp() != nil {
		a.root.Refresh()
	}
	return nil
}

func (a *Adapter) renderScene(scene v3dashboard.Scene) (*fyneui.Container, fyneui.Size, error) {
	sceneRoot := container.NewWithoutLayout()
	objects := []fyneui.CanvasObject{}

	for _, widget := range scene.Widgets {
		widgetObjects, err := a.renderWidget(widget)
		if err != nil {
			return nil, fyneui.Size{}, fmt.Errorf("dashboard %q widget %q: %w", scene.DashboardID, widget.ID, err)
		}
		objects = append(objects, widgetObjects...)
	}

	sceneRoot.Objects = objects
	size := fyneui.NewSize(float32(scene.Size.Width), float32(scene.Size.Height))
	if size.Width <= 0 {
		size.Width = boundsWidth(objects)
	}
	if size.Height <= 0 {
		size.Height = boundsHeight(objects)
	}
	sceneRoot.Resize(size)
	return sceneRoot, size, nil
}

func (a *Adapter) renderWidget(widget v3dashboard.Widget) ([]fyneui.CanvasObject, error) {
	objects := make([]fyneui.CanvasObject, 0, len(widget.Parts))
	baseX, baseY := widgetPosition(widget)

	for index, part := range widget.Parts {
		asset, err := a.loadAsset(part.AssetPath)
		if err != nil {
			return nil, fmt.Errorf("part %d %q asset %q: %w", index, part.Kind, part.AssetPath, err)
		}
		object := canvas.NewImageFromResource(asset.resource)
		object.FillMode = canvas.ImageFillStretch

		x, y := partPosition(baseX, baseY, asset.size, part)
		object.Move(fyneui.NewPos(x, y))
		object.Resize(asset.size)
		objects = append(objects, object)
	}

	return objects, nil
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

func partPosition(baseX, baseY float32, size fyneui.Size, part v3dashboard.Part) (float32, float32) {
	x := baseX
	y := baseY
	if part.Slot > 0 {
		x += float32(part.Slot) * size.Width
	}
	return x, y
}

func boundsWidth(objects []fyneui.CanvasObject) float32 {
	var width float32 = 1
	for _, object := range objects {
		right := object.Position().X + object.Size().Width
		if right > width {
			width = right
		}
	}
	return width
}

func boundsHeight(objects []fyneui.CanvasObject) float32 {
	var height float32 = 1
	for _, object := range objects {
		bottom := object.Position().Y + object.Size().Height
		if bottom > height {
			height = bottom
		}
	}
	return height
}

// RenderedObjectCount is intended for focused adapter tests.
func (a *Adapter) RenderedObjectCount() int {
	if a == nil || a.root == nil {
		return 0
	}
	count := 0
	for _, sceneObject := range a.root.Objects {
		count += countObjects(sceneObject)
	}
	return count
}

func countObjects(object fyneui.CanvasObject) int {
	containerObject, ok := object.(*fyneui.Container)
	if !ok {
		return 1
	}
	count := 0
	for _, child := range containerObject.Objects {
		count += countObjects(child)
	}
	return count
}
