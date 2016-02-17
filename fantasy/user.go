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
	// Guid is the unique ID of the user
	Guid string `xml:"guid"`
}

// UserQueryBuilder formats queries for Yahoo Fantasy Users.
type UserQueryBuilder struct {
	ActiveUser bool
}

// Url returns the api url that the query builder fields create.
func (q UserQueryBuilder) Url() string {
	url := baseUrl + "users"
	if q.ActiveUser {
		url += ";use_login=1"
	}
	url += "?format=xml"

	return url
}

// Get formats the appropriate api url to query and unmarshals the response into a slice of users.
func (q UserQueryBuilder) Get(client *http.Client) ([]User, error) {
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

	var p Parser

	err = xml.Unmarshal(data, &p)

	if err != nil {
		return []User{}, err
	}

	return p.Users, nil
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
