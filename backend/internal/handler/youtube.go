package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Beige0905/recipe-backend/internal/service"
)

type YouTubeHandler struct {
	youtubeService *service.YouTubeService
	recipeService  *service.RecipeService
}

func NewYouTubeHandler(ys *service.YouTubeService, rs *service.RecipeService) *YouTubeHandler {
	return &YouTubeHandler{
		youtubeService: ys,
		recipeService:  rs,
	}
}

type youtubeExtractRequest struct {
	URL string `json:"url"`
}

func (h *YouTubeHandler) Extract(w http.ResponseWriter, r *http.Request) {
	var req youtubeExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "URL이 필요합니다", http.StatusBadRequest)
		return
	}

	// 1. 기존 레시피가 있는지 확인 (URL 중복 캐싱)
	if existing, err := h.recipeService.FindByOriginURL(r.Context(), req.URL); err == nil && existing != nil {
		var ingredients []service.Ingredient
		if err := json.Unmarshal(existing.Ingredients, &ingredients); err == nil {
			resp := service.YouTubeExtractResponse{
				Title:        existing.Title,
				BaseServings: existing.BaseServings,
				Ingredients:  ingredients,
			}
			if existing.ImageURL != nil {
				resp.ImageURL = *existing.ImageURL
			}
			writeJSON(w, resp)
			return
		}
	}

	// 2. 없으면 AI 추출 실행
	recipeData, err := h.youtubeService.ExtractRecipeData(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, recipeData)
}
