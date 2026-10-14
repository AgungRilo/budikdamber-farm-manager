package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type CycleHandler struct {
	svc *service.CycleService
}

func NewCycleHandler(svc *service.CycleService) *CycleHandler {
	return &CycleHandler{svc: svc}
}

func cycleError(c *gin.Context, err error) {
	var capErr *repository.CapacityError
	switch {
	case errors.As(err, &capErr):
		c.JSON(http.StatusConflict, response.Body{
			Success: false,
			Message: capErr.Error(),
			Data:    capErr, // max_stock, used, remaining ikut dikirim
		})
	case errors.Is(err, repository.ErrPondUnavailable),
		errors.Is(err, repository.ErrSpeciesNotFound),
		errors.Is(err, service.ErrStartDateInFuture):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		handleError(c, err, "siklus")
	}
}

func queryInt64(c *gin.Context, key string) int64 {
	v, _ := strconv.ParseInt(c.Query(key), 10, 64)
	return v
}

// List godoc
// @Summary   Daftar siklus budidaya
// @Tags      Cycles
// @Produce   json
// @Security  BearerAuth
// @Param     status      query  string  false  "Status"  Enums(active, closed)
// @Param     pond_id     query  int     false  "Filter wadah"
// @Param     species_id  query  int     false  "Filter jenis ikan"
// @Param     page        query  int     false  "Halaman"  default(1)
// @Param     limit       query  int     false  "Jumlah per halaman"  default(10)
// @Success   200  {object}  response.Body{data=[]model.Cycle,meta=response.Meta}
// @Router    /cycles [get]
func (h *CycleHandler) List(c *gin.Context) {
	f := model.CycleFilter{
		Status:    c.Query("status"),
		PondID:    queryInt64(c, "pond_id"),
		SpeciesID: queryInt64(c, "species_id"),
		Page:      queryInt(c, "page", 1, 1, 100000),
		Limit:     queryInt(c, "limit", 10, 1, 100),
	}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		cycleError(c, err)
		return
	}
	response.Paginated(c, "berhasil mengambil daftar siklus", list, response.NewMeta(f.Page, f.Limit, total))
}

// Get godoc
// @Summary   Detail siklus (termasuk jumlah hidup, mati, terjual)
// @Tags      Cycles
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID siklus"
// @Success   200  {object}  response.Body{data=model.Cycle}
// @Failure   404  {object}  response.Body
// @Router    /cycles/{id} [get]
func (h *CycleHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	cy, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		cycleError(c, err)
		return
	}
	response.OK(c, "berhasil mengambil detail siklus", cy)
}

// Create godoc
// @Summary   Mulai siklus baru / tebar bibit (owner)
// @Tags      Cycles
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     body  body      model.CycleInput  true  "Data siklus"
// @Success   201   {object}  response.Body{data=model.Cycle}
// @Failure   400   {object}  response.Body
// @Failure   409   {object}  response.Body  "Kapasitas wadah tidak cukup"
// @Router    /cycles [post]
func (h *CycleHandler) Create(c *gin.Context) {
	var in model.CycleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	cy, err := h.svc.Create(c.Request.Context(), in, c.GetInt64(middleware.CtxUserID))
	if err != nil {
		cycleError(c, err)
		return
	}
	response.Created(c, "siklus berhasil dibuat", cy)
}

// cycleStateError: error khusus tutup/buka siklus, sisanya diteruskan ke cycleError.
func cycleStateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrCycleNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, repository.ErrCycleClosed), errors.Is(err, repository.ErrCycleNotClosed):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, repository.ErrEndBeforeStart), errors.Is(err, service.ErrLogDateInFuture):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		cycleError(c, err) // CapacityError, ErrPondUnavailable, dst.
	}
}

// Close godoc
// @Summary   Tutup siklus manual (owner)
// @Tags      Cycles
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                    true   "ID siklus"
// @Param     body  body      model.CycleCloseInput  false  "Tanggal selesai (default hari ini)"
// @Success   200   {object}  response.Body{data=model.Cycle}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /cycles/{id}/close [post]
func (h *CycleHandler) Close(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.CycleCloseInput
	if c.Request.ContentLength > 0 { // body boleh kosong
		if err := c.ShouldBindJSON(&in); err != nil {
			response.ValidationError(c, err)
			return
		}
	}
	cy, err := h.svc.Close(c.Request.Context(), id, in, c.GetInt64(middleware.CtxUserID))
	if err != nil {
		cycleStateError(c, err)
		return
	}
	response.OK(c, "siklus berhasil ditutup", cy)
}

// Reopen godoc
// @Summary   Buka kembali siklus yang sudah ditutup (owner)
// @Tags      Cycles
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID siklus"
// @Success   200  {object}  response.Body{data=model.Cycle}
// @Failure   404  {object}  response.Body
// @Failure   409  {object}  response.Body
// @Router    /cycles/{id}/reopen [post]
func (h *CycleHandler) Reopen(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	cy, err := h.svc.Reopen(c.Request.Context(), id)
	if err != nil {
		cycleStateError(c, err)
		return
	}
	response.OK(c, "siklus berhasil dibuka kembali", cy)
}

// Metrics godoc
// @Summary      Metrik bisnis siklus: SR, biaya, HPP, laba, margin (owner)
// @Description  Angka bersifat sementara selama siklus masih aktif (is_final=false).
// @Tags         Cycles
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID siklus"
// @Success      200  {object}  response.Body{data=model.CycleMetrics}
// @Failure      403  {object}  response.Body
// @Failure      404  {object}  response.Body
// @Router       /cycles/{id}/metrics [get]
func (h *CycleHandler) Metrics(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	m, err := h.svc.Metrics(c.Request.Context(), id)
	if err != nil {
		cycleStateError(c, err)
		return
	}
	response.OK(c, "berhasil mengambil metrik siklus", m)
}
