CREATE TABLE product_urls
(
    id         SERIAL PRIMARY KEY,
    code       TEXT NOT NULL,
    url        TEXT NOT NULL UNIQUE,
    image_url  TEXT,
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
