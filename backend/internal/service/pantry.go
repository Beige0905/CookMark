package service

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/repository"
)

type PantryService struct {
	pantryRepo repository.PantryRepository
	recipeRepo repository.RecipeRepository
}

func NewPantryService(pantryRepo repository.PantryRepository, recipeRepo repository.RecipeRepository) *PantryService {
	return &PantryService{pantryRepo: pantryRepo, recipeRepo: recipeRepo}
}

func (s *PantryService) GetAll(ctx context.Context, userID string) ([]model.PantryItem, error) {
	return s.pantryRepo.FindAll(ctx, userID)
}

func (s *PantryService) Add(ctx context.Context, item *model.PantryItem) error {
	return s.pantryRepo.Create(ctx, item)
}

func (s *PantryService) Remove(ctx context.Context, id int, userID string) error {
	return s.pantryRepo.Delete(ctx, id, userID)
}

type RecommendResult struct {
	Recipe       model.Recipe `json:"recipe"`
	MatchedCount int          `json:"matched_count"`
	TotalCount   int          `json:"total_count"`
}

func (s *PantryService) Recommend(ctx context.Context, userID string) ([]RecommendResult, error) {
	items, err := s.pantryRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	recipes, err := s.recipeRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	pantryNames := make([]string, len(items))
	for i, item := range items {
		pantryNames[i] = strings.ToLower(item.Name)
	}

	var results []RecommendResult
	for _, recipe := range recipes {
		var ingredients []struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(recipe.Ingredients, &ingredients); err != nil {
			continue
		}
		total := len(ingredients)
		if total == 0 {
			continue
		}
		matched := 0
		for _, ing := range ingredients {
			ingName := strings.ToLower(ing.Name)
			for _, pantry := range pantryNames {
				if strings.Contains(ingName, pantry) || strings.Contains(pantry, ingName) {
					matched++
					break
				}
			}
		}
		if matched > 0 {
			results = append(results, RecommendResult{
				Recipe:       recipe,
				MatchedCount: matched,
				TotalCount:   total,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		si := float64(results[i].MatchedCount) / float64(results[i].TotalCount)
		sj := float64(results[j].MatchedCount) / float64(results[j].TotalCount)
		return si > sj
	})

	return results, nil
}
