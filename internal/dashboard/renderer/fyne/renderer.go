package fyne

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	fyneui "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/assets"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scene"
)

type Renderer struct {
	root              *fyneui.Container
	assetRegistry     *assets.Registry
	minRenderInterval time.Duration
	lastRender        time.Time
	elements          map[string]*renderedElement
	resources         map[string]fyneui.Resource
}

type renderedElement struct {
	id          string
	elementType string
	object      fyneui.CanvasObject
	signature   string
}

func New(assetRegistry *assets.Registry) *Renderer {
	return &Renderer{
		root:          container.NewWithoutLayout(),
		assetRegistry: assetRegistry,
		elements:      map[string]*renderedElement{},
		resources:     map[string]fyneui.Resource{},
	}
}

func (r *Renderer) SetMinRenderInterval(interval time.Duration) {
	if interval < 0 {
		interval = 0
	}
	r.minRenderInterval = interval
}

func (r *Renderer) CanvasObject() fyneui.CanvasObject {
	return r.root
}

func (r *Renderer) Update(sceneState scene.Scene) error {
	now := time.Now()
	if r.minRenderInterval > 0 && !r.lastRender.IsZero() && now.Sub(r.lastRender) < r.minRenderInterval {
		return nil
	}

	visited := map[string]bool{}
	objects := make([]fyneui.CanvasObject, 0, len(sceneState.Elements))
	for index, element := range sceneState.Elements {
		object, err := r.renderCachedElement(element, rootOccurrenceKey(index, element), visited)
		if err != nil {
			return err
		}
		if object != nil {
			objects = append(objects, object)
		}
	}

	if canvasObjectsChanged(r.root.Objects, objects) {
		r.root.Objects = objects
		r.root.Refresh()
	}
	r.pruneElements(visited)
	r.lastRender = now
	return nil
}

func (r *Renderer) renderCachedElement(element scene.Element, cacheKey string, visited map[string]bool) (fyneui.CanvasObject, error) {
	visited[cacheKey] = true

	if !element.Visible {
		r.elements[cacheKey] = &renderedElement{
			id:          element.ID,
			elementType: element.Type,
			signature:   elementSignature(element),
		}
		return nil, nil
	}

	if element.Type == config.DashboardBlockGroup {
		return r.updateGroupObject(element, cacheKey, visited)
	}

	signature := elementSignature(element)
	cached := r.elements[cacheKey]
	if cached != nil && cached.elementType == element.Type && cached.signature == signature && cached.object != nil {
		return cached.object, nil
	}

	if cached != nil && cached.elementType == element.Type && cached.object != nil {
		if updated, ok, err := r.updateExistingObject(cached.object, element); ok || err != nil {
			if err != nil {
				return nil, err
			}
			cached.signature = signature
			return updated, nil
		}
	}

	object, err := r.buildElementObject(element)
	if err != nil {
		return nil, err
	}
	r.elements[cacheKey] = &renderedElement{
		id:          element.ID,
		elementType: element.Type,
		object:      object,
		signature:   signature,
	}
	return object, nil
}

func (r *Renderer) updateExistingObject(object fyneui.CanvasObject, element scene.Element) (fyneui.CanvasObject, bool, error) {
	switch element.Type {
	case config.DashboardBlockImage:
		asset, err := r.requireImageAsset(element.AssetID)
		if err != nil {
			return nil, true, fmt.Errorf("element %q: %w", element.ID, err)
		}
		image, ok := object.(*canvas.Image)
		if !ok {
			return nil, false, nil
		}
		r.updateImageObject(image, asset.Path, asset.Data, element.Geometry)
		return image, true, nil
	case config.DashboardBlockSpriteFrame:
		if !element.HasFrame {
			return nil, true, fmt.Errorf("element %q has no resolved frame", element.ID)
		}
		image, ok := object.(*canvas.Image)
		if !ok {
			return nil, false, nil
		}
		r.updateImageObject(image, frameResourceName(element), element.Frame.Data, element.Geometry)
		return image, true, nil
	case config.DashboardBlockSpriteText:
		text, ok := object.(*fyneui.Container)
		if !ok {
			return nil, false, nil
		}
		if r.updateSpriteTextObject(text, element) {
			return text, true, nil
		}
		return nil, false, nil
	default:
		return nil, false, nil
	}
}

