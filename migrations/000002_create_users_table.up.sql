CREATE TYPE user_role AS ENUM ('owner', 'worker');

CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL,
    password_hash TEXT         NOT NULL,
    role          user_role    NOT NULL DEFAULT 'worker',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT ck_users_name_not_blank CHECK (length(trim(name)) > 0)
);

-- Email unik tanpa membedakan huruf besar/kecil (Budi@x.com = budi@x.com)
CREATE UNIQUE INDEX uq_users_email ON users (LOWER(email));

CREATE TRIGGER trg_users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();