package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/service"
	authpkg "github.com/Beige0905/recipe-backend/pkg/auth"
)

type PantryHandler struct {
	svc *service.PantryService
}

func NewPantryHandler(svc *service.PantryService) *PantryHandler {
	return &PantryHandler{svc: svc}
}

func (h *PantryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	items, err := h.svc.GetAll(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

func (h *PantryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, "재료 이름이 필요합니다", http.StatusBadRequest)
		return
	}
	item := &model.PantryItem{UserID: userID, Name: body.Name}
	if err := h.svc.Add(r.Context(), item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, item)
}

func (h *PantryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "잘못된 ID", http.StatusBadRequest)
		return
	}
	if err := h.svc.Remove(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PantryHandler) Recommend(w http.ResponseWriter, r *http.Request) {
	userID, ok := authpkg.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	results, err := h.svc.Recommend(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if results == nil {
		results = []service.RecommendResult{}
	}
	writeJSON(w, results)
}
