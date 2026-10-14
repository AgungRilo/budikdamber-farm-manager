package model

import "time"

type Cycle struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	PondID           int64     `json:"pond_id"`
	PondName         string    `json:"pond_name"`
	SpeciesID        int64     `json:"species_id"`
	SpeciesName      string    `json:"species_name"`
	StartDate        string    `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	AgeDays          int       `json:"age_days"`
	InitialQty       int64     `json:"initial_qty"`
	DeadQty          int64     `json:"dead_qty"`
	SoldQty          int64     `json:"sold_qty"`
	AliveQty         int64     `json:"alive_qty"`
	SeedSizeCm       *float64  `json:"seed_size_cm"`
	SeedPricePerFish int64     `json:"seed_price_per_fish"`
	Status           string    `json:"status"`
	Notes            *string   `json:"notes"`
	CreatedBy        int64     `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CycleInput struct {
	PondID           int64    `json:"pond_id" binding:"required,gt=0"`
	SpeciesID        int64    `json:"species_id" binding:"required,gt=0"`
	StartDate        string   `json:"start_date" binding:"required,datetime=2006-01-02"`
	InitialQty       int      `json:"initial_qty" binding:"required,gt=0"`
	SeedSizeCm       *float64 `json:"seed_size_cm" binding:"omitempty,gt=0"`
	SeedPricePerFish int64    `json:"seed_price_per_fish" binding:"gte=0"`
	Notes            *string  `json:"notes" binding:"omitempty,max=500"`
}

type CycleFilter struct {
	Status    string
	PondID    int64
	SpeciesID int64
	Page      int
	Limit     int
}

type CycleCloseInput struct {
	EndDate string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

type CycleMetrics struct {
	CycleID         int64    `json:"cycle_id"`
	Code            string   `json:"code"`
	Status          string   `json:"status"`
	InitialQty      int64    `json:"initial_qty"`
	DeadQty         int64    `json:"dead_qty"`
	AliveQty        int64    `json:"alive_qty"`
	SaleQty         int64    `json:"sale_qty"`
	OwnUseQty       int64    `json:"own_use_qty"`
	SurvivalRatePct float64  `json:"survival_rate_pct"`
	FeedGram        int64    `json:"feed_gram"`
	SeedCost        int64    `json:"seed_cost"`
	FeedCost        int64    `json:"feed_cost"`
	OtherCost       int64    `json:"other_cost"`
	TotalCost       int64    `json:"total_cost"`
	Revenue         int64    `json:"revenue"`
	Receivable      int64    `json:"receivable"`
	Profit          int64    `json:"profit"`
	HppPerFish      *int64   `json:"hpp_per_fish"` // null kalau tidak ada ikan yang bertahan
	MarginPct       *float64 `json:"margin_pct"`   // null kalau belum ada pendapatan
	IsFinal         bool     `json:"is_final"`     // false = siklus masih berjalan, angka sementara
}
