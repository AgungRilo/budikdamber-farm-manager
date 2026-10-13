package model

import "time"

type Feed struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Brand      *string   `json:"brand"`
	ProteinPct *float64  `json:"protein_pct"`
	PricePerKg int64     `json:"price_per_kg"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FeedInput struct {
	Name       string   `json:"name" binding:"required,max=100"`
	Brand      *string  `json:"brand" binding:"omitempty,max=100"`
	ProteinPct *float64 `json:"protein_pct" binding:"omitempty,gte=0,lte=100"`
	PricePerKg int64    `json:"price_per_kg" binding:"required,gt=0"`
}

type FeedFilter struct {
	Search string
	Page   int
	Limit  int
}
