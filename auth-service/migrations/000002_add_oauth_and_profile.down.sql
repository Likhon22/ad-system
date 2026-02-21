DROP INDEX IF EXISTS idx_users_provider_id;
ALTER TABLE users
    ALTER COLUMN password_hash SET NOT NULL,
    DROP COLUMN provider,
    DROP COLUMN provider_id,
    DROP COLUMN name,
    DROP COLUMN avatar_url;
