package ebiten

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
	ebitenui "github.com/hajimehoshi/ebiten/v2"
)

const sceneGap = 12

// Adapter renders v3 dashboard scene output with Ebiten. It consumes the same
// dashboard scene model as the Fyne adapter and deliberately does not read
// sensors, poll endpoints, or own dashboard state.
type Adapter struct {
	repoRoot string
	width    int
	height   int

	ctx context.Context

	mu     sync.RWMutex
	assets map[string]cachedAsset
	parts  []renderedPart
}

type cachedAsset struct {
	image  *ebitenui.Image
	width  int
	height int
}

type loadedPart struct {
	index int
	part  v3dashboard.Part
	asset cachedAsset
}

type renderedPart struct {
	asset cachedAsset
	x     float64
	y     float64
	scale float64

	needle bool
	angle  float64
	pivotX float64
	pivotY float64

	source image.Rectangle
}

// New creates an Ebiten adapter for v3 dashboard scene output.
func New(repoRoot string, width int, height int) (*Adapter, error) {
	root, err := cleanRepoRoot(repoRoot)
	if err != nil {
		return nil, err
	}
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 480
	}
	return &Adapter{
		repoRoot: root,
		width:    width,
		height:   height,
		assets:   map[string]cachedAsset{},
	}, nil
}

// Run starts the Ebiten renderer loop. It blocks until the window closes or the
// supplied context is cancelled.
func (a *Adapter) Run(ctx context.Context, title string) error {
	if a == nil {
		return fmt.Errorf("v3 Ebiten adapter is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	a.ctx = ctx
	if title == "" {
		title = "GoDriveLog v3"
	}
	ebitenui.SetWindowTitle(title)
	ebitenui.SetWindowSize(a.width, a.height)
	return ebitenui.RunGame(a)
}

// UpdateScenes renders the latest v3 dashboard scenes. It keeps decoded Ebiten
// images cached by asset path while replacing the small render-part list for the
// next Draw call.
func (a *Adapter) UpdateScenes(scenes []v3dashboard.Scene) error {
	if a == nil {
		return fmt.Errorf("v3 Ebiten adapter is nil")
	}
	parts, width, height, err := a.renderParts(scenes)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.parts = parts
	if width > 0 {
		a.width = width
	}
	if height > 0 {
		a.height = height
	}
	a.mu.Unlock()
	return nil
}

// Update is Ebiten's game tick hook.
func (a *Adapter) Update() error {
	if a == nil || a.ctx == nil {
		return nil
	}
	select {
	case <-a.ctx.Done():
		return a.ctx.Err()
	default:
		return nil
	}
}

// Draw is Ebiten's render hook.
func (a *Adapter) Draw(screen *ebitenui.Image) {
	if a == nil || screen == nil {
		return
	}
	a.mu.RLock()
	parts := append([]renderedPart(nil), a.parts...)
	a.mu.RUnlock()
	for _, part := range parts {
		drawImage := part.asset.image
		if drawImage == nil {
			continue
		}
		if !part.source.Empty() {
			source := clampSourceRect(part.source, part.asset)
			if source.Empty() {
				continue
			}
			if subImage, ok := drawImage.SubImage(source).(*ebitenui.Image); ok {
				drawImage = subImage
			}
		}
		options := &ebitenui.DrawImageOptions{}
		if part.needle {
			options.GeoM.Translate(-part.pivotX, -part.pivotY)
			options.GeoM.Rotate(part.angle * math.Pi / 180)
			options.GeoM.Scale(part.scale, part.scale)
			options.GeoM.Translate(part.x, part.y)
		} else {
			options.GeoM.Scale(part.scale, part.scale)
			options.GeoM.Translate(part.x, part.y)
		}
		screen.DrawImage(drawImage, options)
	}
}

// Layout is Ebiten's layout hook.
func (a *Adapter) Layout(outsideWidth int, outsideHeight int) (int, int) {
	if a == nil {
		return 1, 1
	}
	a.mu.RLock()
	width, height := a.width, a.height
	a.mu.RUnlock()
	if width <= 0 {
		width = outsideWidth
	}
	if height <= 0 {
		height = outsideHeight
	}
	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}
	return width, height
}

func (a *Adapter) renderParts(scenes []v3dashboard.Scene) ([]renderedPart, int, int, error) {
	parts := []renderedPart{}
	yOffset := 0.0
	maxWidth := 0.0
	for _, scene := range scenes {
		sceneParts, sceneWidth, sceneHeight, err := a.renderSceneParts(scene, yOffset)
		if err != nil {
			return nil, 0, 0, err
		}
		parts = append(parts, sceneParts...)
		yOffset += sceneHeight + sceneGap
		if sceneWidth > maxWidth {
			maxWidth = sceneWidth
		}
	}
	if len(scenes) == 0 {
		return parts, 1, 1, nil
	}
	return parts, int(math.Ceil(maxWidth)), int(math.Ceil(yOffset - sceneGap)), nil
}

func (a *Adapter) renderSceneParts(scene v3dashboard.Scene, yOffset float64) ([]renderedPart, float64, float64, error) {
	parts := []renderedPart{}
	for _, widget := range scene.Widgets {
		widgetParts, err := a.renderWidgetParts(scene.DashboardID, widget, yOffset)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("dashboard %q widget %q: %w", scene.DashboardID, widget.ID, err)
		}
		parts = append(parts, widgetParts...)
	}
	width := float64(scene.Size.Width)
	height := float64(scene.Size.Height)
	if width <= 0 {
		width = partsWidth(parts)
	}
	if height <= 0 {
		height = partsHeight(parts, yOffset)
	}
	return parts, width, height, nil
}

