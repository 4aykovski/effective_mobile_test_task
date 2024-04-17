-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars
(
    registration_number VARCHAR(9) PRIMARY KEY,
    mark VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    year INT,
    owner_id INT REFERENCES owner(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cars;
-- +goose StatementEnd
