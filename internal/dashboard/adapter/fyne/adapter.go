package fyne

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	fyneui "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

const (
	sceneGap                = 12
	needleFrameAngleStep    = 1.0
	needleFrameCount        = int(360 / needleFrameAngleStep)
	needleFrameZeroIndex    = needleFrameCount / 2
)

// Adapter renders v3 dashboard scene output with Fyne. It deliberately consumes
// only dashboard scene data and resolved asset paths; it does not read sensors,
// poll OBD endpoints, or own dashboard state.
type Adapter struct {
	repoRoot        string
	root            *fyneui.Container
	assets          map[string]cachedAsset
	needleFrameSets map[string]needleFrameSet
	images          map[string]*canvas.Image
	refreshRoot     func(*fyneui.Container)
	refreshImage    func(*canvas.Image)
}

type cachedAsset struct {
	resource fyneui.Resource
	size     fyneui.Size
	data     []byte
}

type needleFrameSet struct {
	frames []rotatedNeedleAsset
}

type rotatedNeedleAsset struct {
	asset  cachedAsset
	pivotX float32
	pivotY float32
}

type loadedPart struct {
	index int
	part  v3dashboard.Part
	asset cachedAsset
}

type renderedPart struct {
	key   string
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
		repoRoot:        root,
		root:            container.NewWithoutLayout(),
		assets:          map[string]cachedAsset{},
		needleFrameSets: map[string]needleFrameSet{},
		images:          map[string]*canvas.Image{},
		refreshRoot: func(root *fyneui.Container) {
			root.Refresh()
		},
		refreshImage: func(image *canvas.Image) {
			image.Refresh()
		},
	}, nil
}

// CanvasObject returns the visible Fyne object managed by the adapter.
func (a *Adapter) CanvasObject() fyneui.CanvasObject {
	return a.root
}

