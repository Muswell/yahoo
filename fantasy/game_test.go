package fantasy

import (
	"github.com/muswell/gotest"
	"io/ioutil"
	"os"
	"testing"
)

func TestGameQueryBuilderURL(t *testing.T) {
	var tests = []struct {
		input GameQueryBuilder
		want  string
	}{
		{
			GameQueryBuilder{Available: true},
			baseUrl + "/games;is_available=1?format=xml",
		},
		{
			GameQueryBuilder{ActiveUser: true},
			baseUrl + "/users;use_login=1/games?format=xml",
		},
		{
			GameQueryBuilder{Available: true, ActiveUser: true},
			baseUrl + "/users;use_login=1/games;is_available=1?format=xml",
		},
	}

	for _, test := range tests {
		if got := test.input.Url(); got != test.want {
			t.Errorf("%v.Url = %q, want %q", test.input, got, test.want)
		}
	}
}
func TestGetGamesErrors(t *testing.T) {
	qb := GameQueryBuilder{Available: true}
	client := gotest.NewRegisteredClient()
	url := qb.Url()

	// test bad client request
	_, err := qb.Get(client.Client)
	if err == nil {
		t.Error("Expected GameQueryBuilder.Get to return an error")
	}

	// test non-xml response
	client.Register(url, "get", gotest.NewSimpleRoundTrip([]byte("Hello, world"), nil))
	_, err = ActiveUser(client.Client)
	if err == nil {
		t.Error("Expected GameQueryBuilder.Get to return an error")
	}
}
func TestGetAllAvailableGames(t *testing.T) {
	qb := GameQueryBuilder{Available: true}
	client := gotest.NewRegisteredClient()
	url := qb.Url()

	// test valid xml
	file, err := os.Open("test/all-available-games.xml")
	if err != nil {
		t.Error("Could not open test/all-available-games.xml")
	}
	xml, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Could not read test/all-available-games.xml")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "text/xml"
	client.Register(url, "get", gotest.NewSimpleRoundTrip(xml, headers))

	games, err := qb.Get(client.Client)

	if err != nil {
		t.Errorf("Unexpected GameQueryBuilder.Get error %v", err)
	}

	if len(games) != 5 {
		t.Errorf("GameQueryBuilder.Get returned %d games expected %d", len(games), 5)
	}

	//todo test values of Game
}

func TestGetUserGames(t *testing.T) {
	qb := GameQueryBuilder{ActiveUser: true}
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

	games, err := qb.Get(client.Client)

	if err != nil {
		t.Errorf("Unexpected GameQueryBuilder.Get error %v", err)
	}

	if len(games) != 1 {
		t.Errorf("GameQueryBuilder.Get returned %d games expected %d", len(games), 1)
	}
}
