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

func (r *PgRecipeRepository) FindAll(ctx context.Context, userID string) ([]model.Recipe, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, origin_url, image_url, base_servings, ingredients, instructions, created_at
		 FROM recipes WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`,
		userID,
	)
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

func (r *PgRecipeRepository) FindByID(ctx context.Context, id int, userID string) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, origin_url, image_url, base_servings, ingredients, instructions, created_at
		 FROM recipes WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		id, userID,
	).Scan(&recipe.ID, &recipe.Title, &recipe.OriginURL, &recipe.ImageURL, &recipe.BaseServings, &recipe.Ingredients, &recipe.Instructions, &recipe.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *PgRecipeRepository) FindByOriginURL(ctx context.Context, url string, userID string) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, origin_url, image_url, base_servings, ingredients, instructions, created_at
		 FROM recipes WHERE origin_url = $1 AND user_id = $2 AND deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT 1`,
		url, userID,
	).Scan(&recipe.ID, &recipe.Title, &recipe.OriginURL, &recipe.ImageURL, &recipe.BaseServings, &recipe.Ingredients, &recipe.Instructions, &recipe.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *PgRecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO recipes (title, origin_url, image_url, base_servings, ingredients, instructions, user_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at`,
		recipe.Title, recipe.OriginURL, recipe.ImageURL, recipe.BaseServings, recipe.Ingredients, recipe.Instructions, recipe.UserID,
	).Scan(&recipe.ID, &recipe.CreatedAt)
	if err != nil {
		return fmt.Errorf("레시피 생성 실패: %w", err)
	}
	return nil
}
