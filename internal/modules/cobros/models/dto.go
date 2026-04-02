package models

import "time"

type CreateCobroRequest struct {
	SesionID    uint    `json:"sesion_id" binding:"required"`
	PacienteID  uint    `json:"id_paciente" binding:"required"`
	MedicoID    uint    `json:"id_medico" binding:"required"`
	ClinicaID   uint    `json:"id_clinica" binding:"required"`
	MontoCobrar float64 `json:"monto_cobrar" binding:"required,min=0"`
	Descuento   float64 `json:"descuento" binding:"omitempty,min=0"`
	Recargo     float64 `json:"recargo" binding:"omitempty,min=0"`
	Observacion string  `json:"observacion" binding:"omitempty,max=255"`
}

type RegistrarPagoRequest struct {
	FechaPago   time.Time `json:"fecha_pago" binding:"required"`
	MontoPagado float64   `json:"monto_pagado" binding:"required,min=0.01"`
	MedioPagoID uint      `json:"medio_pago_id" binding:"required"`
	Referencia  string    `json:"referencia" binding:"omitempty,max=100"`
	Observacion string    `json:"observacion" binding:"omitempty,max=255"`
}

type CreateEgresoRequest struct {
	ClinicaID    uint      `json:"id_clinica" binding:"required"`
	TipoEgresoID uint      `json:"tipo_egreso_id" binding:"required"`
	Fecha        time.Time `json:"fecha" binding:"required"`
	Monto        float64   `json:"monto" binding:"required,min=0.01"`
	Descripcion  string    `json:"descripcion" binding:"omitempty,max=255"`
	Proveedor    string    `json:"proveedor" binding:"omitempty,max=150"`
	Referencia   string    `json:"referencia" binding:"omitempty,max=100"`
}
