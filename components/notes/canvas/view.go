package canvas

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"

	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/text"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	/*
		"github.com/spinsrv/browser"
		"github.com/spinsrv/browser/ui"
		"github.com/tdewolff/canvas"
		"github.com/tdewolff/canvas/text"
		"golang.org/x/net/html"
		"golang.org/x/net/html/atom"
	*/)

type State struct {
	Theme *ui.Theme
}

func (s *State) Handle(e browser.Event) {
	switch e.(type) {
	}
}

func View(s *State) *browser.Node {
	return ui.VStack(
		s.Theme.Text("hello canvas"),
		&browser.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Canvas,
			Style: browser.Style{
				Border: browser.Border{
					Color: "red",
					Width: browser.Size{Value: 1, Unit: browser.UnitPX},
					Type:  browser.BorderSolid,
				},
			},
			/*

						CanvasDraw: func(c dom.H.Context) {
							return
				//			draw(c)
						},
			*/
		},
	)
}

func draw(c *canvas.Context) {
	cms := canvas.NewFontFamily("Computer Modern Serif")
	if err := cms.LoadFont(cmunrm, 0, canvas.FontRegular); err != nil {
		panic(err)
	}
	//	log.Print(cms.FindLocalFont("", canvas.FontRegular))
	//	cms = canvas.NewFontFamily("SF Pro")
	//	if err := cms.LoadLocalFont("SF Pro", canvas.FontRegular); err != nil {
	//		panic(err)
	//	}

	draw2(c, cms)
}

func drawText(c *canvas.Context, x, y float64, face *canvas.FontFace, rich *canvas.RichText) {
	metrics := face.Metrics()
	width, height := 90.0, 32.0

	text := rich.ToText(width, height, canvas.Justify, canvas.Top, 0.0, 0.0)

	c.SetFillColor(color.RGBA{192, 0, 64, 255})
	c.DrawPath(x, y, text.Bounds().ToPath())
	c.SetFillColor(color.RGBA{51, 51, 51, 51})
	c.DrawPath(x, y, canvas.Rectangle(width, -metrics.LineHeight))
	c.SetFillColor(color.RGBA{0, 0, 0, 51})
	c.DrawPath(x, y+metrics.CapHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.CapHeight-metrics.Descent))
	c.DrawPath(x, y+metrics.XHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.XHeight))

	c.SetFillColor(canvas.Black)
	c.DrawPath(x, y, canvas.Rectangle(width, -height).Stroke(0.2, canvas.RoundCap, canvas.RoundJoin))
	c.DrawText(x, y, text)
}

func draw2(c *canvas.Context, font *canvas.FontFamily) {
	// Draw a comprehensive text box
	pt := 14.0
	face := font.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	rt := canvas.NewRichText(face)
	rt.Add(face, "Lorem dolor ipsum ❤️")
	rt.Add(font.Face(pt, canvas.White, canvas.FontBold, canvas.FontNormal), "confiscator")
	rt.Add(face, " cur\u200babitur ")
	rt.Add(font.Face(pt, canvas.Black, canvas.FontItalic, canvas.FontNormal), "mattis")
	rt.Add(face, " dui ")
	rt.Add(font.Face(pt, canvas.Black, canvas.FontBold|canvas.FontItalic, canvas.FontNormal), "tellus")
	rt.Add(face, " vel. Proin ")
	rt.Add(font.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontUnderline), "sodales")
	rt.Add(face, " eros vel ")
	rt.Add(font.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontSineUnderline), "nibh")
	rt.Add(face, " fringilla pellen\u200btesque eu cillum. ")

	face = font.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	face.Language = "ru"
	face.Script = text.Cyrillic
	rt.Add(face, "дёжжэнтиюнт холст ")

	/*
		face = fontDevanagari.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
		face.Language = "hi"
		face.Script = text.Devanagari
		rt.Add(face, "हालाँकि प्र ")
	*/

	drawText(c, 9, 95, face, rt)

	// Draw the word Stroke being stroked
	face = font.Face(80.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	p, _, _ := face.ToPath("Stroke")
	c.DrawPath(100, 5, p.Stroke(0.75, canvas.RoundCap, canvas.RoundJoin))

	// Draw an elliptic arc being dashed
	ellipse, err := canvas.ParseSVG(fmt.Sprintf("A10 30 30 1 0 30 0z"))
	if err != nil {
		panic(err)
	}
	c.SetFillColor(canvas.Whitesmoke)
	c.DrawPath(110, 40, ellipse)

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(0.75)
	c.SetStrokeCapper(canvas.RoundCap)
	c.SetStrokeJoiner(canvas.RoundJoin)
	c.SetDashes(0.0, 2.0, 4.0, 2.0, 2.0, 4.0, 2.0)
	//ellipse = ellipse.Dash(0.0, 2.0, 4.0, 2.0).Stroke(0.5, canvas.RoundCap, canvas.RoundJoin)
	c.DrawPath(110, 40, ellipse)
	c.SetStrokeColor(canvas.Transparent)
	c.SetDashes(0.0)

	// Draw a LaTeX formula
	latex, err := canvas.ParseLaTeX(`$y = \sin(\frac{x}{180}\pi)$`)
	if err != nil {
		panic(err)
	}
	latex = latex.Transform(canvas.Identity.Rotate(-30))
	c.SetFillColor(canvas.Black)
	c.DrawPath(135, 85, latex)

	// Draw a raster image
	charbuffer := bytes.NewBuffer(char)
	img, err := canvas.NewPNGImage(charbuffer)
	if err != nil {
		panic(err)
	}
	//	c.Rotate(5)
	c.DrawImage(50.0, 0.0, img, 15)
	c.SetView(canvas.Identity)

	// Draw an closed set of points being smoothed
	polyline := &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(30.0, 0.0)
	polyline.Add(30.0, 15.0)
	polyline.Add(0.0, 30.0)
	polyline.Add(0.0, 0.0)
	c.SetFillColor(canvas.Seagreen)
	c.FillColor.R = byte(float64(c.FillColor.R) * 0.25)
	c.FillColor.G = byte(float64(c.FillColor.G) * 0.25)
	c.FillColor.B = byte(float64(c.FillColor.B) * 0.25)
	c.FillColor.A = byte(float64(c.FillColor.A) * 0.25)
	c.SetStrokeColor(canvas.Seagreen)
	c.DrawPath(155, 35, polyline.Smoothen())

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(0.5)
	c.DrawPath(155, 35, polyline.ToPath())
	c.SetStrokeWidth(0.75)
	for _, coord := range polyline.Coords() {
		c.DrawPath(155, 35, canvas.Circle(2.0).Translate(coord.X, coord.Y))
	}

	// Draw a open set of points being smoothed
	polyline = &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(20.0, 10.0)
	polyline.Add(40.0, 30.0)
	polyline.Add(60.0, 40.0)
	polyline.Add(80.0, 20.0)
	c.SetStrokeColor(canvas.Dodgerblue)
	c.DrawPath(10, 15, polyline.Smoothen())
	c.SetStrokeColor(canvas.Black)
	for _, coord := range polyline.Coords() {
		c.DrawPath(10, 15, canvas.Circle(2.0).Translate(coord.X, coord.Y))
	}

}

////go:embed fonts/cmun-serif/cmunrm.ttf
//go:embed DejaVuSerif.ttf
var cmunrm []byte

//go:embed char.png
var char []byte
