package windowmgt

import (
	"math"

	"github.com/gotk3/gotk3/cairo"
	"github.com/spf13/viper"
)

const (
	ContrastThreshold       uint32 = 128
	RFactorForContrastRatio uint32 = 299
	GFactorForContrastRatio uint32 = 587
	BFactorForContrastRatio uint32 = 114
)

type Pill struct {
	msg   string
	color Color
	cr    *cairo.Context
}

func NewPill(msg string, color Color, cr *cairo.Context) (res *Pill) {
	return &Pill{msg: msg, color: color, cr: cr}
}

func (pill *Pill) GetDimensions() (width, height uint32) {
	pill.cr.SelectFontFace(viper.GetString("group_labels.font"), cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	pill.cr.SetFontSize(viper.GetFloat64("group_labels.font_size"))
	txtExtents := pill.cr.TextExtents(pill.msg)
	fontExtents := pill.cr.FontExtents()
	height = uint32(fontExtents.Height) * 2
	width = uint32(txtExtents.Width) + height

	return
}

func (pill *Pill) draw(x, y float64) {
	pill.cr.Save()
	pill.cr.SelectFontFace(viper.GetString("group_labels.font"), cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	pill.cr.SetFontSize(viper.GetFloat64("group_labels.font_size"))
	txtExtents := pill.cr.TextExtents(pill.msg)
	fontExtents := pill.cr.FontExtents()
	height := fontExtents.Height
	width := txtExtents.Width
	r := ((float64(pill.color.R) * 100) / float64(math.MaxUint8)) / 100 //nolint:gomnd //100 is an obvious value to compute a percentage
	g := ((float64(pill.color.G) * 100) / float64(math.MaxUint8)) / 100 //nolint:gomnd //100 is an obvious value to compute a percentage
	b := ((float64(pill.color.B) * 100) / float64(math.MaxUint8)) / 100 //nolint:gomnd //100 is an obvious value to compute a percentage
	pill.cr.SetSourceRGB(r, g, b)
	pill.cr.Translate(x, y)
	pill.cr.MoveTo(height, height)
	//nolint:gomnd //Obvious value as we need the half of Pi to draw a half circle
	pill.cr.Arc(height, height, height, math.Pi/2, math.Pi+math.Pi/2)
	pill.cr.ClosePath()
	pill.cr.Fill()
	pill.cr.Rectangle(height, 0, width, height*2) //nolint:gomnd //Obvious value to multiply by 2 the height of the pill to center it
	pill.cr.ClosePath()
	pill.cr.Fill()
	pill.cr.Arc(height+width, height, height, -math.Pi/2, math.Pi/2) //nolint:gomnd //Obvious value to draw the second half of the circle
	pill.cr.ClosePath()
	pill.cr.Fill()

	oppositeColor := NewColor(math.MaxUint8, math.MaxUint8, math.MaxUint8)

	//nolint:gomnd //Obvious value of 1000 to have a value < 255
	c := (RFactorForContrastRatio*uint32(pill.color.R) + GFactorForContrastRatio*uint32(pill.color.G) +
		BFactorForContrastRatio*uint32(pill.color.B)) / 1000
	if c > ContrastThreshold {
		oppositeColor = NewColor(0, 0, 0)
	}

	r = ((float64(oppositeColor.R) * 100) / float64(math.MaxUint8)) / 100 //nolint:gomnd //100 is an obvious value to compute a percentage
	g = ((float64(oppositeColor.G) * 100) / float64(math.MaxUint8)) / 100 //nolint:gomnd //100 is an obvious value to compute a percentage
	b = ((float64(oppositeColor.B) * 100) / float64(math.MaxUint8)) / 100 //nolint:gomnd //100 is an obvious value to compute a percentage
	pill.cr.SetSourceRGB(r, g, b)
	pill.cr.MoveTo(height, height*2-txtExtents.Height)
	pill.cr.ShowText(pill.msg)
	pill.cr.Restore()
}
