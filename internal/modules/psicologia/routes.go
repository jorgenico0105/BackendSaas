package psicologia

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/auth/models"
	"saas-medico/internal/modules/psicologia/handlers"
	"saas-medico/internal/modules/psicologia/repositories"
	"saas-medico/internal/modules/psicologia/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewPsicologiaRepository(database.GetDB())
	service := services.NewPsicologiaService(repo)
	handler := handlers.NewPsicologiaHandler(service)

	psicologia := router.Group("/psicologia")
	psicologia.Use(authMiddleware.RequireAuth())
	{
		// Ping para verificar que el módulo funciona (requiere autenticación)
		psicologia.GET("/ping", handler.Ping)

		// Rutas protegidas por roles específicos
		// Solo psicólogos, admins y super_admins pueden acceder
		restricted := psicologia.Group("")
		restricted.Use(authMiddleware.RequireRoles(
			models.RolSuperAdmin,
			models.RolAdmin,
			models.RolPsicologo,
		))
		{
			// Aquí irán las rutas específicas de psicología
			// Ejemplo:
			// restricted.GET("/pacientes", handler.ListPacientes)
			// restricted.GET("/pacientes/:id", handler.GetPaciente)
			// restricted.POST("/pacientes", handler.CreatePaciente)
			// restricted.PUT("/pacientes/:id", handler.UpdatePaciente)
			// restricted.DELETE("/pacientes/:id", handler.DeletePaciente)
		}
	}
}
