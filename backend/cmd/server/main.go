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
	noteRepo := repository.NewPgNoteRepository(pool)
	userRepo := repository.NewPgUserRepository(pool)

	aiService := service.NewAIService()
	authService := service.NewAuthService(userRepo)
	recipeService := service.NewRecipeService(recipeRepo)
	noteService := service.NewNoteService(noteRepo)

	recipeHandler := handler.NewRecipeHandler(recipeService)
	noteHandler := handler.NewNoteHandler(noteService)
	authHandler := handler.NewAuthHandler(authService)
	youtubeHandler := handler.NewYouTubeHandler(
		service.NewYouTubeService(aiService),
		recipeService,
	)
	ocrHandler := handler.NewOCRHandler(
		service.NewOCRService(aiService),
	)

	auth := middleware.Auth(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("GET /api/auth/me", authHandler.Me)

	mux.HandleFunc("GET /api/recipes", auth(recipeHandler.List))
	mux.HandleFunc("GET /api/recipes/{id}", auth(recipeHandler.Get))
	mux.HandleFunc("POST /api/recipes", auth(recipeHandler.Create))
	mux.HandleFunc("GET /api/recipes/{id}/note", auth(noteHandler.Get))
	mux.HandleFunc("PUT /api/recipes/{id}/note", auth(noteHandler.Put))
	mux.HandleFunc("POST /api/youtube/extract", auth(youtubeHandler.Extract))
	mux.HandleFunc("POST /api/recipes/extract-image", auth(ocrHandler.ExtractImage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("서버 시작: :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CORS(mux)))
}
