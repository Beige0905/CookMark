package repository

import (
	"context"
	"errors"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgNoteRepository struct {
	pool *pgxpool.Pool
}

func NewPgNoteRepository(pool *pgxpool.Pool) *PgNoteRepository {
	return &PgNoteRepository{pool: pool}
}

func (r *PgNoteRepository) FindByRecipeID(ctx context.Context, recipeID int, userID string) (*model.RecipeNote, error) {
	note := &model.RecipeNote{RecipeID: recipeID, UserID: userID}
	err := r.pool.QueryRow(ctx,
		`SELECT memo, adjustments, updated_at FROM recipe_notes WHERE recipe_id = $1 AND user_id = $2`,
		recipeID, userID,
	).Scan(&note.Memo, &note.Adjustments, &note.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		note.Adjustments = []byte("{}")
		return note, nil
	}
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *PgNoteRepository) Upsert(ctx context.Context, note *model.RecipeNote) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO recipe_notes (recipe_id, user_id, memo, adjustments, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (recipe_id) DO UPDATE
		SET memo = EXCLUDED.memo,
		    adjustments = EXCLUDED.adjustments,
		    updated_at = NOW()
	`, note.RecipeID, note.UserID, note.Memo, note.Adjustments)
	return err
}
