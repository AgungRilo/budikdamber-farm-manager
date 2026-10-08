CREATE TYPE pond_type AS ENUM ('ember', 'terpal', 'kolam_tanah', 'kolam_beton');

CREATE TABLE ponds (
    id           BIGSERIAL    PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    type         pond_type    NOT NULL,
    volume_liter INT,
    max_stock    INT          NOT NULL,
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    notes        TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT ck_ponds_name_not_blank      CHECK (length(trim(name)) > 0),
    CONSTRAINT ck_ponds_max_stock_positive  CHECK (max_stock > 0),
    CONSTRAINT ck_ponds_volume_positive     CHECK (volume_liter IS NULL OR volume_liter > 0)
);

-- Nama unik (tanpa membedakan huruf besar/kecil), hanya untuk data yang belum dihapus
CREATE UNIQUE INDEX uq_ponds_name ON ponds (LOWER(name)) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_ponds_set_updated_at
    BEFORE UPDATE ON ponds
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();