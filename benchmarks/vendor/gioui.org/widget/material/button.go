// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/internal/f32color"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type ButtonStyle struct {
	Text string
	// Color is the text color.
	Color        color.NRGBA
	Font         text.Font
	TextSize     unit.Value
	Background   color.NRGBA
	CornerRadius unit.Value
	Inset        layout.Inset
	Button       *widget.Clickable
	shaper       text.Shaper
}

type ButtonLayoutStyle struct {
	Background   color.NRGBA
	CornerRadius unit.Value
	Button       *widget.Clickable
}

type IconButtonStyle struct {
	Background color.NRGBA
	// Color is the icon color.
	Color color.NRGBA
	Icon  *widget.Icon
	// Size is the icon size.
	Size   unit.Value
	Inset  layout.Inset
	Button *widget.Clickable
}

func Button(th *Theme, button *widget.Clickable, txt string) ButtonStyle {
	return ButtonStyle{
		Text:         txt,
		Color:        th.Palette.ContrastFg,
		CornerRadius: unit.Dp(4),
		Background:   th.Palette.ContrastBg,
		TextSize:     th.TextSize.Scale(14.0 / 16.0),
		Inset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
		Button: button,
		shaper: th.Shaper,
	}
}

func ButtonLayout(th *Theme, button *widget.Clickable) ButtonLayoutStyle {
	return ButtonLayoutStyle{
		Button:       button,
		Background:   th.Palette.ContrastBg,
		CornerRadius: unit.Dp(4),
	}
}

func IconButton(th *Theme, button *widget.Clickable, icon *widget.Icon) IconButtonStyle {
	return IconButtonStyle{
		Background: th.Palette.ContrastBg,
		Color:      th.Palette.ContrastFg,
		Icon:       icon,
		Size:       unit.Dp(24),
		Inset:      layout.UniformInset(unit.Dp(12)),
		Button:     button,
	}
}

// Clickable lays out a rectangular clickable widget without further
// decoration.
func Clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(button.Layout),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			clip.RRect{
				Rect: f32.Rectangle{Max: layout.FPt(gtx.Constraints.Min)},
			}.Add(gtx.Ops)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}

func (b ButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return ButtonLayoutStyle{
		Background:   b.Background,
		CornerRadius: b.CornerRadius,
		Button:       b.Button,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
			return widget.Label{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
		})
	})
}

func (b ButtonLayoutStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rr := float32(gtx.Px(b.CornerRadius))
			clip.UniformRRect(f32.Rectangle{Max: f32.Point{
				X: float32(gtx.Constraints.Min.X),
				Y: float32(gtx.Constraints.Min.Y),
			}}, rr).Add(gtx.Ops)
			background := b.Background
			switch {
			case gtx.Queue == nil:
				background = f32color.Disabled(b.Background)
			case b.Button.Hovered():
				background = f32color.Hovered(b.Background)
			}
			paint.Fill(gtx.Ops, background)
			for _, c := range b.Button.History() {
				drawInk(gtx, c)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = min
			return layout.Center.Layout(gtx, w)
		}),
		layout.Expanded(b.Button.Layout),
	)
}

func (b IconButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sizex, sizey := gtx.Constraints.Min.X, gtx.Constraints.Min.Y
			sizexf, sizeyf := float32(sizex), float32(sizey)
			rr := (sizexf + sizeyf) * .25
			clip.UniformRRect(f32.Rectangle{
				Max: f32.Point{X: sizexf, Y: sizeyf},
			}, rr).Add(gtx.Ops)
			background := b.Background
			switch {
			case gtx.Queue == nil:
				background = f32color.Disabled(b.Background)
			case b.Button.Hovered():
				background = f32color.Hovered(b.Background)
			}
			paint.Fill(gtx.Ops, background)
			for _, c := range b.Button.History() {
				drawInk(gtx, c)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				size := gtx.Px(b.Size)
				if b.Icon != nil {
					b.Icon.Color = b.Color
					b.Icon.Layout(gtx, unit.Px(float32(size)))
				}
				return layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			pointer.Ellipse(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return b.Button.Layout(gtx)
		}),
	)
}

func drawInk(gtx layout.Context, c widget.Press) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := float32(gtx.Constraints.Min.X)
	if h := float32(gtx.Constraints.Min.Y); h > size {
		size = h
	}
	// Cover the entire constraints min rectangle.
	size *= 2 * float32(math.Sqrt(2))
	// Apply curve values to size and color.
	size *= sizeBezier
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	defer op.Save(gtx.Ops).Load()
	rgba := f32color.MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size * .5
	op.Offset(c.Position.Add(f32.Point{
		X: -rr,
		Y: -rr,
	})).Add(gtx.Ops)
	clip.UniformRRect(f32.Rectangle{Max: f32.Pt(size, size)}, rr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
