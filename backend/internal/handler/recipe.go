package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/service"
)

type RecipeHandler struct {
	svc *service.RecipeService
}

func NewRecipeHandler(svc *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{svc: svc}
}

func (h *RecipeHandler) List(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, recipes)
}

func (h *RecipeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "잘못된 ID", http.StatusBadRequest)
		return
	}
	recipe, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, recipe)
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var recipe model.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(r.Context(), &recipe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, recipe)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
