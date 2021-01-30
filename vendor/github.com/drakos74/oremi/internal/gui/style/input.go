package style

import (
	"github.com/drakos74/oremi/internal/gui/canvas"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Input struct {
	text    string
	editor  *widget.Editor
	trigger chan canvas.Event
}

func NewInput() *Input {
	return &Input{"",
		&widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		make(canvas.Events)}
}

func (i *Input) Draw(gtx *layout.Context, th *material.Theme) error {
	e := th.Editor("Hint")
	e.Font.Style = text.Italic
	e.Layout(gtx, i.editor)
	for _, e := range i.editor.Events(gtx) {
		if _, ok := e.(widget.SubmitEvent); ok {
			i.editor.SetText("")
		}
	}
	if i.text != i.editor.Text() {
		i.text = i.editor.Text()
		i.trigger <- canvas.Event{
			T: canvas.Trigger,
			A: false,
			S: i.text,
		}
	}
	return nil
}

func (i *Input) IsActive() bool {
	panic("implement me")
}

func (i *Input) Set(active bool) {
	panic("implement me")
}

func (i *Input) Trigger() canvas.EventReceiver {
	return i.trigger
}

func (i *Input) Ack() canvas.EventEmitter {
	panic("implement me")
}
