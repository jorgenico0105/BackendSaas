package handlers

import (
	"saas-medico/internal/modules/odontologia/services"
	"saas-medico/internal/shared/responses"

	"github.com/gin-gonic/gin"
)

type OdontologiaHandler struct {
	service *services.OdontologiaService
}

func NewOdontologiaHandler(service *services.OdontologiaService) *OdontologiaHandler {
	return &OdontologiaHandler{service: service}
}

func (h *OdontologiaHandler) Ping(c *gin.Context) {
	result := h.service.Ping()
	responses.Success(c, result, nil)
}

// Aquí irán los handlers específicos de odontología
// Ejemplo:
// func (h *OdontologiaHandler) GetHistorial(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := utils.ParseUint(idStr)
// 	if err != nil {
// 		responses.BadRequest(c, "ID inválido")
// 		return
// 	}
// 	historial, err := h.service.GetHistorial(id)
// 	if err != nil {
// 		responses.NotFound(c, "Historial no encontrado")
// 		return
// 	}
// 	responses.Success(c, "Historial encontrado", historial)
// }
