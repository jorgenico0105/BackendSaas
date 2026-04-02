package models

import (
	adminModels "saas-medico/internal/modules/admin/models"
	"time"
)

type Rol struct {
	ID          uint                          `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre      string                        `gorm:"size:50;uniqueIndex;not null" json:"nombre"`
	Descripcion string                        `gorm:"size:150" json:"descripcion,omitempty"`
	State       string                        `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
	es          *[]adminModels.RolTransaccion `gorm:"many2many:user_languages;"`
}

func (Rol) TableName() string { return "roles" }

// Roles predefinidos del sistema
const (
	RolSuperAdmin    = "super_admin"
	RolAdmin         = "admin_clinica"
	RolMedico        = "medico"
	RolPsicologo     = "psicologo"
	RolNutriologo    = "nutriologo"
	RolOdontologo    = "odontologo"
	RolRecepcionista = "recepcionista"
	RolPaciente      = "paciente"
)
