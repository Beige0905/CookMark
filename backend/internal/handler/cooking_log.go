package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/service"
	authpkg "github.com/Beige0905/recipe-backend/pkg/auth"
)

type CookingLogHandler struct {
	svc *service.CookingLogService
}

func NewCookingLogHandler(svc *service.CookingLogService) *CookingLogHandler {
	return &CookingLogHandler{svc: svc}
}

func (h *CookingLogHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	recipeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "잘못된 ID", http.StatusBadRequest)
		return
	}
	logs, err := h.svc.GetByRecipeID(r.Context(), recipeID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, logs)
}

func (h *CookingLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	recipeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "잘못된 ID", http.StatusBadRequest)
		return
	}
	var body struct {
		Comment  *string `json:"comment"`
		CookedAt *string `json:"cooked_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
		return
	}
	cookedAt := time.Now()
	if body.CookedAt != nil && *body.CookedAt != "" {
		if t, err := time.Parse("2006-01-02", *body.CookedAt); err == nil {
			cookedAt = t
		}
	}
	l := &model.CookingLog{
		UserID:   userID,
		RecipeID: recipeID,
		Comment:  body.Comment,
		CookedAt: cookedAt,
	}
	if err := h.svc.Create(r.Context(), l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, l)
}

func (h *CookingLogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(r.PathValue("logId"))
	if err != nil {
		http.Error(w, "잘못된 ID", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
