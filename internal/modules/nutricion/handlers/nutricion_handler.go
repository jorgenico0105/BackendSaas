package handlers

import (
	"fmt"
	"log"
	"strconv"

	"saas-medico/internal/modules/nutricion/models"
	"saas-medico/internal/modules/nutricion/services"
	"saas-medico/internal/shared/openia"
	"saas-medico/internal/shared/responses"
	"saas-medico/internal/shared/uploads"

	"github.com/gin-gonic/gin"
)

type NutricionHandler struct {
	svc           *services.NutricionService
	openiaService *openia.OpenIaService
}

func NewNutricionHandler(svc *services.NutricionService, openiaService *openia.OpenIaService) *NutricionHandler {
	return &NutricionHandler{svc: svc, openiaService: openiaService}
}

func paramUint(c *gin.Context, key string) (uint, bool) {
	n, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return 0, false
	}
	return uint(n), true
}

// Reportes

func (h *NutricionHandler) CreateMenuReport(c *gin.Context) {
	id, ok := paramUint(c, "menuID")
	if !ok {
		return
	}

	menu, err := h.svc.GetMenu(uint(id))
	if err != nil {
		responses.NotFound(c, "Menú no encontrado")
		return
	}

	dieta, err := h.svc.GetDieta(menu.DietaPacienteID)
	if err != nil {
		responses.NotFound(c, "Dieta no encontrada")
		return
	}

	outputPath, err := h.svc.GetMenuReport(dieta, menu)
	if err != nil {
		responses.InternalError(c, "Error al generar el PDF")
		return
	}

	filename := fmt.Sprintf("menu_%d.pdf", id)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(outputPath)
}

// ─── Grupos de alimento ───────────────────────────────────────────────────────

func (h *NutricionHandler) ListGruposAlimento(c *gin.Context) {
	list, err := h.svc.ListGruposAlimento()
	if err != nil {
		responses.InternalError(c, "Error al listar grupos de alimento")
		return
	}
	responses.Success(c, "Grupos de alimento", list)
}

func (h *NutricionHandler) GetGrupoAlimento(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	g, err := h.svc.GetGrupoAlimento(id)
	if err != nil {
		responses.NotFound(c, "Grupo de alimento no encontrado")
		return
	}
	responses.Success(c, "Grupo de alimento", g)
}

func (h *NutricionHandler) CreateGrupoAlimento(c *gin.Context) {
	var req models.CreateGrupoAlimentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	g, err := h.svc.CreateGrupoAlimento(req)
	if err != nil {
		responses.InternalError(c, "Error al crear grupo de alimento")
		return
	}
	responses.Created(c, "Grupo de alimento creado", g)
}

func (h *NutricionHandler) UpdateGrupoAlimento(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateGrupoAlimentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	g, err := h.svc.UpdateGrupoAlimento(id, req)
	if err != nil {
		responses.NotFound(c, "Grupo de alimento no encontrado")
		return
	}
	responses.Success(c, "Grupo de alimento actualizado", g)
}

func (h *NutricionHandler) DeleteGrupoAlimento(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteGrupoAlimento(id); err != nil {
		responses.NotFound(c, "Grupo de alimento no encontrado")
		return
	}
	responses.Success(c, "Grupo de alimento eliminado", nil)
}

// ─── Alimentos ────────────────────────────────────────────────────────────────

func (h *NutricionHandler) ListAlimentos(c *gin.Context) {
	list, err := h.svc.ListAlimentos(c.Query("categoria"))
	if err != nil {
		responses.InternalError(c, "Error al listar alimentos")
		return
	}
	responses.Success(c, "Alimentos", list)
}

func (h *NutricionHandler) GetAlimento(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	a, err := h.svc.GetAlimento(id)
	if err != nil {
		responses.NotFound(c, "Alimento no encontrado")
		return
	}
	responses.Success(c, "Alimento", a)
}

func (h *NutricionHandler) CreateAlimento(c *gin.Context) {
	var req models.CreateAlimentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	userID := c.GetUint("userID")
	a, err := h.svc.CreateAlimento(req, userID)
	if err != nil {
		responses.InternalError(c, "Error al crear alimento")
		return
	}
	responses.Created(c, "Alimento creado", a)
}

func (h *NutricionHandler) UpdateAlimento(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateAlimentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	a, err := h.svc.UpdateAlimento(id, req)
	if err != nil {
		responses.NotFound(c, "Alimento no encontrado")
		return
	}
	responses.Success(c, "Alimento actualizado", a)
}

