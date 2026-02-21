package postgres

const (
	queryCreateUser = `
    INSERT INTO users (id, email, name, password_hash, provider, role, status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	queryFindByEmail = `
    SELECT id, email, name, password_hash, provider, role, status, created_at, updated_at
    FROM users WHERE email = $1
`
	queryFindByID = `
        SELECT id, email, name, password_hash, role, status, created_at, updated_at
        FROM users WHERE id = $1
    `
)
