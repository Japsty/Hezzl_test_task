-- +goose Up
CREATE TABLE IF NOT EXISTS goods
(
    id         SERIAL ,
    project_id INT ,
    name       VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    priority INT NOT NULL ,
    removed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id, project_id)
);

CREATE INDEX id_and_project_id_index ON goods (id, project_id);

-- +goose Down
DROP TABLE IF EXISTS goods;