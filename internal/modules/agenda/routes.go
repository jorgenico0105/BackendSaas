package agenda

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/agenda/handlers"
	"saas-medico/internal/modules/agenda/repositories"
	"saas-medico/internal/modules/agenda/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewAgendaRepository(database.GetDB())
	svc := services.NewAgendaService(repo)
	h := handlers.NewAgendaHandler(svc)

	g := router.Group("/agenda")
	g.Use(authMiddleware.RequireAuth())
	{
		// Catálogos
		g.GET("/tipos-cita", h.ListTiposCita)
		g.GET("/estados-cita", h.ListEstadosCita)

		// Citas
		g.GET("/citas", h.ListCitas)
		g.POST("/citas", h.CreateCita)
		g.GET("/citas/:id", h.GetCita)
		g.PUT("/citas/:id/estado", h.UpdateEstadoCita)
		g.PATCH("/citas/:id/paciente", h.UpdateCitaPaciente)
		g.DELETE("/citas/:id", h.DeleteCita)

		// Sesiones
		g.POST("/citas/:citaID/sesion", h.CreateSesion)
		g.GET("/sesiones/:id", h.GetSesion)
		g.PUT("/sesiones/:id", h.UpdateSesion)

		// Horarios del médico
		g.GET("/horarios", h.ListHorarios)
		g.POST("/horarios", h.CreateHorario)
		g.DELETE("/horarios/:id", h.DeleteHorario)

		// Bloqueos de agenda
		g.GET("/bloqueos", h.ListBloqueos)
		g.POST("/bloqueos", h.CreateBloqueo)
		g.DELETE("/bloqueos/:id", h.DeleteBloqueo)
	}
}
