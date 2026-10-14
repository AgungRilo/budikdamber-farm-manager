package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) Stock(ctx context.Context, out *model.DashboardStock) error {
	err := r.db.QueryRow(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM cycles WHERE status = 'active'),
		  (SELECT COALESCE(SUM(s.alive_qty), 0)::bigint
		     FROM cycles c JOIN v_cycle_status s ON s.cycle_id = c.id WHERE c.status = 'active'),
		  (SELECT COUNT(*) FROM ponds WHERE deleted_at IS NULL AND is_active),
		  (SELECT COALESCE(SUM(max_stock), 0)::bigint FROM ponds WHERE deleted_at IS NULL AND is_active)`).
		Scan(&out.ActiveCycles, &out.AliveFish, &out.ActivePonds, &out.TotalCapacity)
	if err != nil {
		return err
	}
	if out.TotalCapacity > 0 {
		out.CapacityUsedPct = float64(out.AliveFish*1000/out.TotalCapacity) / 10 // 1 angka desimal
	}
	return nil
}

func (r *DashboardRepository) Finance(ctx context.Context, out *model.DashboardFinance) error {
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_cost), 0)::bigint,
		       COALESCE(SUM(revenue), 0)::bigint,
		       COALESCE(SUM(receivable), 0)::bigint,
		       COALESCE(SUM(profit), 0)::bigint,
		       (SELECT COALESCE(SUM(amount), 0)::bigint FROM expenses WHERE cycle_id IS NULL)
		FROM v_cycle_metrics`).
		Scan(&out.TotalCost, &out.Revenue, &out.Receivable, &out.CycleProfit, &out.GeneralExpenses)
	if err != nil {
		return err
	}
	out.NetProfit = out.CycleProfit - out.GeneralExpenses
	return nil
}

func (r *DashboardRepository) SalesByChannel(ctx context.Context, f model.DashboardFilter, out *[]model.ChannelSales) error {
	rows, err := r.db.Query(ctx, `
		SELECT s.channel::text, COUNT(*), SUM(s.qty)::bigint, SUM(s.total)::bigint,
		       SUM(s.total - s.qty * COALESCE(vm.hpp_per_fish, 0))::bigint
		FROM sales s
		JOIN v_cycle_metrics vm ON vm.cycle_id = s.cycle_id
		WHERE s.sale_type = 'sale'
		  AND ($1::date IS NULL OR s.sale_date >= $1::date)
		  AND ($2::date IS NULL OR s.sale_date <= $2::date)
		GROUP BY s.channel
		ORDER BY 4 DESC`, nilIfEmpty(f.DateFrom), nilIfEmpty(f.DateTo))
	if err != nil {
		return err
	}
	defer rows.Close()

	list := make([]model.ChannelSales, 0)
	for rows.Next() {
		var c model.ChannelSales
		if err := rows.Scan(&c.Channel, &c.Transactions, &c.Qty, &c.Revenue, &c.EstGrossMargin); err != nil {
			return err
		}
		list = append(list, c)
	}
	*out = list
	return rows.Err()
}

func (r *DashboardRepository) Receivables(ctx context.Context, f model.DashboardFilter, out *[]model.Receivable) error {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, c.code, s.buyer_name, to_char(s.sale_date, 'YYYY-MM-DD'), s.total,
		       ($1::date - s.sale_date)::int
		FROM sales s
		JOIN cycles c ON c.id = s.cycle_id
		WHERE NOT s.is_paid
		ORDER BY s.sale_date, s.id
		LIMIT 10`, f.Today)
	if err != nil {
		return err
	}
	defer rows.Close()

	list := make([]model.Receivable, 0)
	for rows.Next() {
		var x model.Receivable
		if err := rows.Scan(&x.SaleID, &x.CycleCode, &x.BuyerName, &x.SaleDate, &x.Total, &x.DaysOutstanding); err != nil {
			return err
		}
		list = append(list, x)
	}
	*out = list
	return rows.Err()
}

func (r *DashboardRepository) ActiveCycles(ctx context.Context, f model.DashboardFilter, out *[]model.ActiveCycleSummary) error {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.code, p.name, fs.name,
		       ($1::date - c.start_date)::int AS age_days,
		       s.alive_qty, fs.target_days,
		       (fs.target_days - ($1::date - c.start_date))::int AS days_to_harvest
		FROM cycles c
		JOIN ponds p           ON p.id = c.pond_id
		JOIN fish_species fs   ON fs.id = c.species_id
		JOIN v_cycle_status s  ON s.cycle_id = c.id
		WHERE c.status = 'active'
		ORDER BY days_to_harvest NULLS LAST, c.id`, f.Today)
	if err != nil {
		return err
	}
	defer rows.Close()

	list := make([]model.ActiveCycleSummary, 0)
	for rows.Next() {
		var a model.ActiveCycleSummary
		if err := rows.Scan(&a.CycleID, &a.Code, &a.PondName, &a.SpeciesName,
			&a.AgeDays, &a.AliveQty, &a.TargetDays, &a.DaysToHarvest); err != nil {
			return err
		}
		list = append(list, a)
	}
	*out = list
	return rows.Err()
}
