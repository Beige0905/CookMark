package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Beige0905/recipe-backend/internal/service"
)

func isSecure() bool {
	return os.Getenv("APP_ENV") == "production"
}

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		http.Error(w, "이메일, 비밀번호, 이름은 필수입니다", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "비밀번호는 8자 이상이어야 합니다", http.StatusBadRequest)
		return
	}
	user, err := h.svc.Register(r.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
		return
	}
	accessToken, refreshToken, user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	setTokenCookies(w, accessToken, refreshToken)
	writeJSON(w, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.svc.Logout(r.Context(), refreshCookie.Value)
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "리프레시 토큰이 없습니다", http.StatusUnauthorized)
		return
	}
	accessToken, user, err := h.svc.Refresh(r.Context(), refreshCookie.Value)
	if err != nil {
		clearTokenCookies(w)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure(),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})
	writeJSON(w, user)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
		return
	}
	userID, err := h.svc.ValidateAccessToken(cookie.Value)
	if err != nil {
		http.Error(w, "유효하지 않은 토큰", http.StatusUnauthorized)
		return
	}
	user, err := h.svc.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "유저를 찾을 수 없습니다", http.StatusNotFound)
		return
	}
	writeJSON(w, user)
}

func setTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	secure := isSecure()
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
}

func clearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
