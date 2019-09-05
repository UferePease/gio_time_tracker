package main

import (
	"gioui.org/ui"
	"gioui.org/ui/f32"
	"gioui.org/ui/input"
	"gioui.org/ui/layout"
	"gioui.org/ui/measure"
	"gioui.org/ui/paint"
	"gioui.org/ui/text"
	"golang.org/x/image/font/sfnt"
	"image"
	"image/color"
)

type TimerUI struct {
	faces        measure.Faces

	descEditor, projectEditor *text.Editor
	descLabel, projectLabel   *text.Label

	//startBtn	*text.Label
	projDrpDwn	*text.Editor
}

type ActionButton struct {
	face    text.Face
	Open    bool
	icons   []*icon
	sendIco *icon
}

type icon struct {
	src  []byte
	size ui.Value

	// Cached values.
	img     image.Image
	imgSize int
}

var theme struct {
	text     ui.MacroOp
	tertText ui.MacroOp
	brand    ui.MacroOp
	white    ui.MacroOp
}

var fonts struct {
	regular *sfnt.Font
	bold    *sfnt.Font
	italic  *sfnt.Font
	mono    *sfnt.Font
}

func colorMaterial(ops *ui.Ops, color color.RGBA) ui.MacroOp {
	var mat ui.MacroOp
	mat.Record(ops)
	paint.ColorOp{Color: color}.Add(ops)
	mat.Stop()
	return mat
}

type fill struct {
	material ui.MacroOp
}

func rgb(c uint32) color.RGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func newUI() *TimerUI {
	//u := new(TimerUI)
	u := &TimerUI{
	}

	u.descEditor = &text.Editor{
		Face: u.face(fonts.italic, 14),
		//Alignment: text.End,
		SingleLine:   false,
		Hint:         "Enter work description here",
		HintMaterial: theme.tertText,
		Material:     theme.text,
	}
	//u.descEditor.SetText("Enter work description here")
	u.projectEditor = &text.Editor{
		Face:     u.face(fonts.regular, 16),
		Material: theme.text,
		//Alignment: text.End,
		//SingleLine: true,
	}
	u.projectEditor.SetText("Select Project")

	u.descLabel = &text.Label{
		Face:     u.face(fonts.regular, 16),
		Material: theme.text,
		Text: "Description",
	}
	u.projectLabel = &text.Label{
		Face:     u.face(fonts.regular, 16),
		Material: theme.text,
		Text: "Project",
	}
	return u
}

func (u *TimerUI) face(f *sfnt.Font, size float32) text.Face {
	var faces measure.Faces
	return faces.For(f, ui.Sp(size))
}

func (f fill) Layout(ops *ui.Ops, cs layout.Constraints) layout.Dimensions {
	d := image.Point{X: cs.Width.Max, Y: cs.Height.Max}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	f.material.Add(ops)
	paint.PaintOp{Rect: dr}.Add(ops)
	return layout.Dimensions{Size: d, Baseline: d.Y}
}

func (u *TimerUI) layoutPage(c ui.Config, q input.Queue, ops *ui.Ops, cs layout.Constraints) layout.Dimensions {
	st := (&layout.Stack{}).Init(ops, cs)
	cs = st.Rigid()
	al := layout.Align{Alignment: layout.SE}
	in := layout.UniformInset(ui.Dp(16))

	dims := in.End(u.descLabel.Layout(ops, in.Begin(c, ops, cs)))
	//dims := in.End(u.descLabel.Layout(ops, cs))
	dims = al.End(in.End(dims))
	c2 := st.End(dims)

	cs = st.Expand()
	{
		f := (&layout.Flex{Axis: layout.Vertical}).Init(ops, cs)

		cs = f.Rigid()
		//dims := u.projectEditor.Layout(c, q, ops, in.Begin(c, ops, cs))
		{
			cs.Width.Min = cs.Width.Max
			in := layout.UniformInset(ui.Dp(16))
			sz := c.Px(ui.Dp(200))
			cs = layout.RigidConstraints(cs.Constrain(image.Point{X: sz, Y: sz}))
			dims := u.projectEditor.Layout(c, q, ops, in.Begin(c, ops, cs))
			dims = in.End(dims)
		}
		c1 := f.End(dims)

		cs = f.Rigid()
		{
			cs.Width.Min = cs.Width.Max
			in := layout.Inset{Bottom: ui.Dp(16), Left: ui.Dp(16), Right: ui.Dp(16)}
			dims = u.descEditor.Layout(c, q, ops, in.Begin(c, ops, cs))
			dims = in.End(dims)
		}
		c2 := f.End(dims)

		cs = f.Rigid()
		{
			cs.Width.Min = cs.Width.Max
			s := layout.Stack{Alignment: layout.Center}
			s.Init(ops, cs)
			cs = s.Rigid()
			in := layout.Inset{Top: ui.Dp(16), Right: ui.Dp(8), Bottom: ui.Dp(8), Left: ui.Dp(8)}
			grey := colorMaterial(ops, rgb(0x888888))
			lbl := text.Label{Material: grey, Face: u.face(fonts.regular, 11), Text: "TIMER"}
			dims = in.End(lbl.Layout(ops, in.Begin(c, ops, cs)))
			c2 := s.End(dims)
			c1 := s.End(fill{colorMaterial(ops, rgb(0xf2f2f2))}.Layout(ops, s.Expand()))
			dims = s.Layout(c1, c2)
		}
		c3 := f.End(dims)

		dims = f.Layout(c1, c2, c3)
	}
	c1 := st.End(dims)
	return st.Layout(c1, c2)
	//return st.Layout(c1)
}

func (u *TimerUI) Layout(c ui.Config, q input.Queue, ops *ui.Ops, cs layout.Constraints) layout.Dimensions {
	u.faces.Reset(c)

	return u.layoutPage(c, q, ops, cs)
}