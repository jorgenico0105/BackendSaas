package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"saas-medico/internal/modules/admin/models"
	"saas-medico/internal/modules/admin/services"
	"saas-medico/internal/shared/responses"
)

type AdminHandler struct {
	svc *services.AdminService
}

func NewAdminHandler(svc *services.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func paramUint(c *gin.Context, key string) (uint, bool) {
	n, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return 0, false
	}
	return uint(n), true
}

func paginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}

// ── Clinicas ──────────────────────────────────────────────────────────────────

func (h *AdminHandler) ListClinicas(c *gin.Context) {
	page, size := paginationParams(c)
	list, total, err := h.svc.ListClinicas(page, size)
	if err != nil {
		responses.InternalError(c, "Error al obtener clínicas")
		return
	}
	responses.Paginated(c, list, page, size, total)
}

func (h *AdminHandler) CreateClinica(c *gin.Context) {
	var req models.CreateClinicaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	clinica, err := h.svc.CreateClinica(req)
	if err != nil {
		responses.InternalError(c, "Error al crear clínica")
		return
	}
	responses.Created(c, "Clínica creada exitosamente", clinica)
}

func (h *AdminHandler) GetClinica(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	clinica, err := h.svc.GetClinica(id)
	if err != nil {
		responses.NotFound(c, "Clínica no encontrada")
		return
	}
	responses.Success(c, "Clínica obtenida", clinica)
}

func (h *AdminHandler) UpdateClinica(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.CreateClinicaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	clinica, err := h.svc.UpdateClinica(id, req)
	if err != nil {
		if errors.Is(err, services.ErrClinicaNotFound) {
			responses.NotFound(c, "Clínica no encontrada")
		} else {
			responses.InternalError(c, "Error al actualizar clínica")
		}
		return
	}
	responses.Success(c, "Clínica actualizada", clinica)
}

func (h *AdminHandler) DeleteClinica(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteClinica(id); err != nil {
		responses.NotFound(c, "Clínica no encontrada")
		return
	}
	responses.Success(c, "Clínica eliminada", nil)
}

// ── Consultorios ──────────────────────────────────────────────────────────────

func (h *AdminHandler) ListConsultorios(c *gin.Context) {
	clinicaID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	list, err := h.svc.ListConsultorios(clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al obtener consultorios")
		return
	}
	responses.Success(c, "Consultorios obtenidos", list)
}

func (h *AdminHandler) CreateConsultorio(c *gin.Context) {
	clinicaID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.CreateConsultorioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	con, err := h.svc.CreateConsultorio(clinicaID, req)
	if err != nil {
		if errors.Is(err, services.ErrClinicaNotFound) {
			responses.NotFound(c, "Clínica no encontrada")
		} else {
			responses.InternalError(c, "Error al crear consultorio")
		}
		return
	}
	responses.Created(c, "Consultorio creado exitosamente", con)
}

func (h *AdminHandler) GetConsultorio(c *gin.Context) {
	id, ok := paramUint(c, "consultorioID")
	if !ok {
		return
	}
	con, err := h.svc.GetConsultorio(id)
	if err != nil {
		responses.NotFound(c, "Consultorio no encontrado")
		return
	}
	responses.Success(c, "Consultorio obtenido", con)
}

func (h *AdminHandler) UpdateConsultorio(c *gin.Context) {
	id, ok := paramUint(c, "consultorioID")
	if !ok {
		return
	}
	var req models.CreateConsultorioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	con, err := h.svc.UpdateConsultorio(id, req)
	if err != nil {
		if errors.Is(err, services.ErrConsultorioNotFound) {
			responses.NotFound(c, "Consultorio no encontrado")
		} else {
			responses.InternalError(c, "Error al actualizar consultorio")
		}
		return
	}
	responses.Success(c, "Consultorio actualizado", con)
}

func (h *AdminHandler) DeleteConsultorio(c *gin.Context) {
	id, ok := paramUint(c, "consultorioID")
	if !ok {
		return
	}
	if err := h.svc.DeleteConsultorio(id); err != nil {
		responses.NotFound(c, "Consultorio no encontrado")
		return
	}
	responses.Success(c, "Consultorio eliminado", nil)
}

// ── UsuarioConsultorio ────────────────────────────────────────────────────────

