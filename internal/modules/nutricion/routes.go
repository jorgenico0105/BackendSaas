package nutricion

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/nutricion/handlers"
	"saas-medico/internal/modules/nutricion/repositories"
	"saas-medico/internal/modules/nutricion/services"

	"saas-medico/internal/shared/openia"
	"saas-medico/internal/shared/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, redis *redis.Client, openiaService *openia.OpenIaService) {
	repo := repositories.NewNutricionRepository(database.GetDB())
	service := services.NewNutricionService(repo, redis)
	h := handlers.NewNutricionHandler(service, openiaService)

	scheduler.StartCron(service.DeactivateOldMenus)

	n := router.Group("/nutricion")

	n.Use(authMiddleware.RequireAuth())
	{
		n.GET("/generate-menu-pdf/:menuID", h.CreateMenuReport)
		// ─── Catálogos ─────────────────────────────────────────────────
		n.GET("/grupos-alimento", h.ListGruposAlimento)
		n.GET("/grupos-alimento/:id", h.GetGrupoAlimento)
		n.POST("/grupos-alimento", h.CreateGrupoAlimento)
		n.PUT("/grupos-alimento/:id", h.UpdateGrupoAlimento)
		n.DELETE("/grupos-alimento/:id", h.DeleteGrupoAlimento)

		n.GET("/alimentos", h.ListAlimentos)
		n.GET("/alimentos/:id", h.GetAlimento)
		n.POST("/alimentos", h.CreateAlimento)
		n.PUT("/alimentos/:id", h.UpdateAlimento)

		n.GET("/dietas-catalogo", h.ListDietasCatalogo)

		n.GET("/ejercicios-catalogo", h.ListEjerciciosCatalogo)
		n.POST("/ejercicios-catalogo", h.CreateEjercicioCatalogo)

		n.GET("/logros-catalogo", h.ListLogrosCatalogo)

		// ─── Catálogo de tipos de recurso ──────────────────────────────
		n.GET("/tipo-recursos", h.ListTipoRecursos)
		n.POST("/tipo-recursos", h.CreateTipoRecurso)
		n.PUT("/tipo-recursos/:id", h.UpdateTipoRecurso)
		n.DELETE("/tipo-recursos/:id", h.DeleteTipoRecurso)

		// ─── Especial: menús que requieren cambio mañana ───────────────
		n.GET("/menus/requieren-cambio", h.ListDietasRequierenCambio)

		// ─── Plantillas de menú semanal ────────────────────────────────
		// Filtros: ?num_comidas=5  ?semana_numero=1
		n.GET("/plantillas-menu", h.ListPlantillas)
		n.POST("/plantillas-menu", h.CreatePlantillaSemana)
		n.GET("/plantillas-menu/:plantillaId", h.GetPlantillaSemana)
		n.PUT("/plantillas-menu/:plantillaId", h.UpdatePlantillaSemana)
		n.DELETE("/plantillas-menu/:plantillaId", h.DeletePlantillaSemana)

		// Detalles (día+comida) de una plantilla
		n.GET("/plantillas-menu/:plantillaId/detalles", h.GetDetallesPlantilla)
		n.POST("/plantillas-menu/:plantillaId/detalles", h.AddDetallePlantilla)
		n.PUT("/plantillas-menu-detalles/:detalleId", h.UpdateDetallePlantilla)
		n.DELETE("/plantillas-menu-detalles/:detalleId", h.DeleteDetallePlantilla)

		// Alimentos de un detalle de plantilla
		n.GET("/plantillas-menu-detalles/:detalleId/alimentos", h.GetAlimentosPlantillaDetalle)
		n.POST("/plantillas-menu-detalles/:detalleId/alimentos", h.AddAlimentoPlantillaDetalle)
		n.PUT("/plantillas-menu-alimentos/:id", h.UpdateAlimentoPlantillaDetalle)
		n.DELETE("/plantillas-menu-alimentos/:id", h.DeleteAlimentoPlantillaDetalle)

		// ─── Archivos PDF (biblioteca de la clínica/paciente) ──────────
		n.GET("/archivos-pdf", h.ListArchivosPDF)
		n.GET("/archivos-pdf/:pacienteID", h.ListArchivosPDFByPaciente)
		n.POST("/archivos-pdf", h.CreateArchivoPDF)
		n.DELETE("/archivos-pdf/:id", h.DeleteArchivoPDF)

		// Cálculo de fórmulas nutricionales (IMC, ICC, Harris-Benedict) — sin persistencia
		n.POST("/formulas", h.CalcularFormulas)

		// ─── Por paciente ──────────────────────────────────────────────
		pac := n.Group("/pacientes/:pacienteId")
		{
			// Plan de dieta
			pac.GET("/dietas", h.ListDietasByPaciente)
			pac.POST("/dietas", h.CreateDieta)
			pac.GET("/dietas/:dietaId", h.GetDieta)
			pac.PUT("/dietas/:dietaId", h.UpdateDieta)

			// Menús semanales de una dieta
			pac.GET("/dietas/:dietaId/menus", h.ListMenusByDieta)
			pac.POST("/dietas/:dietaId/menus", h.CreateMenu)
			pac.GET("/menus/:menuId", h.GetMenu)
			pac.GET("/menus/:menuId/detalles", h.GetDetallesMenu)
			pac.POST("/menus/:menuId/detalles", h.AddDetalleMenu)

			// Alimentos de un detalle de menú
			pac.GET("/menu-detalles/:detalleId/alimentos", h.GetAlimentosMenuDetalle)
			pac.POST("/menu-detalles/:detalleId/alimentos", h.AddAlimentoMenuDetalle)
			pac.PATCH("/menu-detalles/:detalleId", h.UpdateDetalleMenu)
			pac.DELETE("/menu-detalles/:detalleId/alimentos/:id", h.DeleteAlimentoMenuDetalle)
			pac.PUT("/menu-detalles/:detalleId/alimentos/:id", h.UpdateAlimentoMenuDetalle)

			// R24H — Recordatorio 24 horas
			pac.GET("/r24h", h.ListR24H)
			pac.POST("/r24h", h.CreateR24H)
			pac.GET("/r24h/:r24hId/items", h.ListR24HItems)
			pac.POST("/r24h/:r24hId/items", h.AddR24HItem)

			// Preferencias alimentarias
			pac.GET("/preferencias", h.ListPreferencias)
			pac.POST("/preferencias", h.AddPreferencia)
			pac.DELETE("/preferencias/:id", h.DeletePreferencia)

			// Síntomas
			pac.GET("/sintomas", h.ListSintomas)
			pac.POST("/sintomas", h.CreateSintoma)

			// Ejercicios asignados
			pac.GET("/ejercicios", h.ListEjerciciosByPaciente)
			pac.POST("/ejercicios", h.AddEjercicioPaciente)

			// Registros diarios
			pac.GET("/registros-comida", h.ListRegistrosComida)
			pac.POST("/registros-comida", h.CreateRegistroComida)
			pac.POST("/registros-comida/:registroId/alimentos", h.AddRegistroAlimento)
			pac.POST("/registros-comida/:registroId/foto", h.UploadFotoComida)
			pac.PATCH("/registros-comida/:registroId/consumir", h.MarcarRegistroComidaConsumida)
			pac.POST("/registros-comida/fuera-plan")

			pac.GET("/registros-ejercicio", h.ListRegistrosEjercicio)
			pac.POST("/registros-ejercicio", h.CreateRegistroEjercicio)

			pac.GET("/resumen-diario", h.GetResumenDiario)

			pac.GET("/progreso", h.ListProgreso)
			pac.POST("/progreso", h.AddProgreso)

			// XP y logros
			pac.GET("/xp", h.GetXP)
			pac.GET("/logros", h.ListLogros)
			pac.POST("/ask-ia", h.ChatWhitIa)
		}
	}
}
