package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type SaleHandler struct {
	svc *service.SaleService
}

func NewSaleHandler(svc *service.SaleService) *SaleHandler {
	return &SaleHandler{svc: svc}
}

func saleError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrChannelRequired) {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cycleLogError(c, err, "penjualan")
}

// List godoc
// @Summary   Daftar penjualan per siklus
// @Tags      Sales
// @Produce   json
// @Security  BearerAuth
// @Param     id         path   int     true   "ID siklus"
// @Param     date_from  query  string  false  "Dari tanggal (YYYY-MM-DD)"
// @Param     date_to    query  string  false  "Sampai tanggal (YYYY-MM-DD)"
// @Param     sale_type  query  string  false  "Jenis"  Enums(sale, own_use)
// @Param     is_paid    query  bool    false  "Status lunas"
// @Param     page       query  int     false  "Halaman"  default(1)
// @Param     limit      query  int     false  "Jumlah per halaman"  default(20)
// @Success   200  {object}  response.Body{data=[]model.Sale,meta=response.Meta}
// @Failure   404  {object}  response.Body
// @Router    /cycles/{id}/sales [get]
func (h *SaleHandler) List(c *gin.Context) {
	cycleID, ok := parseID(c)
	if !ok {
		return
	}
	lf, ok := parseLogFilter(c, cycleID)
	if !ok {
		return
	}
	f := model.SaleFilter{LogFilter: lf, SaleType: c.Query("sale_type"), IsPaid: queryBool(c, "is_paid")}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		saleError(c, err)
		return
	}
	response.Paginated(c, "berhasil mengambil daftar penjualan", list, response.NewMeta(f.Page, f.Limit, total))
}

// Create godoc
// @Summary   Catat penjualan / konsumsi sendiri (owner & worker)
// @Tags      Sales
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int              true  "ID siklus"
// @Param     body  body      model.SaleInput  true  "Data penjualan"
// @Success   201   {object}  response.Body{data=model.SaleCreateResult}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /cycles/{id}/sales [post]
func (h *SaleHandler) Create(c *gin.Context) {
	cycleID, ok := parseID(c)
	if !ok {
		return
	}
	var in model.SaleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	res, err := h.svc.Create(c.Request.Context(), cycleID, in, c.GetInt64(middleware.CtxUserID))
	if err != nil {
		saleError(c, err)
		return
	}
	msg := "penjualan berhasil dicatat"
	if res.CycleClosed {
		msg += "; ikan hidup habis, siklus otomatis ditutup"
	}
	response.Created(c, msg, res)
}

// SetPaid godoc
// @Summary   Tandai penjualan lunas / belum lunas (owner)
// @Tags      Sales
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                  true  "ID penjualan"
// @Param     body  body      model.SalePaidInput  true  "Status lunas"
// @Success   200   {object}  response.Body{data=model.Sale}
// @Failure   404   {object}  response.Body
// @Router    /sales/{id}/paid [patch]
func (h *SaleHandler) SetPaid(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.SalePaidInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	s, err := h.svc.SetPaid(c.Request.Context(), id, *in.IsPaid)
	if err != nil {
		saleError(c, err)
		return
	}
	msg := "penjualan ditandai lunas"
	if !*in.IsPaid {
		msg = "penjualan ditandai belum lunas"
	}
	response.OK(c, msg, s)
}

// Delete godoc
// @Summary   Hapus penjualan (owner)
// @Tags      Sales
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID penjualan"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Failure   409  {object}  response.Body
// @Router    /sales/{id} [delete]
func (h *SaleHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		saleError(c, err)
		return
	}
	response.OK(c, "penjualan berhasil dihapus", nil)
}
