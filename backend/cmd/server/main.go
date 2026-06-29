package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Beige0905/recipe-backend/internal/handler"
	"github.com/Beige0905/recipe-backend/internal/repository"
	"github.com/Beige0905/recipe-backend/internal/service"
	"github.com/Beige0905/recipe-backend/pkg/database"
	"github.com/Beige0905/recipe-backend/pkg/middleware"
)

func main() {
	ctx := context.Background()

	pool, err := database.Connect(ctx)
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	defer pool.Close()

	recipeRepo := repository.NewPgRecipeRepository(pool)
	userRepo := repository.NewPgUserRepository(pool)
	pantryRepo := repository.NewPgPantryRepository(pool)
	logRepo := repository.NewPgCookingLogRepository(pool)

	aiService := service.NewAIService()
	authService := service.NewAuthService(userRepo)
	recipeService := service.NewRecipeService(recipeRepo)
	pantryService := service.NewPantryService(pantryRepo, recipeRepo)
	logService := service.NewCookingLogService(logRepo)

	recipeHandler := handler.NewRecipeHandler(recipeService)
	authHandler := handler.NewAuthHandler(authService)
	pantryHandler := handler.NewPantryHandler(pantryService)
	logHandler := handler.NewCookingLogHandler(logService)
	youtubeHandler := handler.NewYouTubeHandler(
		service.NewYouTubeService(aiService),
		recipeService,
	)
	ocrHandler := handler.NewOCRHandler(
		service.NewOCRService(aiService),
	)

	auth := middleware.Auth(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("GET /api/auth/me", authHandler.Me)

	mux.HandleFunc("GET /api/recipes", auth(recipeHandler.List))
	mux.HandleFunc("GET /api/recipes/{id}", auth(recipeHandler.Get))
	mux.HandleFunc("POST /api/recipes", auth(recipeHandler.Create))
	mux.HandleFunc("PUT /api/recipes/{id}", auth(recipeHandler.Update))
	mux.HandleFunc("DELETE /api/recipes/{id}", auth(recipeHandler.Delete))
	mux.HandleFunc("GET /api/recipes/{id}/logs", auth(logHandler.List))
	mux.HandleFunc("POST /api/recipes/{id}/logs", auth(logHandler.Create))
	mux.HandleFunc("DELETE /api/recipes/{id}/logs/{logId}", auth(logHandler.Delete))
	mux.HandleFunc("GET /api/recipes/{id}/pantry-matches", auth(pantryHandler.MatchForRecipe))
	mux.HandleFunc("POST /api/youtube/extract", auth(youtubeHandler.Extract))
	mux.HandleFunc("POST /api/recipes/extract-image", auth(ocrHandler.ExtractImage))

	mux.HandleFunc("GET /api/pantry", auth(pantryHandler.List))
	mux.HandleFunc("POST /api/pantry", auth(pantryHandler.Create))
	mux.HandleFunc("DELETE /api/pantry/{id}", auth(pantryHandler.Delete))
	mux.HandleFunc("GET /api/pantry/recommend", auth(pantryHandler.Recommend))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("서버 시작: :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CORS(mux)))
}
