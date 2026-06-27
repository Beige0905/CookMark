-- 1. 사용자 테이블
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name  TEXT NOT NULL,
    avatar_url    TEXT,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. 리프레시 토큰
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens (token_hash);

-- 3. 레시피
CREATE TABLE IF NOT EXISTS recipes (
    id            SERIAL PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title         TEXT NOT NULL,
    origin_url    TEXT,
    image_url     TEXT,
    base_servings INTEGER DEFAULT 1,
    ingredients   JSONB NOT NULL,
    instructions  TEXT[],
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_recipes_user_id ON recipes (user_id);
CREATE INDEX idx_recipes_ingredients ON recipes USING GIN (ingredients);
CREATE INDEX idx_recipes_origin_url ON recipes (origin_url);
CREATE INDEX idx_recipes_deleted_at ON recipes (deleted_at) WHERE deleted_at IS NULL;

-- 4. 개인 메모 & 재료 조정 (deprecated: CookingLog로 대체 예정)
CREATE TABLE IF NOT EXISTS recipe_notes (
    id          SERIAL PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id   INTEGER NOT NULL UNIQUE REFERENCES recipes(id) ON DELETE CASCADE,
    memo        TEXT NOT NULL DEFAULT '',
    adjustments JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_recipe_notes_user_id ON recipe_notes (user_id);

-- 5. 요리 기록
CREATE TABLE IF NOT EXISTS cooking_logs (
    id        SERIAL PRIMARY KEY,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    comment   TEXT,
    cooked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cooking_logs_user_id ON cooking_logs (user_id);

-- 6. 냉장고 재료
CREATE TABLE IF NOT EXISTS my_pantry (
    id           SERIAL PRIMARY KEY,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    category     TEXT,
    amount_num   NUMERIC(10, 2),
    unit         TEXT,
    expiry_date  DATE,
    input_method TEXT DEFAULT 'manual',
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, name)
);

CREATE INDEX idx_my_pantry_user_id ON my_pantry (user_id);
CREATE INDEX idx_my_pantry_name ON my_pantry (name);
