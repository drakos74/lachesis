package graph

import (
	"image/color"

	"github.com/drakos74/oremi/internal/config"

	"github.com/drakos74/oremi/internal/gui"
	"github.com/drakos74/oremi/internal/gui/canvas"
	"github.com/drakos74/oremi/internal/gui/style"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

// AddPoint is a point element
type Point struct {
	gui.Item
	*canvas.RawDynamicElement
	w     float32
	c     f32.Point
	label style.Label
	color color.RGBA
}

// NewPoint creates a new point
func NewPoint(label string, center f32.Point, cc color.RGBA) *Point {
	var w float32 = config.PointWidth
	rect := calculateRect(center, w)
	p := &Point{
		gui.NewRawItem(),
		canvas.NewDynamicElement(&rect),
		w,
		center,
		style.NewLabel(center, label),
		cc,
	}
	return p
}

func calculateRect(center f32.Point, w float32) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: center.X - w, Y: center.Y - w},
		Max: f32.Point{X: center.X + w, Y: center.Y + w},
	}
}

// Draw draws the point on the canvas
func (p *Point) Draw(gtx *layout.Context, th *material.Theme) error {
	r := p.Rect()
	if p.IsActive() {
		r = calculateRect(p.c, 2*p.w)
		err := p.label.Draw(gtx, th)
		if err != nil {
			return err
		}
	}
	paint.ColorOp{Color: p.color}.Add(gtx.Ops)
	paint.PaintOp{Rect: r}.Add(gtx.Ops)
	return nil
}
