package service

import (
	"context"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/repository"
)

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) GetByRecipeID(ctx context.Context, recipeID int) (*model.RecipeNote, error) {
	return s.repo.FindByRecipeID(ctx, recipeID)
}

func (s *NoteService) Upsert(ctx context.Context, note *model.RecipeNote) error {
	return s.repo.Upsert(ctx, note)
}
