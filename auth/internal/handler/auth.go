package handler

import (
	"database/sql"
	"net/http"

	"github.com/DenisOzindzheDev/chatix/auth/internal/github"
	"github.com/DenisOzindzheDev/chatix/auth/internal/repository"
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
	// Проверка state
	if state := c.Query("state"); state != "state-token" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state token"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// Обмен кода на токен
	token, err := github.GetOAuthConfig().Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Получение профиля пользователя
	user, err := github.GetUserProfile(c, token.AccessToken)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Could not retrieve user profile from GitHub",
		})
		return
	}

	// Поиск/создание пользователя
	existingUser, err := repository.FindUserByGithubID(c, user.ID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not check user existence",
		})
		return
	}

	var u *repository.User
	if existingUser == nil {
		u, err = repository.CreateUser(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not create user account",
			})
			return
		}
	} else {
		u = existingUser
	}

	// Возвращаем только необходимые данные
	c.JSON(http.StatusOK, gin.H{
		"username":   u.Username,
		"email":      u.Email,
		"avatar_url": u.AvatarUrl,
	})
}
