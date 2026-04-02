package models

// ─── Formularios CRUD ─────────────────────────────────────────────────────────

type CreateFormularioRequest struct {
	Nombre           string                  `json:"nombre" binding:"required,max=150"`
	Descripcion      string                  `json:"descripcion" binding:"omitempty,max=255"`
	TipoFormularioID uint                    `json:"tipo_formulario_id" binding:"required"`
	ProfesionID      *uint                   `json:"profesion_id"`
	Preguntas        []CreatePreguntaRequest `json:"preguntas"`
}

type UpdateFormularioRequest struct {
	Nombre      string                  `json:"nombre" binding:"omitempty,max=150"`
	Descripcion string                  `json:"descripcion" binding:"omitempty,max=255"`
	Preguntas   []CreatePreguntaRequest `json:"preguntas"`
}

type CreatePreguntaRequest struct {
	Pregunta      string                `json:"pregunta" binding:"required,max=255"`
	TipoRespuesta string                `json:"tipo_respuesta" binding:"required,oneof=TEXT NUMBER DATE SELECT MULTISELECT BOOLEAN"`
	Obligatorio   bool                  `json:"obligatorio"`
	Orden         int                   `json:"orden"`
	Opciones      []CreateOpcionRequest `json:"opciones"`
}

type CreateOpcionRequest struct {
	Valor    string  `json:"valor" binding:"required,max=100"`
	Etiqueta string  `json:"etiqueta" binding:"required,max=150"`
	Orden    int     `json:"orden"`
	Puntos   float64 `json:"puntos"`
}

// ─── Historia Clínica ─────────────────────────────────────────────────────────

type CreateHistoriaClinicaRequest struct {
	MedicoID           uint                             `json:"medico_id" binding:"required"`
	FormularioID       uint                             `json:"formulario_id" binding:"required"`
	Fecha              string                           `json:"fecha" binding:"required"` // ejemplo: "2026-03-12"
	ObservacionGeneral string                           `json:"observacion_general" binding:"omitempty"`
	Preguntas          []CreateHistoriaRespuestaRequest `json:"preguntas" binding:"required,dive"`
}

type CreateHistoriaRespuestaRequest struct {
	PreguntaID      uint     `json:"pregunta_id" binding:"required"`
	RespuestaTexto  string   `json:"respuesta_texto" binding:"omitempty"`
	RespuestaNumero *float64 `json:"respuesta_numero" binding:"omitempty"`
	RespuestaFecha  *string  `json:"respuesta_fecha" binding:"omitempty"`
}
type ResultadoHistoria struct {
	Pregunta          string   `json:"pregunta"`
	OrdenPregunta     int      `json:"orden_pregunta"`
	Multiple          bool     `json:"multiple"`
	RespuestaText     string   `json:"respuesta_text"`
	RespuestaNumero   *float64 `json:"respuesta_numero"`
	IDHistoriaClinica uint     `json:"id_historia_clinica"`
	FechaRegistro     string   `json:"fecha_registro"`
	NombreFormulario  string   `json:"nombre_formulario"`
}
