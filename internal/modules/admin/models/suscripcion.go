package models

import "time"

type Profesion struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre        string    `gorm:"size:100;uniqueIndex;not null" json:"nombre"`
	Descripcion   string    `gorm:"size:255" json:"descripcion,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

type PlanSaas struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo        string    `gorm:"size:30;uniqueIndex;not null" json:"codigo"`
	Nombre        string    `gorm:"size:120;not null" json:"nombre"`
	Descripcion   string    `gorm:"size:255" json:"descripcion,omitempty"`
	PrecioMensual float64   `gorm:"type:decimal(10,2);default:0" json:"precio_mensual"`
	PrecioAnual   float64   `gorm:"type:decimal(10,2);default:0" json:"precio_anual"`
	MaxUsuarios   *int      `json:"max_usuarios,omitempty"`
	MaxPacientes  *int      `json:"max_pacientes,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (PlanSaas) TableName() string { return "planes_saas" }

type EstadoSuscripcion struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:20;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (EstadoSuscripcion) TableName() string { return "estados_suscripcion" }

// Códigos predefinidos de estado de suscripción
const (
	SuscripcionPrueba    = "PRUEBA"
	SuscripcionActiva    = "ACTIVA"
	SuscripcionPausada   = "PAUSADA"
	SuscripcionVencida   = "VENCIDA"
	SuscripcionCancelada = "CANCELADA"
)

type Suscripcion struct {
	ID            uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	PlanID        uint               `gorm:"not null;index" json:"plan_id"`
	Plan          *PlanSaas          `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	ClinicaID     *uint              `gorm:"index" json:"clinica_id,omitempty"`
	UsuarioID     *uint              `gorm:"index" json:"usuario_id,omitempty"`
	EstadoID      uint               `gorm:"not null;index" json:"estado_id"`
	Estado        *EstadoSuscripcion `gorm:"foreignKey:EstadoID" json:"estado,omitempty"`
	Inicio        time.Time          `gorm:"not null" json:"inicio"`
	Fin           *time.Time         `json:"fin,omitempty"`
	ProximoCobro  time.Time          `gorm:"not null" json:"proximo_cobro"`
	GraciaHasta   time.Time          `gorm:"not null" json:"gracia_hasta"`
	State         string             `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time          `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time          `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

type Transaccion struct {
	ID            uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre        string       `gorm:"size:120;not null" json:"nombre"`
	PadreID       *uint        `gorm:"index" json:"padre_id,omitempty"`
	Padre         *Transaccion `gorm:"foreignKey:PadreID" json:"padre,omitempty"`
	Orden         int          `gorm:"default:0" json:"orden"`
	Ruta          string       `gorm:"size:200" json:"ruta,omitempty"`
	Icono         string       `gorm:"size:80" json:"icono,omitempty"`
	Tipo          string       `gorm:"size:10;not null" json:"tipo"`
	Visible       bool         `gorm:"default:true" json:"visible"`
	General       bool         `gorm:"default:false" json:"general"`
	ClinicaID     *uint        `gorm:"index" json:"clinica_id,omitempty"`
	State         string       `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time    `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time    `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

func (Transaccion) TableName() string { return "transacciones" }

type BloqueoAcceso struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID      *uint      `gorm:"index" json:"clinica_id,omitempty"`
	UsuarioID      *uint      `gorm:"index" json:"usuario_id,omitempty"`
	Motivo         string     `gorm:"size:255;not null" json:"motivo"`
	BloqueadoDesde time.Time  `gorm:"not null" json:"bloqueado_desde"`
	BloqueadoHasta *time.Time `json:"bloqueado_hasta,omitempty"`
	State          string     `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (BloqueoAcceso) TableName() string { return "bloqueos_acceso" }
