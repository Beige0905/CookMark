-- 1. 레시피 테이블 (JSONB 사용)
CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    origin_url TEXT,
    image_url TEXT,
    base_servings INTEGER DEFAULT 1,
    ingredients JSONB NOT NULL,
    instructions TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_recipes_ingredients ON recipes USING GIN (ingredients);
CREATE INDEX idx_recipes_origin_url ON recipes (origin_url);
CREATE INDEX idx_recipes_deleted_at ON recipes (deleted_at) WHERE deleted_at IS NULL;

-- 2. "만들었어!" 로그 테이블
CREATE TABLE IF NOT EXISTS cooking_logs (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE, 
    comment TEXT,
    cooked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. 보유 식재료 테이블
CREATE TABLE IF NOT EXISTS my_pantry (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    category TEXT,
    amount_num NUMERIC(10, 2),
    unit TEXT,
    expiry_date DATE,
    input_method TEXT DEFAULT 'manual',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_my_pantry_name ON my_pantry (name);
