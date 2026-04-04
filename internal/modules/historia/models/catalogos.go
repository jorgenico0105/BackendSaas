package models

import "time"

// ─── Catálogos ────────────────────────────────────────────────────────────────

type TipoFormulario struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo      string    `gorm:"type:char(3);uniqueIndex;not null" json:"codigo"`
	Nombre      string    `gorm:"size:100;not null" json:"nombre"`
	Descripcion string    `gorm:"size:255" json:"descripcion,omitempty"`
	State       string    `gorm:"type:char(1);default:'A'" json:"state"`
	CreadoEn    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	CreadoPor   uint      `json:"created_by,omitempty"`
	RoleID      uint      `gorm:"not null;index" json:"rol_id"`
}

func (TipoFormulario) TableName() string { return "tipo_formulario" }

// Códigos de tipo formulario
const (
	TipoFormHC  = "HCL" // Historia Clínica
	TipoFormANM = "ANM" // Anamnesis
	TipoFormSEG = "SEG" // Seguimiento
	TipoFormTST = "TST" // Test
)

type AlergiaCatalogo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"size:150;not null" json:"nombre"`
	Categoria   string `gorm:"size:80" json:"categoria,omitempty"`
	Descripcion string `gorm:"size:255" json:"descripcion,omitempty"`
	State       string `gorm:"type:char(1);default:'A'" json:"state"`
}

func (AlergiaCatalogo) TableName() string { return "alergias_catalogo" }

type TipoAntecedente struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:5;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:100;not null" json:"nombre"`
}

func (TipoAntecedente) TableName() string { return "tipos_antecedente" }

// Códigos de antecedente
const (
	AntecedentePER = "PER" // Personal
	AntecedenteFAM = "FAM" // Familiar
	AntecedenteQUI = "QUI" // Quirúrgico
	AntecedentePAT = "PAT" // Patológico
	AntecedenteFAR = "FAR" // Farmacológico
	AntecedenteOTR = "OTR" // Otro
)

type HabitoCatalogo struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:5;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:100;not null" json:"nombre"`
}

func (HabitoCatalogo) TableName() string { return "habitos_catalogo" }

// Códigos de hábito
const (
	HabitoTAB = "TAB" // Tabaco
	HabitoALC = "ALC" // Alcohol
	HabitoSUE = "SUE" // Sueño
	HabitoEJE = "EJE" // Ejercicio
	HabitoDIE = "DIE" // Dieta
	HabitoCAF = "CAF" // Cafeína
)

type DiagnosticoCatalogo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo      string `gorm:"size:20;uniqueIndex" json:"codigo,omitempty"`
	Nombre      string `gorm:"size:200;not null" json:"nombre"`
	Descripcion string `gorm:"size:255" json:"descripcion,omitempty"`
	Categoria   string `gorm:"size:100" json:"categoria,omitempty"`
	State       string `gorm:"type:char(1);default:'A'" json:"state"`
	CreadoPor   uint   `gorm:"index" json:"creado_por,omitempty"`
}

func (DiagnosticoCatalogo) TableName() string { return "diagnosticos_catalogo" }

type TipoExamen struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"size:150;not null" json:"nombre"`
	State  string `gorm:"type:char(1);default:'A'" json:"state"`
}

func (TipoExamen) TableName() string { return "tipo_examen" }
