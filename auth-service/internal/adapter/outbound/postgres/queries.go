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
)
