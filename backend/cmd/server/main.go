package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Beige0905/recipe-backend/internal/handler"
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

	// TODO: repository 구현체를 만들고 여기에 주입하세요
	// 예: repo := repository.NewPgRecipeRepository(pool)
	recipeHandler := handler.NewRecipeHandler(
		service.NewRecipeService(nil),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/recipes", recipeHandler.List)
	mux.HandleFunc("GET /api/recipes/{id}", recipeHandler.Get)
	mux.HandleFunc("POST /api/recipes", recipeHandler.Create)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("서버 시작: :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CORS(mux)))
}
