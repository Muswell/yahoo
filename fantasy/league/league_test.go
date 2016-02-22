package league

import (
	"github.com/muswell/yahoo/fantasy"
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
	}

	for _, test := range tests {
		if got := test.input.Url(); got != test.want {
			t.Errorf("Url = %q, want %q", got, test.want)
		}
	}
}
