package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type MortalityLogHandler struct {
	svc *service.MortalityLogService
}

func NewMortalityLogHandler(svc *service.MortalityLogService) *MortalityLogHandler {
	return &MortalityLogHandler{svc: svc}
}

// List godoc
// @Summary   Daftar log kematian per siklus
// @Tags      Mortality Logs
// @Produce   json
// @Security  BearerAuth
// @Param     id         path   int     true   "ID siklus"
// @Param     date_from  query  string  false  "Dari tanggal (YYYY-MM-DD)"
// @Param     date_to    query  string  false  "Sampai tanggal (YYYY-MM-DD)"
// @Param     page       query  int     false  "Halaman"  default(1)
// @Param     limit      query  int     false  "Jumlah per halaman"  default(20)
// @Success   200  {object}  response.Body{data=[]model.MortalityLog,meta=response.Meta}
// @Failure   404  {object}  response.Body
// @Router    /cycles/{id}/mortality-logs [get]
func (h *MortalityLogHandler) List(c *gin.Context) {
	cycleID, ok := parseID(c)
	if !ok {
		return
	}
	f, ok := parseLogFilter(c, cycleID)
	if !ok {
		return
	}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		cycleLogError(c, err, "log kematian")
		return
	}
	response.Paginated(c, "berhasil mengambil log kematian", list, response.NewMeta(f.Page, f.Limit, total))
}

// Create godoc
// @Summary   Catat ikan mati (owner & worker); siklus otomatis ditutup jika ikan hidup habis
// @Tags      Mortality Logs
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                      true  "ID siklus"
// @Param     body  body      model.MortalityLogInput  true  "Data log kematian"
// @Success   201   {object}  response.Body{data=model.MortalityCreateResult}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /cycles/{id}/mortality-logs [post]
func (h *MortalityLogHandler) Create(c *gin.Context) {
	cycleID, ok := parseID(c)
	if !ok {
		return
	}
	var in model.MortalityLogInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	res, err := h.svc.Create(c.Request.Context(), cycleID, in, c.GetInt64(middleware.CtxUserID))
	if err != nil {
		cycleLogError(c, err, "log kematian")
		return
	}
	msg := "log kematian berhasil dicatat"
	if res.CycleClosed {
		msg += "; ikan hidup habis, siklus otomatis ditutup"
	}
	response.Created(c, msg, res)
}

// Delete godoc
// @Summary   Hapus log kematian (owner)
// @Tags      Mortality Logs
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID log kematian"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Failure   409  {object}  response.Body
// @Router    /mortality-logs/{id} [delete]
func (h *MortalityLogHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		cycleLogError(c, err, "log kematian")
		return
	}
	response.OK(c, "log kematian berhasil dihapus", nil)
}
