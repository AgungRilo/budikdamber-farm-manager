package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrUserInactive):
		response.Error(c, http.StatusForbidden, err.Error())
	case err != nil:
		log.Printf("login error: %v", err)
		response.Error(c, http.StatusInternalServerError, "terjadi kesalahan server")
	default:
		response.OK(c, "login berhasil", res)
	}
}

func (h *AuthHandler) Me(c *gin.Context) {
	u, err := h.svc.Me(c.Request.Context(), c.GetInt64(middleware.CtxUserID))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan")
	case err != nil:
		log.Printf("me error: %v", err)
		response.Error(c, http.StatusInternalServerError, "terjadi kesalahan server")
	case !u.IsActive:
		response.Error(c, http.StatusForbidden, "akun dinonaktifkan")
	default:
		response.OK(c, "berhasil mengambil profil", u)
	}
}
