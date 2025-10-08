CREATE TABLE short_urls (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    is_deleted BOOLEAN DEFAULT false
);

CREATE INDEX idx_id ON short_urls(id);
CREATE INDEX idx_user_id ON short_urls(user_id);
CREATE UNIQUE INDEX idx_original_url ON short_urls(original_url);