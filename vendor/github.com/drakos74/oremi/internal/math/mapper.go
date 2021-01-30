package math

import (
	"fmt"
	"log"
	"math"
)

type ScaleFactor int

const (
	Normal ScaleFactor = iota + 1
	Inverse
)

type Transform func(x float32) float32

func voidTransform(x float32) float32 {
	return x
}

type Mapper interface {
	ScaleAt(i int, factor ScaleFactor) Transform
	DeScaleAt(i int, factor ScaleFactor) Transform
}

type StackedMapper struct {
	stack []Mapper
}

func NewStackedMapper(stack ...Mapper) *StackedMapper {
	return &StackedMapper{stack: stack}
}

func (s StackedMapper) ScaleAt(i int, factor ScaleFactor) Transform {
	return func(x float32) float32 {
		for _, s := range s.stack {
			x = s.ScaleAt(i, factor)(x)
		}
		return x
	}
}

func (s StackedMapper) DeScaleAt(i int, factor ScaleFactor) Transform {
	return func(sx float32) float32 {
		for _, s := range s.stack {
			sx = s.DeScaleAt(i, factor)(sx)
		}
		return sx
	}
}

type VoidCalcMapper struct {
}

func (v VoidCalcMapper) ScaleAt(i int, factor ScaleFactor) Transform {
	return voidTransform
}

func (v VoidCalcMapper) DeScaleAt(i int, factor ScaleFactor) Transform {
	return voidTransform
}

type CoordinateMapper struct {
	scale float32
	rect  *Rect
}

func (c CoordinateMapper) ScaleAt(i int, factor ScaleFactor) Transform {
	switch factor {
	case Normal:
		return func(x float32) float32 {
			return scaleAt(i, *c.rect, c.scale, x)
		}
	case Inverse:
		return func(x float32) float32 {
			return scaleInvAt(i, *c.rect, c.scale, x)
		}
	default:
		log.Fatalf("scaleFactor not recognised: %v", factor)
		return voidTransform
	}
}

func (c CoordinateMapper) DeScaleAt(i int, factor ScaleFactor) Transform {
	switch factor {
	case Normal:
		return func(x float32) float32 {
			return deScaleAt(i, *c.rect, c.scale, x)
		}
	case Inverse:
		return func(x float32) float32 {
			return deScaleInvAt(i, *c.rect, c.scale, x)
		}
	default:
		log.Fatalf("deScaleFactor not recognised: %v", factor)
		return voidTransform
	}
}

// scaleX calculates the 'real' x - coordinate of a relative value to the grid
func scaleAt(i int, rect Rect, scale, value float32) float32 {
	return rect.Min[i] + ((rect.Max[i] - rect.Min[i]) * value / scale)
}

// deScaleX calculates the 'relative' x - coordinate of a 'real' value
func deScaleAt(i int, rect Rect, scale, value float32) float32 {
	return (value - rect.Min[i]) / safe(rect.Max[i]-rect.Min[i]) * scale
}

// scaleY calculates the 'real' y - coordinate of a relative value to the grid
func scaleInvAt(i int, rect Rect, scale, value float32) float32 {
	return rect.Max[i] - ((rect.Max[i] - rect.Min[i]) * value / scale)
}

// deScaleY calculates the 'relative' y - coordinate of a 'real' value
func deScaleInvAt(i int, rect Rect, scale, value float32) float32 {
	return scale - (value-rect.Min[i])/safe(rect.Max[i]-rect.Min[i])*scale
}

// MonotonicMapper scales values for x and y linearly to certain ranges and vice versa
type MonotonicMapper struct {
	*Rect
	scale float32
}

// NewMonotonicMapper creates a new linearly scale calculation element
func NewMonotonicMapper(scale float32) *MonotonicMapper {
	return newMonotonicMapper(scale, NewV(math.MaxFloat32, math.MaxFloat32), NewV(0, 0))
}

func (l *MonotonicMapper) Max(max V) bool {
	var recalc bool
	for i, v := range max {
		if v > l.Rect.Max[i] {
			l.Rect.Max[i] = v
			recalc = true
		}
	}
	return recalc
}

func (l *MonotonicMapper) Min(min V) bool {
	var recalc bool
	for i, v := range min {
		if v < l.Rect.Min[i] {
			l.Rect.Min[i] = v
			recalc = true
		}
	}
	return recalc
}

// newMonotonicMapper creates a new linearly scale calculation element
func newMonotonicMapper(scale float32, min, max V) *MonotonicMapper {
	return &MonotonicMapper{
		Rect: &Rect{
			Min: min,
			Max: max,
		},
		scale: scale,
	}
}

func (l MonotonicMapper) DeScaleAt(i int, factor ScaleFactor) Transform {
	switch factor {
	case Normal:
		fallthrough
	case Inverse:
		return func(x float32) float32 {
			return x/l.scale*(l.Rect.Max[i]-l.Rect.Min[i]) + l.Rect.Min[i]
		}
	default:
		log.Fatalf("scale factor %v not supported for linear mapper", factor)
		return voidTransform
	}
}

func (l MonotonicMapper) ScaleAt(i int, factor ScaleFactor) Transform {
	switch factor {
	case Normal:
		fallthrough
	case Inverse:
		return func(sx float32) float32 {
			return l.scale * (sx - l.Rect.Min[i]) / safe(l.Rect.Max[i]-l.Rect.Min[i])
		}
	default:
		log.Fatalf("scale factor %v not supported for linear mapper", factor)
		return voidTransform
	}
}

func NewRawCalcElement(rect *Rect, scale float32) *CoordinateMapper {
	return &CoordinateMapper{
		scale: scale,
		rect:  rect,
	}
}

// safe makes sure that we dont encounter NaN when dividing by '0'
func safe(f float32) float32 {
	if f == 0 {
		return 1
	}
	return f
}

const PrecisionThreshold = 1

// Float32 converts to a float32 and panics if there is loss of precision
func Float32(f float64) float32 {
	x := float32(f)
	l := math.Abs(float64(x) - f)
	if l > PrecisionThreshold {
		println(fmt.Sprintf("precision loss for f64:%f vs f32:%f is %v", f, x, l))
	}
	return x
}
