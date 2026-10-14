package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Get godoc
// @Summary      Dashboard usaha (owner)
// @Description  Stok, keuangan, penjualan per kanal, piutang, dan siklus aktif. 5 query dijalankan paralel (goroutine/errgroup).
// @Tags         Dashboard
// @Produce      json
// @Security     BearerAuth
// @Param        date_from  query  string  false  "Periode penjualan per kanal: dari (YYYY-MM-DD)"
// @Param        date_to    query  string  false  "Periode penjualan per kanal: sampai (YYYY-MM-DD)"
// @Success      200  {object}  response.Body{data=model.Dashboard}
// @Failure      403  {object}  response.Body
// @Router       /dashboard [get]
func (h *DashboardHandler) Get(c *gin.Context) {
	f := model.DashboardFilter{DateFrom: c.Query("date_from"), DateTo: c.Query("date_to")}
	for _, d := range []string{f.DateFrom, f.DateTo} {
		if d == "" {
			continue
		}
		if _, err := time.Parse("2006-01-02", d); err != nil {
			response.Error(c, http.StatusBadRequest, "format tanggal filter harus YYYY-MM-DD")
			return
		}
	}
	d, err := h.svc.Get(c.Request.Context(), f)
	if err != nil {
		handleError(c, err, "dashboard")
		return
	}
	response.OK(c, "berhasil mengambil dashboard", d)
}
