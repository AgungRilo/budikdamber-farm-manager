package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type PondHandler struct {
	svc *service.PondService
}

func NewPondHandler(svc *service.PondService) *PondHandler {
	return &PondHandler{svc: svc}
}

// List godoc
// @Summary   Daftar wadah
// @Tags      Ponds
// @Produce   json
// @Security  BearerAuth
// @Param     page       query  int     false  "Halaman"             default(1)
// @Param     limit      query  int     false  "Jumlah per halaman"  default(10)
// @Param     search     query  string  false  "Cari berdasarkan nama"
// @Param     type       query  string  false  "Jenis wadah"  Enums(ember, terpal, kolam_tanah, kolam_beton)
// @Param     is_active  query  bool    false  "Status aktif"
// @Success   200  {object}  response.Body{data=[]model.Pond,meta=response.Meta}
// @Failure   401  {object}  response.Body
// @Router    /ponds [get]
func (h *PondHandler) List(c *gin.Context) {
	f := model.PondFilter{
		Search:   c.Query("search"),
		Type:     c.Query("type"),
		IsActive: queryBool(c, "is_active"),
		Page:     queryInt(c, "page", 1, 1, 100000),
		Limit:    queryInt(c, "limit", 10, 1, 100),
	}
	ponds, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.Paginated(c, "berhasil mengambil daftar wadah", ponds, response.NewMeta(f.Page, f.Limit, total))
}

// Get godoc
// @Summary   Detail wadah
// @Tags      Ponds
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID wadah"
// @Success   200  {object}  response.Body{data=model.Pond}
// @Failure   404  {object}  response.Body
// @Router    /ponds/{id} [get]
func (h *PondHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	p, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.OK(c, "berhasil mengambil detail wadah", p)
}

// Create godoc
// @Summary   Buat wadah (owner)
// @Tags      Ponds
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     body  body      model.PondInput  true  "Data wadah"
// @Success   201   {object}  response.Body{data=model.Pond}
// @Failure   400   {object}  response.Body
// @Failure   403   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /ponds [post]
func (h *PondHandler) Create(c *gin.Context) {
	var in model.PondInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	p, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.Created(c, "wadah berhasil dibuat", p)
}

// Update godoc
// @Summary   Ubah wadah (owner)
// @Tags      Ponds
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int              true  "ID wadah"
// @Param     body  body      model.PondInput  true  "Data wadah"
// @Success   200   {object}  response.Body{data=model.Pond}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /ponds/{id} [put]
func (h *PondHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.PondInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.OK(c, "wadah berhasil diperbarui", p)
}

// Delete godoc
// @Summary   Hapus wadah / soft delete (owner)
// @Tags      Ponds
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID wadah"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Router    /ponds/{id} [delete]
func (h *PondHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.OK(c, "wadah berhasil dihapus", nil)
}

var _ = http.StatusOK // jaga import net/http bila belum terpakai
