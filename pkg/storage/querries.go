package storage

const (
	CreateQuery = `
		INSERT INTO goods (id, project_id, name, description, priority, removed, creadted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, project_id, name, description, priority, removed, created_at
	`
)
