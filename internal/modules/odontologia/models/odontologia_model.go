package models

import (
	"time"
)

// Base model que puedes extender con tus tablas existentes
type OdontologiaBase struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Ejemplo de modelo - reemplazar con tus modelos reales
// type HistorialDental struct {
// 	OdontologiaBase
// 	PacienteID   uint   `gorm:"not null" json:"paciente_id"`
// 	Diagnostico  string `gorm:"type:text" json:"diagnostico"`
// 	Tratamiento  string `gorm:"type:text" json:"tratamiento"`
// 	// ... más campos
// }
