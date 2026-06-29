package service

import (
	"context"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/repository"
)

type CookingLogService struct {
	repo repository.CookingLogRepository
}

func NewCookingLogService(repo repository.CookingLogRepository) *CookingLogService {
	return &CookingLogService{repo: repo}
}

func (s *CookingLogService) GetByRecipeID(ctx context.Context, recipeID int, userID string) ([]model.CookingLog, error) {
	return s.repo.FindByRecipeID(ctx, recipeID, userID)
}

func (s *CookingLogService) Create(ctx context.Context, l *model.CookingLog, pantryIDs []int) error {
	if len(pantryIDs) > 0 {
		return s.repo.CreateAndDeletePantry(ctx, l, pantryIDs, l.UserID)
	}
	return s.repo.Create(ctx, l)
}

func (s *CookingLogService) Delete(ctx context.Context, id int, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
