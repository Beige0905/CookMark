package service

import (
	"context"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/repository"
)

type RecipeService struct {
	repo repository.RecipeRepository
}

func NewRecipeService(repo repository.RecipeRepository) *RecipeService {
	return &RecipeService{repo: repo}
}

func (s *RecipeService) GetAll(ctx context.Context) ([]model.Recipe, error) {
	return s.repo.FindAll(ctx)
}

func (s *RecipeService) GetByID(ctx context.Context, id int) (*model.Recipe, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *RecipeService) Create(ctx context.Context, r *model.Recipe) error {
	return s.repo.Create(ctx, r)
}
