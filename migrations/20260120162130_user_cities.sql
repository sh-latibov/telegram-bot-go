-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_cities
(
    user_id    BIGINT  NOT NULL REFERENCES users (id),
    city_id    INTEGER NOT NULL REFERENCES cities (id),
    is_default BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id, city_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_cities;
-- +goose StatementEnd
