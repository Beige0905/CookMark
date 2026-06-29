package repository

import (
	"context"
	"fmt"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgCookingLogRepository struct {
	pool *pgxpool.Pool
}

func NewPgCookingLogRepository(pool *pgxpool.Pool) *PgCookingLogRepository {
	return &PgCookingLogRepository{pool: pool}
}

func (r *PgCookingLogRepository) FindByRecipeID(ctx context.Context, recipeID int, userID string) ([]model.CookingLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, recipe_id, comment, cooked_at FROM cooking_logs
		 WHERE recipe_id = $1 AND user_id = $2 ORDER BY cooked_at DESC`,
		recipeID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]model.CookingLog, 0)
	for rows.Next() {
		var l model.CookingLog
		if err := rows.Scan(&l.ID, &l.RecipeID, &l.Comment, &l.CookedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *PgCookingLogRepository) Create(ctx context.Context, l *model.CookingLog) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO cooking_logs (user_id, recipe_id, comment, cooked_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		l.UserID, l.RecipeID, l.Comment, l.CookedAt,
	).Scan(&l.ID)
	if err != nil {
		return fmt.Errorf("요리 기록 추가 실패: %w", err)
	}
	return nil
}

func (r *PgCookingLogRepository) CreateAndDeletePantry(ctx context.Context, l *model.CookingLog, pantryIDs []int, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO cooking_logs (user_id, recipe_id, comment, cooked_at) VALUES ($1, $2, $3, $4) RETURNING id`,
		l.UserID, l.RecipeID, l.Comment, l.CookedAt,
	).Scan(&l.ID)
	if err != nil {
		return fmt.Errorf("요리 기록 추가 실패: %w", err)
	}

	if len(pantryIDs) > 0 {
		_, err = tx.Exec(ctx,
			`DELETE FROM my_pantry WHERE id = ANY($1) AND user_id = $2`,
			pantryIDs, userID,
		)
		if err != nil {
			return fmt.Errorf("재료 삭제 실패: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PgCookingLogRepository) Delete(ctx context.Context, id int, userID string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM cooking_logs WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("요리 기록 삭제 실패: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("요리 기록을 찾을 수 없습니다")
	}
	return nil
}
