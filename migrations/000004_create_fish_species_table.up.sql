CREATE TABLE fish_species (
    id                     BIGSERIAL    PRIMARY KEY,
    name                   VARCHAR(100) NOT NULL,
    target_days            INT,
    default_price_per_fish BIGINT       NOT NULL DEFAULT 0,
    default_price_per_kg   BIGINT       NOT NULL DEFAULT 0,
    created_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ,

    CONSTRAINT ck_fish_species_name_not_blank      CHECK (length(trim(name)) > 0),
    CONSTRAINT ck_fish_species_target_days_positive CHECK (target_days IS NULL OR target_days > 0),
    CONSTRAINT ck_fish_species_prices_non_negative  CHECK (default_price_per_fish >= 0 AND default_price_per_kg >= 0)
);

CREATE UNIQUE INDEX uq_fish_species_name ON fish_species (LOWER(name)) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_fish_species_set_updated_at
    BEFORE UPDATE ON fish_species
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();