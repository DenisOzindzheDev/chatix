package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/DenisOzindzheDev/chatix/auth/internal/db"
	"github.com/DenisOzindzheDev/chatix/auth/internal/github"
)

type User struct {
	ID        int64
	GithubID  int64
	Username  string
	Email     string
	AvatarUrl string
	CreatedAt time.Time
}

func FindUserByGithubID(ctx context.Context, githubID int64) (*User, error) {
	query := `
		SELECT id, github_id, username, email, avatar_url, created_at
		FROM users 
		WHERE github_id = $1
	`

	row := db.GetDB().QueryRowContext(ctx, query, githubID)

	var u User
	err := row.Scan(&u.ID, &u.GithubID, &u.Username, &u.Email, &u.AvatarUrl, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func CreateUser(ctx context.Context, profile *github.GitUser) (*User, error) {
	query := `
		INSERT INTO users(github_id, username, email, avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, github_id, username, email, avatar_url, created_at
	`
	row := db.GetDB().QueryRowContext(ctx, query,
		profile.ID,
		profile.Login,
		profile.Email,
		profile.AvatarUrl,
	)

	var u User

	if err := row.Scan(&u.ID, &u.GithubID, &u.Username, &u.Email, &u.AvatarUrl, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
