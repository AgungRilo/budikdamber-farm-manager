CREATE TABLE mortality_logs (
    id          BIGSERIAL    PRIMARY KEY,
    cycle_id    BIGINT       NOT NULL,
    log_date    DATE         NOT NULL,
    qty         INT          NOT NULL,
    cause       VARCHAR(100),           -- contoh: "kualitas air", "penyakit", "tidak diketahui"
    notes       TEXT,
    created_by  BIGINT       NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_mortality_logs_cycle_id   FOREIGN KEY (cycle_id)   REFERENCES cycles(id) ON DELETE RESTRICT,
    CONSTRAINT fk_mortality_logs_created_by FOREIGN KEY (created_by) REFERENCES users(id)  ON DELETE RESTRICT,
    CONSTRAINT ck_mortality_logs_qty_positive CHECK (qty > 0)
);

CREATE INDEX idx_mortality_logs_cycle_id_log_date ON mortality_logs (cycle_id, log_date);