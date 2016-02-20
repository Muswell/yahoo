package fantasy

import (
/*"fmt"
"log"
"strings"*/
)

// FormatAsBool is a helper type that handles the inconsistent bool formatting in the json returned from Yahoo!
type formatAsBool bool

/*func (b *formatAsBool) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	if s == "1" || s == "true" {
		*b = formatAsBool(true)
	} else if s == "0" || s == "false" {
		*b = formatAsBool(false)
	} else {
		log.Println("woops, an error")
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", s)
	}
	return nil
}*/
