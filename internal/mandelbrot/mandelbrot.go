package mandelbrot

import (
	"fmt"
	"math"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	maxIter  = 100
	infinity = 5
)

var xmin, xmax, ymin, ymax = -2.5, 1.5, -1.5, 1.5

type Mandelbrot struct {
	window *gtk.ApplicationWindow
	da     *gtk.DrawingArea
}

func NewMandelbrot(win *gtk.ApplicationWindow, da *gtk.DrawingArea) *Mandelbrot {
	m := &Mandelbrot{win, da}
	da.Connect("draw", m.Draw)
	win.Connect("key-press-event", m.KeyPress)
	return m
}

func (m *Mandelbrot) Draw(da *gtk.DrawingArea, con *cairo.Context) {
	allocation := da.GetAllocation()
	width, height := allocation.GetWidth(), allocation.GetHeight()

	n := 0
	z := 0 + 0i
	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {
			a := mapp(float64(x), 0, float64(width), xmin, xmax)
			b := mapp(float64(y), 0, float64(height), ymin, ymax)
			c := complex(a, b)

			n = 0
			z = 0

			for n = 0; n < maxIter; n++ {
				z = z*z + c

				if math.Abs(real(z)) > infinity {
					break
				}
			}

			bright := mapp(float64(n), 0, maxIter, 0, 1)
			bright = mapp(math.Sqrt(bright), 0, 1, 0, 255)
			if n == maxIter {
				bright = 0
			}
			con.SetSourceRGBA(bright/255, bright/255, bright/255, 1.0)
			con.Rectangle(float64(x), float64(y), 1.0, 1.0)
			con.Fill()
		}
	}
}

func mapp(v, imin, imax, omin, omax float64) float64 {
	return (v-imin)/(imax-imin)*(omax-omin) + omin
}

func (m *Mandelbrot) KeyPress(_ *gtk.ApplicationWindow, e *gdk.Event) {
	ke := gdk.EventKeyNewFromEvent(e)

	switch ke.KeyVal() {
	case gdk.KEY_q:
		m.quit()
	case gdk.KEY_r:
		reset()
	case gdk.KEY_Page_Up:
		zoomIn()
	case gdk.KEY_Page_Down:
		zoomOut()
	case gdk.KEY_Up:
		ymax -= getYMoveDistance()
		ymin -= getYMoveDistance()
	case gdk.KEY_Down:
		ymax += getYMoveDistance()
		ymin += getYMoveDistance()
	case gdk.KEY_Right:
		xmax += getXMoveDistance()
		xmin += getXMoveDistance()
	case gdk.KEY_Left:
		xmax -= getXMoveDistance()
		xmin -= getXMoveDistance()
	}

	m.da.QueueDraw()
	fmt.Printf("(%0.1f,%0.1f) - (%0.1f,%0.1f)\n", xmin, ymax, xmax, ymin)
}

func (m *Mandelbrot) quit() {
	m.window.Destroy()
}

func zoomIn() {
	d := (ymax - ymin) * 0.1
	ymax -= d
	ymin += d
	d = (xmax - xmin) * 0.1
	xmax -= d
	xmin += d
}

func zoomOut() {
	d := (ymax - ymin) * 0.1
	ymax += d
	ymin -= d
	d = (xmax - xmin) * 0.1
	xmax += d
	xmin -= d
}

func reset() {
	xmin, xmax, ymin, ymax = -2.5, 1.5, -1.5, 1.5
}

func getYMoveDistance() float64 {
	return (ymax - ymin) / 30
}

func getXMoveDistance() float64 {
	return (xmax - xmin) / 30
}
