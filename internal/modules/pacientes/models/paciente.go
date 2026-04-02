package models

import (
	"saas-medico/internal/modules/auth/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Paciente struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID          uint       `gorm:"not null;index" json:"clinica_id"`
	Nombres            string     `gorm:"size:100;not null" json:"nombres"`
	Apellidos          string     `gorm:"size:100;not null" json:"apellidos"`
	Sexo               string     `gorm:"size:5" json:"sexo,omitempty"`
	FechaNacimiento    *time.Time `json:"fecha_nacimiento,omitempty"`
	LugarNacimiento    string     `gorm:"size:100" json:"lugar_nacimiento,omitempty"`
	Nacionalidad       *uint      `json:"nacionalidad,omitempty"`
	Direccion          string     `gorm:"size:200" json:"direccion,omitempty"`
	Telefono           string     `gorm:"size:20" json:"telefono,omitempty"`
	Correo             string     `gorm:"size:100" json:"correo,omitempty"`
	ContactoEmergencia string     `gorm:"size:100" json:"contacto_emergencia,omitempty"`
	TelefonoEmergencia string     `gorm:"size:100" json:"telefono_emergencia,omitempty"`
	TipoDocumento      *uint      `json:"tipo_documento,omitempty"`
	NumeroDocumento    string     `gorm:"size:100" json:"numero_documento,omitempty"`
	TipoSangre         string     `gorm:"size:10" json:"tipo_sangre,omitempty"`
	TipoPaciente       uint       `gorm:"default:1" json:"tipo_paciente"`
	Foto               string     `gorm:"size:250" json:"foto,omitempty"`
	State              string     `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreatedBy          uint       `gorm:"not null;index" json:"created_by"`
	Creado             time.Time  `gorm:"autoCreateTime;column:creado" json:"creado"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Paciente) TableName() string { return "pacientes" }

type PrePaciente struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID       uint       `gorm:"not null;index" json:"clinica_id"`
	Nombres         string     `gorm:"size:120;not null" json:"nombres"`
	Apellidos       string     `gorm:"size:120;not null" json:"apellidos"`
	Telefono        string     `gorm:"size:30;not null" json:"telefono"`
	Correo          string     `gorm:"size:150" json:"correo,omitempty"`
	Identificacion  string     `gorm:"size:30" json:"identificacion,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	Sexo            string     `gorm:"type:char(1)" json:"sexo,omitempty"`
	Origen          string     `gorm:"size:50;default:'MANUAL'" json:"origen"`
	Notas           string     `gorm:"size:255" json:"notas,omitempty"`
	State           string     `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn        time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (PrePaciente) TableName() string { return "pre_pacientes" }

// Orígenes de pre-paciente
const (
	OrigenWeb      = "WEB"
	OrigenWhatsApp = "WHATSAPP"
	OrigenManual   = "MANUAL"
)

// PacienteUsuario — credenciales de acceso del paciente a la app de la clínica
type PacienteUsuario struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID uint      `gorm:"not null;index" json:"paciente_id"`
	ClinicaID  uint      `gorm:"not null;index" json:"clinica_id"`
	Username   string    `gorm:"size:100;not null" json:"username"`
	Password   string    `gorm:"size:255;not null" json:"-"`
	State      string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn   time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	Paciente   Paciente
}

func (PacienteUsuario) TableName() string { return "paciente_usuario" }

func (pu *PacienteUsuario) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	pu.Password = string(hashed)
	return nil
}

func (pu *PacienteUsuario) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pu.Password), []byte(password)) == nil
}

// ── Aplicaciones móviles del SaaS ─────────────────────────────────────────────

// Aplicacion — apps móviles habilitadas por clínica
type Aplicacion struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID   uint      `gorm:"not null;index" json:"clinica_id"`
	MedicoID    *uint     `gorm:"index" json:"medico_id,omitempty"`                           // doctor responsable de esta app
	Codigo      string    `gorm:"size:30;not null;uniqueIndex:udx_app_clinica" json:"codigo"` // NUTRICION, PSICOLOGIA, etc.
	Nombre      string    `gorm:"size:100;not null" json:"nombre"`
	Descripcion string    `gorm:"size:255" json:"descripcion,omitempty"`
	State       string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn    time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	Medico      models.User
}

// PacienteAccesoApp — registro de cada vez que un paciente abre/inicia sesión en la app
type PacienteAccesoApp struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID   uint      `gorm:"not null;index:idx_acceso_pac" json:"paciente_id"`
	ClinicaID    uint      `gorm:"not null;index" json:"clinica_id"`
	AplicacionID uint      `gorm:"not null;index" json:"aplicacion_id"`
	Tipo         string    `gorm:"size:20;default:'LOGIN'" json:"tipo"` // LOGIN, ACCESO
	CreadoEn     time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (PacienteAccesoApp) TableName() string { return "paciente_acceso_app" }

func (Aplicacion) TableName() string { return "aplicaciones" }

// PacienteAplicacion — acceso de un paciente a una app, scoped por clínica
type PacienteAplicacion struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint       `gorm:"not null;uniqueIndex:udx_pac_app_clinica;index" json:"paciente_id"`
	AplicacionID  uint       `gorm:"not null;uniqueIndex:udx_pac_app_clinica;index" json:"aplicacion_id"`
	ClinicaID     uint       `gorm:"not null;uniqueIndex:udx_pac_app_clinica;index" json:"clinica_id"`
	State         string     `gorm:"type:char(1);default:'A';not null" json:"state"` // A=activo, I=inactivo
	CreadoPor     uint       `gorm:"not null;index" json:"creado_por"`
	CreadoEn      time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time  `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
	Aplicacion    Aplicacion `gorm:"foreignKey:AplicacionID"`
}

func (PacienteAplicacion) TableName() string { return "paciente_aplicacion" }

// DTOs
type CreateAplicacionRequest struct {
	Codigo      string `json:"codigo" binding:"required,max=30"`
	Nombre      string `json:"nombre" binding:"required,max=100"`
	Descripcion string `json:"descripcion"`
	MedicoID    *uint  `json:"medico_id"` // doctor responsable de esta app
}

type AsignarAplicacionRequest struct {
	AplicacionID uint `json:"aplicacion_id" binding:"required"`
}

type PacienteLoginRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	ClinicaID    uint   `json:"clinica_id" binding:"required"`
	AplicacionID uint   `json:"aplicacion_id" binding:"required"`
}

type CreatePacienteUsuarioRequest struct {
	PacienteID uint   `json:"paciente_id" binding:"required"`
	Username   string `json:"username" binding:"required,min=4,max=100"`
	Password   string `json:"password" binding:"required,min=6"`
}

type PacienteLoginResponse struct {
	AccessToken        string `json:"access_token"`
	TokenType          string `json:"token_type"`
	ExpiresIn          int64  `json:"expires_in"`
	PacienteID         uint   `json:"paciente_id"`
	ClinicaID          uint   `json:"clinica_id"`
	AplicacionID       uint   `json:"aplicacion_id"`
	UserName           string `json:"username"`
	DoctorID           *uint  `json:"doctor_id,omitempty"`
	DoctorNombre       string `json:"doctor_nombre,omitempty"`
	DoctorApellidos    string `json:"doctor_apellidos,omitempty"`
	DoctorEspecialidad string `json:"doctor_especialidad,omitempty"`
	Medico             models.User
}
