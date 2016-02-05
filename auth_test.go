package yahoo

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAuth(t *testing.T) {
	clientId, clientSecret, redirectUrl := "clientId", "clientSecret", "http://redirect.com"
	auth := New(clientId, clientSecret, redirectUrl)

	if auth.Config == nil {
		t.Error("auth is nil")
	}

	if auth.ClientID != clientId {
		t.Errorf("Unexpected clientId: got %s want %s", auth.ClientID, clientId)
	}

	if auth.ClientSecret != clientSecret {
		t.Errorf("Unexpected clientSecret: got %s want %s", auth.ClientSecret, clientSecret)
	}

	if auth.RedirectURL != redirectUrl {
		t.Errorf("Unexpected redirectUrl: got %s want %s", auth.RedirectURL, redirectUrl)
	}

	auth = New(clientId, clientSecret, "")

	if auth.RedirectURL != "oob" {
		t.Errorf("Unexpected redirectUrl: got %s want %s", auth.RedirectURL, "oob")
	}
}

func TestLoginUrl(t *testing.T) {
	auth := New("clientId", "clientSecret", "http://redirect.com")
	state := "abcdefg"
	url := auth.LoginUrl(state)

	if auth.state != state {
		t.Errorf("Unexpected state: got %s want %s", auth.state, state)
	}

	if !strings.HasPrefix(url, AuthURL) {
		t.Errorf("Unexpected login url: got %s want %s", url, AuthURL)
	}

	if !strings.Contains(url, "state="+state) {
		t.Errorf("Login url (%s) does not contain state (%s)", url, "state="+state)
	}
}
func TestConnect(t *testing.T) {
	authHeader := base64.StdEncoding.EncodeToString([]byte("clientId:clientSecret"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Authorization"), authHeader) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"access_token":"Jzxbkqqcvjqik2IMxGFEE1cuaos--",
			"token_type":"bearer",
			"expires_in":3600,
			"refresh_token":"AOiRUlJn_qOmByVGTmUpwcMKW3XDcipToOoHx2wRoyLgJC_RFlA-",
			"xoauth_yahoo_guid":"JT4FACLQZI2OCE"}`))
		}
	}))
	defer server.Close()

	auth := New("clientId", "clientSecret", "http://redirect.com")
	auth.Endpoint.AuthURL = server.URL
	auth.Endpoint.TokenURL = server.URL
	state := "abcdefg"
	auth.LoginUrl(state)

	// check successfull connection
	err := auth.Connect("code", state)
	if err != nil {
		t.Errorf("Could not connect, unexpected error %v", err)
	}

	if auth.token.AccessToken != "Jzxbkqqcvjqik2IMxGFEE1cuaos--" {
		t.Errorf("Access token not set correctly, got %s expected %s", auth.token.AccessToken, "Jzxbkqqcvjqik2IMxGFEE1cuaos--")
	}

	if auth.Client == nil {
		t.Error("auth Client is nil")
	}

	// check bad authorization
	auth.ClientSecret = "BAD"
	// check successfull connection
	err = auth.Connect("code", state)

	if err == nil {
		t.Error("auth Connect unexpected success. Expecting 403 error")
	}
}
func TestConnectMismatchedState(t *testing.T) {
	auth := New("clientId", "clientSecret", "http://redirect.com")
	state := "abcdefg"
	auth.LoginUrl(state)

	err := auth.Connect("code", "hijklmn")

	if err == nil {
		t.Error("Expecting an umatched state error")
	} else if !strings.Contains(err.Error(), "mismatched state") {
		t.Errorf("incorrect error: got %s expected %s", err.Error(), "mismatched state")
	}
}
