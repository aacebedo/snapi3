package internal

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/rotisserie/eris"
)

func HexStringToInt(hexStrToConvert string) (res uint32, err error) {
	regexString := ""

	for i := 0; i < 7; i++ {
		regexString += "([[:xdigit:]]{1})?"
	}

	hexStrRegexp, err := regexp.Compile(fmt.Sprintf("^0x(?P<value>%s[[:xdigit:]]{1})$", regexString))
	if err != nil {
		err = eris.Wrap(InternalError, "Unable to compile the hexstr regex, this is not a normal behavior")

		return
	}

	hexStr := hexStrRegexp.FindStringSubmatch(hexStrToConvert)

	parsedValue, err := strconv.ParseInt(hexStr[1], 16, 0)
	if err != nil {
		err = eris.Wrapf(InvalidArgumentError, "The hex string '%s' does not match the conversion regex", hexStr)

		return
	}

	res = uint32(parsedValue)

	return
}
