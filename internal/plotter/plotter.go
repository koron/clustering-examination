package plotter

import "image"

// Point is 2D point. It is an alias of image.Point.
type Point = image.Point

// Rectangle is 2D rectangle. It is an alias of image.Rectangle.
type Rectangle = image.Rectangle

func rectRegulate(r Rectangle) Rectangle {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func rectWidth(r Rectangle) int {
	return r.Max.X - r.Min.X
}

func rectHeight(r Rectangle) int {
	return r.Max.Y - r.Min.Y
}

func rectSize(r Rectangle) Point {
	return Point{
		X: r.Max.X - r.Min.X,
		Y: r.Max.Y - r.Min.Y,
	}
}

// PointF is 2D point with float64.
type PointF struct {
	X float64
	Y float64
}

func toPointF(p Point) PointF {
	return PointF{X: float64(p.X), Y: float64(p.Y)}
}

// RectangleF is 2D rectangle with float64.
type RectangleF struct {
	Min PointF
	Max PointF
}

func (r RectangleF) width() float64 {
	return r.Max.X - r.Min.X
}

func (r RectangleF) height() float64 {
	return r.Max.Y - r.Min.Y
}

func (r RectangleF) size() PointF {
	return PointF{
		X: r.width(),
		Y: r.height(),
	}
}

func (r RectangleF) regulate() RectangleF {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func (r RectangleF) in(p PointF) bool {
	if p.X < r.Min.X || p.X > r.Max.X || p.Y < r.Min.Y || p.Y > r.Max.Y {
		return false
	}
	return true
}

// Plotter is point converter from float64 (PointF) to int with scaling and
// offsetting.
type Plotter struct {
	bbox   RectangleF // bounding box
	ratio  float64
	offset Point
}

func ratio(dst PointF, src PointF) (float64, bool) {
	ratioX := dst.X / src.X
	ratioY := dst.Y / src.Y
	if ratioY < ratioX {
		return ratioY, false
	}
	return ratioX, true
}

// New creates a Ploatter with destination and source rectangles.
func New(dst Rectangle, src RectangleF) Plotter {
	dstF := rectRegulate(dst)
	srcF := src.regulate()
	ratio, vert := ratio(toPointF(rectSize(dstF)), srcF.size())
	offset := dst.Min
	// centering
	if vert {
		offset.Y += (rectHeight(dst) - int(src.height()*ratio+0.5)) / 2
	} else {
		offset.X += (rectWidth(dst) - int(src.width()*ratio+0.5)) / 2
	}
	return Plotter{
		bbox:   srcF,
		offset: offset,
		ratio:  ratio,
	}
}

// Plot converts PointF to Point.
func (p Plotter) Plot(src PointF) (Point, bool) {
	return Point{
		X: int((src.X-p.bbox.Min.X)*p.ratio+0.5) + p.offset.X,
		Y: int((p.bbox.Max.Y-src.Y)*p.ratio+0.5) + p.offset.Y,
	}, p.bbox.in(src)
}
