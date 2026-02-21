ALTER TABLE users
    ALTER COLUMN password_hash DROP NOT NULL,
    ADD COLUMN provider    VARCHAR(50) NOT NULL DEFAULT 'local',
    ADD COLUMN provider_id TEXT,
    ADD COLUMN name        VARCHAR(255),
    ADD COLUMN avatar_url  TEXT;


CREATE UNIQUE INDEX idx_users_provider_id 
    ON users(provider, provider_id) 
    WHERE provider_id IS NOT NULL;