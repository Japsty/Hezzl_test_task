-- +goose up
CREATE TABLE IF NOT EXISTS logs(
    Id INT ,
    ProjectId INT,
    Name VARCHAR(255),
    Description VARCHAR(255),
    Priority INT,
    Removed BOOLEAN,
    EventTime TIMESTAMP
)
-- +goose down