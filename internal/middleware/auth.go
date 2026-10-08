package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/auth"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
)

const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

func AuthRequired(jm *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || tokenStr == "" {
			response.Abort(c, http.StatusUnauthorized, "token tidak ditemukan")
			return
		}

		claims, err := jm.Parse(tokenStr)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, "token tidak valid atau kedaluwarsa")
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}