// Update renders the latest v3 dashboard scenes. It reuses Fyne image objects by
// stable scene/widget/part identity so fast dashboard updates do not rebuild the
// native canvas object tree when digit resources change.
func (a *Adapter) Update(scenes []v3dashboard.Scene) error {
	if a == nil {
		return fmt.Errorf("v3 Fyne adapter is nil")
	}

	parts, size, err := a.renderParts(scenes)
	if err != nil {
		return err
	}

	changed := a.syncImages(parts)
	if size.Width <= 0 {
		size.Width = 1
	}
	if size.Height <= 0 {
		size.Height = 1
	}
	if a.root.Size() != size {
		a.root.Resize(size)
		changed = true
	}
	if changed && a.refreshRoot != nil {
		a.refreshRoot(a.root)
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
		widgetParts, err := a.renderWidgetParts(scene.DashboardID, widget, yOffset)
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

func (a *Adapter) renderWidgetParts(dashboardID string, widget v3dashboard.Widget, yOffset float32) ([]renderedPart, error) {
	loaded := make([]loadedPart, 0, len(widget.Parts))
	gaugeSize := fyneui.Size{}
	for index, part := range widget.Parts {
		asset, err := a.loadAsset(part.AssetPath)
		if err != nil {
			return nil, fmt.Errorf("part %d %q asset %q: %w", index, part.Kind, part.AssetPath, err)
		}
		loaded = append(loaded, loadedPart{index: index, part: part, asset: asset})
		gaugeSize = radialGaugeSize(gaugeSize, part, asset.size)
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
			rendered, err := a.renderNeedlePart(dashboardID, widget.ID, loadedPart.index, part, asset, gaugeSize, baseX, baseY, widgetScale)
			if err != nil {
				return nil, err
			}
			parts = append(parts, rendered)
			continue
		}
		size := scaledSize(asset.size, widgetScale)
		x, y := partPosition(baseX, baseY, size, widgetScale, part)
		parts = append(parts, renderedPart{key: renderedPartKey(dashboardID, widget.ID, loadedPart.index, part), asset: asset, size: size, x: x, y: y})
	}
	return parts, nil
}

func (a *Adapter) renderNeedlePart(dashboardID string, widgetID string, index int, part v3dashboard.Part, asset cachedAsset, gaugeSize fyneui.Size, baseX float32, baseY float32, widgetScale float64) (renderedPart, error) {
	if gaugeSize.Width <= 0 || gaugeSize.Height <= 0 {
		gaugeSize = asset.size
	}
	rotated, err := a.preparedNeedleFrame(asset, part.Angle, part.NeedlePivot.X, part.NeedlePivot.Y)
	if err != nil {
		return renderedPart{}, err
	}
	scale := float32(widgetScale)
	faceX := baseX + float32(part.FacePivot.X)*gaugeSize.Width*scale
	faceY := baseY + float32(part.FacePivot.Y)*gaugeSize.Height*scale
	size := scaledSize(rotated.asset.size, widgetScale)
	x := faceX - rotated.pivotX*scale
	y := faceY - rotated.pivotY*scale
	return renderedPart{key: renderedPartKey(dashboardID, widgetID, index, part), asset: rotated.asset, size: size, x: x, y: y}, nil
}

func (a *Adapter) syncImages(parts []renderedPart) bool {
	if a.images == nil {
		a.images = map[string]*canvas.Image{}
	}
	changed := false
	objects := make([]fyneui.CanvasObject, 0, len(parts))
	nextImages := make(map[string]*canvas.Image, len(parts))
	for _, part := range parts {
		object, ok := a.images[part.key]
		if !ok {
			object = canvas.NewImageFromResource(part.asset.resource)
			object.FillMode = canvas.ImageFillStretch
			changed = true
		}
		if a.applyImagePart(object, part) {
			changed = true
		}
		nextImages[part.key] = object
		objects = append(objects, object)
	}
	if len(a.images) != len(nextImages) {
		changed = true
	}
	if !sameObjects(a.root.Objects, objects) {
		a.root.Objects = objects
		changed = true
	}
	a.images = nextImages
	return changed
}

func (a *Adapter) applyImagePart(object *canvas.Image, part renderedPart) bool {
	changed := false
	if object.Resource != part.asset.resource {
		object.Resource = part.asset.resource
		if a.refreshImage != nil {
			a.refreshImage(object)
		}
		changed = true
	}
	position := fyneui.NewPos(part.x, part.y)
	if object.Position() != position {
		object.Move(position)
		changed = true
	}
	if object.Size() != part.size {
		object.Resize(part.size)
		changed = true
	}
	return changed
}

func sameObjects(current []fyneui.CanvasObject, next []fyneui.CanvasObject) bool {
	if len(current) != len(next) {
		return false
	}
	for index := range current {
		if current[index] != next[index] {
			return false
		}
	}
	return true
}

func renderedPartKey(dashboardID string, widgetID string, index int, part v3dashboard.Part) string {
	if part.Kind == v3dashboard.PartKindLayer && part.Layer != "" {
		return fmt.Sprintf("%s/%s/layer/%s", dashboardID, widgetID, part.Layer)
	}
	switch part.Kind {
	case v3dashboard.PartKindBackground, v3dashboard.PartKindCharacter, v3dashboard.PartKindDecimalPoint, v3dashboard.PartKindForeground, v3dashboard.PartKindCell:
		return fmt.Sprintf("%s/%s/%s/%d", dashboardID, widgetID, part.Kind, part.Slot)
	case v3dashboard.PartKindState:
		return fmt.Sprintf("%s/%s/%s/%s", dashboardID, widgetID, part.Kind, part.State)
	case v3dashboard.PartKindFrame, v3dashboard.PartKindNeedle:
		return fmt.Sprintf("%s/%s/%s", dashboardID, widgetID, part.Kind)
	default:
		return fmt.Sprintf("%s/%s/%s/%d", dashboardID, widgetID, part.Kind, index)
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
		data:     data,
	}
	a.assets[cacheKey] = cached
	return cached, nil
}

func (a *Adapter) preparedNeedleFrame(asset cachedAsset, angle float64, pivotX float64, pivotY float64) (rotatedNeedleAsset, error) {
	key := needleFrameSetKey(asset, pivotX, pivotY)
	if a.needleFrameSets == nil {
		a.needleFrameSets = map[string]needleFrameSet{}
	}
	set, ok := a.needleFrameSets[key]
	if !ok {
		prepared, err := prepareNeedleFrameSet(asset, pivotX, pivotY, key)
		if err != nil {
			return rotatedNeedleAsset{}, err
		}
		a.needleFrameSets[key] = prepared
		set = prepared
	}
	index := needleFrameIndex(angle)
	if index < 0 || index >= len(set.frames) {
		return rotatedNeedleAsset{}, fmt.Errorf("radial needle frame index %d out of range", index)
	}
	return set.frames[index], nil
}

func needleFrameSetKey(asset cachedAsset, pivotX float64, pivotY float64) string {
	return strings.Join([]string{
		asset.resource.Name(),
		strconv.FormatFloat(pivotX, 'f', -1, 64),
		strconv.FormatFloat(pivotY, 'f', -1, 64),
	}, "|")
}

func prepareNeedleFrameSet(asset cachedAsset, pivotX float64, pivotY float64, key string) (needleFrameSet, error) {
	width := float64(asset.size.Width)
	height := float64(asset.size.Height)
	pivotPX := pivotX * width
	pivotPY := pivotY * height
	decoded, _, err := image.Decode(bytes.NewReader(asset.data))
	if err != nil {
		return needleFrameSet{}, err
	}
	frames := make([]rotatedNeedleAsset, needleFrameCount)
	for index := 0; index < needleFrameCount; index++ {
		angle := needleFrameAngle(index)
		if math.Abs(angle) < 0.000001 {
			frames[index] = rotatedNeedleAsset{asset: asset, pivotX: float32(pivotPX), pivotY: float32(pivotPY)}
			continue
		}
		rotatedImage, rotatedPivotX, rotatedPivotY := rotateImageAroundPivot(decoded, angle, pivotPX, pivotPY)
		var buf bytes.Buffer
		if err := png.Encode(&buf, rotatedImage); err != nil {
			return needleFrameSet{}, err
		}
		data := buf.Bytes()
		rotatedAsset := cachedAsset{
			resource: fyneui.NewStaticResource(fmt.Sprintf("%s@radial-needle|%s|%03d", asset.resource.Name(), key, index), data),
			size:     fyneui.NewSize(float32(rotatedImage.Bounds().Dx()), float32(rotatedImage.Bounds().Dy())),
			data:     data,
		}
		frames[index] = rotatedNeedleAsset{asset: rotatedAsset, pivotX: rotatedPivotX, pivotY: rotatedPivotY}
	}
	return needleFrameSet{frames: frames}, nil
}

func needleFrameIndex(angle float64) int {
	step := int(math.Round(angle / needleFrameAngleStep))
	step = ((step + needleFrameZeroIndex) % needleFrameCount + needleFrameCount) % needleFrameCount
	return step
}

func needleFrameAngle(index int) float64 {
	return float64(index-needleFrameZeroIndex) * needleFrameAngleStep
}

func rotateImageAroundPivot(source image.Image, angle float64, pivotX float64, pivotY float64) (*image.NRGBA, float32, float32) {
	bounds := source.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())
	radians := angle * math.Pi / 180
	cosAngle := math.Cos(radians)
	sinAngle := math.Sin(radians)
	corners := [][2]float64{{0, 0}, {width, 0}, {0, height}, {width, height}}
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, corner := range corners {
		x, y := rotatePoint(corner[0]-pivotX, corner[1]-pivotY, cosAngle, sinAngle)
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	rotatedWidth := int(math.Ceil(maxX - minX))
	rotatedHeight := int(math.Ceil(maxY - minY))
	if rotatedWidth < 1 {
		rotatedWidth = 1
	}
	if rotatedHeight < 1 {
		rotatedHeight = 1
	}
	rotated := image.NewNRGBA(image.Rect(0, 0, rotatedWidth, rotatedHeight))
	for y := 0; y < rotatedHeight; y++ {
		for x := 0; x < rotatedWidth; x++ {
			rotatedX := float64(x) + minX
			rotatedY := float64(y) + minY
			sourceX := rotatedX*cosAngle + rotatedY*sinAngle + pivotX
			sourceY := -rotatedX*sinAngle + rotatedY*cosAngle + pivotY
			sampleX := int(math.Floor(sourceX + 0.5))
			sampleY := int(math.Floor(sourceY + 0.5))
			if sampleX < 0 || sampleY < 0 || sampleX >= bounds.Dx() || sampleY >= bounds.Dy() {
				continue
			}
			c := color.NRGBAModel.Convert(source.At(bounds.Min.X+sampleX, bounds.Min.Y+sampleY)).(color.NRGBA)
			rotated.SetNRGBA(x, y, c)
		}
	}
	return rotated, float32(-minX), float32(-minY)
}

func rotatePoint(x float64, y float64, cosAngle float64, sinAngle float64) (float64, float64) {
	return x*cosAngle - y*sinAngle, x*sinAngle + y*cosAngle
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

func radialGaugeSize(current fyneui.Size, part v3dashboard.Part, size fyneui.Size) fyneui.Size {
	if part.Kind == v3dashboard.PartKindNeedle {
		return current
	}
	if part.Kind == v3dashboard.PartKindLayer && part.Layer == "face" {
		return size
	}
	if current.Width <= 0 || current.Height <= 0 {
		return size
	}
	if part.Kind == v3dashboard.PartKindLayer && current.Width*current.Height < size.Width*size.Height {
		return size
	}
	return current
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
