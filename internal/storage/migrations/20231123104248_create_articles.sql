-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    articles (
        id SERIAL PRIMARY KEY,
        source_id INT NOT NULL,
        title VARCHAR(255) NOT NULL,
        link VARCHAR(255) NOT NULL UNIQUE,
        summary TEXT NOT NULL,
        published_at TIMESTAMP NOT NULL, -- in original resource
        created_at TIMESTAMP NOT NULL DEFAULT NOW (), -- in our DB
        posted_at TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS articles;

-- +goose StatementEnd