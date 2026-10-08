CREATE TABLE feeds (
    id           BIGSERIAL    PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    brand        VARCHAR(100),
    protein_pct  NUMERIC(4,1),
    price_per_kg BIGINT       NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT ck_feeds_name_not_blank      CHECK (length(trim(name)) > 0),
    CONSTRAINT ck_feeds_protein_range       CHECK (protein_pct IS NULL OR protein_pct BETWEEN 0 AND 100),
    CONSTRAINT ck_feeds_price_non_negative  CHECK (price_per_kg >= 0)
);

-- Kombinasi nama + merek unik (brand kosong dianggap '')
CREATE UNIQUE INDEX uq_feeds_name_brand
    ON feeds (LOWER(name), LOWER(COALESCE(brand, '')))
    WHERE deleted_at IS NULL;

CREATE TRIGGER trg_feeds_set_updated_at
    BEFORE UPDATE ON feeds
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();