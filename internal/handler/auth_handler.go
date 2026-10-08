package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "email (format valid) dan password wajib diisi"})
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrUserInactive):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case err != nil:
		log.Printf("login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan server"})
	default:
		c.JSON(http.StatusOK, res)
	}
}

func (h *AuthHandler) Me(c *gin.Context) {
	u, err := h.svc.Me(c.Request.Context(), c.GetInt64(middleware.CtxUserID))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan"})
	case err != nil:
		log.Printf("me error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan server"})
	case !u.IsActive:
		c.JSON(http.StatusForbidden, gin.H{"error": "akun dinonaktifkan"})
	default:
		c.JSON(http.StatusOK, u)
	}
}
