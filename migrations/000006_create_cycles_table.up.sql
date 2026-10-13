CREATE TYPE cycle_status AS ENUM ('active', 'closed');

CREATE TABLE cycles (
    id                  BIGSERIAL    PRIMARY KEY,
    code                VARCHAR(30)  NOT NULL,
    pond_id             BIGINT       NOT NULL,
    species_id          BIGINT       NOT NULL,
    start_date          DATE         NOT NULL,
    end_date            DATE,
    initial_qty         INT          NOT NULL,
    seed_size_cm        NUMERIC(4,1),
    seed_price_per_fish BIGINT       NOT NULL DEFAULT 0,
    status              cycle_status NOT NULL DEFAULT 'active',
    notes               TEXT,
    created_by          BIGINT       NOT NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_cycles_code                    UNIQUE (code),
    CONSTRAINT fk_cycles_pond_id                 FOREIGN KEY (pond_id)    REFERENCES ponds(id)        ON DELETE RESTRICT,
    CONSTRAINT fk_cycles_species_id              FOREIGN KEY (species_id) REFERENCES fish_species(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cycles_created_by              FOREIGN KEY (created_by) REFERENCES users(id)        ON DELETE RESTRICT,
    CONSTRAINT ck_cycles_initial_qty_positive    CHECK (initial_qty > 0),
    CONSTRAINT ck_cycles_seed_price_non_negative CHECK (seed_price_per_fish >= 0),
    CONSTRAINT ck_cycles_seed_size_positive      CHECK (seed_size_cm IS NULL OR seed_size_cm > 0),
    CONSTRAINT ck_cycles_end_after_start         CHECK (end_date IS NULL OR end_date >= start_date),
    CONSTRAINT ck_cycles_closed_has_end          CHECK (status = 'active' OR end_date IS NOT NULL)
);

CREATE INDEX idx_cycles_pond_id_status ON cycles (pond_id, status);
CREATE INDEX idx_cycles_species_id     ON cycles (species_id);

CREATE TRIGGER trg_cycles_set_updated_at
    BEFORE UPDATE ON cycles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();