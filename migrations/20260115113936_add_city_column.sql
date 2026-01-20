-- +goose Up
-- +goose StatementBegin
CREATE TABLE cities
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255)   NOT NULL,
    lat        DECIMAL(10, 8) NOT NULL,
    lon        DECIMAL(11, 8) NOT NULL,
    country    VARCHAR(2),
    state      VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cities;
-- +goose StatementEnd
