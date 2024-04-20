-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS owners
(
    name       VARCHAR(255) NOT NULL,
    surname    VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    PRIMARY KEY (name, surname)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS owner;
-- +goose StatementEnd
