package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type ExpenseHandler struct {
	svc *service.ExpenseService
}

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

// List godoc
// @Summary   Daftar pengeluaran
// @Tags      Expenses
// @Produce   json
// @Security  BearerAuth
// @Param     cycle_id      query  int     false  "Filter per siklus"
// @Param     general_only  query  bool    false  "Hanya biaya umum (tanpa siklus)"
// @Param     category      query  string  false  "Kategori"  Enums(obat_probiotik, peralatan, utilitas, tenaga_kerja, lainnya)
// @Param     date_from     query  string  false  "Dari tanggal (YYYY-MM-DD)"
// @Param     date_to       query  string  false  "Sampai tanggal (YYYY-MM-DD)"
// @Param     page          query  int     false  "Halaman"  default(1)
// @Param     limit         query  int     false  "Jumlah per halaman"  default(20)
// @Success   200  {object}  response.Body{data=[]model.Expense,meta=response.Meta}
// @Router    /expenses [get]
func (h *ExpenseHandler) List(c *gin.Context) {
	f := model.ExpenseFilter{
		Category: c.Query("category"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     queryInt(c, "page", 1, 1, 100000),
		Limit:    queryInt(c, "limit", 20, 1, 100),
	}
	if v, err := strconv.ParseInt(c.Query("cycle_id"), 10, 64); err == nil && v > 0 {
		f.CycleID = &v
	}
	if b := queryBool(c, "general_only"); b != nil {
		f.GeneralOnly = *b
	}
	for _, d := range []string{f.DateFrom, f.DateTo} {
		if d == "" {
			continue
		}
		if _, err := time.Parse("2006-01-02", d); err != nil {
			response.Error(c, http.StatusBadRequest, "format tanggal filter harus YYYY-MM-DD")
			return
		}
	}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		cycleLogError(c, err, "pengeluaran")
		return
	}
	response.Paginated(c, "berhasil mengambil daftar pengeluaran", list, response.NewMeta(f.Page, f.Limit, total))
}

// Get godoc
// @Summary   Detail pengeluaran
// @Tags      Expenses
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID pengeluaran"
// @Success   200  {object}  response.Body{data=model.Expense}
// @Failure   404  {object}  response.Body
// @Router    /expenses/{id} [get]
func (h *ExpenseHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	e, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		cycleLogError(c, err, "pengeluaran")
		return
	}
	response.OK(c, "berhasil mengambil detail pengeluaran", e)
}

// Create godoc
// @Summary   Catat pengeluaran (owner)
// @Tags      Expenses
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     body  body      model.ExpenseInput  true  "Data pengeluaran"
// @Success   201   {object}  response.Body{data=model.Expense}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /expenses [post]
func (h *ExpenseHandler) Create(c *gin.Context) {
	var in model.ExpenseInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	e, err := h.svc.Create(c.Request.Context(), in, c.GetInt64(middleware.CtxUserID))
	if err != nil {
		cycleLogError(c, err, "pengeluaran")
		return
	}
	response.Created(c, "pengeluaran berhasil dicatat", e)
}

// Update godoc
// @Summary   Ubah pengeluaran (owner)
// @Tags      Expenses
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                 true  "ID pengeluaran"
// @Param     body  body      model.ExpenseInput  true  "Data pengeluaran"
// @Success   200   {object}  response.Body{data=model.Expense}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /expenses/{id} [put]
func (h *ExpenseHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.ExpenseInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	e, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		cycleLogError(c, err, "pengeluaran")
		return
	}
	response.OK(c, "pengeluaran berhasil diperbarui", e)
}

// Delete godoc
// @Summary   Hapus pengeluaran (owner)
// @Tags      Expenses
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID pengeluaran"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Failure   409  {object}  response.Body
// @Router    /expenses/{id} [delete]
func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		cycleLogError(c, err, "pengeluaran")
		return
	}
	response.OK(c, "pengeluaran berhasil dihapus", nil)
}
