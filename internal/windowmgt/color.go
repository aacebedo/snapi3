package windowmgt

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/rotisserie/eris"

	"github.com/aacebedo/snapi3/internal"
)

const NumberOfComponentsInHTMLCode uint32 = 4

type Color struct {
	R uint8
	G uint8
	B uint8
}

func NewColor(r, g, b uint8) (res *Color) {
	res = &Color{R: r, G: g, B: b}

	return
}

func IsStringAnHTMLColorCode(htmlColorCodeStr string) (res bool) {
	res = false

	htmlColorCodeRegexp, err := regexp.Compile("^#[[:xdigit:]]{2}[[:xdigit:]]{2}[[:xdigit:]]{2}$")
	if err != nil {
		internal.VerboseLogger.Errorf("Impossible to compile regex '%s'", htmlColorCodeRegexp)

		return
	}

	return htmlColorCodeRegexp.MatchString(htmlColorCodeStr)
}

func NewColorFromHTMLCode(htmlCode string) (res *Color, err error) {
	//nolint:gocritic //MustCompile make the app immediately crash, I want to handle the error gracefully
	colorRegexp, err := regexp.Compile("^#(?P<R>[[:xdigit:]]{2})(?P<G>[[:xdigit:]]{2})(?P<B>[[:xdigit:]]{2})$")
	if err != nil {
		err = eris.Wrapf(err, "Impossible to compile regex '%s'", colorRegexp)

		return
	}

	colorComponents := colorRegexp.FindStringSubmatch(htmlCode)
	if len(colorComponents) != int(NumberOfComponentsInHTMLCode) {
		err = eris.Wrapf(internal.InvalidArgumentError, "HTML color code '%s' is invalid", htmlCode)

		return
	}

	rValue, rErr := hex.DecodeString(colorComponents[1])
	gValue, gErr := hex.DecodeString(colorComponents[2])
	bValue, bErr := hex.DecodeString(colorComponents[3])

	if rErr != nil || gErr != nil || bErr != nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Unable to decode a color component from the HTML color code '%s'", htmlCode)

		return
	}

	res = NewColor(rValue[0], gValue[0], bValue[0])

	return
}

func (c *Color) ToHTMLCode() (res string) {
	res = fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)

	return
}

func (c *Color) String() (res string) {
	return c.ToHTMLCode()
}
