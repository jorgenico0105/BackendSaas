package models

import "time"

// ─── Formularios dinámicos ─────────────────────────────────────────────────────

type Formulario struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre           string    `gorm:"size:150;not null" json:"nombre"`
	Descripcion      string    `gorm:"size:255" json:"descripcion,omitempty"`
	ProfesionID      *uint     `gorm:"index" json:"profesion_id,omitempty"`
	ClinicaID        *uint     `gorm:"index" json:"clinica_id,omitempty"`
	UsuarioID        uint      `gorm:"not null;index" json:"usuario_id"`
	TipoFormularioID uint      `gorm:"not null;index" json:"tipo_formulario_id"`
	State            string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn         time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	Preguntas        []FormularioPregunta
}

func (Formulario) TableName() string { return "formularios" }

// Tipos de respuesta posibles para FormularioPregunta
const (
	TipoRespuestaText        = "TEXT"
	TipoRespuestaNumber      = "NUMBER"
	TipoRespuestaDate        = "DATE"
	TipoRespuestaSelect      = "SELECT"
	TipoRespuestaMultiselect = "MULTISELECT"
	TipoRespuestaBoolean     = "BOOLEAN"
)

type FormularioPregunta struct {
	ID            uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	FormularioID  uint               `gorm:"not null;index:idx_form_orden" json:"formulario_id"`
	Pregunta      string             `gorm:"size:255;not null" json:"pregunta"`
	TipoRespuesta string             `gorm:"size:30;not null" json:"tipo_respuesta"`
	Obligatorio   bool               `gorm:"default:false" json:"obligatorio"`
	Orden         int                `gorm:"default:0;index:idx_form_orden" json:"orden"`
	State         string             `gorm:"type:char(1);default:'A';not null" json:"state"`
	Puntua        bool               `gorm:"default:false" json:"puntua"`
	Peso          float64            `gorm:"type:decimal(10,2);default:1" json:"peso"`
	MinVal        *float64           `gorm:"type:decimal(10,2)" json:"min_val,omitempty"`
	MaxVal        *float64           `gorm:"type:decimal(10,2)" json:"max_val,omitempty"`
	PermiteMulti  bool               `gorm:"default:false" json:"permite_multi"`
	Opciones      []FormularioOpcion `gorm:"foreignKey:PreguntaID"`
}

func (FormularioPregunta) TableName() string { return "formulario_preguntas" }

type FormularioOpcion struct {
	ID         uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	PreguntaID uint    `gorm:"not null;index" json:"pregunta_id"`
	Valor      string  `gorm:"size:100;not null" json:"valor"`
	Etiqueta   string  `gorm:"size:150;not null" json:"etiqueta"`
	Orden      int     `gorm:"default:0" json:"orden"`
	Puntos     float64 `gorm:"type:decimal(10,2);default:0" json:"puntos"`
	State      string  `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (FormularioOpcion) TableName() string { return "formulario_opciones" }
