-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS versions (
    id SERIAL PRIMARY KEY,
    store_id INT NOT NULL,
    version_number INT NOT NULL,
    creator VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    open_at VARCHAR(255) NOT NULL,
    close_at VARCHAR(255) NOT NULL,
    created_at VARCHAR(255) NOT NULL,
    is_deleted boolean default false,
    FOREIGN KEY (store_id) REFERENCES stores (id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS versions;
-- +goose StatementEnd
