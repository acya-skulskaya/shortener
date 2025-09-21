CREATE TABLE short_urls (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(255) NOT NULL
);

CREATE INDEX idx_id ON short_urls(id);
CREATE UNIQUE INDEX idx_original_url ON short_urls(original_url);