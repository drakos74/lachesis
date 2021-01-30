package style

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Label struct {
	Properties
	position f32.Point
	reset    f32.Point
	text     string
}

func NewLabel(center f32.Point, text string) Label {
	c := f32.Point{
		X: center.X + 10,
		Y: center.Y - 40,
	}
	return Label{
		position: c,
		reset:    neg(c),
		text:     text,
	}
}

func neg(p f32.Point) f32.Point {
	return f32.Point{
		X: -1 * p.X,
		Y: -1 * p.Y,
	}
}

// Text adjusts the labels text
func (l *Label) Text(text string) {
	// TODO : given that we might to center always to the pointer,
	// we might in the future adjust the position as well
	l.text = text
}

// Position adjusts the label position
func (l *Label) Position(position f32.Point) {
	l.position = position
	l.reset = neg(position)
}

func (l Label) Draw(gtx *layout.Context, th *material.Theme) error {
	// TODO : try with StackOp
	dim := gtx.Dimensions
	op.TransformOp{}.Offset(l.position).Add(gtx.Ops)
	th.Label(unit.Px(30), l.text).Layout(gtx)
	op.TransformOp{}.Offset(l.reset).Add(gtx.Ops)
	gtx.Dimensions = dim
	return nil
}
