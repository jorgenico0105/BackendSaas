package handlers

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"saas-medico/internal/modules/pacientes/models"
	"saas-medico/internal/modules/pacientes/services"
	"saas-medico/internal/shared/responses"
)

type PacienteHandler struct {
	svc *services.PacienteService
}

func NewPacienteHandler(svc *services.PacienteService) *PacienteHandler {
	return &PacienteHandler{svc: svc}
}

func paramUint(c *gin.Context, key string) (uint, bool) {
	n, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return 0, false
	}
	return uint(n), true
}

func (h *PacienteHandler) List(c *gin.Context) {
	search := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	clinicaID := int(c.GetUint("clinicaID"))
	userID := int(c.GetUint("userID"))
	list, total, err := h.svc.List(search, page, size, clinicaID, userID)
	if err != nil {
		responses.InternalError(c, "Error al obtener pacientes")
		return
	}
	responses.Paginated(c, list, page, size, total)
}

func (h *PacienteHandler) Create(c *gin.Context) {
	var req models.CreatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	userID := c.GetUint("userID")
	clinicaID := c.GetUint("clinicaID")
	p, err := h.svc.Create(req, userID, clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al crear paciente")
		return
	}
	responses.Created(c, "Paciente creado exitosamente", p)
}

func (h *PacienteHandler) Get(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	p, err := h.svc.GetByID(id)
	if err != nil {
		responses.NotFound(c, "Paciente no encontrado")
		return
	}
	responses.Success(c, "Paciente obtenido", p)
}

func (h *PacienteHandler) Update(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	p, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrPacienteNotFound) {
			responses.NotFound(c, "Paciente no encontrado")
		} else {
			responses.InternalError(c, "Error al actualizar paciente")
		}
		return
	}
	responses.Success(c, "Paciente actualizado", p)
}

func (h *PacienteHandler) Delete(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		responses.NotFound(c, "Paciente no encontrado")
		return
	}
	responses.Success(c, "Paciente eliminado", nil)
}

// ── PrePacientes ──────────────────────────────────────────────────────────────

func (h *PacienteHandler) ListPrePacientes(c *gin.Context) {
	clinicaID, ok := paramUint(c, "clinicaID")
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	list, total, err := h.svc.ListPrePacientes(clinicaID, page, size)
	if err != nil {
		responses.InternalError(c, "Error al obtener pre-pacientes")
		return
	}
	responses.Paginated(c, list, page, size, total)
}

func (h *PacienteHandler) CreatePrePaciente(c *gin.Context) {
	var req models.CreatePrePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	pp, err := h.svc.CreatePrePaciente(req)
	if err != nil {
		responses.InternalError(c, "Error al crear pre-paciente")
		return
	}
	responses.Created(c, "Pre-paciente creado exitosamente", pp)
}

func (h *PacienteHandler) DeletePrePaciente(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeletePrePaciente(id); err != nil {
		responses.NotFound(c, "Pre-paciente no encontrado")
		return
	}
	responses.Success(c, "Pre-paciente eliminado", nil)
}

// ── Aplicaciones ──────────────────────────────────────────────────────────────

func (h *PacienteHandler) ListAplicaciones(c *gin.Context) {
	log.Println(c.GetUint("clinicaID"))
	clinicaID := c.GetUint("clinicaID")

	list, err := h.svc.ListAplicaciones(clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al listar aplicaciones")
		return
	}
	responses.Success(c, "Aplicaciones disponibles", list)
}

func (h *PacienteHandler) CreateAplicacion(c *gin.Context) {
	var req models.CreateAplicacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	clinicaID := c.GetUint("clinicaID")
	a, err := h.svc.CreateAplicacion(req, clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al crear aplicación")
		return
	}
	responses.Created(c, "Aplicación creada", a)
}

func (h *PacienteHandler) ListAplicacionesPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	clinicaID := c.GetUint("clinicaID")
	list, err := h.svc.ListAplicacionesPaciente(pacienteID, clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al listar aplicaciones del paciente")
		return
	}
	responses.Success(c, "Aplicaciones del paciente", list)
}

func (h *PacienteHandler) AsignarAplicacion(c *gin.Context) {
	pacienteID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.AsignarAplicacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	clinicaID := c.GetUint("clinicaID")
	userID := c.GetUint("userID")
	pa, err := h.svc.AsignarAplicacion(pacienteID, req, clinicaID, userID)
	if err != nil {
		if errors.Is(err, services.ErrAplicacionYaAsignada) {
			responses.BadRequest(c, "El paciente ya tiene acceso a esta aplicación")
		} else {
			responses.InternalError(c, "Error al asignar aplicación")
		}
		return
	}
	responses.Created(c, "Aplicación asignada al paciente", pa)
}

func (h *PacienteHandler) RevocarAplicacion(c *gin.Context) {
	pacienteID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	aplicacionID, ok := paramUint(c, "aplicacionId")
	if !ok {
		return
	}
	clinicaID := c.GetUint("clinicaID")
	if err := h.svc.RevocarAplicacion(pacienteID, aplicacionID, clinicaID); err != nil {
		responses.NotFound(c, "Acceso no encontrado")
		return
	}
	responses.Success(c, "Acceso revocado", nil)
}

// ── PacienteUsuario ───────────────────────────────────────────────────────────

func (h *PacienteHandler) LoginPaciente(c *gin.Context) {

	var req models.PacienteLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	result, err := h.svc.LoginPaciente(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPacienteCredenciales):
			responses.Unauthorized(c, "Credenciales inválidas")
		case errors.Is(err, services.ErrSinAccesoAplicacion):
			responses.Forbidden(c, "No tienes acceso a esta aplicación")
		default:
			responses.InternalError(c, "Error al iniciar sesión")
		}
		return
	}
	responses.Success(c, "Login exitoso", result)
}

// GetAccesoStats — frecuencia de uso del app para el dashboard web
func (h *PacienteHandler) GetAccesoStats(c *gin.Context) {
	pacienteID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	desde := c.Query("desde")
	hasta := c.Query("hasta")

	count  := h.svc.CountAccesos(pacienteID, desde, hasta)
	ultimo := h.svc.UltimoAcceso(pacienteID)

	responses.Success(c, "Estadísticas de acceso", map[string]any{
		"paciente_id":   pacienteID,
		"total_accesos": count,
		"ultimo_acceso": ultimo,
	})
}

func (h *PacienteHandler) CreatePacienteUsuario(c *gin.Context) {
	var req models.CreatePacienteUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	clinicaID := c.GetUint("clinicaID")
	pu, err := h.svc.CreatePacienteUsuario(req, clinicaID)
	if err != nil {
		if errors.Is(err, services.ErrPacienteUsuarioExists) {
			responses.BadRequest(c, "El paciente ya tiene usuario en esta clínica")
		} else {
			responses.InternalError(c, "Error al crear usuario de paciente")
		}
		return
	}
	responses.Created(c, "Usuario de paciente creado exitosamente", pu)
}
