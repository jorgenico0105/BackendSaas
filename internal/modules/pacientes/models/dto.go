package models

import "time"

// ParseFechaNacimiento parsea "YYYY-MM-DD" y retorna nil si el string está vacío.
func ParseFechaNacimiento(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type CreatePacienteRequest struct {
	Nombres            string `json:"nombres" binding:"required,min=2,max=100"`
	Apellidos          string `json:"apellidos" binding:"required,min=2,max=100"`
	Sexo               string `json:"sexo" binding:"omitempty,oneof=M F"`
	FechaNacimiento    string `json:"fecha_nacimiento"`  // "YYYY-MM-DD"
	LugarNacimiento    string `json:"lugar_nacimiento" binding:"omitempty,max=100"`
	Direccion          string `json:"direccion" binding:"omitempty,max=200"`
	Telefono           string `json:"telefono" binding:"omitempty,max=20"`
	Correo             string `json:"correo" binding:"omitempty,email"`
	ContactoEmergencia string `json:"contacto_emergencia" binding:"omitempty,max=100"`
	TelefonoEmergencia string `json:"telefono_emergencia" binding:"omitempty,max=100"`
	NumeroDocumento    string `json:"numero_documento" binding:"omitempty,max=100"`
	TipoSangre         string `json:"tipo_sangre" binding:"omitempty,max=10"`
	TipoPaciente       uint   `json:"tipo_paciente"`
}

type UpdatePacienteRequest struct {
	Nombres            string `json:"nombres" binding:"omitempty,min=2,max=100"`
	Apellidos          string `json:"apellidos" binding:"omitempty,min=2,max=100"`
	Sexo               string `json:"sexo" binding:"omitempty,oneof=M F"`
	FechaNacimiento    string `json:"fecha_nacimiento"`  // "YYYY-MM-DD"
	LugarNacimiento    string `json:"lugar_nacimiento" binding:"omitempty,max=100"`
	Direccion          string `json:"direccion" binding:"omitempty,max=200"`
	Telefono           string `json:"telefono" binding:"omitempty,max=20"`
	Correo             string `json:"correo" binding:"omitempty,email"`
	ContactoEmergencia string `json:"contacto_emergencia" binding:"omitempty,max=100"`
	TelefonoEmergencia string `json:"telefono_emergencia" binding:"omitempty,max=100"`
	NumeroDocumento    string `json:"numero_documento" binding:"omitempty,max=100"`
	TipoSangre         string `json:"tipo_sangre" binding:"omitempty,max=10"`
	Foto               string `json:"foto" binding:"omitempty"`
}

type CreatePrePacienteRequest struct {
	ClinicaID       uint   `json:"clinica_id" binding:"required"`
	Nombres         string `json:"nombres" binding:"required,min=2,max=120"`
	Apellidos       string `json:"apellidos" binding:"required,min=2,max=120"`
	Telefono        string `json:"telefono" binding:"required,max=30"`
	Correo          string `json:"correo" binding:"omitempty,email"`
	Identificacion  string `json:"identificacion" binding:"omitempty,max=30"`
	FechaNacimiento string `json:"fecha_nacimiento"`  // "YYYY-MM-DD"
	Sexo            string `json:"sexo" binding:"omitempty,oneof=M F"`
	Origen          string `json:"origen" binding:"omitempty,oneof=WEB WHATSAPP MANUAL"`
	Notas           string `json:"notas" binding:"omitempty,max=255"`
}
