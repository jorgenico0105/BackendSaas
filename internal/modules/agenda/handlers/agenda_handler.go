package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"saas-medico/internal/modules/agenda/models"
	"saas-medico/internal/modules/agenda/services"
	"saas-medico/internal/shared/responses"
)

type AgendaHandler struct {
	svc *services.AgendaService
}

func NewAgendaHandler(svc *services.AgendaService) *AgendaHandler {
	return &AgendaHandler{svc: svc}
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

// ── Citas ─────────────────────────────────────────────────────────────────────

func (h *AgendaHandler) ListCitas(c *gin.Context) {
	medicoID := queryUint(c, "medico_id")
	clinicaID := queryUint(c, "clinica_id")
	fecha := c.Query("fecha")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := h.svc.ListCitas(medicoID, clinicaID, fecha, page, size)
	if err != nil {
		responses.InternalError(c, "Error al obtener citas")
		return
	}
	responses.Paginated(c, list, page, size, total)
}

func (h *AgendaHandler) CreateCita(c *gin.Context) {
	var req models.CreateCitaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	cita, err := h.svc.CreateCita(req)
	if err != nil {
		responses.InternalError(c, "Error al crear cita: "+err.Error())
		return
	}
	responses.Created(c, "Cita creada exitosamente", cita)
}

func (h *AgendaHandler) GetCita(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	cita, err := h.svc.GetCita(id)
	if err != nil {
		responses.NotFound(c, "Cita no encontrada")
		return
	}
	responses.Success(c, "Cita obtenida", cita)
}

func (h *AgendaHandler) UpdateEstadoCita(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateEstadoCitaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	cita, err := h.svc.UpdateEstadoCita(id, req.EstadoCodigo)
	if err != nil {
		if errors.Is(err, services.ErrCitaNotFound) {
			responses.NotFound(c, "Cita no encontrada")
		} else {
			responses.BadRequest(c, err.Error())
		}
		return
	}
	responses.Success(c, "Estado actualizado", cita)
}

func (h *AgendaHandler) UpdateCitaPaciente(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateCitaPacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	cita, err := h.svc.UpdateCitaPaciente(id, req.PacienteID)
	if err != nil {
		if errors.Is(err, services.ErrCitaNotFound) {
			responses.NotFound(c, "Cita no encontrada")
		} else {
			responses.InternalError(c, "Error al actualizar cita")
		}
		return
	}
	responses.Success(c, "Cita actualizada", cita)
}

func (h *AgendaHandler) DeleteCita(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteCita(id); err != nil {
		responses.NotFound(c, "Cita no encontrada")
		return
	}
	responses.Success(c, "Cita eliminada", nil)
}

// ── Sesiones ──────────────────────────────────────────────────────────────────

func (h *AgendaHandler) CreateSesion(c *gin.Context) {
	citaID, ok := paramUint(c, "citaID")
	if !ok {
		return
	}
	var req models.CreateSesionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	sesion, err := h.svc.CreateSesion(citaID, req)
	if err != nil {
		if errors.Is(err, services.ErrCitaNotFound) {
			responses.NotFound(c, "Cita no encontrada")
		} else if errors.Is(err, services.ErrSesionExiste) {
			responses.BadRequest(c, "Ya existe una sesión para esta cita")
		} else {
			responses.InternalError(c, "Error al crear sesión")
		}
		return
	}
	responses.Created(c, "Sesión creada exitosamente", sesion)
}

func (h *AgendaHandler) GetSesion(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	sesion, err := h.svc.GetSesion(id)
	if err != nil {
		responses.NotFound(c, "Sesión no encontrada")
		return
	}
	responses.Success(c, "Sesión obtenida", sesion)
}

func (h *AgendaHandler) UpdateSesion(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateSesionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	sesion, err := h.svc.UpdateSesion(id, req)
	if err != nil {
		if errors.Is(err, services.ErrSesionNotFound) {
			responses.NotFound(c, "Sesión no encontrada")
		} else {
			responses.InternalError(c, "Error al actualizar sesión")
		}
		return
	}
	responses.Success(c, "Sesión actualizada", sesion)
}

// ── Horarios ──────────────────────────────────────────────────────────────────

func (h *AgendaHandler) ListHorarios(c *gin.Context) {
	medicoID := queryUint(c, "medico_id")
	list, err := h.svc.ListHorarios(medicoID)
	if err != nil {
		responses.InternalError(c, "Error al obtener horarios")
		return
	}
	responses.Success(c, "Horarios obtenidos", list)
}

func (h *AgendaHandler) CreateHorario(c *gin.Context) {
	var req models.CreateHorarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	h2, err := h.svc.CreateHorario(req)
	if err != nil {
		responses.InternalError(c, "Error al crear horario")
		return
	}
	responses.Created(c, "Horario creado", h2)
}

func (h *AgendaHandler) DeleteHorario(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteHorario(id); err != nil {
		responses.NotFound(c, "Horario no encontrado")
		return
	}
	responses.Success(c, "Horario eliminado", nil)
}

// ── Bloqueos ──────────────────────────────────────────────────────────────────

func (h *AgendaHandler) ListBloqueos(c *gin.Context) {
	medicoID := queryUint(c, "medico_id")
	list, err := h.svc.ListBloqueos(medicoID)
	if err != nil {
		responses.InternalError(c, "Error al obtener bloqueos")
		return
	}
	responses.Success(c, "Bloqueos obtenidos", list)
}

func (h *AgendaHandler) CreateBloqueo(c *gin.Context) {
	var req models.CreateBloqueoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	b, err := h.svc.CreateBloqueo(req)
	if err != nil {
		responses.InternalError(c, "Error al crear bloqueo")
		return
	}
	responses.Created(c, "Bloqueo creado", b)
}

func (h *AgendaHandler) DeleteBloqueo(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteBloqueo(id); err != nil {
		responses.InternalError(c, "Error al eliminar bloqueo")
		return
	}
	responses.Success(c, "Bloqueo eliminado", nil)
}

// ── Catálogos ─────────────────────────────────────────────────────────────────

func (h *AgendaHandler) ListTiposCita(c *gin.Context) {
	// Staff: use rolID from JWT; fallback to ?rol_id query param (used by patient app)
	rolID := c.GetUint("rolID")
	if rolID == 0 {
		if v, err := strconv.ParseUint(c.Query("rol_id"), 10, 64); err == nil {
			rolID = uint(v)
		}
	}
	list, err := h.svc.ListTiposCita(rolID)
	if err != nil {
		responses.InternalError(c, "Error al obtener tipos de cita")
		return
	}
	responses.Success(c, "Tipos de cita obtenidos", list)
}

func (h *AgendaHandler) ListEstadosCita(c *gin.Context) {
	list, err := h.svc.ListEstadosCita()
	if err != nil {
		responses.InternalError(c, "Error al obtener estados de cita")
		return
	}
	responses.Success(c, "Estados de cita obtenidos", list)
}
