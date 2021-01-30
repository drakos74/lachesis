package model

import (
	"github.com/drakos74/oremi/internal/gui/style"

	"github.com/drakos74/oremi/internal/data/model"
	"github.com/drakos74/oremi/internal/math"

	"gioui.org/f32"
)

type LabeledPoint struct {
	f32.Point
	Label []string
}

// TODO : embed the data model collection into the interface
type Collection interface {
	Bounds() *f32.Rectangle
	Next() (point *LabeledPoint, ok, next bool)
	Size() int
	Reset()
	Labels() []string
	Style() style.Properties
}

type Series struct {
	model.Collection
	style style.Properties
}

func (s Series) Style() style.Properties {
	return s.style
}

func NewSeries(collection model.Collection, style style.Properties) Collection {
	return Series{Collection: collection, style: style}
}

func (s Series) Bounds() *f32.Rectangle {
	min, max := s.Edge()
	return &f32.Rectangle{
		Min: f32.Point{
			X: float32(min.Coords[0]),
			Y: float32(min.Coords[1]),
		},
		Max: f32.Point{
			X: float32(max.Coords[0]),
			Y: float32(max.Coords[1]),
		},
	}
}

func (s Series) Next() (point *LabeledPoint, ok, next bool) {
	if p, ok, next := s.Collection.Next(); ok {
		return &LabeledPoint{
			// TODO : make the coordinate choice connected to the labels and the graph options in general
			Point: f32.Point{
				X: math.Float32(p.Coords[0]),
				Y: math.Float32(p.Coords[1]),
			},
			Label: p.Label,
		}, true, next
	}
	return nil, false, false
}

func (s Series) Reset() {
	s.Collection.Reset()
}
