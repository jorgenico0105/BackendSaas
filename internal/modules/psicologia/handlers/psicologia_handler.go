package handlers

import (
	"saas-medico/internal/modules/psicologia/services"
	"saas-medico/internal/shared/responses"

	"github.com/gin-gonic/gin"
)

type PsicologiaHandler struct {
	service *services.PsicologiaService
}

func NewPsicologiaHandler(service *services.PsicologiaService) *PsicologiaHandler {
	return &PsicologiaHandler{service: service}
}

func (h *PsicologiaHandler) Ping(c *gin.Context) {
	result := h.service.Ping()
	responses.Success(c, result, nil)
}

// Aquí irán los handlers específicos de psicología
// Ejemplo:
// func (h *PsicologiaHandler) GetPaciente(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := utils.ParseUint(idStr)
// 	if err != nil {
// 		responses.BadRequest(c, "ID inválido")
// 		return
// 	}
// 	paciente, err := h.service.GetPaciente(id)
// 	if err != nil {
// 		responses.NotFound(c, "Paciente no encontrado")
// 		return
// 	}
// 	responses.Success(c, "Paciente encontrado", paciente)
// }
