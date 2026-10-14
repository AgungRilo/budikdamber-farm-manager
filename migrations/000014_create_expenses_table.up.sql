CREATE TYPE expense_category AS ENUM ('obat_probiotik', 'peralatan', 'utilitas', 'tenaga_kerja', 'lainnya');

CREATE TABLE expenses (
    id            BIGSERIAL         PRIMARY KEY,
    cycle_id      BIGINT,                          -- NULL = biaya umum (tidak terikat siklus)
    expense_date  DATE              NOT NULL,
    category      expense_category  NOT NULL,
    description   VARCHAR(200)      NOT NULL,
    amount        BIGINT            NOT NULL,
    created_by    BIGINT            NOT NULL,
    created_at    TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ       NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_expenses_cycle_id   FOREIGN KEY (cycle_id)   REFERENCES cycles(id) ON DELETE RESTRICT,
    CONSTRAINT fk_expenses_created_by FOREIGN KEY (created_by) REFERENCES users(id)  ON DELETE RESTRICT,
    CONSTRAINT ck_expenses_amount_positive      CHECK (amount > 0),
    CONSTRAINT ck_expenses_description_not_blank CHECK (length(trim(description)) > 0)
);

CREATE INDEX idx_expenses_cycle_id_expense_date ON expenses (cycle_id, expense_date);
CREATE INDEX idx_expenses_expense_date ON expenses (expense_date);

CREATE TRIGGER trg_expenses_set_updated_at
    BEFORE UPDATE ON expenses
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();