package handlers

import (
	"net/http"
	"saas-medico/internal/modules/historia/models"
	"saas-medico/internal/modules/historia/services"
	"saas-medico/internal/shared/responses"
	"saas-medico/internal/shared/uploads"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HistoriaHandler struct {
	service *services.HistoriaService
}

func NewHistoriaHandler(service *services.HistoriaService) *HistoriaHandler {
	return &HistoriaHandler{service: service}
}

func paramUint(c *gin.Context, key string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		responses.BadRequest(c, "ID inválido")
		return 0, false
	}
	return uint(v), true
}

// ─── Formularios ──────────────────────────────────────────────────────────────

func (h *HistoriaHandler) ListTiposFormulario(c *gin.Context) {
	list, err := h.service.ListTiposFormulario()
	if err != nil {
		responses.InternalError(c, "Error al listar tipos de formulario")
		return
	}
	responses.Success(c, "Tipos de formulario", list)
}

func (h *HistoriaHandler) CreateFormulario(c *gin.Context) {
	var req models.CreateFormularioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	usuarioID := c.GetUint("userID")
	clinicaID := c.GetUint("clinicaID")
	f, err := h.service.CreateFormularioCompleto(req, usuarioID, clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al crear formulario: "+err.Error())
		return
	}
	responses.Created(c, "Formulario creado", f)
}

func (h *HistoriaHandler) UpdateFormulario(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	var req models.UpdateFormularioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	if err := h.service.UpdateFormularioCompleto(id, req); err != nil {
		responses.InternalError(c, "Error al actualizar formulario: "+err.Error())
		return
	}
	responses.Success(c, "Formulario actualizado", nil)
}

func (h *HistoriaHandler) DeleteFormulario(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	if err := h.service.DeleteFormulario(id); err != nil {
		responses.NotFound(c, "Formulario no encontrado")
		return
	}
	responses.Success(c, "Formulario eliminado", nil)
}

func (h *HistoriaHandler) ListFormularios(c *gin.Context) {
	var tipoID uint
	if t := c.Query("tipo_id"); t != "" {
		v, _ := strconv.ParseUint(t, 10, 64)
		tipoID = uint(v)
	}
	usuarioID := c.GetUint("userID")
	clinicaID := c.GetUint("clinicaID")
	list, err := h.service.ListFormularios(tipoID, usuarioID, clinicaID)
	if err != nil {
		responses.InternalError(c, "Error al listar formularios")
		return
	}
	responses.Success(c, "Formularios", list)
}

func (h *HistoriaHandler) GetFormulario(c *gin.Context) {
	id, ok := paramUint(c, "id")
	if !ok {
		return
	}
	f, err := h.service.GetFormulario(id)
	if err != nil {
		responses.NotFound(c, "Formulario no encontrado")
		return
	}
	preguntas, opciones, err := h.service.GetPreguntasConOpciones(id)
	if err != nil {
		responses.InternalError(c, "Error al obtener preguntas")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"formulario": f,
		"preguntas":  preguntas,
		"opciones":   opciones,
	})
}

// ─── Historias ────────────────────────────────────────────────────────────────
func (h *HistoriaHandler) CreateHistoriaPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	var req models.CreateHistoriaClinicaRequest
	if !ok {
		responses.BadRequest(c, "Error parsenado")
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}
	respuestas, err := h.service.CreateHistoria(&req, int(pacienteID))
	if err != nil {
		responses.BadRequest(c, "Error creando la hitoria clinica")
	}
	responses.Success(c, "Historia clinica guardada!", respuestas)
}

func (h *HistoriaHandler) GetHistoriaClinicaByUser(c *gin.Context) {
	userId := c.GetUint("userID")
	clinicaId := c.GetUint("clinicaID")
	tipoForm := 1
	list, err := h.service.GetHistoriaClinicaByUser(int(userId), int(clinicaId), tipoForm)
	if err != nil {
		responses.InternalError(c, "Error al traer el formulario historia clinica")
		return
	}
	responses.Success(c, "Hisotria clinica formulario", list)
}

func (h *HistoriaHandler) ListHistoriasByPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListHistoriasByPaciente(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar historias")
		return
	}
	responses.Success(c, "Historias clínicas", list)
}
func (h *HistoriaHandler) GetHistoriasByPaciente(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	res, err := h.service.FindHistoriasClinicasByPaciente(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar alergias")
		return
	}
	responses.Success(c, "Hisotrias Clinicas", res)
}

// ─── Alergias ─────────────────────────────────────────────────────────────────

func (h *HistoriaHandler) ListAlergias(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListAlergias(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar alergias")
		return
	}
	responses.Success(c, "Alergias", list)
}

// ─── Antecedentes ─────────────────────────────────────────────────────────────

func (h *HistoriaHandler) ListAntecedentes(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListAntecedentes(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar antecedentes")
		return
	}
	responses.Success(c, "Antecedentes", list)
}

// ─── Hábitos ──────────────────────────────────────────────────────────────────

func (h *HistoriaHandler) ListHabitos(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListHabitos(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar hábitos")
		return
	}
	responses.Success(c, "Hábitos", list)
}

// ─── Diagnósticos ─────────────────────────────────────────────────────────────

func (h *HistoriaHandler) ListDiagnosticos(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListDiagnosticos(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar diagnósticos")
		return
	}
	responses.Success(c, "Diagnósticos", list)
}

// ─── Imágenes del paciente ────────────────────────────────────────────────────

func (h *HistoriaHandler) UploadPacienteImagen(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	medicoID := c.GetUint("userID")

	fileHeader, err := c.FormFile("imagen")
	if err != nil {
		responses.BadRequest(c, "Se requiere el archivo de imagen")
		return
	}

	result, err := uploads.SaveFile(c, fileHeader, "paciente_imagenes", uploads.AllowedImageTypes)
	if err != nil {
		responses.BadRequest(c, "Error al guardar imagen: "+err.Error())
		return
	}

	descripcion := c.PostForm("descripcion")

	img := &models.PacienteImagen{
		PacienteID:    pacienteID,
		MedicoID:      &medicoID,
		NombreArchivo: result.FileName,
		UrlArchivo:    result.FilePath,
		TipoImagen:    1,
		Descripcion:   descripcion,
	}
	if err := h.service.AddPacienteImagen(img); err != nil {
		responses.InternalError(c, "Error al registrar imagen")
		return
	}
	responses.Created(c, "Imagen registrada", img)
}

func (h *HistoriaHandler) ListPacienteImagenes(c *gin.Context) {
	pacienteID, ok := paramUint(c, "pacienteId")
	if !ok {
		return
	}
	list, err := h.service.ListImagenesPaciente(pacienteID)
	if err != nil {
		responses.InternalError(c, "Error al listar imágenes")
		return
	}
	responses.Success(c, "Imágenes del paciente", list)
}
