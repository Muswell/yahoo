package fantasy

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Game represents a Yahoo Fantasy Sport.
type Game struct {
	XMLName xml.Name `xml:"game"`
	Key     int64    `xml:"game_key"`
	ID      int64    `xml:"game_id"`
	Name    string   `xml:"name"`
	Code    string   `xml:"code"`
	URL     url.URL  `xml:"url"`
	Season  int64    `xml:"season"`
}

type GameQueryBuilder struct {
	ActiveUser bool
	Available  bool
}

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

type xmlGameParser interface {
	parseXML([]byte) ([]Game, error)
}

type defaultXMLGameParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Games is where a GameBuilderQuery stores its results.
		Games []Game `xml:"games>game"`
	}
}

func (p defaultXMLGameParser) parseXML(data []byte) ([]Game, error) {
	err := xml.Unmarshal(data, &p.result)

	if err != nil {
		return []Game{}, err
	}

	return p.result.Games, nil
}

type userXMLGameParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Games is where a GameBuilderQuery stores its results.
		Games []Game `xml:"users>user>games>game"`
	}
}

func (p userXMLGameParser) parseXML(data []byte) ([]Game, error) {
	err := xml.Unmarshal(data, &p.result)

	if err != nil {
		return []Game{}, err
	}

	return p.result.Games, nil
}

func (qb *GameQueryBuilder) xmlParser() xmlGameParser {
	if qb.ActiveUser {
		return new(userXMLGameParser)
	}
	return new(defaultXMLGameParser)
}

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
