package models

import "time"

type TipoCita struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre       string `gorm:"size:100;not null" json:"nombre"`
	State        string `gorm:"type:char(1);default:'A'" json:"state"`
	RolID        uint   `gorm:"not null;index;column:id_rol" json:"id_rol"`
	NeedHistoria bool   `gorm:"default:false" json:"need_historia"`
}

func (TipoCita) TableName() string { return "tipo_citas" }

type EstadoCita struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"size:5;uniqueIndex;not null" json:"codigo"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (EstadoCita) TableName() string { return "estado_citas" }

// Códigos de estado de cita
const (
	CitaPendiente  = "PE"
	CitaConfirmada = "CF"
	CitaAtendida   = "AT"
	CitaCancelada  = "CA"
	CitaNoAsistio  = "NA"
)

// PacienteRef — lightweight read-only reference to the pacientes table, used for preloading in Cita.
type PacienteRef struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Nombres   string `gorm:"column:nombres" json:"nombres"`
	Apellidos string `gorm:"column:apellidos" json:"apellidos"`
	Telefono  string `gorm:"column:telefono" json:"telefono,omitempty"`
}
type PrePacienteRef struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Nombres   string `gorm:"column:nombres" json:"nombres"`
	Apellidos string `gorm:"column:apellidos" json:"apellidos"`
	Telefono  string `gorm:"column:telefono" json:"telefono,omitempty"`
}

func (PacienteRef) TableName() string    { return "pacientes" }
func (PrePacienteRef) TableName() string { return "pre_pacientes" }

type Cita struct {
	ID            uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Fecha         time.Time       `gorm:"not null" json:"fecha"`
	Hora          string          `gorm:"size:8;not null" json:"hora"`
	DuracionMin   int             `gorm:"default:30" json:"duracion_min"`
	MedicoID      uint            `gorm:"not null;index;column:id_medico" json:"id_medico"`
	PacienteID    uint            `gorm:"not null;index;column:id_paciente" json:"id_paciente"`
	Paciente      *PacienteRef    `gorm:"foreignKey:PacienteID;references:ID" json:"paciente,omitempty"`
	ClinicaID     uint            `gorm:"not null;index;column:id_clinica" json:"id_clinica"`
	TipoCitaID    uint            `gorm:"not null;index" json:"tipo_cita_id"`
	TipoCita      *TipoCita       `gorm:"foreignKey:TipoCitaID" json:"tipo_cita,omitempty"`
	EstadoCitaID  uint            `gorm:"not null;index" json:"estado_cita_id"`
	EstadoCita    *EstadoCita     `gorm:"foreignKey:EstadoCitaID" json:"estado_cita,omitempty"`
	SucursalID    *uint           `gorm:"index" json:"sucursal_id,omitempty"`
	ConsultorioID *uint           `gorm:"index" json:"consultorio_id,omitempty"`
	PrePacienteID *uint           `gorm:"index" json:"pre_paciente_id,omitempty"`
	PrePaciente   *PrePacienteRef `gorm:"foreignKey:PrePacienteID;references:ID" json:"pre-paciente,omitempty"`
	Motivo        string          `gorm:"size:255" json:"motivo,omitempty"`
	UrlSesion     string          `gorm:"size:250" json:"url_sesion,omitempty"`
	Notificado    bool            `gorm:"default:false" json:"notificado"`
	State         string          `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time       `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time       `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

type Sesion struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CitaID        uint       `gorm:"not null;uniqueIndex" json:"cita_id"`
	Cita          *Cita      `gorm:"foreignKey:CitaID;references:ID" json:"cita,omitempty"`
	Inicio        time.Time  `gorm:"not null" json:"inicio"`
	Fin           *time.Time `json:"fin,omitempty"`
	Resumen       string     `gorm:"type:text" json:"resumen,omitempty"`
	Conclusiones  string     `gorm:"type:text" json:"conclusiones,omitempty"`
	State         string     `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time  `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time  `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

type HorarioMedico struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	MedicoID      uint   `gorm:"not null;index" json:"medico_id"`
	ClinicaID     *uint  `gorm:"index" json:"clinica_id,omitempty"`
	ConsultorioID *uint  `gorm:"index" json:"consultorio_id,omitempty"`
	DiaSemana     int    `gorm:"not null" json:"dia_semana"`
	HoraInicio    string `gorm:"size:8;not null" json:"hora_inicio"`
	HoraFin       string `gorm:"size:8;not null" json:"hora_fin"`
	IntervaloMin  int    `gorm:"default:30" json:"intervalo_min"`
	State         string `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (HorarioMedico) TableName() string { return "horarios_medico" }

type BloqueoAgenda struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID     uint      `gorm:"not null;index" json:"clinica_id"`
	SucursalID    *uint     `gorm:"index" json:"sucursal_id,omitempty"`
	ConsultorioID *uint     `gorm:"index" json:"consultorio_id,omitempty"`
	MedicoID      uint      `gorm:"not null;index" json:"medico_id"`
	FechaInicio   time.Time `gorm:"not null" json:"fecha_inicio"`
	FechaFin      time.Time `gorm:"not null" json:"fecha_fin"`
	Motivo        string    `gorm:"size:255" json:"motivo,omitempty"`
	TipoBloqueo   string    `gorm:"size:30" json:"tipo_bloqueo,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (BloqueoAgenda) TableName() string { return "bloqueos_agenda" }
