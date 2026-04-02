package models

import "time"

// Base es el struct base para todos los modelos del sistema
type Base struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

const (
	StateActivo   = "A"
	StateInactivo = "I"
)
