-- Миграция создающая таблицу projects в PostgreSQL
-- Для индексирования используется btree поскольку для интовых значений он быстрее производит поиск в большом дипазоне

-- +goose Up
CREATE TABLE IF NOT EXISTS projects
(
    id         SERIAL PRIMARY KEY ,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX id_index ON projects USING btree (id);

INSERT INTO projects(name)
VALUES ('Запись 1');

-- +goose Down
DROP TABLE IF EXISTS projects;