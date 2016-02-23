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
			baseUrl + "games;is_available=1?format=xml",
		},
		{
			GameQueryBuilder{UserQB: &UserQueryBuilder{ActiveUser: true}},
			baseUrl + "users;use_login=1/games?format=xml",
		},
		{
			GameQueryBuilder{Available: true, UserQB: &UserQueryBuilder{ActiveUser: true}},
			baseUrl + "users;use_login=1/games;is_available=1?format=xml",
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
	game := games[0]

	if game.Code != "nfl" {
		t.Errorf("Incorrect Game Code got %s, expected %s", game.Code, "nfl")
	}
	if game.ID != 348 {
		t.Errorf("Incorrect Game ID got %d, expected %d", game.ID, 348)
	}

	if !game.IsRegistrationOver {
		t.Errorf("Game registration was expected to be over %s.", game.IsRegistrationOver)
	}
}

func TestGetUserGames(t *testing.T) {
	qb := GameQueryBuilder{UserQB: &UserQueryBuilder{ActiveUser: true}}
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
