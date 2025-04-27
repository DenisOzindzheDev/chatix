package handler

import (
	"net/http"

	"github.com/DenisOzindzheDev/chatix/auth/internal/github"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func RegisterAuthRoutes(router *gin.Engine) {
	router.GET("/login/github", githubLoginHandler)
	router.GET("/callback/github", githubCallbackHandler)
}

func githubLoginHandler(c *gin.Context) {
	url := github.GetOAuthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func githubCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MissingCode"})
		return
	}

	token, err := github.GetOAuthConfig().Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token exchange failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken,
	})
}
