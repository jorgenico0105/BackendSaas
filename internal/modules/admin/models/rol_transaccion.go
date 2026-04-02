package models

import "time"

// RolTransaccion define qué transacciones (rutas/funciones del sistema)
// tiene acceso un rol. Es la tabla pivot del sistema de permisos basado en
// transacciones (menú/acciones del frontend).
type RolTransaccion struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RolID          uint      `gorm:"not null;index;uniqueIndex:udx_rol_transaccion" json:"rol_id"`
	TransaccionID  uint      `gorm:"not null;index;uniqueIndex:udx_rol_transaccion" json:"transaccion_id"`
	State          string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn       time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (RolTransaccion) TableName() string { return "rol_transaccion" }
