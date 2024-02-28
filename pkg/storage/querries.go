package storage

const (
	// CreateGoodQuerries -----
	SelectMaxPriority = `
		SELECT COALESCE(MAX(priority), 0) 
		FROM goods;
	`
	CreateQuery = `
		INSERT INTO goods (project_id, name, description, priority, removed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, project_id, name, description, priority, removed, created_at;
	`
	// ------------------------

	// UpdateGoodsQuerries -----
	UpdateQuery = `
		UPDATE goods
		SET name = $3, description = $4
		WHERE id = $1 AND project_id = $2
		RETURNING id, project_id, name, description, priority, removed, created_at;
	`
	CheckRecord = `
		SELECT EXISTS (
			SELECT 1
			FROM goods
			WHERE id = $1 AND project_id = $2
		);
	`
	// --------------------------
)
