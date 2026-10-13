package model

import "time"

type FishSpecies struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	TargetDays          *int      `json:"target_days"`
	DefaultPricePerFish int64     `json:"default_price_per_fish"`
	DefaultPricePerKg   int64     `json:"default_price_per_kg"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type FishSpeciesInput struct {
	Name                string `json:"name" binding:"required,max=100"`
	TargetDays          *int   `json:"target_days" binding:"omitempty,gt=0"`
	DefaultPricePerFish int64  `json:"default_price_per_fish" binding:"gte=0"`
	DefaultPricePerKg   int64  `json:"default_price_per_kg" binding:"gte=0"`
}

type FishSpeciesFilter struct {
	Search string
	Page   int
	Limit  int
}
