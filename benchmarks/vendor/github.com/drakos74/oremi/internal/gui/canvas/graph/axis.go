package graph

import (
	"fmt"

	"github.com/drakos74/oremi/internal/gui"
	"github.com/drakos74/oremi/internal/gui/canvas"
	"github.com/drakos74/oremi/internal/gui/style"
	"github.com/drakos74/oremi/internal/math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

// Axis is an axis element for graphs
type Axis struct {
	gui.Item
	gui.Area
	start  f32.Point
	length float32
	layout.Axis
	label style.Label
}

// NewAxisX creates a new x axis
func NewAxisX(label string, start f32.Point, length float32) *Axis {
	rect := &f32.Rectangle{
		Min: start,
		Max: f32.Point{
			X: start.X + length,
			Y: start.Y + 1,
		},
	}
	axis := &Axis{
		gui.NewRawItem(),
		gui.Rect(rect),
		start,
		length,
		layout.Horizontal,
		style.NewLabel(f32.Point{
			X: start.X + length,
			Y: start.Y + 20,
		}, label),
	}

	return axis
}

func (axis Axis) Delimiters(delim int, calc math.Mapper) []*Delimiter {
	delimiters := make([]*Delimiter, delim+1)
	for i := 0; i <= delim; i++ {
		d := float32(i) / float32(delim)
		switch axis.Axis {
		case layout.Horizontal:
			delimiters[i] = NewDelimiterX(
				f32.Point{
					X: axis.start.X + axis.length*d,
					Y: axis.start.Y,
				},
				calc.DeScaleAt(0, math.Normal))
		case layout.Vertical:
			delimiters[i] = NewDelimiterY(
				f32.Point{
					X: axis.start.X,
					Y: axis.start.Y + axis.length*d,
				},
				calc.DeScaleAt(1, math.Inverse))
		}
	}
	return delimiters
}

// NewAxisY creates a new y axis
func NewAxisY(label string, start f32.Point, length float32) *Axis {
	rect := &f32.Rectangle{
		Min: start,
		Max: f32.Point{
			X: start.X + 1,
			Y: start.Y + length,
		},
	}
	axis := &Axis{
		gui.NewRawItem(),
		gui.Rect(rect),
		start,
		length,
		layout.Vertical,
		style.NewLabel(f32.Point{
			X: start.X - 20,
			Y: start.Y - 2,
		}, label),
	}
	return axis
}

// Draw draws the axis
func (a *Axis) Draw(gtx *layout.Context, th *material.Theme) error {
	paint.ColorOp{Color: style.Black}.Add(gtx.Ops)
	paint.PaintOp{Rect: a.Rect()}.Add(gtx.Ops)
	return a.label.Draw(gtx, th)
}

// Delimiter is an axis child element representing a value on the respective axis
type Delimiter struct {
	gui.Item
	canvas.DynamicElement
	label     style.Label
	transform func() float32
}

// NewDelimiterX creates a new delimiter for an x axis
func NewDelimiterX(p f32.Point, transform math.Transform) *Delimiter {
	rect := &f32.Rectangle{
		Min: f32.Point{
			X: p.X,
			Y: p.Y - 10,
		},
		Max: f32.Point{
			X: p.X + 1,
			Y: p.Y + 10,
		},
	}
	return &Delimiter{
		gui.NewRawItem(),
		canvas.NewDynamicElement(rect),
		style.NewLabel(p, ""),
		func() float32 {
			return transform(p.X)
		},
	}
}

// NewDelimiterY creates a new delimiter for an x axis
func NewDelimiterY(p f32.Point, transform math.Transform) *Delimiter {
	rect := &f32.Rectangle{
		Min: f32.Point{
			X: p.X - 10,
			Y: p.Y,
		},
		Max: f32.Point{
			X: p.X + 10,
			Y: p.Y + 1,
		},
	}
	return &Delimiter{
		gui.NewRawItem(),
		canvas.NewDynamicElement(rect),
		style.NewLabel(p, ""),
		func() float32 {
			return transform(p.Y)
		},
	}
}

// Draw draws the delimiter
func (d *Delimiter) Draw(gtx *layout.Context, th *material.Theme) error {
	if d.IsActive() {
		d.label.Text(fmt.Sprintf("%v", d.transform()))
		err := d.label.Draw(gtx, th)
		if err != nil {
			return err
		}
	}
	paint.ColorOp{Color: style.Black}.Add(gtx.Ops)
	paint.PaintOp{Rect: d.Rect()}.Add(gtx.Ops)
	return nil
}
