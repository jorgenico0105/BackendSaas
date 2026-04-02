package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	adminModels "saas-medico/internal/modules/admin/models"
	"saas-medico/internal/modules/admin/services"
	authModels "saas-medico/internal/modules/auth/models"
	"saas-medico/internal/shared/responses"
)

type RolHandler struct {
	svc *services.RolService
}

func NewRolHandler(svc *services.RolService) *RolHandler {
	return &RolHandler{svc: svc}
}

func (h *RolHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		responses.InternalError(c, "Error al obtener roles")
		return
	}
	responses.Success(c, "Roles obtenidos", list)
}

func (h *RolHandler) Get(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	rol, err := h.svc.GetByID(id)
	if err != nil {
		responses.NotFound(c, "Rol no encontrado")
		return
	}
	responses.Success(c, "Rol obtenido", rol)
}

func (h *RolHandler) Create(c *gin.Context) {
	var req authModels.CreateRolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	rol, err := h.svc.Create(req)
	if err != nil {
		responses.InternalError(c, "Error al crear rol")
		return
	}
	responses.Created(c, "Rol creado exitosamente", rol)
}

func (h *RolHandler) Update(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req authModels.CreateRolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	rol, err := h.svc.Update(id, req)
	if err != nil {
		responses.NotFound(c, "Rol no encontrado")
		return
	}
	responses.Success(c, "Rol actualizado", rol)
}

// func (h *RolHandler) Delete(c *gin.Context) {
// 	id, ok := paramUint(c, "id")
// 	if !ok {
// 		return
// 	}
// 	if err := h.svc.Delete(id); err != nil {
// 		responses.NotFound(c, "Rol no encontrado")
// 		return
// 	}
// 	responses.Success(c, "Rol eliminado", nil)
// }

// ─── RolTransaccion ───────────────────────────────────────────────────────────

func (h *RolHandler) ListTransacciones(c *gin.Context) {
	list, err := h.svc.ListTransacciones(nil)
	if err != nil {
		responses.InternalError(c, "Error al obtener transacciones")
		return
	}
	responses.Success(c, "Transacciones", list)
}

func (h *RolHandler) ListTransaccionesByRol(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	list, err := h.svc.ListTransaccionesByRol(id)
	if err != nil {
		responses.InternalError(c, "Error al obtener transacciones del rol")
		return
	}
	responses.Success(c, "Transacciones del rol", list)
}

func (h *RolHandler) AsignarTransacciones(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req adminModels.AsignarTransaccionesRolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	if err := h.svc.AsignarTransacciones(id, req.TransaccionIDs); err != nil {
		responses.InternalError(c, "Error al asignar transacciones")
		return
	}
	responses.Success(c, "Transacciones asignadas", nil)
}

func (h *RolHandler) RevocarTransaccion(c *gin.Context) {
	rolID, ok := paramUint(c, "id")
	if !ok {
		return
	}
	transaccionID, err := strconv.ParseUint(c.Param("transaccionId"), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID de transacción inválido")
		return
	}
	if err := h.svc.RevocarTransaccion(rolID, uint(transaccionID)); err != nil {
		responses.InternalError(c, "Error al revocar transacción")
		return
	}
	responses.Success(c, "Transacción revocada", nil)
}

// ─── UsuarioRol ───────────────────────────────────────────────────────────────

func (h *RolHandler) ListRolesByUsuario(c *gin.Context) {
	usuarioID, ok := paramUint(c, "usuarioId")
	if !ok {
		return
	}
	clinicaID, err := strconv.ParseUint(c.Query("clinica_id"), 10, 64)
	if err != nil || clinicaID == 0 {
		responses.BadRequest(c, "Se requiere clinica_id")
		return
	}
	list, err := h.svc.ListRolesByUsuario(usuarioID, uint(clinicaID))
	if err != nil {
		responses.InternalError(c, "Error al obtener roles del usuario")
		return
	}
	responses.Success(c, "Roles del usuario", list)
}

func (h *RolHandler) AsignarRolAUsuario(c *gin.Context) {
	usuarioID, ok := paramUint(c, "usuarioId")
	if !ok {
		return
	}
	var req adminModels.AsignarRolUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	if err := h.svc.AsignarRolAUsuario(usuarioID, req.RolID, req.ClinicaID); err != nil {
		responses.InternalError(c, "Error al asignar rol")
		return
	}
	responses.Success(c, "Rol asignado", nil)
}

func (h *RolHandler) RevocarRolDeUsuario(c *gin.Context) {
	usuarioID, ok := paramUint(c, "usuarioId")
	if !ok {
		return
	}
	rolID, err := strconv.ParseUint(c.Param("rolId"), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID de rol inválido")
		return
	}
	clinicaID, err := strconv.ParseUint(c.Query("clinica_id"), 10, 64)
	if err != nil || clinicaID == 0 {
		responses.BadRequest(c, "Se requiere clinica_id")
		return
	}
	if err := h.svc.RevocarRolDeUsuario(usuarioID, uint(rolID), uint(clinicaID)); err != nil {
		responses.InternalError(c, "Error al revocar rol")
		return
	}
	responses.Success(c, "Rol revocado", nil)
}
