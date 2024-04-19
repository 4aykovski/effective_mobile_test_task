-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars
(
    registration_number VARCHAR(9) PRIMARY KEY,
    mark VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    year INT,
    owner_name VARCHAR(255) NOT NULL,
    owner_surname VARCHAR(255) NOT NULL,
    FOREIGN KEY (owner_name, owner_surname) REFERENCES owner(name, surname)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cars;
-- +goose StatementEnd
