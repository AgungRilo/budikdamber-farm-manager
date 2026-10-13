// ===== internal/handler/fish_species_handler.go =====
package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type FishSpeciesHandler struct {
	svc *service.FishSpeciesService
}

func NewFishSpeciesHandler(svc *service.FishSpeciesService) *FishSpeciesHandler {
	return &FishSpeciesHandler{svc: svc}
}

// List godoc
// @Summary   Daftar jenis ikan
// @Tags      Fish Species
// @Produce   json
// @Security  BearerAuth
// @Param     page    query  int     false  "Halaman"  default(1)
// @Param     limit   query  int     false  "Jumlah per halaman"  default(10)
// @Param     search  query  string  false  "Cari nama"
// @Success   200  {object}  response.Body{data=[]model.FishSpecies,meta=response.Meta}
// @Router    /fish-species [get]
func (h *FishSpeciesHandler) List(c *gin.Context) {
	f := model.FishSpeciesFilter{
		Search: c.Query("search"),
		Page:   queryInt(c, "page", 1, 1, 100000),
		Limit:  queryInt(c, "limit", 10, 1, 100),
	}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		handleError(c, err, "jenis ikan")
		return
	}
	response.Paginated(c, "berhasil mengambil daftar jenis ikan", list, response.NewMeta(f.Page, f.Limit, total))
}

// Get godoc
// @Summary   Detail jenis ikan
// @Tags      Fish Species
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID jenis ikan"
// @Success   200  {object}  response.Body{data=model.FishSpecies}
// @Failure   404  {object}  response.Body
// @Router    /fish-species/{id} [get]
func (h *FishSpeciesHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	s, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err, "jenis ikan")
		return
	}
	response.OK(c, "berhasil mengambil detail jenis ikan", s)
}

// Create godoc
// @Summary   Buat jenis ikan (owner)
// @Tags      Fish Species
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     body  body      model.FishSpeciesInput  true  "Data jenis ikan"
// @Success   201   {object}  response.Body{data=model.FishSpecies}
// @Failure   400   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /fish-species [post]
func (h *FishSpeciesHandler) Create(c *gin.Context) {
	var in model.FishSpeciesInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	s, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		handleError(c, err, "jenis ikan")
		return
	}
	response.Created(c, "jenis ikan berhasil dibuat", s)
}

// Update godoc
// @Summary   Ubah jenis ikan (owner)
// @Tags      Fish Species
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                     true  "ID jenis ikan"
// @Param     body  body      model.FishSpeciesInput  true  "Data jenis ikan"
// @Success   200   {object}  response.Body{data=model.FishSpecies}
// @Failure   404   {object}  response.Body
// @Router    /fish-species/{id} [put]
func (h *FishSpeciesHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.FishSpeciesInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	s, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		handleError(c, err, "jenis ikan")
		return
	}
	response.OK(c, "jenis ikan berhasil diperbarui", s)
}

// Delete godoc
// @Summary   Hapus jenis ikan / soft delete (owner)
// @Tags      Fish Species
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID jenis ikan"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Router    /fish-species/{id} [delete]
func (h *FishSpeciesHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err, "jenis ikan")
		return
	}
	response.OK(c, "jenis ikan berhasil dihapus", nil)
}
