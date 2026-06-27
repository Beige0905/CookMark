package repository

import (
	"context"
	"fmt"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPantryRepository struct {
	pool *pgxpool.Pool
}

func NewPgPantryRepository(pool *pgxpool.Pool) *PgPantryRepository {
	return &PgPantryRepository{pool: pool}
}

func (r *PgPantryRepository) FindAll(ctx context.Context, userID string) ([]model.PantryItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, created_at FROM my_pantry WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.PantryItem, 0)
	for rows.Next() {
		var item model.PantryItem
		if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *PgPantryRepository) Create(ctx context.Context, item *model.PantryItem) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO my_pantry (user_id, name) VALUES ($1, $2) RETURNING id, created_at`,
		item.UserID, item.Name,
	).Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return fmt.Errorf("재료 추가 실패: %w", err)
	}
	return nil
}

func (r *PgPantryRepository) Delete(ctx context.Context, id int, userID string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM my_pantry WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("재료 삭제 실패: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("재료를 찾을 수 없습니다")
	}
	return nil
}
