package fantasy

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	userUrl = "https://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/?format=xml"
)

// User type represents a single Yahoo fantasy user
type User struct {
	XMLName xml.Name `xml:"user"`
	// Guid is the unique ID of the user
	Guid string `xml:"guid"`
}

// NewUser takes an ouath ready http.client and returns a User object
func NewUser(client *http.Client) (User, error) {
	user := User{}
	resp, err := client.Get(userUrl)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	var p Parser

	err = xml.Unmarshal(data, &p)

	if err != nil {
		return user, err
	}

	if len(p.Users) != 1 {
		err = fmt.Errorf("Error the request returned multiple users.")
		return user, err
	}

	return p.Users[0], nil
}
