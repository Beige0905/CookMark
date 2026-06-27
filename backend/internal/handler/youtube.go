package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Beige0905/recipe-backend/internal/service"
	authpkg "github.com/Beige0905/recipe-backend/pkg/auth"
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
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}

	var req youtubeExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "URL이 필요합니다", http.StatusBadRequest)
		return
	}

	if existing, err := h.recipeService.FindByOriginURL(r.Context(), req.URL, userID); err == nil && existing != nil {
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

	recipeData, err := h.youtubeService.ExtractRecipeData(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, recipeData)
}
