-- Миграция создающая таблицу goods в PostgreSQL
-- Для индексирования используется btree, поскольку для интовых значений он быстрее производит поиск в большом дипазоне
-- Для name используется hash индекс, поскольку для стринговых значений он быстрее

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

CREATE INDEX id_and_project_id_index ON goods USING btree (id, project_id);
CREATE INDEX name_index ON goods USING hash (name);


-- +goose Down
DROP TABLE IF EXISTS goods;