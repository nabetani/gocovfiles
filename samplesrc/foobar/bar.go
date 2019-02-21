package foobar

import "math"

// このソースは、カバレッジ測定の対象外にしたいという設定。

// Vec 点a から 点b へのベクトル
func Vec(a, b Point) *Cart {
	return &Cart{
		x: b.GetX() - a.GetX(),
		y: b.GetY() - a.GetY(),
	}
}

// TriangleSize1 三角形の面積その1
func TriangleSize1(a, b, c Point) float64 {
	ab := Vec(a, b)
	ac := Vec(a, c)
	return 0.5 * math.Abs(ab.GetX()*ac.GetY()-ab.GetY()*ac.GetX())
}

// TriangleSize2 三角形の面積その2
func TriangleSize2(pa, pb, pc Point) float64 {
	a := Distance(pb, pc)
	b := Distance(pc, pa)
	c := Distance(pa, pb)
	s := (a + b + c) / 2
	return math.Sqrt(s * (s - a) * (s - b) * (s - c))
}
