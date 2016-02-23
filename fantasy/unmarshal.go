package fantasy

import (
	/*"fmt"
	  "log"
	  "strings"*/
	"encoding/xml"
	"fmt"
	"time"
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

type calendarDate time.Time

func (c *calendarDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02" // yyyy-mm-dd date format
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = calendarDate(parse)
	return nil
}
