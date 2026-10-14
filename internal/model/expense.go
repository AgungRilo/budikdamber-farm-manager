package model

import "time"

type Expense struct {
	ID            int64     `json:"id"`
	CycleID       *int64    `json:"cycle_id"`
	CycleCode     *string   `json:"cycle_code"`
	ExpenseDate   string    `json:"expense_date"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	Amount        int64     `json:"amount"`
	CreatedBy     int64     `json:"created_by"`
	CreatedByName string    `json:"created_by_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ExpenseInput struct {
	CycleID     *int64 `json:"cycle_id" binding:"omitempty,gt=0"`
	ExpenseDate string `json:"expense_date" binding:"required,datetime=2006-01-02"`
	Category    string `json:"category" binding:"required,oneof=obat_probiotik peralatan utilitas tenaga_kerja lainnya"`
	Description string `json:"description" binding:"required,max=200"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
}

type ExpenseFilter struct {
	CycleID     *int64
	GeneralOnly bool // true → hanya biaya umum (cycle_id NULL)
	Category    string
	DateFrom    string
	DateTo      string
	Page        int
	Limit       int
}