// ─── Catálogo dietas ──────────────────────────────────────────────────────────

func (h *NutricionHandler) ListDietasCatalogo(c *gin.Context) {
	list, err := h.svc.ListDietasCatalogo()
	if err != nil {
		responses.InternalError(c, "Error al listar catálogo de dietas")
		return
	}
	responses.Success(c, "Catálogo de dietas", list)
}

// ─── Dietas del paciente ──────────────────────────────────────────────────────

func (h *NutricionHandler) ListDietasByPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListDietasByPaciente(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar dietas")
		return
	}
	responses.Success(c, "Dietas del paciente", list)
}

func (h *NutricionHandler) GetDieta(c *gin.Context) {
	id, ok := paramUint(c, "dietaId")
	if !ok {
		return
	}
	d, err := h.svc.GetDieta(id)
	if err != nil {
		responses.NotFound(c, "Dieta no encontrada")
		return
	}
	responses.Success(c, "Dieta", d)
}

func (h *NutricionHandler) CreateDieta(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateDietaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	userID := c.GetUint("userID")
	d, err := h.svc.CreateDieta(pacienteID, userID, req)
	if err != nil {
		responses.InternalError(c, "Error al crear dieta")
		return
	}
	responses.Created(c, "Dieta creada", d)
}

func (h *NutricionHandler) UpdateDieta(c *gin.Context) {
	id, ok := paramUint(c, "dietaId")
	if !ok {
		return
	}
	var req models.UpdateDietaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	d, err := h.svc.UpdateDieta(id, req)
	if err != nil {
		responses.NotFound(c, "Dieta no encontrada")
		return
	}
	responses.Success(c, "Dieta actualizada", d)
}

// ─── Menús ────────────────────────────────────────────────────────────────────

func (h *NutricionHandler) CreateMenu(c *gin.Context) {
	dietaID, ok := paramUint(c, "dietaId")
	if !ok {
		return
	}
	var req models.AssignMenuFromPlantillaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	m, err := h.svc.AssignMenuFromPlantilla(dietaID, req)
	if err != nil {
		responses.InternalError(c, "Error al asignar menú al paciente")
		return
	}
	responses.Created(c, "Menú asignado al paciente", m)
}

func (h *NutricionHandler) AddDetalleMenu(c *gin.Context) {
	menuID, ok := paramUint(c, "menuId")
	if !ok {
		return
	}
	var req models.AddDetalleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	d, err := h.svc.AddDetalleMenu(menuID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar detalle al menú")
		return
	}
	responses.Created(c, "Detalle agregado", d)
}

func (h *NutricionHandler) UpdateDetalleMenu(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	var req models.UpdateDetalleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	d, err := h.svc.UpdateDetalleReceta(detalleID, req.NombreReceta)
	if err != nil {
		responses.NotFound(c, "Detalle no encontrado")
		return
	}
	responses.Success(c, "Receta actualizada", d)
}

func (h *NutricionHandler) ListMenusByDieta(c *gin.Context) {
	dietaID, ok := paramUint(c, "dietaId")
	if !ok {
		return
	}
	list, err := h.svc.ListMenusByDieta(dietaID)
	if err != nil {
		responses.InternalError(c, "Error al listar menús")
		return
	}
	responses.Success(c, "Menús de la dieta", list)
}

func (h *NutricionHandler) GetMenu(c *gin.Context) {
	menuID, ok := paramUint(c, "menuId")
	if !ok {
		return
	}
	m, err := h.svc.GetMenu(menuID)
	if err != nil {
		responses.NotFound(c, "Menú no encontrado")
		return
	}
	responses.Success(c, "Menú", m)
}

func (h *NutricionHandler) GetDetallesMenu(c *gin.Context) {
	menuID, ok := paramUint(c, "menuId")
	if !ok {
		return
	}
	list, err := h.svc.GetDetallesMenu(menuID)
	if err != nil {
		responses.InternalError(c, "Error al obtener detalles del menú")
		return
	}
	responses.Success(c, "Detalles del menú", list)
}

func (h *NutricionHandler) GetAlimentosMenuDetalle(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	list, err := h.svc.GetAlimentosMenuDetalle(detalleID)
	if err != nil {
		responses.InternalError(c, "Error al obtener alimentos del detalle")
		return
	}
	responses.Success(c, "Alimentos del detalle de menú", list)
}

