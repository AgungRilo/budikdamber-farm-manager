-- kembali ke versi 000007
CREATE OR REPLACE VIEW v_cycle_status AS
SELECT
    c.id                  AS cycle_id,
    c.initial_qty::bigint AS initial_qty,
    0::bigint             AS dead_qty,
    0::bigint             AS sold_qty,
    c.initial_qty::bigint AS alive_qty
FROM cycles c;