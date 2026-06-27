package middleware

import (
	"net/http"

	authpkg "github.com/Beige0905/recipe-backend/pkg/auth"
	"github.com/Beige0905/recipe-backend/internal/service"
)

func Auth(authSvc *service.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
				return
			}
			userID, err := authSvc.ValidateAccessToken(cookie.Value)
			if err != nil {
				http.Error(w, "유효하지 않은 토큰", http.StatusUnauthorized)
				return
			}
			ctx := authpkg.WithUserID(r.Context(), userID)
			next(w, r.WithContext(ctx))
		}
	}
}
