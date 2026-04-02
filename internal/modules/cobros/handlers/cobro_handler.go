package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"saas-medico/internal/modules/cobros/models"
	"saas-medico/internal/modules/cobros/services"
	"saas-medico/internal/shared/responses"
)

type CobroHandler struct {
	svc *services.CobroService
}

func NewCobroHandler(svc *services.CobroService) *CobroHandler {
	return &CobroHandler{svc: svc}
}

func paramUint(c *gin.Context, key string) (uint, bool) {
	n, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return 0, false
	}
	return uint(n), true
}

func queryUint(c *gin.Context, key string) uint {
	n, _ := strconv.ParseUint(c.Query(key), 10, 64)
	return uint(n)
}

func (h *CobroHandler) CreateCobro(c *gin.Context) {
	var req models.CreateCobroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	cobro, err := h.svc.CreateCobro(req)
	if err != nil {
		responses.InternalError(c, "Error al crear cobro: "+err.Error())
		return
	}
	responses.Created(c, "Cobro creado exitosamente", cobro)
}

func (h *CobroHandler) GetCobro(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	cobro, err := h.svc.GetCobro(id)
	if err != nil {
		responses.NotFound(c, "Cobro no encontrado")
		return
	}
	responses.Success(c, "Cobro obtenido", cobro)
}

func (h *CobroHandler) ListCobrosPaciente(c *gin.Context) {
	pacienteID := queryUint(c, "paciente_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	list, total, err := h.svc.ListCobrosPaciente(pacienteID, page, size)
	if err != nil {
		responses.InternalError(c, "Error al obtener cobros")
		return
	}
	responses.Paginated(c, list, page, size, total)
}

func (h *CobroHandler) RegistrarPago(c *gin.Context) {
	cobroID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.RegistrarPagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	// pacienteID vendría del cobro, pero lo permitimos desde el body o del cobro mismo
	pacienteID := queryUint(c, "paciente_id")
	pago, err := h.svc.RegistrarPago(cobroID, req, pacienteID)
	if err != nil {
		if errors.Is(err, services.ErrCobroNotFound) {
			responses.NotFound(c, "Cobro no encontrado")
		} else {
			responses.InternalError(c, "Error al registrar pago")
		}
		return
	}
	responses.Created(c, "Pago registrado exitosamente", pago)
}

func (h *CobroHandler) CreateEgreso(c *gin.Context) {
	var req models.CreateEgresoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	e, err := h.svc.CreateEgreso(req)
	if err != nil {
		responses.InternalError(c, "Error al crear egreso")
		return
	}
	responses.Created(c, "Egreso creado exitosamente", e)
}

func (h *CobroHandler) ListEgresos(c *gin.Context) {
	clinicaID := queryUint(c, "clinica_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	list, total, err := h.svc.ListEgresos(clinicaID, page, size)
	if err != nil {
		responses.InternalError(c, "Error al obtener egresos")
		return
	}
	responses.Paginated(c, list, page, size, total)
}

func (h *CobroHandler) ListMediosPago(c *gin.Context) {
	list, _ := h.svc.ListMediosPago()
	responses.Success(c, "Medios de pago", list)
}

func (h *CobroHandler) ListEstadosCobro(c *gin.Context) {
	list, _ := h.svc.ListEstadosCobro()
	responses.Success(c, "Estados de cobro", list)
}

func (h *CobroHandler) ListTiposEgreso(c *gin.Context) {
	list, _ := h.svc.ListTiposEgreso()
	responses.Success(c, "Tipos de egreso", list)
}
