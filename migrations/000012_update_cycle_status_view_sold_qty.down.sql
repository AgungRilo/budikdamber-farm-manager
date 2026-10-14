-- kembali ke versi 000010 (dead_qty saja)
CREATE OR REPLACE VIEW v_cycle_status AS
SELECT
    c.id                                        AS cycle_id,
    c.initial_qty::bigint                       AS initial_qty,
    COALESCE(m.dead_qty, 0)::bigint             AS dead_qty,
    0::bigint                                   AS sold_qty,
    (c.initial_qty - COALESCE(m.dead_qty, 0))::bigint AS alive_qty
FROM cycles c
LEFT JOIN (SELECT cycle_id, SUM(qty) AS dead_qty FROM mortality_logs GROUP BY cycle_id) m ON m.cycle_id = c.id;