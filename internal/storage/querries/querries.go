package querries

const (
	CheckRecord = `
		SELECT EXISTS (
			SELECT 1
			FROM goods
			WHERE id = $1 AND project_id = $2
		);
	`

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
	// --------------------------

	// RemoveGoodsQuerries -----
	RemoveQuery = `
		UPDATE goods
		SET removed = true
		WHERE id = $1 AND project_id = $2
		RETURNING id, project_id, removed;
	`
	// --------------------------

	// ListGoodsQuerries ------
	ListQuery = `
			SELECT id, project_id, name, description, priority, removed, created_at
			FROM goods
			ORDER BY created_at
			LIMIT $1 OFFSET $2;
		`

	CountTotalQuery = `
		SELECT COUNT(id) FROM goods;
    `
	CountTotalRemovedQuery = `
		SELECT COUNT(id) FROM goods WHERE removed = TRUE;

    `
	// --------------------------

	// ReprioritiizeQuerries ------
	UpdatePriority = `
		UPDATE goods
		SET priority = $3
		WHERE id = $1 AND project_id = $2
    `

	RepriotiizeQuery = `
		UPDATE goods
		SET priority = priority + 1
		WHERE project_id = $1 AND priority < $2
		RETURNING id, priority;
    `

	RepriotiizeSelectQuery = `
		SELECT id, priority
		FROM goods;
    `
	// --------------------------

	InsetIntoClickhouse = ` 
		INSERT INTO clickhouse 
		(id, project_id, name, description, priority, removed, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`
)
