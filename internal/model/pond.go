package model

import "time"

type Pond struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	VolumeLiter *int      `json:"volume_liter"`
	MaxStock    int       `json:"max_stock"`
	IsActive    bool      `json:"is_active"`
	Notes       *string   `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PondInput struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Type        string  `json:"type" binding:"required,oneof=ember terpal kolam_tanah kolam_beton"`
	VolumeLiter *int    `json:"volume_liter" binding:"omitempty,gt=0"`
	MaxStock    int     `json:"max_stock" binding:"required,gt=0"`
	IsActive    *bool   `json:"is_active"`
	Notes       *string `json:"notes" binding:"omitempty,max=500"`
}

type PondFilter struct {
	Search   string
	Type     string
	IsActive *bool
	Page     int
	Limit    int
}
