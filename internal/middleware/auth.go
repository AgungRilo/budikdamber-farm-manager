package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/auth"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
)

const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

// UserFinder cukup butuh FindByID; dipenuhi oleh repository.UserRepository.
// Pakai interface supaya middleware tidak bergantung langsung ke repository.
type UserFinder interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

func AuthRequired(jm *auth.JWTManager, users UserFinder) gin.HandlerFunc {
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

		// Ambil status & role TERBARU dari database, bukan dari token,
		// supaya nonaktif / turun role langsung berlaku.
		u, err := users.FindByID(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, "akun tidak ditemukan")
			return
		}
		if !u.IsActive {
			response.Abort(c, http.StatusUnauthorized, "akun nonaktif atau belum diaktifkan owner")
			return
		}

		c.Set(CtxUserID, u.ID)
		c.Set(CtxRole, u.Role)
		c.Next()
	}
}
