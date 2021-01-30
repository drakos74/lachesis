package canvas

import (
	"github.com/drakos74/oremi/internal/gui"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

// Container represents a ui scene
type Container struct {
	*gui.InteractiveElement
	CompoundElement
}

func (c *Container) Layout(ops *op.Ops) layout.Dimensions {
	return layout.Dimensions{
		Size:     c.Size(gui.Inset),
		Baseline: 0,
	}
}

// Draw propagates the draw call to all the scene chldren
func (c *Container) Draw(gtx *layout.Context, th *material.Theme) error {
	gtx.Dimensions = c.Layout(gtx.Ops)
	_, err := c.Elements(gtx, DrawAction(gtx, th))
	return err
}

// Event propagates a pointer event to all the scene chldren
func (c *Container) Event(gtx *layout.Context, e *pointer.Event) (f32.Point, bool, error) {
	p, ok, err := c.InteractiveElement.Event(gtx, e)
	if err != nil {
		return p, false, err
	}
	if ok {
		kk, err := c.Elements(gtx, EventAction(p))
		return p, kk, err
	}
	return p, false, nil
}

// NewContainer creates a new scene
func NewContainer(rect *f32.Rectangle) *Container {
	return &Container{
		gui.NewInteractiveElement(rect),
		NewCompoundElement(),
	}
}
