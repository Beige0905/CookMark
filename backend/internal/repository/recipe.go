package repository

import (
	"context"

	"github.com/Beige0905/recipe-backend/internal/model"
)

type RecipeRepository interface {
	FindAll(ctx context.Context) ([]model.Recipe, error)
	FindByID(ctx context.Context, id int) (*model.Recipe, error)
	Create(ctx context.Context, r *model.Recipe) error
}

type CookingLogRepository interface {
	FindByRecipeID(ctx context.Context, recipeID int) ([]model.CookingLog, error)
	Create(ctx context.Context, log *model.CookingLog) error
}

type PantryRepository interface {
	FindAll(ctx context.Context) ([]model.PantryItem, error)
	Create(ctx context.Context, item *model.PantryItem) error
	Delete(ctx context.Context, id int) error
}
