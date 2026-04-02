package models

import (
	"time"
)

// Base model que puedes extender con tus tablas existentes
type PsicologiaBase struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Ejemplo de modelo - reemplazar con tus modelos reales
// type Paciente struct {
// 	PsicologiaBase
// 	Nombre    string `gorm:"size:100;not null" json:"nombre"`
// 	Apellido  string `gorm:"size:100;not null" json:"apellido"`
// 	// ... más campos
// }
