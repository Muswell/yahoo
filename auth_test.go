package yahoo

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	clientId, clientSecret, redirectUrl := "clientId", "clientSecret", "http://redirect.com"
	auth := NewConfig(clientId, clientSecret, redirectUrl)

	if auth.ClientID != clientId {
		t.Errorf("Unexpected clientId: got %s want %s", auth.ClientID, clientId)
	}

	if auth.ClientSecret != clientSecret {
		t.Errorf("Unexpected clientSecret: got %s want %s", auth.ClientSecret, clientSecret)
	}

	if auth.RedirectURL != redirectUrl {
		t.Errorf("Unexpected redirectUrl: got %s want %s", auth.RedirectURL, redirectUrl)
	}

	auth = NewConfig(clientId, clientSecret, "")

	if auth.RedirectURL != "oob" {
		t.Errorf("Unexpected redirectUrl: got %s want %s", auth.RedirectURL, "oob")
	}
}
