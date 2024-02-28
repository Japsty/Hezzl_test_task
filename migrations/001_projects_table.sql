-- +goose Up
CREATE TABLE IF NOT EXISTS projects
(
    id         INT PRIMARY KEY ,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO projects(id, name)
VALUES (1,'Запись 1')

-- +goose Down