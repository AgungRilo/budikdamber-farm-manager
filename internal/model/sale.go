package model

import "time"

type Sale struct {
	ID            int64     `json:"id"`
	CycleID       int64     `json:"cycle_id"`
	SaleDate      string    `json:"sale_date"`
	SaleType      string    `json:"sale_type"`
	Channel       *string   `json:"channel"`
	Qty           int       `json:"qty"`
	WeightKg      *float64  `json:"weight_kg"`
	PricePerFish  int64     `json:"price_per_fish"`
	Total         int64     `json:"total"`
	BuyerName     *string   `json:"buyer_name"`
	IsPaid        bool      `json:"is_paid"`
	Notes         *string   `json:"notes"`
	CreatedBy     int64     `json:"created_by"`
	CreatedByName string    `json:"created_by_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SaleInput struct {
	SaleDate     string   `json:"sale_date" binding:"required,datetime=2006-01-02"`
	SaleType     string   `json:"sale_type" binding:"omitempty,oneof=sale own_use"`
	Channel      *string  `json:"channel" binding:"omitempty,oneof=langsung pasar pengepul online"`
	Qty          int      `json:"qty" binding:"required,gt=0"`
	WeightKg     *float64 `json:"weight_kg" binding:"omitempty,gt=0"`
	PricePerFish *int64   `json:"price_per_fish" binding:"omitempty,gte=0"` // kosong → harga default jenis ikan
	BuyerName    *string  `json:"buyer_name" binding:"omitempty,max=100"`
	IsPaid       *bool    `json:"is_paid"`
	Notes        *string  `json:"notes" binding:"omitempty,max=500"`
}

type SalePaidInput struct {
	IsPaid *bool `json:"is_paid" binding:"required"`
}

type SaleFilter struct {
	LogFilter
	SaleType string
	IsPaid   *bool
}

type SaleCreateResult struct {
	Sale
	AliveQtyAfter int64 `json:"alive_qty_after"`
	CycleClosed   bool  `json:"cycle_closed"`
}
