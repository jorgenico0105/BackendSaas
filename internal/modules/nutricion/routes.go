package nutricion

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/nutricion/handlers"
	"saas-medico/internal/modules/nutricion/repositories"
	"saas-medico/internal/modules/nutricion/services"

	"saas-medico/internal/shared/openia"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, redis *redis.Client, openiaService *openia.OpenIaService) {
	repo := repositories.NewNutricionRepository(database.GetDB())
	service := services.NewNutricionService(repo, redis)
	h := handlers.NewNutricionHandler(service, openiaService)

	n := router.Group("/nutricion")
	n.Use(authMiddleware.RequireAuth())
	{
		// ─── Catálogos ─────────────────────────────────────────────────
		n.GET("/alimentos", h.ListAlimentos)
		n.GET("/alimentos/:id", h.GetAlimento)
		n.POST("/alimentos", h.CreateAlimento)

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

		// ─── Archivos PDF (biblioteca de la clínica/paciente) ──────────
		n.GET("/archivos-pdf", h.ListArchivosPDF)
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

			// Generar dierta (!! importante)

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

			pac.GET("/registros-ejercicio", h.ListRegistrosEjercicio)
			pac.POST("/registros-ejercicio", h.CreateRegistroEjercicio)

			// Resumen diario unificado (comidas + ejercicios + progreso) para móvil y web
			pac.GET("/resumen-diario", h.GetResumenDiario)

			// Progreso (historial de peso)
			pac.GET("/progreso", h.ListProgreso)
			pac.POST("/progreso", h.AddProgreso)

			// XP y logros
			pac.GET("/xp", h.GetXP)
			pac.GET("/logros", h.ListLogros)

			pac.POST("/ask-ia", h.ChatWhitIa)
		}
	}
}
