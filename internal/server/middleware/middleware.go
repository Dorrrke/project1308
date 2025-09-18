package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Dorrrke/project1308/internal/domain"
	"github.com/Dorrrke/project1308/internal/server/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func AuthMiddleware(signer auth.HS256Signer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := signer.ParseAccessToken(tokenStr, auth.ParseOptions{
			ExpectedIssuer:   signer.Issuer,
			ExpectedAudience: signer.Audience,
			AllowedMethods:   []string{"HS256"},
			Leeway:           domain.LeewayTimeout,
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Next()
	}
}

func ZerologMiddleware(log *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		log.Info().
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Int("status", ctx.Writer.Status()).
			Str("client_ip", ctx.ClientIP()).
			Dur("duration", duration).
			Send()
	}
}
