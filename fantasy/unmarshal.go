package fantasy

import (
	/*"fmt"
	  "log"
	  "strings"*/
	"encoding/xml"
	"fmt"
)

// IntAsBool reads integer xml node values and converts them to true or false,
// 0 = false, 1 = true, everything else errors.
type intAsBool bool

// UnmarshalXML takes an xml element, reads its content as an int64 and converts that to a bool.
func (b *intAsBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var i int64
	d.DecodeElement(&i, &start)

	if i == 0 {
		*b = intAsBool(false)
		return nil
	}

	if i == 1 {
		*b = intAsBool(true)
		return nil
	}
	return fmt.Errorf("Bad intAsBool value %d (0 = false 1, = true)", i)
}
