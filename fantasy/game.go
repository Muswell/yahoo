package fantasy

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Game represents a Yahoo Fantasy Sport.
type Game struct {
	// XMLName is the name of the game xml node.
	XMLName xml.Name `xml:"game"`
	// Key, not sure what key is. It seems to duplicate id.
	Key int64 `xml:"game_key"`
	// ID, the unique identifier for this game.
	ID int64 `xml:"game_id"`
	// The name of this game e.g. Baseball
	Name string `xml:"name"`
	// Unique code which acts like an id for the current season of the game e.g. mlb.
	Code string `xml:"code"`
	// The api url associated with this set of games
	URL string `xml:"url"`
	// Season is a 4 digit year in which the season is played.
	Season int64 `xml:"season"`
	//IsRegistrationOver determines if the game is still accepting new signups.
	IsRegistrationOver intAsBool `xml:"is_registration_over"`
}

//GameQueryBuilder contains properties which are used to generate yahoo api game requests.
type GameQueryBuilder struct {
	// ActiveUser set to true sets the query builder to get games only for the logged in user.
	ActiveUser bool
	// Available sets the query builder to only return available games.
	Available bool
}

// Url generates the url needed for a request of the query builder's settings.
func (q *GameQueryBuilder) Url() string {
	url := baseUrl
	if q.ActiveUser {
		url += "/users;use_login=1"
	}

	url += "/games"
	if q.Available {
		url += ";is_available=1"
	}
	url += "?format=xml"

	return url
}

// XmlParser must be able to parse a byte slice of xml data and return a slice of Games.
type xmlGameParser interface {
	parseXML([]byte) ([]Game, error)
}

// DefaultXMLGameParser parses xml with a games node as a direct child of the fantasy_content node.
type defaultXMLGameParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Games is where a GameBuilderQuery stores its results.
		Games []Game `xml:"games>game"`
	}
}

// ParseXML actually does the transformation from xml to Game slice.
func (p defaultXMLGameParser) parseXML(data []byte) ([]Game, error) {
	err := xml.Unmarshal(data, &p.result)
	if err != nil {
		return []Game{}, err
	}
	return p.result.Games, nil
}

//UserXMLGameParser parses xml with a structure fantasy_conent>users>user>games>game.
type userXMLGameParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Games is where a GameBuilderQuery stores its results.
		Games []Game `xml:"users>user>games>game"`
	}
}

// todo create generic function for parseXML.

// ParseXML actually does the transformation from xml to Game slice.
func (p userXMLGameParser) parseXML(data []byte) ([]Game, error) {
	err := xml.Unmarshal(data, &p.result)

	if err != nil {
		return []Game{}, err
	}

	return p.result.Games, nil
}

// XmlParser returns the appropriate xml parser based on the query builder settings.
func (qb *GameQueryBuilder) xmlParser() xmlGameParser {
	if qb.ActiveUser {
		return new(userXMLGameParser)
	}
	return new(defaultXMLGameParser)
}

// Get sends a request to the appropriate url based on the query builder settings.
// It then sends that response to a parser and returns the resulting Game slice.
func (qb *GameQueryBuilder) Get(client *http.Client) ([]Game, error) {
	url := qb.Url()

	resp, err := client.Get(url)
	if err != nil {
		return []Game{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Game{}, err
	}

	return qb.xmlParser().parseXML(data)
}
