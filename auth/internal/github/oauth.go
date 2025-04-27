package github

import (
	"github.com/DenisOzindzheDev/chatix/auth/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oauthConfig *oauth2.Config

func InitOAuth(cfg *config.Config) {
	oauthConfig = &oauth2.Config{
		ClientID:     cfg.Github.ClientID,
		ClientSecret: cfg.Github.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  cfg.Github.RedirectURL,
		Scopes:       []string{"read:user", "user:email"},
	}
}

func GetOAuthConfig() *oauth2.Config {
	return oauthConfig
}
