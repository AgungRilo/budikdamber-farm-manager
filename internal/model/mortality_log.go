package model

import "time"

type MortalityLog struct {
	ID            int64     `json:"id"`
	CycleID       int64     `json:"cycle_id"`
	LogDate       string    `json:"log_date"`
	Qty           int       `json:"qty"`
	Cause         *string   `json:"cause"`
	Notes         *string   `json:"notes"`
	CreatedBy     int64     `json:"created_by"`
	CreatedByName string    `json:"created_by_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type MortalityLogInput struct {
	LogDate string  `json:"log_date" binding:"required,datetime=2006-01-02"`
	Qty     int     `json:"qty" binding:"required,gt=0"`
	Cause   *string `json:"cause" binding:"omitempty,max=100"`
	Notes   *string `json:"notes" binding:"omitempty,max=500"`
}

// MortalityCreateResult = log + info dampaknya ke siklus.
// MortalityLog di-embed, jadi field-nya "naik" ke level atas di JSON.
type MortalityCreateResult struct {
	MortalityLog
	AliveQtyAfter int64 `json:"alive_qty_after"`
	CycleClosed   bool  `json:"cycle_closed"`
}
