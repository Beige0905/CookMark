package repository

import (
	"context"

	"github.com/Beige0905/recipe-backend/internal/model"
)

type RecipeRepository interface {
	FindAll(ctx context.Context, userID string) ([]model.Recipe, error)
	FindByID(ctx context.Context, id int, userID string) (*model.Recipe, error)
	FindByOriginURL(ctx context.Context, url string, userID string) (*model.Recipe, error)
	Create(ctx context.Context, r *model.Recipe) error
	Update(ctx context.Context, r *model.Recipe) error
	Delete(ctx context.Context, id int, userID string) error
}

type CookingLogRepository interface {
	FindByRecipeID(ctx context.Context, recipeID int) ([]model.CookingLog, error)
	Create(ctx context.Context, log *model.CookingLog) error
}

type PantryRepository interface {
	FindAll(ctx context.Context, userID string) ([]model.PantryItem, error)
	Create(ctx context.Context, item *model.PantryItem) error
	Delete(ctx context.Context, id int, userID string) error
}

type NoteRepository interface {
	FindByRecipeID(ctx context.Context, recipeID int, userID string) (*model.RecipeNote, error)
	Upsert(ctx context.Context, note *model.RecipeNote) error
}
