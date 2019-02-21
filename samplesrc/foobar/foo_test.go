package foobar

import (
	"math"
	"testing"
)

func TestCartGetX(t *testing.T) {
	c := Cart{x: 1, y: 2}
	if c.GetX() != 1 {
		t.Errorf("c.GetX()==%v, want 1", c.GetX())
	}
}

func TestCartGetY(t *testing.T) {
	c := Cart{x: 1, y: 2}
	if c.GetY() != 2 {
		t.Errorf("c.GetY()==%v, want 2", c.GetY())
	}
}

func TestCartGetR(t *testing.T) {
	c := Cart{x: 3, y: 4}
	expected := float64(5)
	if c.GetR() != expected {
		t.Errorf("c.GetR()==%v, want %v", c.GetR(), expected)
	}
}

func TestCartGetθ(t *testing.T) {
	c := Cart{x: 3, y: 4}
	expected := math.Atan2(4, 3)
	if c.Getθ() != expected {
		t.Errorf("c.Getθ()==%v, want %v", c.Getθ(), expected)
	}
}

func TestPolarGetXY(t *testing.T) {
	type dataType struct {
		src  Polar
		x, y float64
	}
	data := []dataType{
		dataType{x: 1, y: 0, src: Polar{r: 1, θ: 0}},
		dataType{x: 0, y: 1, src: Polar{r: 1, θ: math.Pi / 2}},
		dataType{x: -1, y: 0, src: Polar{r: 1, θ: math.Pi}},
		dataType{x: 0, y: -1, src: Polar{r: 1, θ: math.Pi * 3 / 2}},
		dataType{x: 1, y: 0, src: Polar{r: 1, θ: math.Pi * 2}},
	}
	for _, d := range data {
		x := d.src.GetX()
		if 1e-10 < math.Abs(x-d.x) {
			t.Errorf("d.src.GetX()=%v, want almost %v", x, d.x)
		}
		y := d.src.GetY()
		if 1e-10 < math.Abs(y-d.y) {
			t.Errorf("d.src.GetY()=%v, want almost %v", y, d.y)
		}
	}
}

func TestPolarGetR(t *testing.T) {
	p := Polar{r: 1, θ: 2}
	if p.GetR() != 1 {
		t.Errorf("p.GetR()==%v, want 1", p.GetR())
	}
}

func TestPolarGetθ(t *testing.T) {
	p := Polar{r: 1, θ: 2}
	if p.Getθ() != 2 {
		t.Errorf("p.Getθ()==%v, want 2", p.Getθ())
	}
}
