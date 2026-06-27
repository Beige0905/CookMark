package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/service"
	authpkg "github.com/Beige0905/recipe-backend/pkg/auth"
)

type NoteHandler struct {
	svc *service.NoteService
}

func NewNoteHandler(svc *service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	note, err := h.svc.GetByRecipeID(r.Context(), recipeID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, note)
}

func (h *NoteHandler) Put(w http.ResponseWriter, r *http.Request) {
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
	var note model.RecipeNote
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
		return
	}
	note.RecipeID = recipeID
	note.UserID = userID
	if err := h.svc.Upsert(r.Context(), &note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
