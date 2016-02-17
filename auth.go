// Yahoo package provides an OAuth 2.0 config for Yahoo Fantasy Sports users
package yahoo

import (
	"golang.org/x/oauth2"
)

const (
	AuthURL  = "https://api.login.yahoo.com/oauth2/request_auth"
	TokenURL = "https://api.login.yahoo.com/oauth2/get_token"
)

// New creates a ready to use Auth instance.
// if the redirectUrl is empty, this is assumed to be an installed application.
func NewConfig(clientId, clientSecret, redirectUrl string) *oauth2.Config {
	if redirectUrl == "" {
		redirectUrl = "oob"
	}

	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
		RedirectURL: redirectUrl,
	}

	return config
}
