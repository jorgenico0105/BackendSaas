package historia

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	authModels "saas-medico/internal/modules/auth/models"
	"saas-medico/internal/modules/historia/handlers"
	"saas-medico/internal/modules/historia/repositories"
	"saas-medico/internal/modules/historia/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewHistoriaRepository(database.GetDB())
	service := services.NewHistoriaService(repo)
	handler := handlers.NewHistoriaHandler(service)

	historia := router.Group("/historia")
	historia.Use(authMiddleware.RequireAuth())
	historia.Use(authMiddleware.RequireRoles(authModels.RolAdmin, authModels.RolNutriologo, authModels.RolPsicologo, authModels.RolSuperAdmin))
	{
		// Catálogos
		historia.GET("/tipos-formulario", handler.ListTiposFormulario)

		// Formularios CRUD
		historia.GET("/formularios", handler.ListFormularios)
		historia.POST("/formularios", handler.CreateFormulario)
		historia.GET("/formularios/:id", handler.GetFormulario)
		historia.PUT("/formularios/en l:id", handler.UpdateFormulario)
		historia.DELETE("/formularios/:id", handler.DeleteFormulario)
		historia.GET("/historia-clinica-form", handler.GetHistoriaClinicaByUser)

		// Por paciente
		paciente := historia.Group("/pacientes/:pacienteId")
		{
			paciente.POST("/historia", handler.CreateHistoriaPaciente)
			paciente.GET("/historias", handler.ListHistoriasByPaciente)
			paciente.GET("/historias-respuestas", handler.GetHistoriasByPaciente)

			paciente.GET("/alergias", handler.ListAlergias)
			paciente.GET("/antecedentes", handler.ListAntecedentes)
			paciente.GET("/habitos", handler.ListHabitos)
			paciente.GET("/diagnosticos", handler.ListDiagnosticos)

			paciente.POST("/imagenes", handler.UploadPacienteImagen)
			paciente.GET("/imagenes", handler.ListPacienteImagenes)
		}
	}
}
