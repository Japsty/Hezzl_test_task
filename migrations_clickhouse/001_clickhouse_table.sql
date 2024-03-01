-- +goose Up
CREATE TABLE IF NOT EXISTS clickhouse
(
    id         INT ,
    project_id INT ,
    name       VARCHAR(255),
    description VARCHAR(255),
    priority INT,
    removed BOOLEAN,
    created_at TIMESTAMP
);

CREATE INDEX id_and_project_id_index ON clickhouse USING btree (id, project_id);
CREATE INDEX name_index ON clickhouse USING hash (name);


-- +goose Down
DROP TABLE IF EXISTS clickhouse;