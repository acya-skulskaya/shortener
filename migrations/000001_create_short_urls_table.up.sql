CREATE TABLE short_urls (
    id VARCHAR(10) NOT NULL PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(255) NOT NULL
);

CREATE INDEX idx_id ON short_urls(id);