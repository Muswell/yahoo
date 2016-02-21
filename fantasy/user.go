package fantasy

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// User type represents a single Yahoo fantasy user
type User struct {
	XMLName xml.Name `xml:"user"`
	// Guid is the unique ID of the user.
	Guid string `xml:"guid"`
	// A list of games for this user.
	Games []Game `xml:"games>game,omitempty"`
}

// UserQueryBuilder formats queries for Yahoo Fantasy Users.
type UserQueryBuilder struct {
	ActiveUser bool
	GameQB     *GameQueryBuilder
}

//Path returns the yahoo api path for the query excluding the host and query string.
func (q *UserQueryBuilder) Path() string {
	path := "users"
	if q.ActiveUser {
		path += ";use_login=1"
	}

	if q.GameQB != nil {
		// todo strip out any redundant user path.
		path += "/" + q.GameQB.Path()

	}
	return path
}

// Url returns the api url that the query builder fields create.
func (q *UserQueryBuilder) Url() string {
	return baseUrl + q.Path() + "?format=xml"
}

// XmlUserParser must be able to parse a byte slice of xml data and return a slice of Users.
type xmlUserParser interface {
	parseXML([]byte) ([]User, error)
}

// DefaultXMLUserParser parses xml with a games node as a direct child of the fantasy_content node.
type defaultXMLUserParser struct {
	result struct {
		// XMLName fantasy_content is the main wrapper tag in a Yahoo api response.
		XMLName xml.Name `xml:"fantasy_content"`
		//Users is where a UserBuilderQuery stores its results.
		Users []User `xml:"users>user"`
	}
}

// ParseXML actually does the transformation from xml to Game slice.
func (p defaultXMLUserParser) parseXML(data []byte) ([]User, error) {
	err := xml.Unmarshal(data, &p.result)
	if err != nil {
		return []User{}, err
	}
	return p.result.Users, nil
}

// XmlParser returns the appropriate xml parser based on the query builder settings.
func (qb *UserQueryBuilder) xmlParser() xmlUserParser {
	/*if qb.GameQB != nil {
		return new(gameXMLUserParser)
	}*/
	return new(defaultXMLUserParser)
}

// Get formats the appropriate api url to query and unmarshals the response into a slice of users.
func (q *UserQueryBuilder) Get(client *http.Client) ([]User, error) {
	url := q.Url()

	resp, err := client.Get(url)
	if err != nil {
		return []User{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []User{}, err
	}

	return q.xmlParser().parseXML(data)
}

// ActiveUser takes an ouath ready http.client and returns a User object representing the authorized user.
func ActiveUser(client *http.Client) (User, error) {
	qb := UserQueryBuilder{ActiveUser: true}
	users, err := qb.Get(client)

	if err != nil {
		return User{}, err
	}

	if len(users) != 1 {
		err = fmt.Errorf("Error the request returned multiple users.")
		return User{}, err
	}

	return users[0], nil
}
