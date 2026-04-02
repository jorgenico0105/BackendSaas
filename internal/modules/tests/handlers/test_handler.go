package handlers

import (
	"strconv"

	"saas-medico/internal/modules/tests/services"
	"saas-medico/internal/shared/responses"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	service *services.TestService
}

func NewTestHandler(service *services.TestService) *TestHandler {
	return &TestHandler{service: service}
}

func paramUint(c *gin.Context, key string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return 0, false
	}
	return uint(v), true
}

func (h *TestHandler) ListReglas(c *gin.Context) {
	formularioID, ok := paramUint(c, "formularioId")
	if !ok {
		return
	}
	list, err := h.service.ListReglas(formularioID)
	if err != nil {
		responses.InternalError(c, "Error al listar reglas")
		return
	}
	responses.Success(c, "Reglas", list)
}

func (h *TestHandler) ListTestsByPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListTestsByPaciente(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar tests")
		return
	}
	responses.Success(c, "Tests", list)
}

func (h *TestHandler) GetTest(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	t, respuestas, err := h.service.GetTest(id)
	if err != nil {
		responses.NotFound(c, "Test no encontrado")
		return
	}
	responses.Success(c, "Test", gin.H{"test": t, "respuestas": respuestas})
}

func (h *TestHandler) ListTestsBySesion(c *gin.Context) {
	sesionID, ok := paramUint(c, "sesionId")
	if !ok {
		return
	}
	list, err := h.service.GetTestsBySesion(sesionID)
	if err != nil {
		responses.InternalError(c, "Error al listar tests de sesión")
		return
	}
	responses.Success(c, "Tests de sesión", list)
}
