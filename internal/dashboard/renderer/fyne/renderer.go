package fyne

import (
	"fmt"

	fyneui "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/assets"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scene"
)

type Renderer struct {
	root          *fyneui.Container
	assetRegistry *assets.Registry
}

func New(assetRegistry *assets.Registry) *Renderer {
	return &Renderer{
		root:          container.NewWithoutLayout(),
		assetRegistry: assetRegistry,
	}
}

func (r *Renderer) CanvasObject() fyneui.CanvasObject {
	return r.root
}

func (r *Renderer) Update(sceneState scene.Scene) error {
	objects := make([]fyneui.CanvasObject, 0, len(sceneState.Elements))
	for _, element := range sceneState.Elements {
		object, err := r.renderElement(element)
		if err != nil {
			return err
		}
		if object != nil {
			objects = append(objects, object)
		}
	}

	r.root.Objects = objects
	r.root.Refresh()
	return nil
}

func (r *Renderer) renderElement(element scene.Element) (fyneui.CanvasObject, error) {
	if !element.Visible {
		return nil, nil
	}

	switch element.Type {
	case config.DashboardBlockImage:
		asset, err := r.requireImageAsset(element.AssetID)
		if err != nil {
			return nil, fmt.Errorf("element %q: %w", element.ID, err)
		}
		return imageObject(asset.Path, asset.Data, element.Geometry), nil
	case config.DashboardBlockSpriteFrame:
		if !element.HasFrame {
			return nil, fmt.Errorf("element %q has no resolved frame", element.ID)
		}
		return imageObject(frameResourceName(element), element.Frame.Data, element.Geometry), nil
	case config.DashboardBlockSpriteText:
		return spriteTextObject(element)
	case config.DashboardBlockGroup:
		return r.groupObject(element)
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

func (r *Renderer) groupObject(element scene.Element) (fyneui.CanvasObject, error) {
	children := make([]fyneui.CanvasObject, 0, len(element.Children))
	for _, child := range element.Children {
		object, err := r.renderElement(child)
		if err != nil {
			return nil, err
		}
		if object != nil {
			children = append(children, object)
		}
	}

	group := container.NewWithoutLayout(children...)
	applyGeometry(group, element.Geometry)
	return group, nil
}

func imageObject(name string, data []byte, geometry config.RectConfig) fyneui.CanvasObject {
	image := canvas.NewImageFromResource(fyneui.NewStaticResource(name, data))
	image.FillMode = canvas.ImageFillStretch
	applyGeometry(image, geometry)
	return image
}

func spriteTextObject(element scene.Element) (fyneui.CanvasObject, error) {
	glyphCount := len(element.Glyphs)
	if glyphCount == 0 {
		return container.NewWithoutLayout(), nil
	}

	glyphWidth := element.Geometry.Width / float64(glyphCount)
	objects := make([]fyneui.CanvasObject, 0, glyphCount)
	for index, glyph := range element.Glyphs {
		glyphObject := imageObject(glyphResourceName(element, index, glyph), glyph.Data, config.RectConfig{
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

func applyGeometry(object fyneui.CanvasObject, geometry config.RectConfig) {
	object.Move(fyneui.NewPos(float32(geometry.X), float32(geometry.Y)))
	if geometry.Width > 0 && geometry.Height > 0 {
		object.Resize(fyneui.NewSize(float32(geometry.Width), float32(geometry.Height)))
	}
}

func frameResourceName(element scene.Element) string {
	if element.Frame.Path != "" {
		return element.Frame.Path
	}
	return fmt.Sprintf("%s-frame-%d", element.ID, element.Frame.Index)
}

func glyphResourceName(element scene.Element, index int, glyph assets.Glyph) string {
	if glyph.Path != "" {
		return glyph.Path
	}
	return fmt.Sprintf("%s-glyph-%d-%s", element.ID, index, glyph.Char)
}
