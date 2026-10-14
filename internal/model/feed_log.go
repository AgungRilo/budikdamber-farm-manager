package model

import "time"

type FeedLog struct {
	ID            int64     `json:"id"`
	CycleID       int64     `json:"cycle_id"`
	FeedID        int64     `json:"feed_id"`
	FeedName      string    `json:"feed_name"`
	LogDate       string    `json:"log_date"`
	QtyGram       int       `json:"qty_gram"`
	PricePerKg    int64     `json:"price_per_kg"`
	Cost          int64     `json:"cost"`
	Notes         *string   `json:"notes"`
	CreatedBy     int64     `json:"created_by"`
	CreatedByName string    `json:"created_by_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type FeedLogInput struct {
	FeedID  int64   `json:"feed_id" binding:"required,gt=0"`
	LogDate string  `json:"log_date" binding:"required,datetime=2006-01-02"`
	QtyGram int     `json:"qty_gram" binding:"required,gt=0,lte=100000"`
	Notes   *string `json:"notes" binding:"omitempty,max=500"`
}

// LogFilter dipakai bersama oleh log pakan & log kematian.
type LogFilter struct {
	CycleID  int64
	DateFrom string
	DateTo   string
	Page     int
	Limit    int
}
