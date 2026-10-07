# Konvensi Database

## Environment
- Lokal: PostgreSQL 16 via Docker, port **5433** (`docker compose up -d`)
- Production: PostgreSQL di Railway (`make migrate-up-prod`)

## Penamaan
| Objek | Aturan | Contoh |
|---|---|---|
| Tabel | snake_case, jamak | `fish_species`, `feed_logs` |
| Kolom | snake_case | `seed_price_per_fish` |
| Primary key | `id BIGSERIAL` | `id` |
| Foreign key | `<tabel tunggal>_id` | `cycle_id`, `pond_id` |
| Constraint FK | `fk_<tabel>_<kolom>` | `fk_feed_logs_cycle_id` |
| Unique | `uq_<tabel>_<kolom>` | `uq_users_email` |
| Check | `ck_<tabel>_<aturan>` | `ck_cycles_initial_qty_positive` |
| Index | `idx_<tabel>_<kolom>` | `idx_feed_logs_cycle_id_fed_at` |
| Enum type | snake_case, tunggal | `cycle_status`, `sales_channel` |
| View | `v_<nama>` | `v_cycle_status` |
| File migrasi | `NNNNNN_<kata kerja>_<objek>` | `000002_create_users_table` |
| Trigger | `trg_<tabel>_<aksi>` | `trg_users_set_updated_at` |

## Tipe data
- Waktu kejadian: `TIMESTAMPTZ`; tanggal log: `DATE`
- Uang: `BIGINT` (rupiah, tanpa desimal)
- Berat/volume: `NUMERIC`
- Boolean: awalan `is_` / `has_`
- Setiap tabel punya `created_at`, `updated_at` (trigger `set_updated_at()`)

## Aturan migrasi
- Satu migrasi = satu perubahan
- Selalu tulis file `.down.sql` dan uji `up → down → up`
- Migrasi yang sudah jalan di Railway tidak boleh diedit; buat migrasi baru

