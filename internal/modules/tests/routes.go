package tests

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/tests/handlers"
	"saas-medico/internal/modules/tests/repositories"
	"saas-medico/internal/modules/tests/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewTestRepository(database.GetDB())
	service := services.NewTestService(repo)
	handler := handlers.NewTestHandler(service)

	t := router.Group("/tests")
	t.Use(authMiddleware.RequireAuth())
	{
		// Reglas de formularios
		t.GET("/formularios/:formularioId/reglas", handler.ListReglas)

		// Tests por paciente
		t.GET("/pacientes/:pacienteId", handler.ListTestsByPaciente)
		t.GET("/:id", handler.GetTest)

		// Tests de una sesión
		t.GET("/sesiones/:sesionId", handler.ListTestsBySesion)
	}
}
