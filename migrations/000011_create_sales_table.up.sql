CREATE TYPE sale_type    AS ENUM ('sale', 'own_use');
CREATE TYPE sale_channel AS ENUM ('langsung', 'pasar', 'pengepul', 'online');

CREATE TABLE sales (
    id              BIGSERIAL     PRIMARY KEY,
    cycle_id        BIGINT        NOT NULL,
    sale_date       DATE          NOT NULL,
    sale_type       sale_type     NOT NULL DEFAULT 'sale',
    channel         sale_channel,                 -- NULL untuk own_use
    qty             INT           NOT NULL,
    weight_kg       NUMERIC(8,2),                 -- info tambahan, opsional
    price_per_fish  BIGINT        NOT NULL DEFAULT 0,
    total           BIGINT        NOT NULL,       -- dihitung server: qty * price_per_fish
    buyer_name      VARCHAR(100),
    is_paid         BOOLEAN       NOT NULL DEFAULT TRUE,
    notes           TEXT,
    created_by      BIGINT        NOT NULL,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sales_cycle_id   FOREIGN KEY (cycle_id)   REFERENCES cycles(id) ON DELETE RESTRICT,
    CONSTRAINT fk_sales_created_by FOREIGN KEY (created_by) REFERENCES users(id)  ON DELETE RESTRICT,
    CONSTRAINT ck_sales_qty_positive       CHECK (qty > 0),
    CONSTRAINT ck_sales_weight_positive    CHECK (weight_kg IS NULL OR weight_kg > 0),
    CONSTRAINT ck_sales_money_non_negative CHECK (price_per_fish >= 0 AND total >= 0),
    CONSTRAINT ck_sales_total_matches      CHECK (total = qty * price_per_fish),
    -- R07: konsumsi sendiri selalu gratis & tanpa kanal
    CONSTRAINT ck_sales_own_use_free       CHECK (sale_type = 'sale' OR (price_per_fish = 0 AND channel IS NULL))
);

CREATE INDEX idx_sales_cycle_id_sale_date ON sales (cycle_id, sale_date);
CREATE INDEX idx_sales_unpaid ON sales (sale_date) WHERE is_paid = FALSE;

CREATE TRIGGER trg_sales_set_updated_at
    BEFORE UPDATE ON sales
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();