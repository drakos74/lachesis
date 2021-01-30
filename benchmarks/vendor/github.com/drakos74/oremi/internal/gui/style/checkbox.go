package style

import (
	"image/color"

	"github.com/drakos74/oremi/internal/gui"
	"github.com/drakos74/oremi/internal/gui/canvas"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// CheckboxControl is a checkbox to be used as a canvas.Control
type CheckboxControl struct {
	gui.RawItem
	label    string
	color    color.RGBA
	checkbox *widget.CheckBox
	active   bool
	trigger  chan canvas.Event
	ack      chan canvas.Event
	enabled  bool
}

func NewCheckBox(label string, active bool, color color.RGBA) *CheckboxControl {
	cb := &CheckboxControl{
		*gui.NewRawItem(),
		label,
		color,
		new(widget.CheckBox),
		active,
		make(chan canvas.Event),
		make(chan canvas.Event),
		true,
	}
	cb.checkbox.SetChecked(active)
	return cb
}

func (c *CheckboxControl) Draw(gtx *layout.Context, th *material.Theme) error {
	if c.enabled {
		theme := material.NewTheme()
		theme.Color.Text = c.color
		theme.CheckBox(c.label).Layout(gtx, c.checkbox)
	}
	active := c.active
	c.active = c.checkbox.Checked(gtx)
	if c.active != active {
		c.trigger <- canvas.Event{canvas.Trigger, c.active, ""}
	}
	return nil
}

func (c *CheckboxControl) Disable() {
	c.enabled = false
}

func (c *CheckboxControl) Enable() {
	c.enabled = true
}

func (c *CheckboxControl) Label() string {
	return c.label
}

func (c *CheckboxControl) Color() color.RGBA {
	return c.color
}

func (c *CheckboxControl) Set(active bool) {
	c.checkbox.SetChecked(active)
}

func (c *CheckboxControl) IsActive() bool {
	return c.active
}

func (c *CheckboxControl) Trigger() canvas.EventReceiver {
	return c.trigger
}

func (c *CheckboxControl) Ack() canvas.EventEmitter {
	return c.ack
}

type CheckboxControlGroup struct {
	CheckboxControl
	cboxes []canvas.Control
}

func NewCheckboxControlGroup(active bool, control ...canvas.Control) *CheckboxControlGroup {
	cb := NewCheckBox("all", active, color.RGBA{})
	group := &CheckboxControlGroup{
		CheckboxControl: *cb,
		cboxes:          control,
	}

	go func() {
		for {
			select {
			case <-cb.Trigger():
				for _, checkbox := range group.cboxes {
					checkbox.Set(group.CheckboxControl.active)
				}
			}
		}
	}()

	return group
}
