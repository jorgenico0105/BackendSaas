package cobros

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/cobros/handlers"
	"saas-medico/internal/modules/cobros/repositories"
	"saas-medico/internal/modules/cobros/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewCobroRepository(database.GetDB())
	svc := services.NewCobroService(repo)
	h := handlers.NewCobroHandler(svc)

	g := router.Group("/cobros")
	g.Use(authMiddleware.RequireAuth())
	{
		// Catálogos
		g.GET("/medios-pago", h.ListMediosPago)
		g.GET("/estados-cobro", h.ListEstadosCobro)
		g.GET("/tipos-egreso", h.ListTiposEgreso)

		// Cobros
		g.GET("", h.ListCobrosPaciente)
		g.POST("", h.CreateCobro)
		g.GET("/:id", h.GetCobro)
		g.POST("/:id/pagos", h.RegistrarPago)

		// Egresos
		g.GET("/egresos", h.ListEgresos)
		g.POST("/egresos", h.CreateEgreso)
	}
}
