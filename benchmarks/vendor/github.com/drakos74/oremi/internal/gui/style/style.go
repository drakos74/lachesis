package style

import "image/color"

var Black = color.RGBA{0, 0, 0, 255}

type Properties struct {
	color.RGBA
}

func (p *Properties) Color(rgba color.RGBA) *Properties {
	p.Color(rgba)
	return p
}
