package postgres

const (
	queryCreateUser = `
        INSERT INTO users (id, email, name, password_hash, provider, provider_id, avatar_url, role, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	queryFindByEmail = `
        SELECT id, email, name, password_hash, provider, provider_id, avatar_url, role, status, created_at, updated_at
        FROM users WHERE email = $1
    `

	queryFindByProviderID = `
        SELECT id, email, name, password_hash, provider, provider_id, avatar_url, role, status, created_at, updated_at
        FROM users WHERE provider = $1 AND provider_id = $2
    `
	queryUpdatePassword = `
    UPDATE users 
    SET password_hash = $1, updated_at = $2 
    WHERE id = $3
`
	queryCreateResetToken = `
    INSERT INTO password_reset_tokens (user_id, token_hash, expires_at, created_at)
    VALUES ($1, $2, $3, $4)
`

	queryFindByTokenHash = `
    SELECT user_id FROM password_reset_tokens
    WHERE token_hash = $1 AND expires_at > NOW()
`

	queryDeleteByTokenHash = `
    DELETE FROM password_reset_tokens WHERE token_hash = $1
`

	queryDeleteAllForUser = `
    DELETE FROM password_reset_tokens WHERE user_id = $1
`
)
