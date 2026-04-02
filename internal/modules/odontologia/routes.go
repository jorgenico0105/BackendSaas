package odontologia

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/auth/models"
	"saas-medico/internal/modules/odontologia/handlers"
	"saas-medico/internal/modules/odontologia/repositories"
	"saas-medico/internal/modules/odontologia/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewOdontologiaRepository(database.GetDB())
	service := services.NewOdontologiaService(repo)
	handler := handlers.NewOdontologiaHandler(service)

	odontologia := router.Group("/odontologia")
	odontologia.Use(authMiddleware.RequireAuth())
	{
		// Ping para verificar que el módulo funciona (requiere autenticación)
		odontologia.GET("/ping", handler.Ping)

		// Rutas protegidas por roles específicos
		// Solo odontólogos, admins y super_admins pueden acceder
		restricted := odontologia.Group("")
		restricted.Use(authMiddleware.RequireRoles(
			models.RolSuperAdmin,
			models.RolAdmin,
			models.RolOdontologo,
		))
		{
			// Aquí irán las rutas específicas de odontología
			// Ejemplo:
			// restricted.GET("/historiales", handler.ListHistoriales)
			// restricted.GET("/historiales/:id", handler.GetHistorial)
			// restricted.POST("/historiales", handler.CreateHistorial)
			// restricted.PUT("/historiales/:id", handler.UpdateHistorial)
			// restricted.DELETE("/historiales/:id", handler.DeleteHistorial)
		}
	}
}
