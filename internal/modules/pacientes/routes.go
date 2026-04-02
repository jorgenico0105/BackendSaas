package pacientes

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	authRoutes "saas-medico/internal/modules/auth"
	authModels "saas-medico/internal/modules/auth/models"
	"saas-medico/internal/modules/pacientes/handlers"
	"saas-medico/internal/modules/pacientes/repositories"
	"saas-medico/internal/modules/pacientes/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewPacienteRepository(database.GetDB())
	svc := services.NewPacienteService(repo, authRoutes.GetJWTService())
	h := handlers.NewPacienteHandler(svc)

	// Ruta pública: login de paciente
	router.POST("/pacientes/login", h.LoginPaciente)

	g := router.Group("/pacientes")
	g.Use(authMiddleware.RequireAuth())
	g.Use(authMiddleware.RequireRoles(authModels.RolMedico,
		authModels.RolNutriologo,
		authModels.RolPsicologo, authModels.RolOdontologo, authModels.RolAdmin))
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.GET("/:id", h.Get)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
		g.POST("/usuario", h.CreatePacienteUsuario)

		// Aplicaciones de la clínica (ruta estática — Gin la resuelve antes que /:id/...)
		g.GET("/aplicaciones", h.ListAplicaciones)
		g.POST("/aplicaciones", h.CreateAplicacion)

		// Aplicaciones del paciente
		g.GET("/:id/aplicaciones", h.ListAplicacionesPaciente)
		g.POST("/:id/aplicaciones", h.AsignarAplicacion)
		g.DELETE("/:id/aplicaciones/:aplicacionId", h.RevocarAplicacion)

		// Métricas de acceso al app (frecuencia de uso)
		g.GET("/:id/acceso-stats", h.GetAccesoStats)
	}

	// Pre-pacientes por clínica
	pre := router.Group("/pre-pacientes")
	pre.Use(authMiddleware.RequireAuth())
	{
		pre.GET("/clinica/:clinicaID", h.ListPrePacientes)
		pre.POST("", h.CreatePrePaciente)
		pre.DELETE("/:id", h.DeletePrePaciente)
	}
}
