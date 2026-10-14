package model

import "time"

type Dashboard struct {
	Stock          DashboardStock       `json:"stock"`
	Finance        DashboardFinance     `json:"finance"`
	SalesByChannel []ChannelSales       `json:"sales_by_channel"`
	Receivables    []Receivable         `json:"receivables"`
	ActiveCycles   []ActiveCycleSummary `json:"active_cycles"`
	GeneratedAt    time.Time            `json:"generated_at"`
}

type DashboardStock struct {
	ActiveCycles    int64   `json:"active_cycles"`
	AliveFish       int64   `json:"alive_fish"`
	ActivePonds     int64   `json:"active_ponds"`
	TotalCapacity   int64   `json:"total_capacity"`
	CapacityUsedPct float64 `json:"capacity_used_pct"`
}

type DashboardFinance struct {
	TotalCost       int64 `json:"total_cost"`
	Revenue         int64 `json:"revenue"`
	Receivable      int64 `json:"receivable"`
	CycleProfit     int64 `json:"cycle_profit"`
	GeneralExpenses int64 `json:"general_expenses"`
	NetProfit       int64 `json:"net_profit"`
}

type ChannelSales struct {
	Channel        string `json:"channel"`
	Transactions   int64  `json:"transactions"`
	Qty            int64  `json:"qty"`
	Revenue        int64  `json:"revenue"`
	EstGrossMargin int64  `json:"est_gross_margin"`
}

type Receivable struct {
	SaleID          int64   `json:"sale_id"`
	CycleCode       string  `json:"cycle_code"`
	BuyerName       *string `json:"buyer_name"`
	SaleDate        string  `json:"sale_date"`
	Total           int64   `json:"total"`
	DaysOutstanding int     `json:"days_outstanding"`
}

type ActiveCycleSummary struct {
	CycleID       int64  `json:"cycle_id"`
	Code          string `json:"code"`
	PondName      string `json:"pond_name"`
	SpeciesName   string `json:"species_name"`
	AgeDays       int    `json:"age_days"`
	AliveQty      int64  `json:"alive_qty"`
	TargetDays    *int   `json:"target_days"`
	DaysToHarvest *int   `json:"days_to_harvest"` // negatif = lewat target
}

type DashboardFilter struct {
	DateFrom string
	DateTo   string
	Today    string // "hari ini" versi WIB, diisi service
}
