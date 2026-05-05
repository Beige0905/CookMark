-- 1. 레시피 정보를 담는 테이블
CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    origin_url TEXT,
    base_servings INTEGER DEFAULT 1,
    ingredients JSONB NOT NULL,         -- JSONB로 재료들 보관
    instructions TEXT[],                -- 조리 순서
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. "만들었어!" 버튼을 누를 때마다 기록되는 로그 테이블
CREATE TABLE IF NOT EXISTS cooking_logs (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    comment TEXT,                       -- 코멘트 (선택 사항)
    cooked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- 만든 날짜/시간
);

CREATE INDEX idx_recipes_ingredients ON recipes USING GIN (ingredients);

-- 보유 식재료 (단순 버전)
CREATE TABLE IF NOT EXISTS my_pantry (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,       -- 재료명 (예: 양파)
    category TEXT,                   -- 채소, 고기, 양념 등
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);