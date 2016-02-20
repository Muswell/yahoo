// Fantasy wraps the Yahoo fantasy API into Go types and provides functions to generate proper request URLs.
package fantasy

import "encoding/xml"

const (
	// BaseUrl is the root yahoo fantasy api url.
	baseUrl = "https://fantasysports.yahooapis.com/fantasy/v2/"
)

// Parser is the main struct for unmarshaling xml from http.client.Get requests.
type Parser struct {
	// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
	XMLName xml.Name `xml:"fantasy_content"`
	// Users is where a user(s) request stores its results.
	Users []User `xml:"users>user"`
	// Games is where a GameBuilderQuery stores its results.
	Games []Game `xml:"games>game"`
}
