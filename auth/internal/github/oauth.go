package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DenisOzindzheDev/chatix/auth/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

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

func GetUserProfile(ctx context.Context, accessToken string) (*GitUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status error: %v", resp.StatusCode)
	}

	var user GitUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error decoding request body: %w", err)
	}

	return &user, nil
}
