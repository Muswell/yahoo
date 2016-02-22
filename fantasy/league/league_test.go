package league

import (
	"github.com/muswell/gotest"
	"github.com/muswell/yahoo/fantasy"
	"io/ioutil"
	"os"
	"testing"
	"time"
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
func TestQueryLeaugeErrors(t *testing.T) {
	qb := QueryBuilder{Keys: []string{"abc"}}
	client := gotest.NewRegisteredClient()
	url := qb.Url()

	// test bad client request
	_, err := qb.Get(client.Client)
	if err == nil {
		t.Error("Expected QueryBuilder.Get to return an error")
	}

	// test non-xml response
	client.Register(url, "get", gotest.NewSimpleRoundTrip([]byte("Hello, world"), nil))
	_, err = qb.Get(client.Client)
	if err == nil {
		t.Error("Expected QueryBuilder.Get to return an error")
	}
}

type queryTest struct {
	qb     *QueryBuilder
	client *gotest.RegisteredClient
	want   int
	next   func([]League, *testing.T)
}

func TestLeagueQueries(t *testing.T) {
	var tests = []queryTest{
		getSingleMetaTestSet(t),
		getUserSingleMetaTestSet(t),
	}

	for _, test := range tests {

		leagues, err := test.qb.Get(test.client.Client)
		if err != nil {
			t.Errorf("Unexpected QueryBuilder.Get error: %s", err)
		}
		if len(leagues) != test.want {
			t.Errorf("Unexpected League len got %d, expected %d", len(leagues), test.want)
		}

		if test.next != nil {
			test.next(leagues, t)
		}
	}
}

func getSingleMetaTestSet(t *testing.T) queryTest {
	q := QueryBuilder{Keys: []string{"357.l.86753"}}
	url := q.Url()

	return queryTest{
		qb:     &q,
		client: getXMLClient(url, "single-league-meta.xml", t),
		want:   1,
		next: func(l []League, t *testing.T) {
			league := l[0]

			if league.GameCode != "mlb" {
				t.Errorf("League unmarshaled incorrectley. GameCode: %s, expected %s", league.GameCode, "mlb")
			}
			start := time.Time(league.StartDate).Format("2006-02-01")
			if start != "2016-03-04" {
				t.Errorf("League unmarshaled incorrectley. StartDate: %s, expected %s", start, "2016-03-04")
			}
		},
	}
}

func getUserSingleMetaTestSet(t *testing.T) queryTest {
	q := QueryBuilder{UserQB: &fantasy.UserQueryBuilder{ActiveUser: true}}
	url := q.Url()

	return queryTest{
		qb:     &q,
		client: getXMLClient(url, "user-leagues-meta.xml", t),
		want:   1,
	}
}

func getXMLClient(url, filename string, t *testing.T) *gotest.RegisteredClient {
	client := gotest.NewRegisteredClient()

	// test valid xml
	file, err := os.Open("test/" + filename)
	if err != nil {
		t.Error("Could not open test/"+filename, err)
	}
	xml, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Could not read test/"+filename, err)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "text/xml"
	client.Register(url, "get", gotest.NewSimpleRoundTrip(xml, headers))

	return client
}
