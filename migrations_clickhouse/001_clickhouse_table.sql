-- +goose Up
CREATE TABLE IF NOT EXISTS clickhouse
(
    id          INT,
    project_id  INT,
    name        VARCHAR(255),
    description VARCHAR(255),
    priority    INT,
    removed     UInt8,
    event_time  TIMESTAMP
) ENGINE = MergeTree()
ORDER BY (id, project_id, name);

-- +goose Down
DROP TABLE IF EXISTS clickhouse;