func (r *Renderer) buildElementObject(element scene.Element) (fyneui.CanvasObject, error) {
	switch element.Type {
	case config.DashboardBlockImage:
		asset, err := r.requireImageAsset(element.AssetID)
		if err != nil {
			return nil, fmt.Errorf("element %q: %w", element.ID, err)
		}
		return r.imageObject(asset.Path, asset.Data, element.Geometry), nil
	case config.DashboardBlockSpriteFrame:
		if !element.HasFrame {
			return nil, fmt.Errorf("element %q has no resolved frame", element.ID)
		}
		return r.imageObject(frameResourceName(element), element.Frame.Data, element.Geometry), nil
	case config.DashboardBlockSpriteText:
		return r.spriteTextObject(element)
	default:
		return nil, fmt.Errorf("element %q type %q is not supported by Fyne scene renderer", element.ID, element.Type)
	}
}

func (r *Renderer) requireImageAsset(assetID string) (assets.Asset, error) {
	if r.assetRegistry == nil {
		return assets.Asset{}, fmt.Errorf("asset registry must not be nil")
	}
	asset, err := r.assetRegistry.MustGet(assetID)
	if err != nil {
		return assets.Asset{}, err
	}
	if asset.Type != assets.TypeImage {
		return assets.Asset{}, fmt.Errorf("asset %q type is %q, want %q", assetID, asset.Type, assets.TypeImage)
	}
	return asset, nil
}

func (r *Renderer) updateGroupObject(element scene.Element, cacheKey string, visited map[string]bool) (fyneui.CanvasObject, error) {
	signature := elementSignature(element)
	cached := r.elements[cacheKey]
	group, ok := cachedContainer(cached)
	if !ok {
		group = container.NewWithoutLayout()
	}

	children := make([]fyneui.CanvasObject, 0, len(element.Children))
	for index, child := range element.Children {
		object, err := r.renderCachedElement(child, childOccurrenceKey(cacheKey, index, child), visited)
		if err != nil {
			return nil, err
		}
		if object != nil {
			children = append(children, object)
		}
	}

	if canvasObjectsChanged(group.Objects, children) {
		group.Objects = children
		group.Refresh()
	}
	applyGeometry(group, element.Geometry)

	r.elements[cacheKey] = &renderedElement{
		id:          element.ID,
		elementType: element.Type,
		object:      group,
		signature:   signature,
	}
	return group, nil
}

func cachedContainer(cached *renderedElement) (*fyneui.Container, bool) {
	if cached == nil || cached.object == nil {
		return nil, false
	}
	group, ok := cached.object.(*fyneui.Container)
	return group, ok
}

func (r *Renderer) imageObject(name string, data []byte, geometry config.RectConfig) fyneui.CanvasObject {
	image := canvas.NewImageFromResource(r.resource(name, data))
	image.FillMode = canvas.ImageFillStretch
	applyGeometry(image, geometry)
	return image
}

func (r *Renderer) updateImageObject(image *canvas.Image, name string, data []byte, geometry config.RectConfig) {
	resource := r.resource(name, data)
	if image.Resource == nil || image.Resource.Name() != resource.Name() {
		image.Resource = resource
		image.Refresh()
	}
	image.FillMode = canvas.ImageFillStretch
	applyGeometry(image, geometry)
}

func (r *Renderer) resource(name string, data []byte) fyneui.Resource {
	if resource, ok := r.resources[name]; ok {
		return resource
	}
	resource := fyneui.NewStaticResource(name, data)
	r.resources[name] = resource
	return resource
}

func (r *Renderer) spriteTextObject(element scene.Element) (fyneui.CanvasObject, error) {
	glyphCount := len(element.Glyphs)
	if glyphCount == 0 {
		return container.NewWithoutLayout(), nil
	}

	glyphWidth := element.Geometry.Width / float64(glyphCount)
	objects := make([]fyneui.CanvasObject, 0, glyphCount)
	for index, glyph := range element.Glyphs {
		glyphObject := r.imageObject(glyphResourceName(element, index, glyph), glyph.Data, config.RectConfig{
			X:      glyphWidth * float64(index),
			Y:      0,
			Width:  glyphWidth,
			Height: element.Geometry.Height,
		})
		objects = append(objects, glyphObject)
	}

	text := container.NewWithoutLayout(objects...)
	applyGeometry(text, element.Geometry)
	return text, nil
}

