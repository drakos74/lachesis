package model

import (
	"log"
	"math"
)

// TODO : unify with math package

// Vector defines a point in n dimensional space
type Vector struct {
	Label  []string
	Coords []float64
}

// NewVector creates a new point at the specified coordinates
func NewVector(label []string, x ...float64) Vector {
	return Vector{Label: label, Coords: x}
}

// TODO : consider removing this abstraction (see usages)
type V interface {
	ToVector() (Vector, error)
}

// Norm returns the norm of the point
// e.g. the distance to the start of the coordinate system
func (p Vector) Norm() float64 {
	var x2 float64
	for _, x := range p.Coords {
		x2 += math.Pow(x, 2)
	}
	return math.Sqrt(x2)
}

// Dim returns the dimensions of the point
func (p Vector) Dim() int {
	return len(p.Coords)
}

// Distance calculates the distance of the point to another vector
func (p Vector) Distance(element Vector) float64 {
	dimensionValidator(p, element)
	var d float64
	for i, x := range p.Coords {
		d += math.Pow(x-p.Coords[i], 2)
	}
	return math.Sqrt(d)
}

// Iterator iterates over all vectors it is used to read the records of a collection
type Iterator interface {
	Next() (vector Vector, ok, hasNext bool)
	Reset()
}

// Collection represents a collection of data points
type Collection interface {
	Iterator
	Size() int
	Edge() (min, max Vector)
	Labels() []string
}

// TODO : consider removing this abstraction (see usages)
type C interface {
	ToCollection() (Collection, error)
}

var dimensionValidator = func(e1, e2 Vector) {
	if e1.Dim() != e2.Dim() {
		log.Fatalf("vectors don't have the same dimension %v vs %v", e1, e2)
	}
}

type Filter func(Collection) bool

var Size Filter = func(collection Collection) bool {
	return collection.Size() > 0
}
