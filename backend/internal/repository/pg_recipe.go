package repository

import (
	"context"
	"fmt"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRecipeRepository struct {
	pool *pgxpool.Pool
}

func NewPgRecipeRepository(pool *pgxpool.Pool) *PgRecipeRepository {
	return &PgRecipeRepository{pool: pool}
}

func (r *PgRecipeRepository) FindAll(ctx context.Context) ([]model.Recipe, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, origin_url, image_url, base_servings, ingredients, instructions, created_at FROM recipes WHERE deleted_at IS NULL ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := make([]model.Recipe, 0)
	for rows.Next() {
		var recipe model.Recipe
		err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.OriginURL, &recipe.ImageURL, &recipe.BaseServings, &recipe.Ingredients, &recipe.Instructions, &recipe.CreatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *PgRecipeRepository) FindByID(ctx context.Context, id int) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.pool.QueryRow(ctx, "SELECT id, title, origin_url, image_url, base_servings, ingredients, instructions, created_at FROM recipes WHERE id = $1 AND deleted_at IS NULL", id).
		Scan(&recipe.ID, &recipe.Title, &recipe.OriginURL, &recipe.ImageURL, &recipe.BaseServings, &recipe.Ingredients, &recipe.Instructions, &recipe.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *PgRecipeRepository) FindByOriginURL(ctx context.Context, url string) (*model.Recipe, error) {
	var recipe model.Recipe
	query := `
		SELECT id, title, origin_url, image_url, base_servings, ingredients, instructions, created_at
		FROM recipes
		WHERE origin_url = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.pool.QueryRow(ctx, query, url).
		Scan(&recipe.ID, &recipe.Title, &recipe.OriginURL, &recipe.ImageURL, &recipe.BaseServings, &recipe.Ingredients, &recipe.Instructions, &recipe.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *PgRecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	query := `
		INSERT INTO recipes (title, origin_url, image_url, base_servings, ingredients, instructions)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query, recipe.Title, recipe.OriginURL, recipe.ImageURL, recipe.BaseServings, recipe.Ingredients, recipe.Instructions).
		Scan(&recipe.ID, &recipe.CreatedAt)
	if err != nil {
		return fmt.Errorf("레시피 생성 실패: %w", err)
	}
	return nil
}