func (r *Renderer) updateSpriteTextObject(text *fyneui.Container, element scene.Element) bool {
	glyphCount := len(element.Glyphs)
	if len(text.Objects) != glyphCount {
		return false
	}
	if glyphCount == 0 {
		applyGeometry(text, element.Geometry)
		return true
	}

	glyphWidth := element.Geometry.Width / float64(glyphCount)
	for index, glyph := range element.Glyphs {
		image, ok := text.Objects[index].(*canvas.Image)
		if !ok {
			return false
		}
		r.updateImageObject(image, glyphResourceName(element, index, glyph), glyph.Data, config.RectConfig{
			X:      glyphWidth * float64(index),
			Y:      0,
			Width:  glyphWidth,
			Height: element.Geometry.Height,
		})
	}
	applyGeometry(text, element.Geometry)
	return true
}

func (r *Renderer) pruneElements(visited map[string]bool) {
	for id := range r.elements {
		if !visited[id] {
			delete(r.elements, id)
		}
	}
}

func canvasObjectsChanged(current []fyneui.CanvasObject, next []fyneui.CanvasObject) bool {
	if len(current) != len(next) {
		return true
	}
	for i := range current {
		if current[i] != next[i] {
			return true
		}
	}
	return false
}

func applyGeometry(object fyneui.CanvasObject, geometry config.RectConfig) {
	object.Move(fyneui.NewPos(float32(geometry.X), float32(geometry.Y)))
	if geometry.Width > 0 && geometry.Height > 0 {
		object.Resize(fyneui.NewSize(float32(geometry.Width), float32(geometry.Height)))
	}
}

func rootOccurrenceKey(index int, element scene.Element) string {
	return "root[" + strconv.Itoa(index) + "]/" + element.ID
}

func childOccurrenceKey(parentKey string, index int, element scene.Element) string {
	return parentKey + "/children[" + strconv.Itoa(index) + "]/" + element.ID
}

func elementSignature(element scene.Element) string {
	var builder strings.Builder
	builder.WriteString(element.ID)
	builder.WriteByte('|')
	builder.WriteString(element.Type)
	builder.WriteByte('|')
	builder.WriteString(element.LayerID)
	builder.WriteByte('|')
	builder.WriteString(strconv.Itoa(element.Z))
	builder.WriteByte('|')
	builder.WriteString(strconv.FormatBool(element.Visible))
	builder.WriteByte('|')
	writeRectSignature(&builder, element.Geometry)

	if !element.Visible {
		return builder.String()
	}

	switch element.Type {
	case config.DashboardBlockImage:
		builder.WriteString("|asset=")
		builder.WriteString(element.AssetID)
	case config.DashboardBlockSpriteFrame:
		builder.WriteString("|frame=")
		builder.WriteString(strconv.FormatBool(element.HasFrame))
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(element.Frame.Index))
		builder.WriteByte(':')
		builder.WriteString(element.Frame.Path)
	case config.DashboardBlockSpriteText:
		builder.WriteString("|text=")
		builder.WriteString(element.Text)
		builder.WriteString("|glyphs=")
		for index, glyph := range element.Glyphs {
			if index > 0 {
				builder.WriteByte(',')
			}
			builder.WriteString(strconv.Itoa(index))
			builder.WriteByte(':')
			builder.WriteString(glyph.Char)
			builder.WriteByte(':')
			builder.WriteString(glyph.Path)
		}
	case config.DashboardBlockGroup:
		builder.WriteString("|children=")
		for index, child := range element.Children {
			if index > 0 {
				builder.WriteByte(',')
			}
			builder.WriteString(child.ID)
			builder.WriteByte(':')
			builder.WriteString(strconv.FormatBool(child.Visible))
		}
	}

	return builder.String()
}

func writeRectSignature(builder *strings.Builder, geometry config.RectConfig) {
	builder.WriteString(strconv.FormatFloat(geometry.X, 'f', -1, 64))
	builder.WriteByte(',')
	builder.WriteString(strconv.FormatFloat(geometry.Y, 'f', -1, 64))
	builder.WriteByte(',')
	builder.WriteString(strconv.FormatFloat(geometry.Width, 'f', -1, 64))
	builder.WriteByte(',')
	builder.WriteString(strconv.FormatFloat(geometry.Height, 'f', -1, 64))
}

func frameResourceName(element scene.Element) string {
	if element.Frame.Path != "" {
		return element.Frame.Path
	}
	return element.ID + "-frame-" + strconv.Itoa(element.Frame.Index)
}

func glyphResourceName(element scene.Element, index int, glyph assets.Glyph) string {
	if glyph.Path != "" {
		return glyph.Path
	}
	return element.ID + "-glyph-" + strconv.Itoa(index) + "-" + glyph.Char
}
