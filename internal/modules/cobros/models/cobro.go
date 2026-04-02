package models

import "time"

type EstadoCobro struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:5;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (EstadoCobro) TableName() string { return "estados_cobro" }

// Códigos de estado de cobro
const (
	CobroPendiente = "PE"
	CobroParcial   = "PA"
	CobradoCobrado = "CO"
	CobroAnulado   = "AN"
)

type MedioPago struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:5;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (MedioPago) TableName() string { return "medios_pago" }

// Medios de pago
const (
	MedioEfectivo      = "EFE"
	MedioTransferencia = "TRA"
	MedioTarjeta       = "TAR"
	MedioDeposito      = "DEP"
	MedioDeUna         = "DNA"
	MedioAhorita       = "AHO"
)

type TipoEgreso struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:5;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (TipoEgreso) TableName() string { return "tipo_egreso" }

// Tipos de egreso
const (
	EgresoMateriales = "MAT"
	EgresoServicios  = "SER"
	EgresoAlquiler   = "ALQ"
	EgresoConsultas  = "CON"
	EgresoOtros      = "OTR"
)

type CobroSesion struct {
	ID            uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	SesionID      uint         `gorm:"not null;index" json:"sesion_id"`
	PacienteID    uint         `gorm:"not null;index;column:id_paciente" json:"id_paciente"`
	MedicoID      uint         `gorm:"not null;index;column:id_medico" json:"id_medico"`
	ClinicaID     uint         `gorm:"not null;index;column:id_clinica" json:"id_clinica"`
	MontoCobrar   float64      `gorm:"type:decimal(10,2);not null" json:"monto_cobrar"`
	Descuento     float64      `gorm:"type:decimal(10,2);default:0" json:"descuento"`
	Recargo       float64      `gorm:"type:decimal(10,2);default:0" json:"recargo"`
	MontoTotal    float64      `gorm:"type:decimal(10,2);not null" json:"monto_total"`
	EstadoCobroID uint         `gorm:"not null;index" json:"estado_cobro_id"`
	EstadoCobro   *EstadoCobro `gorm:"foreignKey:EstadoCobroID" json:"estado_cobro,omitempty"`
	Observacion   string       `gorm:"size:255" json:"observacion,omitempty"`
	State         string       `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time    `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	Pagos         []Pago       `gorm:"foreignKey:CobroID" json:"pagos,omitempty"`
}

func (CobroSesion) TableName() string { return "cobros_sesion" }

type Pago struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CobroID     uint       `gorm:"not null;index" json:"cobro_id"`
	PacienteID  uint       `gorm:"not null;index;column:id_paciente" json:"id_paciente"`
	FechaPago   time.Time  `gorm:"not null" json:"fecha_pago"`
	MontoPagado float64    `gorm:"type:decimal(10,2);not null" json:"monto_pagado"`
	MedioPagoID uint       `gorm:"not null;index" json:"medio_pago_id"`
	MedioPago   *MedioPago `gorm:"foreignKey:MedioPagoID" json:"medio_pago,omitempty"`
	Referencia  string     `gorm:"size:100" json:"referencia,omitempty"`
	Observacion string     `gorm:"size:255" json:"observacion,omitempty"`
	State       string     `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn    time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

type Egreso struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID    uint        `gorm:"not null;index;column:id_clinica" json:"id_clinica"`
	TipoEgresoID uint        `gorm:"not null;index" json:"tipo_egreso_id"`
	TipoEgreso   *TipoEgreso `gorm:"foreignKey:TipoEgresoID" json:"tipo_egreso,omitempty"`
	Fecha        time.Time   `gorm:"not null" json:"fecha"`
	Monto        float64     `gorm:"type:decimal(10,2);not null" json:"monto"`
	Descripcion  string      `gorm:"size:255" json:"descripcion,omitempty"`
	Proveedor    string      `gorm:"size:150" json:"proveedor,omitempty"`
	Referencia   string      `gorm:"size:100" json:"referencia,omitempty"`
	State        string      `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn     time.Time   `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}
