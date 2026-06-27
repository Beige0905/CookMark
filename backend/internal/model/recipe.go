package model

import (
	"encoding/json"
	"time"
)

type Recipe struct {
	ID           int             `json:"id"`
	UserID       string          `json:"-"`
	Title        string          `json:"title"`
	OriginURL    *string         `json:"origin_url,omitempty"`
	ImageURL     *string         `json:"image_url,omitempty"`
	BaseServings int             `json:"base_servings"`
	Ingredients  json.RawMessage `json:"ingredients"`
	Instructions []string        `json:"instructions"`
	CreatedAt    time.Time       `json:"created_at"`
}

type CookingLog struct {
	ID       int       `json:"id"`
	RecipeID int       `json:"recipe_id"`
	Comment  *string   `json:"comment,omitempty"`
	CookedAt time.Time `json:"cooked_at"`
}

type PantryItem struct {
	ID        int       `json:"id"`
	UserID    string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RecipeNote struct {
	RecipeID    int             `json:"recipe_id"`
	UserID      string          `json:"-"`
	Memo        string          `json:"memo"`
	Adjustments json.RawMessage `json:"adjustments"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
