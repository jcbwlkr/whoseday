// An app that says which of my kids gets to pick the first bedtime story for the day.
package main

import (
	"image"
	"log"
	"time"

	_ "image/jpeg"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
)

const (
	width  = 100
	height = 100
)

var (
	startTime = time.Now()

	images *glutil.Images
	eng    sprite.Engine
	scene  *sprite.Node

	sz size.Event
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx)
				a.Publish()
				a.Send(paint.Event{}) // keep animating
			}
		}
	})
}

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	loadScene()
}

func onStop() {
	eng.Release()
	images.Release()
}

func onPaint(glctx gl.Context) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, sz)
}

func newNode() *sprite.Node {
	n := &sprite.Node{}
	eng.Register(n)
	scene.AppendChild(n)
	return n
}

func loadScene() {
	kid := loadKid()
	scene = &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	var x, y float32
	dx, dy := float32(1), float32(1)

	n := newNode()
	// TODO do I need this arranger thing?
	n.Arranger = arrangerFunc(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, kid)

		if x < 0 {
			dx = 1
		}
		if y < 0 {
			dy = 1
		}
		if x+width > float32(sz.WidthPt) {
			dx = -1
		}
		if y+height > float32(sz.HeightPt) {
			dy = -1
		}

		x += dx
		y += dy
		eng.SetTransform(n, f32.Affine{
			{width, 0, x},
			{0, height, y},
		})
	})
}

func loadKid() sprite.SubTex {
	file := picForDay(time.Now())

	a, err := asset.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	img, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(img)
	if err != nil {
		log.Fatal(err)
	}

	return sprite.SubTex{T: t, R: image.Rect(0, 0, 466, 466)}
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

func picForDay(day time.Time) string {
	d := day.Sub(time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local))

	if int(d.Hours()/24)%2 == 0 {
		return "carter.jpeg"
	}

	return "kell.jpeg"
}
