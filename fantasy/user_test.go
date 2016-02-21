package fantasy

import (
	"github.com/muswell/gotest"
	"io/ioutil"
	"os"
	"testing"
)

const guid = "JT4FACLQZI2OCE"

func TestActiveUser(t *testing.T) {
	client := gotest.NewRegisteredClient()
	qb := UserQueryBuilder{ActiveUser: true}
	url := qb.Url()

	// test bad client request
	_, err := ActiveUser(client.Client)
	if err == nil {
		t.Error("Expected ActiveUser to return an error")
	}

	// test non-xml response
	client.Register(url, "get", gotest.NewSimpleRoundTrip([]byte("Hello, world"), nil))
	_, err = ActiveUser(client.Client)
	if err == nil {
		t.Error("Expected ActiveUser to return an error")
	}

	// test well formatted xml response
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<fantasy_content xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/" time="23.27299118042ms" copyright="Data provided by Yahoo! and STATS, LLC" refresh_rate="31" xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
		 <users count="1">
		  <user>
		   <guid>` + guid + `</guid>
		  </user>
		 </users>
		</fantasy_content>`)
	headers := make(map[string]string)
	headers["Content-Type"] = "text/xml"
	client.Register(url, "get", gotest.NewSimpleRoundTrip(xml, headers))

	user, err := ActiveUser(client.Client)
	if err != nil {
		t.Errorf("ActiveUser returned an error %v", err)
	}

	if user.Guid != guid {
		t.Errorf("ActiveUser returned incorrect user Guid was %s expected %s", user.Guid, guid)
	}

	if user.Games != nil {
		t.Error("Expected Games to be nil.")
	}

	// test incorrect xml
	xml = []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<fantasy_content xml:lang="en-US" yahoo:uri="http://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/" time="23.27299118042ms" copyright="Data provided by Yahoo! and STATS, LLC" refresh_rate="31" xmlns:yahoo="http://www.yahooapis.com/v1/base.rng" xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
		 <users count="2">
		  <user>
		   <guid>` + guid + `</guid>
		  </user>
		  <user>
		   <guid>ABCDEFGHI</guid>
		  </user>
		 </users>
		</fantasy_content>`)
	client.Register(url, "get", gotest.NewSimpleRoundTrip(xml, headers))
	_, err = ActiveUser(client.Client)
	if err == nil {
		t.Error("Expected ActiveUser to return an error")
	}
}

func TestGetUserWithGames(t *testing.T) {
	qb := UserQueryBuilder{
		ActiveUser: true,
		GameQB:     &GameQueryBuilder{Available: true},
	}

	client := gotest.NewRegisteredClient()
	url := qb.Url()

	// test valid xml
	file, err := os.Open("test/user-games.xml")
	if err != nil {
		t.Error("Could not open test/user-games.xml")
	}
	xml, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Could not read test/user-games.xml")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "text/xml"
	client.Register(url, "get", gotest.NewSimpleRoundTrip(xml, headers))

	users, err := qb.Get(client.Client)

	if err != nil {
		t.Errorf("UserQueryBuilder.Get returned an error %v", err)
	}

	user := users[0]

	if user.Guid != guid {
		t.Errorf("UserQueryBuilder.Get returned incorrect user Guid was %s expected %s", user.Guid, guid)
	}

	if len(user.Games) != 1 {
		t.Errorf("Incorrect User.Games length got %d, expected %d", len(user.Games), 1)
	}
}
