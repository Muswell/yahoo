package fantasy

import (
	"github.com/muswell/gotest"
	"testing"
)

const guid = "JT4FACLQZI2OCE"

func TestNewUser(t *testing.T) {
	client := gotest.NewRegisteredClient()

	// test bad client request
	_, err := NewUser(client.Client)
	if err == nil {
		t.Error("Expected NewUser to return an error")
	}

	// test non-xml response
	client.Register(userUrl, "get", gotest.NewSimpleRoundTrip([]byte("Hello, world"), nil))
	_, err = NewUser(client.Client)
	if err == nil {
		t.Error("Expected NewUser to return an error")
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
	client.Register(userUrl, "get", gotest.NewSimpleRoundTrip(xml, headers))

	user, err := NewUser(client.Client)
	if err != nil {
		t.Errorf("NewUser returned an error %v", err)
	}

	if user.Guid != guid {
		t.Errorf("NewUser returned incorrect user Guid was %s expected %s", user.Guid, guid)
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
	client.Register(userUrl, "get", gotest.NewSimpleRoundTrip(xml, headers))
	_, err = NewUser(client.Client)
	if err == nil {
		t.Error("Expected NewUser to return an error")
	}
}
