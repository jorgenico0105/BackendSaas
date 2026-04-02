package models

import "time"

// ─── Reglas de puntuación ──────────────────────────────────────────────────────

type TestRegla struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FormularioID uint      `gorm:"not null;uniqueIndex:udx_form_version" json:"formulario_id"`
	Version      int       `gorm:"default:1;uniqueIndex:udx_form_version" json:"version"`
	Nombre       string    `gorm:"size:150;not null" json:"nombre"`
	Descripcion  string    `gorm:"size:255" json:"descripcion,omitempty"`
	State        string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn     time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (TestRegla) TableName() string { return "test_reglas" }

type TestReglaDetalle struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReglaID   uint    `gorm:"not null;index" json:"regla_id"`
	MinVal    float64 `gorm:"type:decimal(10,2);not null" json:"min_val"`
	MaxVal    float64 `gorm:"type:decimal(10,2);not null" json:"max_val"`
	Resultado string  `gorm:"size:150;not null" json:"resultado"`
	Mensaje   string  `gorm:"size:255" json:"mensaje,omitempty"`
	Orden     int     `gorm:"default:0" json:"orden"`
	State     string  `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (TestReglaDetalle) TableName() string { return "test_reglas_detalle" }

// ─── Tests aplicados ──────────────────────────────────────────────────────────

type Test struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID   uint      `gorm:"not null;index" json:"paciente_id"`
	MedicoID     uint      `gorm:"not null;index" json:"medico_id"`
	FormularioID uint      `gorm:"not null;index" json:"formulario_id"`
	ReglaID      uint      `gorm:"not null;index" json:"regla_id"`
	Fecha        time.Time `gorm:"not null" json:"fecha"`
	PuntajeTotal *float64  `gorm:"type:decimal(10,2)" json:"puntaje_total,omitempty"`
	Resultado    string    `gorm:"size:150" json:"resultado,omitempty"`
	Observacion  string    `gorm:"type:text" json:"observacion,omitempty"`
	State        string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn     time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (Test) TableName() string { return "tests" }

type TestRespuesta struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TestID           uint      `gorm:"not null;index" json:"test_id"`
	PreguntaID       uint      `gorm:"not null;index" json:"pregunta_id"`
	OpcionID         *uint     `gorm:"index" json:"opcion_id,omitempty"`
	RespuestaTexto   string    `gorm:"type:text" json:"respuesta_texto,omitempty"`
	RespuestaNumero  *float64  `gorm:"type:decimal(10,2)" json:"respuesta_numero,omitempty"`
	CreadoEn         time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (TestRespuesta) TableName() string { return "test_respuestas" }

type TestArchivo struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TestID        uint      `gorm:"not null;index" json:"test_id"`
	NombreArchivo string    `gorm:"size:255;not null" json:"nombre_archivo"`
	TipoArchivo   string    `gorm:"size:150" json:"tipo_archivo,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (TestArchivo) TableName() string { return "test_archivos" }

// ─── Sesión ↔ Tests ───────────────────────────────────────────────────────────

type SesionTest struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SesionID uint      `gorm:"not null;uniqueIndex:udx_sesion_test" json:"sesion_id"`
	TestID   uint      `gorm:"not null;uniqueIndex:udx_sesion_test" json:"test_id"`
	State    string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (SesionTest) TableName() string { return "sesion_tests" }
