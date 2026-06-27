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

func (s *RecipeService) GetAll(ctx context.Context, userID string) ([]model.Recipe, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *RecipeService) GetByID(ctx context.Context, id int, userID string) (*model.Recipe, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *RecipeService) FindByOriginURL(ctx context.Context, url string, userID string) (*model.Recipe, error) {
	return s.repo.FindByOriginURL(ctx, url, userID)
}

func (s *RecipeService) Create(ctx context.Context, r *model.Recipe) error {
	return s.repo.Create(ctx, r)
}

func (s *RecipeService) Update(ctx context.Context, r *model.Recipe) error {
	return s.repo.Update(ctx, r)
}

func (s *RecipeService) Delete(ctx context.Context, id int, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
