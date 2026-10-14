CREATE TABLE feed_logs (
    id            BIGSERIAL   PRIMARY KEY,
    cycle_id      BIGINT      NOT NULL,
    feed_id       BIGINT      NOT NULL,
    log_date      DATE        NOT NULL,
    qty_gram      INT         NOT NULL,
    price_per_kg  BIGINT      NOT NULL,  -- snapshot harga pakan saat dicatat
    cost          BIGINT      NOT NULL,  -- qty_gram / 1000 * price_per_kg
    notes         TEXT,
    created_by    BIGINT      NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_feed_logs_cycle_id   FOREIGN KEY (cycle_id)   REFERENCES cycles(id) ON DELETE RESTRICT,
    CONSTRAINT fk_feed_logs_feed_id    FOREIGN KEY (feed_id)    REFERENCES feeds(id)  ON DELETE RESTRICT,
    CONSTRAINT fk_feed_logs_created_by FOREIGN KEY (created_by) REFERENCES users(id)  ON DELETE RESTRICT,
    CONSTRAINT ck_feed_logs_qty_positive    CHECK (qty_gram > 0),
    CONSTRAINT ck_feed_logs_money_non_negative CHECK (price_per_kg >= 0 AND cost >= 0)
);

CREATE INDEX idx_feed_logs_cycle_id_log_date ON feed_logs (cycle_id, log_date);
CREATE INDEX idx_feed_logs_feed_id ON feed_logs (feed_id);