func (h *NutricionHandler) AddAlimentoMenuDetalle(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	var req models.AddAlimentoMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	a, err := h.svc.AddAlimentoMenuDetalle(detalleID, req.AlimentoID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar alimento al detalle de menú")
		return
	}
	responses.Created(c, "Alimento agregado al detalle de menú", a)
}

func (h *NutricionHandler) DeleteAlimentoMenuDetalle(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteAlimentoMenuDetalle(id); err != nil {
		responses.NotFound(c, "Alimento no encontrado")
		return
	}
	responses.Success(c, "Alimento eliminado", nil)
}

func (h *NutricionHandler) UpdateAlimentoMenuDetalle(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req struct {
		GramosAsignados float64 `json:"gramos_asignados" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: gramos_asignados requerido y mayor a 0")
		return
	}
	a, err := h.svc.UpdateAlimentoMenuDetalle(id, req.GramosAsignados)
	if err != nil {
		responses.NotFound(c, "Alimento no encontrado")
		return
	}
	responses.Success(c, "Gramaje actualizado", a)
}

func (h *NutricionHandler) ListDietasRequierenCambio(c *gin.Context) {
	list, err := h.svc.ListDietasRequierenCambio()
	if err != nil {
		responses.InternalError(c, "Error al listar dietas que requieren cambio")
		return
	}
	responses.Success(c, "Dietas que requieren cambio", list)
}

// ─── Fórmulas nutricionales ───────────────────────────────────────────────────

func (h *NutricionHandler) CalcularFormulas(c *gin.Context) {
	var req models.CalcularFormulasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	result := h.svc.CalcularFormulas(req)
	responses.Success(c, "Fórmulas calculadas", result)
}

// ─── R24H ─────────────────────────────────────────────────────────────────────

func (h *NutricionHandler) CreateR24H(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateR24HRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	userID := c.GetUint("userID")
	r, err := h.svc.CreateR24H(pacienteID, userID, req)
	if err != nil {
		responses.InternalError(c, "Error al crear recordatorio 24h")
		return
	}
	responses.Created(c, "Recordatorio 24h creado", r)
}

func (h *NutricionHandler) ListR24H(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListR24H(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar recordatorios 24h")
		return
	}
	responses.Success(c, "Recordatorios 24h", list)
}

func (h *NutricionHandler) AddR24HItem(c *gin.Context) {
	r24hID, ok := paramUint(c, "r24hId")
	if !ok {
		return
	}
	var req models.AddR24HItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	item, err := h.svc.AddR24HItem(r24hID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar item al recordatorio 24h")
		return
	}
	responses.Created(c, "Item agregado al recordatorio 24h", item)
}

func (h *NutricionHandler) ListR24HItems(c *gin.Context) {
	r24hID, ok := paramUint(c, "r24hId")
	if !ok {
		return
	}
	list, err := h.svc.ListR24HItems(r24hID)
	if err != nil {
		responses.InternalError(c, "Error al listar items del recordatorio 24h")
		return
	}
	responses.Success(c, "Items del recordatorio 24h", list)
}

// ─── Preferencias ─────────────────────────────────────────────────────────────

func (h *NutricionHandler) AddPreferencia(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreatePreferenciaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	p, err := h.svc.AddPreferencia(pacienteID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar preferencia")
		return
	}
	responses.Created(c, "Preferencia registrada", p)
}

func (h *NutricionHandler) ListPreferencias(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListPreferencias(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar preferencias")
		return
	}
	responses.Success(c, "Preferencias alimentarias", list)
}

func (h *NutricionHandler) DeletePreferencia(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeletePreferencia(id); err != nil {
		responses.InternalError(c, "Error al eliminar preferencia")
		return
	}
	responses.Success(c, "Preferencia eliminada", nil)
}

// ─── Síntomas ─────────────────────────────────────────────────────────────────

func (h *NutricionHandler) CreateSintoma(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateSintomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	s, err := h.svc.CreateSintoma(pacienteID, req)
	if err != nil {
		responses.InternalError(c, "Error al registrar síntoma")
		return
	}
	responses.Created(c, "Síntoma registrado", s)
}

func (h *NutricionHandler) ListSintomas(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListSintomas(pacienteID, c.Query("desde"), c.Query("hasta"))
	if err != nil {
		responses.InternalError(c, "Error al listar síntomas")
		return
	}
	responses.Success(c, "Síntomas", list)
}

// ─── Ejercicios catálogo ──────────────────────────────────────────────────────

func (h *NutricionHandler) ListEjerciciosCatalogo(c *gin.Context) {
	list, err := h.svc.ListEjerciciosCatalogo(c.Query("categoria"))
	if err != nil {
		responses.InternalError(c, "Error al listar ejercicios")
		return
	}
	responses.Success(c, "Catálogo de ejercicios", list)
}

func (h *NutricionHandler) CreateEjercicioCatalogo(c *gin.Context) {
	var req models.CreateEjercicioCatalogoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	userID := c.GetUint("userID")
	e, err := h.svc.CreateEjercicioCatalogo(req, userID)
	if err != nil {
		responses.InternalError(c, "Error al crear ejercicio en catálogo")
		return
	}
	responses.Created(c, "Ejercicio creado en catálogo", e)
}

// ─── Ejercicios paciente ──────────────────────────────────────────────────────

func (h *NutricionHandler) ListEjerciciosByPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListEjerciciosByPaciente(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar ejercicios del paciente")
		return
	}
	responses.Success(c, "Ejercicios del paciente", list)
}

func (h *NutricionHandler) AddEjercicioPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateEjercicioPacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	userID := c.GetUint("userID")
	e, err := h.svc.AddEjercicioPaciente(pacienteID, userID, req)
	if err != nil {
		responses.InternalError(c, "Error al asignar ejercicio al paciente")
		return
	}
	responses.Created(c, "Ejercicio asignado al paciente", e)
}

// ─── Registros comida ─────────────────────────────────────────────────────────

func (h *NutricionHandler) ListRegistrosComida(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListRegistrosComida(pacienteID, c.Query("fecha"), c.Query("desde"), c.Query("hasta"))
	if err != nil {
		responses.InternalError(c, "Error al listar registros de comida")
		return
	}
	responses.Success(c, "Registros de comida", list)
}

func (h *NutricionHandler) CreateRegistroComida(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateRegistroComidaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	rc, err := h.svc.CreateRegistroComida(pacienteID, req)
	if err != nil {
		responses.InternalError(c, "Error al registrar comida")
		return
	}
	responses.Created(c, "Registro de comida creado", rc)
}

func (h *NutricionHandler) AddRegistroAlimento(c *gin.Context) {
	registroID, ok := paramUint(c, "registroId")
	if !ok {
		return
	}
	var req models.AddRegistroAlimentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	ra, err := h.svc.AddRegistroAlimento(registroID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar alimento al registro")
		return
	}
	responses.Created(c, "Alimento agregado al registro", ra)
}

func (h *NutricionHandler) UploadFotoComida(c *gin.Context) {
	registroID, ok := paramUint(c, "registroId")
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("foto")
	if err != nil {
		responses.BadRequest(c, "Se requiere el campo 'foto'")
		return
	}

	result, err := uploads.SaveFile(c, fileHeader, "comidas", uploads.AllowedImageTypes)
	if err != nil {
		responses.BadRequest(c, err.Error())
		return
	}

	rc, oldFoto, err := h.svc.UpdateFotoComida(registroID, result.FilePath)
	if err != nil {
		uploads.DeleteFile(result.FilePath)
		responses.NotFound(c, "Registro de comida no encontrado")
		return
	}

	if oldFoto != "" {
		uploads.DeleteFile(oldFoto)
	}

	responses.Success(c, "Foto actualizada", rc)
}

// MarcarRegistroComidaConsumida cambia el estado del registro a 'C' (consumida)
func (h *NutricionHandler) MarcarRegistroComidaConsumida(c *gin.Context) {
	registroID, ok := paramUint(c, "registroId")
	if !ok {
		return
	}
	if err := h.svc.MarcarConsumida(registroID); err != nil {
		responses.NotFound(c, "Registro no encontrado")
		return
	}
	responses.Success(c, "Comida marcada como consumida", nil)
}

// GetResumenDiario devuelve resumen consolidado del día (comidas + ejercicios + progreso)
func (h *NutricionHandler) GetResumenDiario(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	fecha := c.Query("fecha")
	resumen, err := h.svc.GetResumenDiario(pacienteID, fecha)
	if err != nil {
		responses.InternalError(c, "Error al obtener resumen diario")
		return
	}
	responses.Success(c, "Resumen diario", resumen)
}

// ─── Registros ejercicio ──────────────────────────────────────────────────────

func (h *NutricionHandler) ListRegistrosEjercicio(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListRegistrosEjercicio(pacienteID, c.Query("fecha"), c.Query("desde"), c.Query("hasta"))
	if err != nil {
		responses.InternalError(c, "Error al listar registros de ejercicio")
		return
	}
	responses.Success(c, "Registros de ejercicio", list)
}

func (h *NutricionHandler) CreateRegistroEjercicio(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateRegistroEjercicioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	re, err := h.svc.CreateRegistroEjercicio(pacienteID, req)
	if err != nil {
		responses.InternalError(c, "Error al registrar ejercicio")
		return
	}
	responses.Created(c, "Registro de ejercicio creado", re)
}

// ─── Progreso ─────────────────────────────────────────────────────────────────

func (h *NutricionHandler) ListProgreso(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListProgreso(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar progreso")
		return
	}
	responses.Success(c, "Progreso del paciente", list)
}

func (h *NutricionHandler) AddProgreso(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.CreateProgresoRequest
	if err := c.ShouldBind(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	// Optional photo upload
	if fileHeader, err := c.FormFile("foto"); err == nil {
		result, err := uploads.SaveFile(c, fileHeader, "progreso", uploads.AllowedImageTypes)
		if err == nil {
			req.FotoProgreso = result.FilePath
		}
	}
	userID := c.GetUint("userID")
	p, err := h.svc.AddProgreso(pacienteID, userID, req)
	if err != nil {
		responses.InternalError(c, "Error al registrar progreso")
		return
	}
	responses.Created(c, "Progreso registrado", p)
}

// ─── Tipo de Recurso ──────────────────────────────────────────────────────────

func (h *NutricionHandler) ListTipoRecursos(c *gin.Context) {
	list, err := h.svc.ListTipoRecursos()
	if err != nil {
		responses.InternalError(c, "Error al listar tipos de recursos")
		return
	}
	responses.Success(c, "Tipos de recursos", list)
}

func (h *NutricionHandler) CreateTipoRecurso(c *gin.Context) {
	var req models.CreateTipoRecursoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "El nombre es requerido")
		return
	}
	t, err := h.svc.CreateTipoRecurso(req)
	if err != nil {
		responses.InternalError(c, "Error al crear tipo de recurso")
		return
	}
	responses.Created(c, "Tipo de recurso creado", t)
}

func (h *NutricionHandler) UpdateTipoRecurso(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateTipoRecursoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "El nombre es requerido")
		return
	}
	t, err := h.svc.UpdateTipoRecurso(id, req)
	if err != nil {
		responses.NotFound(c, "Tipo de recurso no encontrado")
		return
	}
	responses.Success(c, "Tipo de recurso actualizado", t)
}

func (h *NutricionHandler) DeleteTipoRecurso(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteTipoRecurso(id); err != nil {
		responses.InternalError(c, "Error al eliminar tipo de recurso")
		return
	}
	responses.Success(c, "Tipo de recurso eliminado", nil)
}

// ─── Archivos PDF ─────────────────────────────────────────────────────────────

func (h *NutricionHandler) CreateArchivoPDF(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		responses.BadRequest(c, "Se requiere un archivo")
		return
	}

	titulo := c.PostForm("titulo")
	if titulo == "" {
		responses.BadRequest(c, "El título es requerido")
		return
	}

	tipoStr := c.PostForm("tipo_recurso_id")
	if tipoStr == "" {
		responses.BadRequest(c, "El tipo de recurso es requerido")
		return
	}
	tipoID, err := strconv.ParseUint(tipoStr, 10, 64)
	if err != nil {
		responses.BadRequest(c, "tipo_recurso_id inválido")
		return
	}

	result, err := uploads.SaveFile(c, fileHeader, "recursos", append(uploads.AllowedDocTypes, uploads.AllowedImageTypes...))
	if err != nil {
		responses.BadRequest(c, err.Error())
		return
	}

	req := models.CreateArchivoPDFRequest{
		TipoRecursoID: uint(tipoID),
		Titulo:        titulo,
		Descripcion:   c.PostForm("descripcion"),
		RutaArchivo:   result.FilePath,
	}
	if pid := c.PostForm("paciente_id"); pid != "" {
		if n, err2 := strconv.ParseUint(pid, 10, 64); err2 == nil {
			v := uint(n)
			req.PacienteID = &v
		}
	}

	clinicaID := c.GetUint("clinicaID")
	userID := c.GetUint("userID")
	a, err := h.svc.CreateArchivoPDF(clinicaID, userID, req)
	if err != nil {
		uploads.DeleteFile(result.FilePath)
		responses.InternalError(c, "Error al guardar archivo PDF")
		return
	}
	responses.Created(c, "Archivo PDF guardado", a)
}

func (h *NutricionHandler) ListArchivosPDF(c *gin.Context) {
	clinicaID := c.GetUint("clinicaID")

	var pacienteID *uint
	if pid := c.Query("paciente_id"); pid != "" {
		if n, err := strconv.ParseUint(pid, 10, 64); err == nil {
			v := uint(n)
			pacienteID = &v
		}
	}

	var tipoRecursoID *uint
	if tid := c.Query("tipo_recurso_id"); tid != "" {
		if n, err := strconv.ParseUint(tid, 10, 64); err == nil {
			v := uint(n)
			tipoRecursoID = &v
		}
	}

	list, err := h.svc.ListArchivosPDF(clinicaID, pacienteID, tipoRecursoID)
	if err != nil {
		responses.InternalError(c, "Error al listar archivos PDF")
		return
	}
	responses.Success(c, "Archivos PDF", list)
}

func (h *NutricionHandler) ListArchivosPDFByPaciente(c *gin.Context) {
	clinicaID := c.GetUint("clinicaID")
	pacienteID, ok := paramUint(c, "pacienteID")
	if !ok {
		return
	}
	list, err := h.svc.ListArchivosPDFByUser(clinicaID, pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar archivos PDF")
		return
	}
	responses.Success(c, "Archivos PDF", list)
}

func (h *NutricionHandler) DeleteArchivoPDF(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteArchivoPDF(id); err != nil {
		responses.InternalError(c, "Error al eliminar archivo PDF")
		return
	}
	responses.Success(c, "Archivo PDF eliminado", nil)
}

// ─── XP y logros ──────────────────────────────────────────────────────────────

func (h *NutricionHandler) GetXP(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	xp, err := h.svc.GetXP(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al obtener XP del paciente")
		return
	}
	responses.Success(c, "XP del paciente", xp)
}

func (h *NutricionHandler) ListLogros(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.svc.ListLogros(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar logros del paciente")
		return
	}
	responses.Success(c, "Logros del paciente", list)
}

func (h *NutricionHandler) ListLogrosCatalogo(c *gin.Context) {
	list, err := h.svc.ListLogrosCatalogo()
	if err != nil {
		responses.InternalError(c, "Error al listar catálogo de logros")
		return
	}
	responses.Success(c, "Catálogo de logros", list)
}
// ─── Plantillas de menú semanal ───────────────────────────────────────────────

func (h *NutricionHandler) ListPlantillas(c *gin.Context) {
	clinicaID := c.GetUint("clinicaID")

	var numComidas *int
	if v := c.Query("num_comidas"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			numComidas = &n
		}
	}
	var semanaNumero *int
	if v := c.Query("semana_numero"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			semanaNumero = &n
		}
	}

	list, err := h.svc.ListPlantillas(clinicaID, numComidas, semanaNumero)
	if err != nil {
		responses.InternalError(c, "Error al listar plantillas de menú")
		return
	}
	responses.Success(c, "Plantillas de menú", list)
}

func (h *NutricionHandler) CreatePlantillaSemana(c *gin.Context) {
	var req models.CreatePlantillaSemanaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	clinicaID := c.GetUint("clinicaID")
	medicoID := c.GetUint("userID")
	p, err := h.svc.CreatePlantillaSemana(clinicaID, medicoID, req)
	if err != nil {
		responses.InternalError(c, "Error al crear plantilla de menú")
		return
	}
	responses.Created(c, "Plantilla de menú creada", p)
}

func (h *NutricionHandler) GetPlantillaSemana(c *gin.Context) {
	id, ok := paramUint(c, "plantillaId")
	if !ok {
		return
	}
	p, err := h.svc.GetPlantillaSemana(id)
	if err != nil {
		responses.NotFound(c, "Plantilla de menú no encontrada")
		return
	}
	responses.Success(c, "Plantilla de menú", p)
}

func (h *NutricionHandler) UpdatePlantillaSemana(c *gin.Context) {
	id, ok := paramUint(c, "plantillaId")
	if !ok {
		return
	}
	var req models.UpdatePlantillaSemanaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	p, err := h.svc.UpdatePlantillaSemana(id, req)
	if err != nil {
		responses.NotFound(c, "Plantilla de menú no encontrada")
		return
	}
	responses.Success(c, "Plantilla de menú actualizada", p)
}

func (h *NutricionHandler) DeletePlantillaSemana(c *gin.Context) {
	id, ok := paramUint(c, "plantillaId")
	if !ok {
		return
	}
	if err := h.svc.DeletePlantillaSemana(id); err != nil {
		responses.NotFound(c, "Plantilla de menú no encontrada")
		return
	}
	responses.Success(c, "Plantilla de menú eliminada", nil)
}

func (h *NutricionHandler) AddDetallePlantilla(c *gin.Context) {
	plantillaID, ok := paramUint(c, "plantillaId")
	if !ok {
		return
	}
	var req models.AddDetallePlantillaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	d, err := h.svc.AddDetallePlantilla(plantillaID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar detalle a la plantilla")
		return
	}
	responses.Created(c, "Detalle agregado a la plantilla", d)
}

func (h *NutricionHandler) GetDetallesPlantilla(c *gin.Context) {
	plantillaID, ok := paramUint(c, "plantillaId")
	if !ok {
		return
	}
	list, err := h.svc.GetDetallesPlantilla(plantillaID)
	if err != nil {
		responses.InternalError(c, "Error al obtener detalles de la plantilla")
		return
	}
	responses.Success(c, "Detalles de la plantilla", list)
}

func (h *NutricionHandler) UpdateDetallePlantilla(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	var req models.UpdateDetallePlantillaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos")
		return
	}
	d, err := h.svc.UpdateDetallePlantilla(detalleID, req)
	if err != nil {
		responses.NotFound(c, "Detalle de plantilla no encontrado")
		return
	}
	responses.Success(c, "Detalle de plantilla actualizado", d)
}

func (h *NutricionHandler) DeleteDetallePlantilla(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	if err := h.svc.DeleteDetallePlantilla(detalleID); err != nil {
		responses.NotFound(c, "Detalle de plantilla no encontrado")
		return
	}
	responses.Success(c, "Detalle de plantilla eliminado", nil)
}

func (h *NutricionHandler) AddAlimentoPlantillaDetalle(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	var req models.AddAlimentoPlantillaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	a, err := h.svc.AddAlimentoPlantillaDetalle(detalleID, req)
	if err != nil {
		responses.InternalError(c, "Error al agregar alimento a la plantilla")
		return
	}
	responses.Created(c, "Alimento agregado a la plantilla", a)
}

func (h *NutricionHandler) GetAlimentosPlantillaDetalle(c *gin.Context) {
	detalleID, ok := paramUint(c, "detalleId")
	if !ok {
		return
	}
	list, err := h.svc.GetAlimentosPlantillaDetalle(detalleID)
	if err != nil {
		responses.InternalError(c, "Error al obtener alimentos del detalle de plantilla")
		return
	}
	responses.Success(c, "Alimentos del detalle de plantilla", list)
}

func (h *NutricionHandler) UpdateAlimentoPlantillaDetalle(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req struct {
		GramosAsignados float64 `json:"gramos_asignados" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: gramos_asignados requerido y mayor a 0")
		return
	}
	a, err := h.svc.UpdateAlimentoPlantillaDetalle(id, req.GramosAsignados)
	if err != nil {
		responses.NotFound(c, "Alimento de plantilla no encontrado")
		return
	}
	responses.Success(c, "Gramaje actualizado en plantilla", a)
}

func (h *NutricionHandler) DeleteAlimentoPlantillaDetalle(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteAlimentoPlantillaDetalle(id); err != nil {
		responses.NotFound(c, "Alimento de plantilla no encontrado")
		return
	}
	responses.Success(c, "Alimento eliminado de la plantilla", nil)
}

func (h *NutricionHandler) ChatWhitIa(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	var req models.AskIaNutricionQuestion
	if err := c.ShouldBind(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	dieta, err := h.svc.ListDietasByPaciente(pacienteID)

	alimentos, err := h.svc.ListAlimentos("")
	if err != nil {
		responses.InternalError(c, "Error fatal")
	}
	log.Printf("[infromacion] %v", dieta[0].Paciente.Apellidos)
	resp, err := h.openiaService.AskModelIa(alimentos, req.Prompt, dieta)

	responses.Success(c, "Respuesta ia", resp)
}
