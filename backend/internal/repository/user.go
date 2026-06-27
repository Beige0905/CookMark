package repository

import (
	"context"
	"time"

	"github.com/Beige0905/recipe-backend/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User, passwordHash string) error
	FindByEmail(ctx context.Context, email string) (*model.User, string, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
}
