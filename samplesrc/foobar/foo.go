package foobar

import "math"

// Cart 特に意味のない構造体
type Cart struct {
	x, y float64
}

// GetX x 座標を返す
func (c *Cart) GetX() float64 {
	return c.x
}

// GetY y 座標を返す
func (c *Cart) GetY() float64 {
	return c.y
}

// GetR 中心からの距離を返す
func (c *Cart) GetR() float64 {
	return math.Sqrt(c.x*c.x + c.y*c.y)
}

// Getθ 偏角を返す
func (c *Cart) Getθ() float64 {
	return math.Atan2(c.y, c.x)
}

// Polar 特に意味のない構造体
type Polar struct {
	r, θ float64
}

// GetX x 座標を返す
func (p *Polar) GetX() float64 {
	return p.r * math.Cos(p.θ)
}

// GetY y 座標を返す
func (p *Polar) GetY() float64 {
	return p.r * math.Sin(p.θ)
}

// GetR 中心からの距離を返す
func (p *Polar) GetR() float64 {
	return p.r
}

// Getθ 偏角を返す
func (p *Polar) Getθ() float64 {
	return p.θ
}

// Point 点の interface
type Point interface {
	GetX() float64
	GetY() float64
	GetR() float64
	Getθ() float64
}

// Distance2 二点間の距離の自乗
func Distance2(a, b Point) float64 {
	dx := a.GetX() - b.GetX()
	dy := a.GetY() - b.GetY()
	return dx*dx + dy*dy
}

// Distance 二点間の距離
func Distance(a, b Point) float64 {
	return math.Sqrt(Distance2(a, b))
}
