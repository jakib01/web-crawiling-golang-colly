-- 1. Categories (for breadcrumb hierarchy)
CREATE TABLE categories (
  id            SERIAL PRIMARY KEY,
  name          VARCHAR(255) NOT NULL,
  parent_id     INT REFERENCES categories(id)      -- self-join for hierarchy
);

-- 2. Products
CREATE TABLE products (
  id                  SERIAL PRIMARY KEY,
  product_number      VARCHAR(50)  NOT NULL UNIQUE,
  name                VARCHAR(500) NOT NULL,
  category_id         INT          NOT NULL REFERENCES categories(id),
  price_yen           NUMERIC(10,2) NOT NULL,      -- e.g. 7990.00
  sense_of_size       VARCHAR(100),                -- e.g. “Loose fit”
  details_url         TEXT         NOT NULL,
  total_reviews       INT          DEFAULT 0,      -- aggregate count
  recommended_rate    NUMERIC(5,2) DEFAULT 0.0     -- e.g. 87.50 (%)
);

-- index to quickly find products by category
CREATE INDEX idx_products_category ON products(category_id);


-- 3. Images
CREATE TABLE product_images (
  id            SERIAL PRIMARY KEY,
  product_id    INT     NOT NULL REFERENCES products(id),
  url           TEXT    NOT NULL,
  is_main       BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_images_product ON product_images(product_id);


-- 4. Detailed Size Variants
CREATE TABLE product_sizes (
  id                 SERIAL PRIMARY KEY,
  product_id         INT     NOT NULL REFERENCES products(id),
  size_label         VARCHAR(20) NOT NULL,         -- e.g. S, M, L
  chest_cm           NUMERIC(5,2),                 -- e.g.  92.00
  back_length_cm     NUMERIC(5,2),                 -- e.g.  70.00
  other_measurements JSONB,                        -- free-form for extra dims
  special_functions  JSONB                         -- e.g. [{"name":"moisture-wicking","desc":"…"}]
);
CREATE INDEX idx_sizes_product ON product_sizes(product_id);


-- 5. Coordinated (Suggested) Products
-- We capture them even if you don't crawl them later as full products
CREATE TABLE coordinated_items (
  id                   SERIAL PRIMARY KEY,
  source_product_id    INT     NOT NULL REFERENCES products(id),
  coord_product_number VARCHAR(50)  NOT NULL,
  name                 VARCHAR(500) NOT NULL,
  price_yen            NUMERIC(10,2) NOT NULL,
  image_url            TEXT         NOT NULL,
  page_url             TEXT         NOT NULL
);
CREATE INDEX idx_coordinated_source ON coordinated_items(source_product_id);


-- 6. Keywords / Tags
CREATE TABLE keywords (
  id    SERIAL PRIMARY KEY,
  kw    VARCHAR(100) UNIQUE NOT NULL
);
CREATE TABLE product_keywords (
  product_id   INT NOT NULL REFERENCES products(id),
  keyword_id   INT NOT NULL REFERENCES keywords(id),
  PRIMARY KEY (product_id, keyword_id)
);
CREATE INDEX idx_prod_kw_product ON product_keywords(product_id);


-- 7. Reviews
CREATE TABLE reviews (
  id               SERIAL PRIMARY KEY,
  product_id       INT       NOT NULL REFERENCES products(id),
  reviewer_id      VARCHAR(100) NOT NULL,
  review_date      DATE      NOT NULL,
  overall_rating   NUMERIC(3,2) NOT NULL,        -- e.g. 4.50 out of 5
  title            VARCHAR(255),
  body             TEXT
);
CREATE INDEX idx_reviews_product ON reviews(product_id);


-- 8. Aspect Ratings per Review
CREATE TABLE review_aspect_ratings (
  id         SERIAL PRIMARY KEY,
  review_id  INT    NOT NULL REFERENCES reviews(id),
  aspect     VARCHAR(100) NOT NULL,              -- e.g. “Comfort”
  rating     NUMERIC(3,2) NOT NULL
);
CREATE INDEX idx_aspect_review ON review_aspect_ratings(review_id);
-- 1. Categories (for breadcrumb hierarchy)
CREATE TABLE categories (
  id            SERIAL PRIMARY KEY,
  name          VARCHAR(255) NOT NULL,
  parent_id     INT REFERENCES categories(id)      -- self-join for hierarchy
);

-- 2. Products
CREATE TABLE products (
  id                  SERIAL PRIMARY KEY,
  product_number      VARCHAR(50)  NOT NULL UNIQUE,
  name                VARCHAR(500) NOT NULL,
  category_id         INT          NOT NULL REFERENCES categories(id),
  price_yen           NUMERIC(10,2) NOT NULL,      -- e.g. 7990.00
  sense_of_size       VARCHAR(100),                -- e.g. “Loose fit”
  details_url         TEXT         NOT NULL,
  total_reviews       INT          DEFAULT 0,      -- aggregate count
  recommended_rate    NUMERIC(5,2) DEFAULT 0.0     -- e.g. 87.50 (%)
);

-- index to quickly find products by category
CREATE INDEX idx_products_category ON products(category_id);


-- 3. Images
CREATE TABLE product_images (
  id            SERIAL PRIMARY KEY,
  product_id    INT     NOT NULL REFERENCES products(id),
  url           TEXT    NOT NULL,
  is_main       BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_images_product ON product_images(product_id);


-- 4. Detailed Size Variants
CREATE TABLE product_sizes (
  id                 SERIAL PRIMARY KEY,
  product_id         INT     NOT NULL REFERENCES products(id),
  size_label         VARCHAR(20) NOT NULL,         -- e.g. S, M, L
  chest_cm           NUMERIC(5,2),                 -- e.g.  92.00
  back_length_cm     NUMERIC(5,2),                 -- e.g.  70.00
  other_measurements JSONB,                        -- free-form for extra dims
  special_functions  JSONB                         -- e.g. [{"name":"moisture-wicking","desc":"…"}]
);
CREATE INDEX idx_sizes_product ON product_sizes(product_id);


-- 5. Coordinated (Suggested) Products
-- We capture them even if you don't crawl them later as full products
CREATE TABLE coordinated_items (
  id                   SERIAL PRIMARY KEY,
  source_product_id    INT     NOT NULL REFERENCES products(id),
  coord_product_number VARCHAR(50)  NOT NULL,
  name                 VARCHAR(500) NOT NULL,
  price_yen            NUMERIC(10,2) NOT NULL,
  image_url            TEXT         NOT NULL,
  page_url             TEXT         NOT NULL
);
CREATE INDEX idx_coordinated_source ON coordinated_items(source_product_id);


-- 6. Keywords / Tags
CREATE TABLE keywords (
  id    SERIAL PRIMARY KEY,
  kw    VARCHAR(100) UNIQUE NOT NULL
);
CREATE TABLE product_keywords (
  product_id   INT NOT NULL REFERENCES products(id),
  keyword_id   INT NOT NULL REFERENCES keywords(id),
  PRIMARY KEY (product_id, keyword_id)
);
CREATE INDEX idx_prod_kw_product ON product_keywords(product_id);


-- 7. Reviews
CREATE TABLE reviews (
  id               SERIAL PRIMARY KEY,
  product_id       INT       NOT NULL REFERENCES products(id),
  reviewer_id      VARCHAR(100) NOT NULL,
  review_date      DATE      NOT NULL,
  overall_rating   NUMERIC(3,2) NOT NULL,        -- e.g. 4.50 out of 5
  title            VARCHAR(255),
  body             TEXT
);
CREATE INDEX idx_reviews_product ON reviews(product_id);


-- 8. Aspect Ratings per Review
CREATE TABLE review_aspect_ratings (
  id         SERIAL PRIMARY KEY,
  review_id  INT    NOT NULL REFERENCES reviews(id),
  aspect     VARCHAR(100) NOT NULL,              -- e.g. “Comfort”
  rating     NUMERIC(3,2) NOT NULL
);
CREATE INDEX idx_aspect_review ON review_aspect_ratings(review_id);
