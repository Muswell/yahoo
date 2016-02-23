// League Contains Query Builders for Yahoo fantasy leagues
package fantasy

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

// League contains the metadata for a Yahoo fantasy league as well as pointers to related information.

// When users join a Fantasy Football, Baseball, Basketball, or Hockey draft and trade game,
// they are organized into leagues with a limited number of friends or other Yahoo! users,
// with each user managing a Team. With the League API, you can obtain the league related information,
// like the league name, the number of teams, the draft status, et cetera. Leagues only exist in the context
// of a particular Game.

// A particular user can only retrieve data for private leagues of which they are a member, or for public leagues.
type League struct {
	XMLName xml.Name `xml:"league"`
	// ID is the League ID.
	ID int64 `xml:"league_id"`
	// Key is an id which matches a game with a league.
	Key string `xml:"league_key"`
	// The name for this league.
	Name string `xml:"name"`
	// The url associated with this league.
	URL string `xml:"url"`
	// The id needed to initiate chat within the league
	ChatID string `xml:"league_chat_id"`
	// What status the draft for this league currently has.
	DraftStatus string `xml:"dratf_status"`
	// The number of teams signed up for this league.
	NumTeams int64 `xml:"num_teams"`
	// The stlye of scoring the league uses.
	ScoringType string `xml:"scoring_type"`
	// Public or private.
	LeagueType string `xml:"league_type"`
	// The beginning date of the season
	StartDate calendarDate `xml:"start_date"`
	// The end date of the season
	EndDate calendarDate `xml:"end_date"`
	//The code of the associated game type.
	GameCode string `xml:"game_code"`
	// 4 digit year
	Season int64 `xml:"season"`
	// A pointer to the League Settings.
	//*Settings `xml:"settings"`
	// A pointer to the League Standings.
	//*Standings `xml:"standings"`
	// A pointer to the League Scoreboard.
	//*Scoreboard `xml:"scoreboard"`
	// A pointer to the League Teams.
	//*Teams `xml:"teams"`
	// A pointer to the League Settings.
	//*Settings `xml:"settings"`
	// A pointer to the League's eligible Players.
	//*Players `xml:"players"`
	// A pointer to the League Draft.
	//*Draft `xml:"draftresults"`
	// A pointer to the League Transactions.
	//*Transactions `xml:"transactions"`
}

//LeagueQueryBuilder contains properties which are used to generate yahoo api league requests.
type LeagueQueryBuilder struct {
	// Add a UserQueryBuilder to filter results by user info.
	UserQB *UserQueryBuilder
	// Add League Keys to return specific leagues.
	Keys []string
	// todo include settings, standings...
}

//Path returns the yahoo api path for the query excluding the host and query string.
func (q *LeagueQueryBuilder) Path() string {
	var path string

	if q.UserQB != nil {
		path += "/" + q.UserQB.Path()
		// a league only exists within a game.
		if !strings.Contains(path, "/games") {
			path += "/games"
		}
	}

	path += "/leagues"

	if q.Keys != nil {
		path += ";league_keys=" + strings.Join(q.Keys, ",")
	}

	return strings.TrimLeft(path, "/")
}

// Url generates the url needed for a request of the query builder's settings.
func (q *LeagueQueryBuilder) Url() string {
	return baseUrl + q.Path() + "?format=xml"
}

// XmlLeagueParser must be able to parse a byte slice of xml data and return a slice of Leagues.
type xmlLeagueParser interface {
	parseXML([]byte) ([]League, error)
}

// DefaultXMLLeagueParser parses xml with a leagues node as a direct child of the fantasy_content node.
type defaultXMLLeagueParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Games is where a GameBuilderQuery stores its results.
		Leagues []League `xml:"leagues>league"`
	}
}

// ParseXML actually does the transformation from xml to League slice.
func (p defaultXMLLeagueParser) parseXML(data []byte) ([]League, error) {
	err := xml.Unmarshal(data, &p.result)
	if err != nil {
		return []League{}, err
	}
	return p.result.Leagues, nil
}

//UserXMLLeagueParser parses xml with a structure fantasy_conent>users>user>games>game.
type userXMLLeagueParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Games is where a GameBuilderQuery stores its results.
		Leagues []League `xml:"users>user>games>game>leagues>league"`
	}
}

// todo create generic function for parseXML.
// ParseXML actually does the transformation from xml to League slice.
func (p userXMLLeagueParser) parseXML(data []byte) ([]League, error) {
	err := xml.Unmarshal(data, &p.result)

	if err != nil {
		return []League{}, err
	}

	return p.result.Leagues, nil
}

// XmlParser returns the appropriate xml parser based on the query builder settings.
func (q *LeagueQueryBuilder) xmlParser() xmlLeagueParser {
	if q.UserQB != nil {
		return new(userXMLLeagueParser)
	}
	return new(defaultXMLLeagueParser)
}

// Get sends a request to the appropriate url based on the query builder settings.
// It then sends that response to a parser and returns the resulting League slice.
func (q *LeagueQueryBuilder) Get(client *http.Client) ([]League, error) {
	url := q.Url()

	resp, err := client.Get(url)
	if err != nil {
		return []League{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []League{}, err
	}

	return q.xmlParser().parseXML(data)
}
