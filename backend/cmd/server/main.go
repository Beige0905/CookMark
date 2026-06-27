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

	// repository 구현체 주입
	recipeRepo := repository.NewPgRecipeRepository(pool)
	noteRepo := repository.NewPgNoteRepository(pool)
	aiService := service.NewAIService()
	recipeService := service.NewRecipeService(recipeRepo)
	noteService := service.NewNoteService(noteRepo)

	recipeHandler := handler.NewRecipeHandler(recipeService)
	noteHandler := handler.NewNoteHandler(noteService)
	youtubeHandler := handler.NewYouTubeHandler(
		service.NewYouTubeService(aiService),
		recipeService,
	)
	ocrHandler := handler.NewOCRHandler(
		service.NewOCRService(aiService),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/recipes", recipeHandler.List)
	mux.HandleFunc("GET /api/recipes/{id}", recipeHandler.Get)
	mux.HandleFunc("POST /api/recipes", recipeHandler.Create)
	mux.HandleFunc("GET /api/recipes/{id}/note", noteHandler.Get)
	mux.HandleFunc("PUT /api/recipes/{id}/note", noteHandler.Put)
	mux.HandleFunc("POST /api/youtube/extract", youtubeHandler.Extract)
	mux.HandleFunc("POST /api/recipes/extract-image", ocrHandler.ExtractImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("서버 시작: :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CORS(mux)))
}
