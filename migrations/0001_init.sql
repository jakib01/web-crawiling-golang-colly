-- 1. Products
CREATE TABLE products
(
    id                           SERIAL PRIMARY KEY,
    product_code                 VARCHAR(50)    NOT NULL UNIQUE,
    name                         VARCHAR(500)   NOT NULL,
    category                     VARCHAR(500)   NOT NULL,
    price_yen                    NUMERIC(10, 2) NOT NULL,
    sense_of_size                VARCHAR(100),
    details_url                  TEXT           NOT NULL,
    total_reviews                INT DEFAULT 0,
    overall_rating               NUMERIC(3, 2)  NOT NULL,
    title_description            TEXT           NOT NULL,
    general_description          TEXT           NOT NULL,
    item_general_description     TEXT           NOT NULL,
    special_function_description TEXT           NOT NULL
);

-- 3. Images
CREATE TABLE product_images
(
    id         SERIAL PRIMARY KEY,
    product_id INT     REFERENCES products (id),
    url        TEXT    NOT NULL,
    is_main    BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_images_product ON product_images (product_id);

-- 4. Detailed Size Variants
CREATE TABLE product_sizes
(
    id                 SERIAL PRIMARY KEY,
    product_id         INT         REFERENCES products (id),
    size_label         VARCHAR(20) NOT NULL,
    chest_cm           NUMERIC(5, 2),
    availability       NUMERIC(5, 2),
    back_length_cm     NUMERIC(5, 2),
    other_measurements TEXT NULL,
    special_functions  TEXT NULL
);
CREATE INDEX idx_sizes_product ON product_sizes (product_id);

-- 5. Coordinated Products
CREATE TABLE coordinated_items
(
    id                SERIAL PRIMARY KEY,
    source_product_id INT            REFERENCES products (id),
    product_number    VARCHAR(50)    NOT NULL,
    name              VARCHAR(500)   NOT NULL,
    price_yen         NUMERIC(10, 2) NOT NULL,
    image_url         TEXT           NOT NULL,
    product_page_url  TEXT           NOT NULL
);
CREATE INDEX idx_coordinated_source ON coordinated_items (source_product_id);

-- 6. Keywords / Tags
CREATE TABLE keywords
(
    id SERIAL PRIMARY KEY,
    kw VARCHAR(100) UNIQUE NOT NULL
);
CREATE TABLE product_keywords
(
    product_id INT REFERENCES products (id),
    keyword_id INT NOT NULL REFERENCES keywords (id),
    PRIMARY KEY (product_id, keyword_id)
);
CREATE INDEX idx_prod_kw_product ON product_keywords (product_id);

-- 7. Reviews
CREATE TABLE reviews
(
    id             SERIAL PRIMARY KEY,
    product_id     INT           REFERENCES products (id),
    review_date    DATE          NOT NULL,
    rating         NUMERIC(3, 2) NOT NULL,
    overall_rating NUMERIC(3, 2) NOT NULL,
    title          VARCHAR(255),
    body           TEXT
);
CREATE INDEX idx_reviews_product ON reviews (product_id);

-- 8. Aspect Ratings per Review
CREATE TABLE review_aspect_ratings
(
    id        SERIAL PRIMARY KEY,
    review_id INT           NOT NULL REFERENCES reviews (id),
    aspect    VARCHAR(100)  NOT NULL,
    rating    NUMERIC(3, 2) NOT NULL
);
CREATE INDEX idx_aspect_review ON review_aspect_ratings (review_id);