func (a *Adapter) renderWidgetParts(dashboardID string, widget v3dashboard.Widget, yOffset float64) ([]renderedPart, error) {
	loaded := make([]loadedPart, 0, len(widget.Parts))
	gaugeWidth := 0.0
	gaugeHeight := 0.0
	for index, part := range widget.Parts {
		asset, err := a.loadAsset(part.AssetPath)
		if err != nil {
			return nil, fmt.Errorf("part %d %q asset %q: %w", index, part.Kind, part.AssetPath, err)
		}
		loaded = append(loaded, loadedPart{index: index, part: part, asset: asset})
		gaugeWidth, gaugeHeight = radialGaugeSize(gaugeWidth, gaugeHeight, part, asset)
	}

	parts := make([]renderedPart, 0, len(loaded))
	baseX, baseY := widgetPosition(widget)
	baseY += yOffset
	widgetScale := widget.Scale
	if widgetScale <= 0 {
		widgetScale = 1
	}

	for _, loadedPart := range loaded {
		part := loadedPart.part
		asset := loadedPart.asset
		if part.Kind == v3dashboard.PartKindNeedle {
			if gaugeWidth <= 0 || gaugeHeight <= 0 {
				gaugeWidth = float64(asset.width)
				gaugeHeight = float64(asset.height)
			}
			faceX := baseX + part.FacePivot.X*gaugeWidth*widgetScale
			faceY := baseY + part.FacePivot.Y*gaugeHeight*widgetScale
			parts = append(parts, renderedPart{
				asset:  asset,
				x:      faceX,
				y:      faceY,
				scale:  widgetScale,
				needle: true,
				angle:  part.Angle,
				pivotX: part.NeedlePivot.X * float64(asset.width),
				pivotY: part.NeedlePivot.Y * float64(asset.height),
			})
			continue
		}
		x, y := partPosition(baseX, baseY, asset, widgetScale, part)
		parts = append(parts, renderedPart{asset: asset, x: x, y: y, scale: widgetScale, source: partSourceRect(part)})
	}
	return parts, nil
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
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return cachedAsset{}, err
	}
	bounds := decoded.Bounds()
	cached := cachedAsset{
		image:  ebitenui.NewImageFromImage(decoded),
		width:  bounds.Dx(),
		height: bounds.Dy(),
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

func widgetPosition(widget v3dashboard.Widget) (float64, float64) {
	if len(widget.Position) < 2 {
		return 0, 0
	}
	return float64(widget.Position[0]), float64(widget.Position[1])
}

func partPosition(baseX, baseY float64, asset cachedAsset, scale float64, part v3dashboard.Part) (float64, float64) {
	x := baseX
	y := baseY
	if len(part.Position) >= 2 {
		x += float64(part.Position[0]) * scale
		y += float64(part.Position[1]) * scale
		return x, y
	}
	if part.Slot > 0 {
		x += float64(part.Slot) * float64(asset.width) * scale
	}
	return x, y
}

func partSourceRect(part v3dashboard.Part) image.Rectangle {
	if (part.Kind != v3dashboard.PartKindWheelStrip && part.Kind != v3dashboard.PartKindBar) || part.Window.Width <= 0 || part.Window.Height <= 0 {
		return image.Rectangle{}
	}
	x, y := 0, 0
	if len(part.Source) >= 2 {
		x = part.Source[0]
		y = part.Source[1]
	}
	return image.Rect(x, y, x+part.Window.Width, y+part.Window.Height)
}

func clampSourceRect(rect image.Rectangle, asset cachedAsset) image.Rectangle {
	if asset.width <= 0 || asset.height <= 0 {
		return image.Rectangle{}
	}
	if rect.Min.X < 0 {
		rect.Min.X = 0
	}
	if rect.Min.Y < 0 {
		rect.Min.Y = 0
	}
	width := rect.Dx()
	height := rect.Dy()
	if width <= 0 || height <= 0 {
		return image.Rectangle{}
	}
	if width > asset.width {
		width = asset.width
	}
	if height > asset.height {
		height = asset.height
	}
	if rect.Min.X+width > asset.width {
		rect.Min.X = asset.width - width
	}
	if rect.Min.Y+height > asset.height {
		rect.Min.Y = asset.height - height
	}
	rect.Max.X = rect.Min.X + width
	rect.Max.Y = rect.Min.Y + height
	return rect
}

func radialGaugeSize(currentWidth float64, currentHeight float64, part v3dashboard.Part, asset cachedAsset) (float64, float64) {
	if part.Kind == v3dashboard.PartKindNeedle {
		return currentWidth, currentHeight
	}
	width := float64(asset.width)
	height := float64(asset.height)
	if part.Kind == v3dashboard.PartKindLayer && part.Layer == "face" {
		return width, height
	}
	if currentWidth <= 0 || currentHeight <= 0 {
		return width, height
	}
	if part.Kind == v3dashboard.PartKindLayer && currentWidth*currentHeight < width*height {
		return width, height
	}
	return currentWidth, currentHeight
}

func partsWidth(parts []renderedPart) float64 {
	max := 0.0
	for _, part := range parts {
		width := part.x + float64(part.asset.width)*part.scale
		if width > max {
			max = width
		}
	}
	return max
}

func partsHeight(parts []renderedPart, yOffset float64) float64 {
	max := yOffset
	for _, part := range parts {
		height := part.y + float64(part.asset.height)*part.scale
		if height > max {
			max = height
		}
	}
	return max
}
