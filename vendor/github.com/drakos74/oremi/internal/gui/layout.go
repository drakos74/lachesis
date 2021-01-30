package gui

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type View struct {
	itemsList
	*layout.List
	height float32
	width  float32
}

func NewView(orientation layout.Axis) *View {
	return &View{
		itemsList: itemsList{
			items: make(map[uint32]item),
			index: make([]uint32, 0),
		},
		List: &layout.List{
			Axis: orientation,
		},
	}
}

func (v *View) Draw(gtx *layout.Context, th *material.Theme) error {
	// TODO : do recover
	//
	//children := make([]layout.FlexChild, len(v.items))
	//
	//for i := 0; i < len(v.items); i++ {
	//	println(fmt.Sprintf("i = %v", i))
	//	children[i] = layout.Rigid(func(j int) func() {
	//		return func() {
	//			if v.height > 0 {
	//				gtx.Constraints.Height.Max = gtx.Px(unit.Dp(v.height))
	//			}
	//			if v.width > 0 {
	//				gtx.Constraints.Width.Max = gtx.Px(unit.Dp(v.width))
	//			}
	//			layout.UniformInset(unit.Dp(0)).Layout(gtx, v.get(j).draw(gtx, th))
	//		}
	//	}(i))
	//}
	//
	//layout.Flex{Alignment: layout.Start}.Layout(gtx, children...)

	v.Layout(gtx, len(v.items), func(i int) {
		if v.height > 0 {
			gtx.Constraints.Height.Max = gtx.Px(unit.Dp(v.height))
		}
		if v.width > 0 {
			gtx.Constraints.Width.Max = gtx.Px(unit.Dp(v.width))
		}
		layout.UniformInset(unit.Dp(0)).Layout(gtx, v.get(i).draw(gtx, th))
	})

	return nil
}

func (v *View) Event(gtx *layout.Context, e *pointer.Event) (redraw bool, err error) {
	for i := 0; i < len(v.items); i++ {
		if v.get(i).event(gtx, e) {
			redraw = true
		}
	}
	return redraw, nil
}

func (v *View) WithMaxHeight(height float32) *View {
	v.height = height
	return v
}

func (v *View) WithMaxWidth(width float32) *View {
	v.width = width
	return v
}
