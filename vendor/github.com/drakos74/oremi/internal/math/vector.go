package math

import (
	"log"
	"math"
)

// Vector defines a point in n dimensional space
type V []float32

func NewV(x ...float32) V {
	return x
}

// Norm returns the norm of the point
// e.g. the distance to the start of the coordinate system
func (v V) Norm() float32 {
	var x2 float64
	for _, x := range v {
		x2 += math.Pow(float64(x), 2)
	}
	return Float32(math.Sqrt(x2))
}

// Dim returns the dimensions of the point
func (v V) Dim() int {
	return len(v)
}

// Distance calculates the distance of the point to another V
func (v V) Distance(p V) float32 {
	dimensionValidator(v, p)
	var d float64
	for i, x := range v {
		d += math.Pow(float64(x-p[i]), 2)
	}
	return Float32(math.Sqrt(d))
}

type Rect struct {
	Min V
	Max V
}

var dimensionValidator = func(e1, e2 V) {
	if e1.Dim() != e2.Dim() {
		log.Fatalf("vectors don't have the same dimension %v vs %v", e1, e2)
	}
}
