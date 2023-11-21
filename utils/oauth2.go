package utils

import (
	"Tally/config"
	"golang.org/x/oauth2"
)

func NewOauth2() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Config.Oauth2.ClientID,
		ClientSecret: config.Config.Oauth2.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Config.Oauth2.AuthURL,
			TokenURL: config.Config.Oauth2.TokenURL,
		},
		RedirectURL: config.Config.Oauth2.RedirectURL,
		Scopes:      []string{config.Config.Oauth2.Scopes},
	}
}
