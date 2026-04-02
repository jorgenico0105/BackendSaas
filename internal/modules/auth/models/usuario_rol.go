package models

import "time"

// UsuarioRol define el rol de un usuario dentro de una clínica específica.
// Reemplaza el sistema de permisos directos — los accesos se determinan por
// las transacciones asignadas al rol (ver RolTransaccion en admin/models).
type UsuarioRol struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UsuarioID uint      `gorm:"not null;index;uniqueIndex:udx_usuario_rol_clinica" json:"usuario_id"`
	RolID     uint      `gorm:"not null;index;uniqueIndex:udx_usuario_rol_clinica" json:"rol_id"`
	ClinicaID uint      `gorm:"not null;index;uniqueIndex:udx_usuario_rol_clinica" json:"clinica_id"`
	State     string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn  time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (UsuarioRol) TableName() string { return "usuario_rol" }
