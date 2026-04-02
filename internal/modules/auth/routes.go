package auth

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	repositoriesAdmin "saas-medico/internal/modules/admin/repositories"
	"saas-medico/internal/modules/auth/handlers"
	"saas-medico/internal/modules/auth/repositories"
	"saas-medico/internal/modules/auth/services"

	"github.com/gin-gonic/gin"
)

var (
	authMiddleware *middleware.AuthMiddleware
	jwtService     *services.JWTService
)

func Setup() {
	jwtService = services.NewJWTService()
	authMiddleware = middleware.NewAuthMiddleware(jwtService)
}

func GetAuthMiddleware() *middleware.AuthMiddleware {
	return authMiddleware
}

func GetJWTService() *services.JWTService {
	return jwtService
}

func RegisterRoutes(router *gin.RouterGroup) {
	db := database.GetDB()

	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	rolRepo := repositories.NewRolRepository(db)
	adminRepo := repositoriesAdmin.NewAdminRepository(db)

	authSvc := services.NewAuthService(userRepo, tokenRepo, rolRepo, adminRepo, jwtService)
	handler := handlers.NewAuthHandler(authSvc)

	auth := router.Group("/auth")
	{
		// Rutas públicas
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
		auth.POST("/refresh", handler.RefreshToken)

		// Rutas protegidas
		protected := auth.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/permisosMenusPaleta", handler.BuildMenuPaleta)
			protected.GET("/estilos", handler.GetEstiloClinica)
			protected.POST("/logout", handler.Logout)
			protected.GET("/me", handler.Me)
			protected.PUT("/me", handler.UpdateProfile)
			protected.POST("/me/foto", handler.UploadFoto)
			protected.POST("/change-password", handler.ChangePassword)
		}
	}
}
