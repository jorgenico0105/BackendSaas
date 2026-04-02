package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"saas-medico/internal/modules/auth/services"
	"saas-medico/internal/shared/responses"
)

type AuthMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthMiddleware(jwtService *services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			responses.Unauthorized(c, "Token de autorización requerido")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			responses.Unauthorized(c, "Formato de token inválido. Use: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(parts[1])

		if err != nil {
			if err == services.ErrExpiredToken {
				responses.Unauthorized(c, "Token expirado")
			} else {
				responses.Unauthorized(c, "Token inválido")
			}
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("rolID", claims.RolID)
		c.Set("rolName", claims.RolName)
		c.Set("clinicaID", claims.ClinicaID)
		c.Set("aplicacionID", claims.AplicacionID)

		c.Next()
	}
}

// RequireRoles verifica que el usuario tenga uno de los roles especificados
func (m *AuthMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolName, exists := c.Get("rolName")
		if !exists {
			responses.Unauthorized(c, "Usuario no autenticado")
			c.Abort()
			return
		}

		userRol := rolName.(string)
		for _, rol := range roles {
			if userRol == rol {
				c.Next()
				return
			}
		}

		responses.Forbidden(c, "No tiene permisos para acceder a este recurso")
		c.Abort()
	}
}
