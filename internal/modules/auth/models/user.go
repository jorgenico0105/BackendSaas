package models

import (
	adminModels "saas-medico/internal/modules/admin/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre          string                `gorm:"size:100;not null" json:"nombre"`
	Apellidos       string                `gorm:"size:100;not null" json:"apellidos"`
	Cedula          string                `gorm:"size:100;not null" json:"cedula"`
	Username        string                `gorm:"size:20;uniqueIndex" json:"username,omitempty"`
	Email           string                `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password        string                `gorm:"size:255;not null" json:"-"`
	Sexo            string                `gorm:"size:10" json:"sexo,omitempty"`
	CodigoProfesion int                   `json:"codigo_profesion,omitempty"`
	Universidad     string                `gorm:"size:50" json:"universidad,omitempty"`
	Celular         string                `gorm:"size:20" json:"celular,omitempty"`
	Foto            string                `gorm:"size:250" json:"foto,omitempty"`
	RolID           uint                  `gorm:"not null" json:"rol_id"`
	Rol             Rol                   `gorm:"foreignKey:RolID" json:"rol,omitempty"`
	ClinicaID       *uint                 `gorm:"index" json:"clinica_id,omitempty"`
	Clinicas        []adminModels.Clinica `gorm:"many2many:usuarios_clinicas;joinForeignKey:UsuarioID;joinReferences:ClinicaID" json:"clinicas"`
	State           string                `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreatedBy       uint                  `json:"created_by,omitempty"`
	CreatedAt       time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time             `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt        `gorm:"index" json:"-"`
	Roles           []*UsuarioRol         `gorm:"many2many:userio_rol;"`
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) FullName() string {
	return u.Nombre + " " + u.Apellidos
}

func (u *User) IsActive() bool {
	return u.State == "A"
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Nombre    string    `json:"nombre"`
	Apellidos string    `json:"apellidos"`
	Email     string    `json:"email"`
	Celular   string    `json:"celular,omitempty"`
	Foto      string    `json:"foto,omitempty"`
	State     string    `json:"state"`
	Rol       string    `json:"rol"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Nombre:    u.Nombre,
		Apellidos: u.Apellidos,
		Email:     u.Email,
		Celular:   u.Celular,
		Foto:      u.Foto,
		State:     u.State,
		Rol:       u.Rol.Nombre,
		CreatedAt: u.CreatedAt,
	}
}
