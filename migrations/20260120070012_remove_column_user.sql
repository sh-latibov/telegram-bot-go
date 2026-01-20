-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS city;
-- +goose StatementEnd

