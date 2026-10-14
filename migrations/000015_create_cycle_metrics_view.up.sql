CREATE VIEW v_cycle_metrics AS
SELECT
    b.*,
    (b.seed_cost + b.feed_cost + b.other_cost)                  AS total_cost,
    (b.revenue - (b.seed_cost + b.feed_cost + b.other_cost))    AS profit,
    ROUND((b.initial_qty - b.dead_qty) * 100.0 / b.initial_qty, 1)::float8 AS survival_rate_pct,
    CASE WHEN b.initial_qty - b.dead_qty > 0
         THEN ROUND((b.seed_cost + b.feed_cost + b.other_cost)::numeric / (b.initial_qty - b.dead_qty))::bigint
    END                                                         AS hpp_per_fish,
    CASE WHEN b.revenue > 0
         THEN ROUND((b.revenue - (b.seed_cost + b.feed_cost + b.other_cost)) * 100.0 / b.revenue, 1)::float8
    END                                                         AS margin_pct
FROM (
    SELECT
        c.id                                              AS cycle_id,
        c.initial_qty::bigint                             AS initial_qty,
        st.dead_qty,
        st.alive_qty,
        COALESCE(sl.sale_qty, 0)::bigint                  AS sale_qty,
        COALESCE(sl.own_use_qty, 0)::bigint               AS own_use_qty,
        (c.initial_qty::bigint * c.seed_price_per_fish)   AS seed_cost,
        COALESCE(fl.feed_cost, 0)::bigint                 AS feed_cost,
        COALESCE(fl.feed_gram, 0)::bigint                 AS feed_gram,
        COALESCE(ex.other_cost, 0)::bigint                AS other_cost,
        COALESCE(sl.revenue, 0)::bigint                   AS revenue,
        COALESCE(sl.receivable, 0)::bigint                AS receivable
    FROM cycles c
    JOIN v_cycle_status st ON st.cycle_id = c.id
    LEFT JOIN (
        SELECT cycle_id, SUM(cost) AS feed_cost, SUM(qty_gram) AS feed_gram
        FROM feed_logs GROUP BY cycle_id
    ) fl ON fl.cycle_id = c.id
    LEFT JOIN (
        SELECT cycle_id, SUM(amount) AS other_cost
        FROM expenses WHERE cycle_id IS NOT NULL GROUP BY cycle_id
    ) ex ON ex.cycle_id = c.id
    LEFT JOIN (
        SELECT cycle_id,
               SUM(qty)   FILTER (WHERE sale_type = 'sale')    AS sale_qty,
               SUM(qty)   FILTER (WHERE sale_type = 'own_use') AS own_use_qty,
               SUM(total)                                      AS revenue,
               SUM(total) FILTER (WHERE NOT is_paid)           AS receivable
        FROM sales GROUP BY cycle_id
    ) sl ON sl.cycle_id = c.id
) b;