package models

import "time"

// ParseFecha intenta parsear una fecha en formato "YYYY-MM-DD" o RFC3339.
// Usa time.Local para evitar desfase cuando el DSN de MySQL usa loc=Local.
func ParseFecha(s string) (time.Time, error) {
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

type CreateCitaRequest struct {
	// Fecha acepta "YYYY-MM-DD" o RFC3339; se parsea en el service
	Fecha         string `json:"fecha" binding:"required"`
	Hora          string `json:"hora" binding:"required,max=8"`
	DuracionMin   int    `json:"duracion_min"`
	MedicoID      uint   `json:"id_medico" binding:"required"`
	// PacienteID puede ser 0 para citas anónimas (pre_paciente_id lleva el contacto)
	PacienteID    uint   `json:"id_paciente"`
	ClinicaID     uint   `json:"id_clinica" binding:"required"`
	TipoCitaID    uint   `json:"tipo_cita_id" binding:"required"`
	SucursalID    *uint  `json:"sucursal_id"`
	ConsultorioID *uint  `json:"consultorio_id"`
	PrePacienteID *uint  `json:"pre_paciente_id"`
	Motivo        string `json:"motivo" binding:"omitempty,max=255"`
	UrlSesion     string `json:"url_sesion" binding:"omitempty,max=250"`
}

type UpdateEstadoCitaRequest struct {
	EstadoCodigo string `json:"estado_codigo" binding:"required,oneof=PE CF AT CA NA"`
}

type UpdateCitaPacienteRequest struct {
	PacienteID uint `json:"paciente_id" binding:"required"`
}

type CreateSesionRequest struct {
	Inicio      time.Time  `json:"inicio" binding:"required"`
	Fin         *time.Time `json:"fin"`
	Resumen     string     `json:"resumen"`
	Conclusiones string    `json:"conclusiones"`
}

type UpdateSesionRequest struct {
	Fin          *time.Time `json:"fin"`
	Resumen      string     `json:"resumen"`
	Conclusiones string     `json:"conclusiones"`
}

type CreateHorarioRequest struct {
	MedicoID      uint   `json:"medico_id" binding:"required"`
	ClinicaID     *uint  `json:"clinica_id"`
	ConsultorioID *uint  `json:"consultorio_id"`
	DiaSemana     int    `json:"dia_semana" binding:"min=0,max=6"`
	HoraInicio    string `json:"hora_inicio" binding:"required"`
	HoraFin       string `json:"hora_fin" binding:"required"`
	IntervaloMin  int    `json:"intervalo_min"`
}

type CreateBloqueoRequest struct {
	ClinicaID     uint      `json:"clinica_id" binding:"required"`
	SucursalID    *uint     `json:"sucursal_id"`
	ConsultorioID *uint     `json:"consultorio_id"`
	MedicoID      uint      `json:"medico_id" binding:"required"`
	FechaInicio   time.Time `json:"fecha_inicio" binding:"required"`
	FechaFin      time.Time `json:"fecha_fin" binding:"required"`
	Motivo        string    `json:"motivo" binding:"omitempty,max=255"`
	TipoBloqueo   string    `json:"tipo_bloqueo" binding:"omitempty,max=30"`
}
