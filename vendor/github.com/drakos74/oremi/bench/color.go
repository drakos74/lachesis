package bench

import (
	"encoding/binary"
	"hash/fnv"
	"image/color"
)

type Color struct {
	colors map[string]color.RGBA
}

func Palette(num int) *Color {
	colors := Color{
		colors: make(map[string]color.RGBA),
	}

	return &colors

}

func (c *Color) Get(label string) color.RGBA {
	if col, ok := c.colors[label]; ok {
		return col
	}

	cc := newColor(label)
	c.colors[label] = cc
	return c.colors[label]
}

var Divergence uint8 = 80

const (
	max uint8 = 255
)

func newColor(label string) color.RGBA {
	b := hash(label)
	return color.RGBA{clamp(b[1]), clamp(b[2]), clamp(b[3]), 255}
}

func hash(s string) []uint8 {
	h := fnv.New32a()
	h.Write([]byte(s))
	x := h.Sum32()
	b := make([]uint8, 8)
	binary.PutVarint(b, int64(x))
	return b[:4]
}

func clamp(value uint8) uint8 {
	if value < Divergence {
		return 0
	}
	if value > max {
		return max - Divergence
	}
	return value - Divergence
}
