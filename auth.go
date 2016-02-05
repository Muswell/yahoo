// Yahoo sports package provides an OAuth 2.0 conntection mechanism for Yahoo Fantasy Sports users
// and provides interfaces to read and write from
package yahoo

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	AuthURL  = "https://api.login.yahoo.com/oauth2/request_auth"
	TokenURL = "https://api.login.yahoo.com/oauth2/get_token"
)

// The auth type configures and and manages a user connection to Yahoo.
type Auth struct {
	*oauth2.Config
	state string
	*http.Client
	token *oauth2.Token
}

// New creates a ready to use Auth instance.
// if the redirectUrl is empty, this is assumed to be an installed application.
func New(clientId, clientSecret, redirectUrl string) *Auth {
	if redirectUrl == "" {
		redirectUrl = "oob"
	}

	auth := &Auth{
		Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  AuthURL,
				TokenURL: TokenURL,
			},
			RedirectURL: redirectUrl,
		},
	}

	return auth
}

// LoginUrl generates a URL which must be visited to begin the authorization process

// State is a token used to protect from CSRF attacks.
// It is stored internally and checked when trying to check the authoriztion_code

// Opts may include AccessTypeOnline or AccessTypeOffline, as well as ApprovalForce.
func (auth *Auth) LoginUrl(state string, opts ...oauth2.AuthCodeOption) string {
	loginUrl := auth.Config.AuthCodeURL(state, opts...)
	auth.state = state
	return loginUrl
}

// Connect takes the authorization code and state returned By Yahoo and creates the http Client and token
func (auth *Auth) Connect(code, state string) error {
	if state != auth.state {
		return fmt.Errorf("mismatched state %s expecting %s", state, auth.state)
	}
	if token, err := auth.Config.Exchange(oauth2.NoContext, code); err == nil {
		//guid := token.Extra("xoauth_yahoo_guid")
		auth.token = token
		auth.Client = auth.Config.Client(oauth2.NoContext, token)
		return nil
	} else {
		return err
	}
}
