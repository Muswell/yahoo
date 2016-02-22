package league

import (
	"github.com/muswell/gotest"
	"github.com/muswell/yahoo/fantasy"
	"io/ioutil"
	"os"
	"testing"
)

func TestGameQueryBuilderURL(t *testing.T) {
	var tests = []struct {
		input QueryBuilder
		want  string
	}{
		{
			QueryBuilder{UserQB: &fantasy.UserQueryBuilder{ActiveUser: true}},
			fantasy.BaseUrl + "users;use_login=1/games/leagues?format=xml",
		},
		{
			QueryBuilder{Keys: []string{"357.l.37903", "357.l.37825"}},
			fantasy.BaseUrl + "leagues;league_keys=357.l.37903,357.l.37825?format=xml",
		},
		{
			QueryBuilder{
				UserQB: &fantasy.UserQueryBuilder{
					ActiveUser: true,
					GameQB: &fantasy.GameQueryBuilder{
						Available:  true,
						ActiveUser: true,
					},
				},
			},
			fantasy.BaseUrl + "users;use_login=1/games;is_available=1/leagues?format=xml",
		},
		{
			QueryBuilder{
				UserQB: &fantasy.UserQueryBuilder{
					ActiveUser: true,
					GameQB: &fantasy.GameQueryBuilder{
						Available:  true,
						ActiveUser: true,
					},
				},
				Keys: []string{"357.l.37903"},
			},
			fantasy.BaseUrl + "users;use_login=1/games;is_available=1/leagues;league_keys=357.l.37903?format=xml",
		},
	}

	for _, test := range tests {
		if got := test.input.Url(); got != test.want {
			t.Errorf("Url = %q, want %q", got, test.want)
		}
	}
}

func TestQuerySingleLeagueMetaData(t *testing.T) {
	q := QueryBuilder{Keys: []string{"357.l.86753"}}
	client := gotest.NewRegisteredClient()
	url := q.Url()

	// test valid xml
	file, err := os.Open("test/single-league-meta.xml")
	if err != nil {
		t.Error("Could not open test/single-league-meta.xml", err)
	}
	xml, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Could not read test/single-league-meta.xml", err)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "text/xml"
	client.Register(url, "get", gotest.NewSimpleRoundTrip(xml, headers))
	leagues, err := q.Get(client.Client)
	if err != nil {
		t.Errorf("Unexpected QueryBuilder.Get error: %s", err)
	}
	if len(leagues) != 1 {
		t.Errorf("Unexpected League len got %d, expected %d", len(leagues), 1)

	}

	league := leagues[0]

	if league.GameCode != "mlb" {
		t.Errorf("League unmarshaled incorrectley. GameCode: %s, expected %s", league.GameCode, "mlb")
	}
}
