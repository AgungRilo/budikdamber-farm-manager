-- Data demo Budikdamber Farm Manager.
-- Dijalankan oleh cmd/seed-demo. Aman diulang: dilewati jika data [demo] sudah ada.
DO $$
DECLARE
    d        date := CURRENT_DATE;
    v_owner  bigint;
    v_worker bigint;
    p1 bigint; p2 bigint; p3 bigint;
    s_lele bigint; s_nila bigint;
    f1 bigint; f2 bigint;
    c1 bigint; c2 bigint; c3 bigint;
BEGIN
    IF EXISTS (SELECT 1 FROM cycles WHERE notes LIKE '[demo]%') THEN
        RAISE NOTICE 'data demo sudah ada, dilewati';
        RETURN;
    END IF;

    SELECT id INTO v_owner FROM users WHERE role = 'owner' AND is_active ORDER BY id LIMIT 1;
    IF v_owner IS NULL THEN
        RAISE EXCEPTION 'akun owner belum ada, jalankan make seed terlebih dahulu';
    END IF;
    SELECT id INTO v_worker FROM users WHERE LOWER(email) = 'pekerja@budikdamber.local';
    v_worker := COALESCE(v_worker, v_owner);

    -- ===== Data master =====
    INSERT INTO ponds (name, type, volume_liter, max_stock, is_active, notes) VALUES
        ('Ember Demo 1', 'ember',   80,   50,  TRUE, '[demo] ember 80 liter'),
        ('Ember Demo 2', 'ember',   80,   50,  TRUE, '[demo] ember 80 liter'),
        ('Terpal Demo',  'terpal',  1000, 300, TRUE, '[demo] kolam terpal bulat')
    ON CONFLICT DO NOTHING;
    SELECT id INTO p1 FROM ponds WHERE LOWER(name) = 'ember demo 1' AND deleted_at IS NULL;
    SELECT id INTO p2 FROM ponds WHERE LOWER(name) = 'ember demo 2' AND deleted_at IS NULL;
    SELECT id INTO p3 FROM ponds WHERE LOWER(name) = 'terpal demo'  AND deleted_at IS NULL;

    INSERT INTO fish_species (name, target_days, default_price_per_fish, default_price_per_kg) VALUES
        ('Lele',   90,  5000,  22000),
        ('Nila',   120, 6000,  30000),
        ('Gurame', 180, 15000, 45000)
    ON CONFLICT DO NOTHING;
    SELECT id INTO s_lele FROM fish_species WHERE LOWER(name) = 'lele' AND deleted_at IS NULL;
    SELECT id INTO s_nila FROM fish_species WHERE LOWER(name) = 'nila' AND deleted_at IS NULL;

    INSERT INTO feeds (name, brand, protein_pct, price_per_kg) VALUES
        ('Pelet Apung 2mm', 'Demo Feed', 31, 15000),
        ('Pelet Apung 3mm', 'Demo Feed', 30, 14000)
    ON CONFLICT DO NOTHING;
    SELECT id INTO f1 FROM feeds
     WHERE LOWER(name) = 'pelet apung 2mm' AND LOWER(COALESCE(brand, '')) = 'demo feed' AND deleted_at IS NULL;
    SELECT id INTO f2 FROM feeds
     WHERE LOWER(name) = 'pelet apung 3mm' AND LOWER(COALESCE(brand, '')) = 'demo feed' AND deleted_at IS NULL;

    -- ===== Siklus 1: SELESAI (panen habis) =====
    c1 := nextval('cycles_id_seq');
    INSERT INTO cycles (id, code, pond_id, species_id, start_date, initial_qty, seed_size_cm, seed_price_per_fish, notes, created_by)
    VALUES (c1, 'SK-' || to_char(d - 120, 'YYMMDD') || '-' || c1, p1, s_lele, d - 120, 45, 8, 500,
            '[demo] siklus selesai, panen habis', v_owner);

    INSERT INTO feed_logs (cycle_id, feed_id, log_date, qty_gram, price_per_kg, cost, created_by)
    SELECT c1, f.id, g.ts::date, x.q, f.price_per_kg, ROUND(x.q * f.price_per_kg / 1000.0)::bigint, v_worker
    FROM generate_series((d - 120)::timestamp, (d - 31)::timestamp, interval '1 day') AS g(ts)
    CROSS JOIN LATERAL (SELECT (30 + (g.ts::date - (d - 120)))::int AS q) x
    JOIN feeds f ON f.id = f1;

    INSERT INTO mortality_logs (cycle_id, log_date, qty, cause, created_by) VALUES
        (c1, d - 110, 2, 'adaptasi bibit', v_worker),
        (c1, d - 80,  2, 'kualitas air',   v_worker);

    -- 45 - 4 mati = 41 = 3 konsumsi + 15 + 10 + 13
    INSERT INTO sales (cycle_id, sale_date, sale_type, channel, qty, price_per_fish, total, buyer_name, is_paid, created_by) VALUES
        (c1, d - 31, 'own_use', NULL,       3,  0,    0,     NULL,                 TRUE, v_owner),
        (c1, d - 31, 'sale',    'langsung', 15, 5000, 75000, 'Tetangga RT 03',     TRUE, v_owner),
        (c1, d - 30, 'sale',    'pasar',    10, 5000, 50000, 'Bu Siti',            TRUE, v_owner),
        (c1, d - 30, 'sale',    'pengepul', 13, 4500, 58500, 'Pengepul Pak Darto', TRUE, v_owner);

    INSERT INTO expenses (cycle_id, expense_date, category, description, amount, created_by) VALUES
        (c1, d - 90, 'obat_probiotik', 'EM4 dan garam ikan', 25000, v_owner);

    UPDATE cycles SET status = 'closed', end_date = d - 30, closed_by = v_owner WHERE id = c1;

    -- ===== Siklus 2: AKTIF, hampir panen =====
    c2 := nextval('cycles_id_seq');
    INSERT INTO cycles (id, code, pond_id, species_id, start_date, initial_qty, seed_size_cm, seed_price_per_fish, notes, created_by)
    VALUES (c2, 'SK-' || to_char(d - 75, 'YYMMDD') || '-' || c2, p3, s_lele, d - 75, 250, 7, 450,
            '[demo] siklus berjalan, hampir panen', v_owner);

    INSERT INTO feed_logs (cycle_id, feed_id, log_date, qty_gram, price_per_kg, cost, created_by)
    SELECT c2, f.id, g.ts::date, x.q, f.price_per_kg, ROUND(x.q * f.price_per_kg / 1000.0)::bigint, v_worker
    FROM generate_series((d - 75)::timestamp, (d - 1)::timestamp, interval '1 day') AS g(ts)
    CROSS JOIN LATERAL (SELECT (200 + (g.ts::date - (d - 75)) * 10)::int AS q) x
    JOIN feeds f ON f.id = f2;

    INSERT INTO mortality_logs (cycle_id, log_date, qty, cause, created_by) VALUES
        (c2, d - 70, 8, 'adaptasi bibit',  v_worker),
        (c2, d - 45, 5, 'penyakit kulit',  v_worker),
        (c2, d - 15, 3, 'tidak diketahui', v_worker);

    INSERT INTO sales (cycle_id, sale_date, sale_type, channel, qty, price_per_fish, total, buyer_name, is_paid, created_by) VALUES
        (c2, d - 5, 'sale',    'online',   20, 6000, 120000, 'Pak Budi',      FALSE, v_owner),
        (c2, d - 3, 'sale',    'langsung', 30, 5000, 150000, 'Warung Bu Ani', TRUE,  v_owner),
        (c2, d - 2, 'own_use', NULL,       4,  0,    0,      NULL,            TRUE,  v_owner);

    INSERT INTO expenses (cycle_id, expense_date, category, description, amount, created_by) VALUES
        (c2, d - 60, 'obat_probiotik', 'Probiotik air kolam', 40000, v_owner),
        (c2, d - 40, 'utilitas',       'Listrik pompa air',   30000, v_owner);

    -- ===== Siklus 3: AKTIF, baru mulai =====
    c3 := nextval('cycles_id_seq');
    INSERT INTO cycles (id, code, pond_id, species_id, start_date, initial_qty, seed_size_cm, seed_price_per_fish, notes, created_by)
    VALUES (c3, 'SK-' || to_char(d - 20, 'YYMMDD') || '-' || c3, p2, s_nila, d - 20, 40, 5, 700,
            '[demo] siklus baru', v_owner);

    INSERT INTO feed_logs (cycle_id, feed_id, log_date, qty_gram, price_per_kg, cost, created_by)
    SELECT c3, f.id, g.ts::date, x.q, f.price_per_kg, ROUND(x.q * f.price_per_kg / 1000.0)::bigint, v_worker
    FROM generate_series((d - 20)::timestamp, (d - 1)::timestamp, interval '1 day') AS g(ts)
    CROSS JOIN LATERAL (SELECT (40 + (g.ts::date - (d - 20)) * 4)::int AS q) x
    JOIN feeds f ON f.id = f1;

    INSERT INTO mortality_logs (cycle_id, log_date, qty, cause, created_by) VALUES
        (c3, d - 12, 1, 'adaptasi bibit', v_worker);

    -- ===== Biaya umum (tidak terikat siklus) =====
    INSERT INTO expenses (cycle_id, expense_date, category, description, amount, created_by) VALUES
        (NULL, d - 60, 'peralatan', 'Aerator kecil untuk ember', 85000, v_owner);

    RAISE NOTICE 'data demo berhasil dibuat: siklus % (selesai), % dan % (aktif)', c1, c2, c3;
END $$;