func (h *AdminHandler) ListUsuariosConsultorio(c *gin.Context) {
	consultorioID, ok := paramUint(c, "consultorioID")
	if !ok {
		return
	}
	list, err := h.svc.ListUsuariosByConsultorio(consultorioID)
	if err != nil {
		responses.InternalError(c, "Error al obtener usuarios del consultorio")
		return
	}
	responses.Success(c, "Usuarios del consultorio", list)
}

func (h *AdminHandler) AsignarUsuarioConsultorio(c *gin.Context) {
	consultorioID, ok := paramUint(c, "consultorioID")
	if !ok {
		return
	}
	var req models.AsignarUsuarioConsultorioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	if err := h.svc.AsignarUsuarioAConsultorio(consultorioID, req.UsuarioID); err != nil {
		responses.InternalError(c, "Error al asignar usuario al consultorio")
		return
	}
	responses.Success(c, "Usuario asignado al consultorio", nil)
}

func (h *AdminHandler) RemoverUsuarioConsultorio(c *gin.Context) {
	consultorioID, ok := paramUint(c, "consultorioID")
	if !ok {
		return
	}
	usuarioID, ok := paramUint(c, "usuarioID")
	if !ok {
		return
	}
	if err := h.svc.RemoverUsuarioDeConsultorio(consultorioID, usuarioID); err != nil {
		responses.InternalError(c, "Error al remover usuario del consultorio")
		return
	}
	responses.Success(c, "Usuario removido del consultorio", nil)
}

// ── Profesiones ───────────────────────────────────────────────────────────────

func (h *AdminHandler) ListProfesiones(c *gin.Context) {
	list, err := h.svc.ListProfesiones()
	if err != nil {
		responses.InternalError(c, "Error al obtener profesiones")
		return
	}
	responses.Success(c, "Profesiones obtenidas", list)
}

func (h *AdminHandler) CreateProfesion(c *gin.Context) {
	var req models.CreateProfesionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	p, err := h.svc.CreateProfesion(req)
	if err != nil {
		responses.InternalError(c, "Error al crear profesión")
		return
	}
	responses.Created(c, "Profesión creada", p)
}

func (h *AdminHandler) UpdateProfesion(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.CreateProfesionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	p, err := h.svc.UpdateProfesion(id, req)
	if err != nil {
		responses.NotFound(c, "Profesión no encontrada")
		return
	}
	responses.Success(c, "Profesión actualizada", p)
}

func (h *AdminHandler) DeleteProfesion(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteProfesion(id); err != nil {
		responses.NotFound(c, "Profesión no encontrada")
		return
	}
	responses.Success(c, "Profesión eliminada", nil)
}

// ── Planes SaaS ───────────────────────────────────────────────────────────────

func (h *AdminHandler) ListPlanes(c *gin.Context) {
	list, err := h.svc.ListPlanes()
	if err != nil {
		responses.InternalError(c, "Error al obtener planes")
		return
	}
	responses.Success(c, "Planes obtenidos", list)
}

func (h *AdminHandler) CreatePlan(c *gin.Context) {
	var req models.CreatePlanSaasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	p, err := h.svc.CreatePlan(req)
	if err != nil {
		responses.InternalError(c, "Error al crear plan")
		return
	}
	responses.Created(c, "Plan creado", p)
}

// ── Usuarios por Clínica ──────────────────────────────────────────────────────

func (h *AdminHandler) ListUsuariosClinica(c *gin.Context) {
	clinicaID, ok := paramUint(c, "clinicaID")
	if !ok {
		return
	}
	list, err := h.svc.ListUsuariosClinica(clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al obtener usuarios")
		return
	}
	responses.Success(c, "Usuarios de la clínica", list)
}

func (h *AdminHandler) AsignarUsuario(c *gin.Context) {
	clinicaID, ok := paramUint(c, "clinicaID")
	if !ok {
		return
	}
	var req models.AsignarUsuarioClinicaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	uc, err := h.svc.AsignarUsuario(clinicaID, req)
	if err != nil {
		responses.InternalError(c, "Error al asignar usuario")
		return
	}
	responses.Created(c, "Usuario asignado a la clínica", uc)
}

func (h *AdminHandler) RemoverUsuario(c *gin.Context) {
	clinicaID, ok := paramUint(c, "clinicaID")
	if !ok {
		return
	}
	usuarioID, ok := paramUint(c, "usuarioID")
	if !ok {
		return
	}
	if err := h.svc.RemoverUsuario(clinicaID, usuarioID); err != nil {
		responses.InternalError(c, "Error al remover usuario")
		return
	}
	responses.Success(c, "Usuario removido de la clínica", nil)
}
