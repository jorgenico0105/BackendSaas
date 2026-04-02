package models

import "time"

// ─── Historia Clínica ──────────────────────────────────────────────────────────

type HistoriaClinica struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID         uint      `gorm:"not null;index" json:"paciente_id"`
	MedicoID           uint      `gorm:"not null;index" json:"medico_id"`
	FormularioID       uint      `gorm:"not null;index" json:"formulario_id"`
	Fecha              time.Time `gorm:"not null" json:"fecha"`
	ObservacionGeneral string               `gorm:"type:text" json:"observacion_general,omitempty"`
	State              string               `gorm:"type:char(1);default:'A';not null" json:"state"`
	Respuestas         []HistoriaRespuesta  `gorm:"foreignKey:HistoriaID" json:"respuestas,omitempty"`
}

func (HistoriaClinica) TableName() string { return "historias_clinicas" }

type HistoriaRespuesta struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	HistoriaID      uint       `gorm:"not null;index" json:"historia_id"`
	PreguntaID      uint       `gorm:"not null;index" json:"pregunta_id"`
	RespuestaTexto  string     `gorm:"type:text" json:"respuesta_texto,omitempty"`
	RespuestaNumero *float64   `gorm:"type:decimal(10,2)" json:"respuesta_numero,omitempty"`
	RespuestaFecha  *time.Time `gorm:"type:date" json:"respuesta_fecha,omitempty"`
	CreadoEn        time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (HistoriaRespuesta) TableName() string { return "historia_respuestas" }

// ─── Alergias del paciente ─────────────────────────────────────────────────────

type PacienteAlergia struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint      `gorm:"not null;uniqueIndex:udx_pac_alergia" json:"paciente_id"`
	AlergiaID     uint      `gorm:"not null;uniqueIndex:udx_pac_alergia" json:"alergia_id"`
	Severidad     string    `gorm:"size:50" json:"severidad,omitempty"`
	Reaccion      string    `gorm:"size:255" json:"reaccion,omitempty"`
	Observacion   string    `gorm:"size:255" json:"observacion,omitempty"`
	FechaRegistro time.Time `gorm:"not null" json:"fecha_registro"`
	MedicoID      *uint     `gorm:"index" json:"medico_id,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (PacienteAlergia) TableName() string { return "paciente_alergias" }

// ─── Antecedentes del paciente ─────────────────────────────────────────────────

type PacienteAntecedente struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID        uint      `gorm:"not null;index" json:"paciente_id"`
	TipoAntecedenteID uint      `gorm:"not null;index" json:"tipo_antecedente_id"`
	Descripcion       string    `gorm:"type:text;not null" json:"descripcion"`
	FechaRegistro     time.Time `gorm:"not null" json:"fecha_registro"`
	MedicoID          *uint     `gorm:"index" json:"medico_id,omitempty"`
	State             string    `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (PacienteAntecedente) TableName() string { return "paciente_antecedentes" }

// ─── Hábitos del paciente ──────────────────────────────────────────────────────

type PacienteHabito struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint      `gorm:"not null;uniqueIndex:udx_pac_habito" json:"paciente_id"`
	HabitoID      uint      `gorm:"not null;uniqueIndex:udx_pac_habito" json:"habito_id"`
	Valor         string    `gorm:"size:120" json:"valor,omitempty"`
	Frecuencia    string    `gorm:"size:120" json:"frecuencia,omitempty"`
	Observacion   string    `gorm:"size:255" json:"observacion,omitempty"`
	FechaRegistro time.Time `gorm:"not null" json:"fecha_registro"`
	MedicoID      *uint     `gorm:"index" json:"medico_id,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (PacienteHabito) TableName() string { return "paciente_habitos" }

// ─── Diagnósticos del paciente ─────────────────────────────────────────────────

const (
	EstadoClinicoDiagActivo        = "ACTIVO"
	EstadoClinicoDiagResuelto      = "RESUELTO"
	EstadoClinicoDiagCronico       = "CRONICO"
	EstadoClinicoDiagEnSeguimiento = "EN_SEGUIMIENTO"
)

type PacienteDiagnostico struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID       uint       `gorm:"not null;index" json:"paciente_id"`
	DiagnosticoID    uint       `gorm:"not null;index" json:"diagnostico_id"`
	MedicoID         uint       `gorm:"not null;index" json:"medico_id"`
	SesionID         *uint      `gorm:"index" json:"sesion_id,omitempty"`
	CitaID           *uint      `gorm:"index" json:"cita_id,omitempty"`
	EstadoClinico    string     `gorm:"size:30;default:'ACTIVO'" json:"estado_clinico"`
	FechaDiagnostico time.Time  `gorm:"type:date;not null" json:"fecha_diagnostico"`
	FechaResolucion  *time.Time `gorm:"type:date" json:"fecha_resolucion,omitempty"`
	Observaciones    string     `gorm:"type:text" json:"observaciones,omitempty"`
	State            string     `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn         time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn    time.Time  `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

func (PacienteDiagnostico) TableName() string { return "paciente_diagnosticos" }

// ─── Exámenes e imágenes del paciente ─────────────────────────────────────────

type PacienteExamenResultado struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint       `gorm:"not null;index" json:"paciente_id"`
	TipoExamenID  uint       `gorm:"not null;index" json:"tipo_examen_id"`
	NombreArchivo string     `gorm:"size:255;not null" json:"nombre_archivo"`
	FechaExamen   *time.Time `gorm:"type:date" json:"fecha_examen,omitempty"`
	CreadoEn      time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (PacienteExamenResultado) TableName() string { return "paciente_examenes_resultados" }

type PacienteImagen struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint      `gorm:"not null;index" json:"paciente_id"`
	MedicoID      *uint     `gorm:"index" json:"medico_id,omitempty"`
	NombreArchivo string    `gorm:"size:255;not null" json:"nombre_archivo"`
	UrlArchivo    string    `gorm:"size:500;not null" json:"url_archivo"`
	TipoImagen    int       `gorm:"not null;default:1" json:"tipo_imagen"`
	Descripcion   string    `gorm:"size:255" json:"descripcion,omitempty"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (PacienteImagen) TableName() string { return "paciente_imagenes" }

type PacienteCertificado struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint      `gorm:"not null;index" json:"paciente_id"`
	NombreArchivo string    `gorm:"size:255;not null" json:"nombre_archivo"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (PacienteCertificado) TableName() string { return "paciente_certificados" }
