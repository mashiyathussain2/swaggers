package app

import (
	"go-app/server/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth interface {
	AuthCodeURL(opts ...oauth2.AuthCodeOption) string
	Exchange(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type GoogleOAuthImpl struct {
	Config *config.GoogleOAuth
	Client *oauth2.Config
}

type GoogleOAuthOpts struct {
	Config *config.GoogleOAuth
}

func NewGoogleOAuth(opts *GoogleOAuthOpts) GoogleOAuth {
	client := oauth2.Config{
		ClientID:     opts.Config.ClientID,
		ClientSecret: opts.Config.ClientSecret,
		RedirectURL:  opts.Config.RedirectURL,
		Scopes:       opts.Config.Scopes,
		Endpoint:     google.Endpoint,
	}
	gl := GoogleOAuthImpl{
		Config: opts.Config,
		Client: &client,
	}
	return &gl
}

func (gi *GoogleOAuthImpl) AuthCodeURL(opts ...oauth2.AuthCodeOption) string {
	return gi.Client.AuthCodeURL(gi.Config.State, opts...)
}

func (gi *GoogleOAuthImpl) Exchange(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return gi.Client.Exchange(oauth2.NoContext, code, opts...)
}
