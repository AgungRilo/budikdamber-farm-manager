ALTER TABLE cycles ADD COLUMN closed_by BIGINT;
ALTER TABLE cycles ADD CONSTRAINT fk_cycles_closed_by
    FOREIGN KEY (closed_by) REFERENCES users(id) ON DELETE RESTRICT;