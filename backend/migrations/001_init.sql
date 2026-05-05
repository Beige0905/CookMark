CREATE TABLE IF NOT EXISTS recipes (
    id            SERIAL  PRIMARY KEY,
    title         TEXT    NOT NULL,
    origin_url    TEXT,
    base_servings INTEGER DEFAULT 1,
    ingredients   JSONB   NOT NULL,
    instructions  TEXT[],
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cooking_logs (
    id        SERIAL  PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    comment   TEXT,
    cooked_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_recipes_ingredients ON recipes USING GIN (ingredients);

CREATE TABLE IF NOT EXISTS my_pantry (
    id         SERIAL PRIMARY KEY,
    name       TEXT   NOT NULL UNIQUE,
    category   TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
