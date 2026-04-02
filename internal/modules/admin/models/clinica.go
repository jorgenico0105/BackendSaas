package models

import "time"

type Clinica struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre             string    `gorm:"size:150;not null" json:"nombre"`
	Ruc                string    `gorm:"size:13;uniqueIndex" json:"ruc,omitempty"`
	RazonSocial        string    `gorm:"size:200" json:"razon_social,omitempty"`
	Direccion          string    `gorm:"size:250" json:"direccion,omitempty"`
	Ciudad             string    `gorm:"size:100" json:"ciudad,omitempty"`
	Provincia          string    `gorm:"size:100" json:"provincia,omitempty"`
	Pais               string    `gorm:"size:100;default:'Ecuador'" json:"pais"`
	Telefono           string    `gorm:"size:20" json:"telefono,omitempty"`
	Correo             string    `gorm:"size:150" json:"correo,omitempty"`
	SitioWeb           string    `gorm:"size:150" json:"sitio_web,omitempty"`
	RepresentanteLegal string    `gorm:"size:150" json:"representante_legal,omitempty"`
	TipoClinica        string    `gorm:"size:100" json:"tipo_clinica,omitempty"`
	State              string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn           time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn      time.Time `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
	Estilo             EstiloClinica
}
type EstiloClinica struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID     uint   `gorm:"not null;index" json:"clinica_id"`
	UsuarioID     uint   `gorm:"not null;index" json:"usuario_id"`
	NombreArchivo string `gorm:"size:255" json:"nombre_archivo,omitempty"`
	UrlArchivo    string `gorm:"size:500" json:"url_archivo,omitempty"`
	LogoPath      string `gorm:"size:500" json:"logo_path,omitempty"`
	TipoLogo      string `gorm:"size:30" json:"tipo_logo,omitempty"`

	PrimaryColor   string `gorm:"size:20" json:"primary_color,omitempty"`
	SecondaryColor string `gorm:"size:20" json:"secondary_color,omitempty"`
	ThirdColor     string `gorm:"size:20" json:"third_color,omitempty"`

	Surface1 string `gorm:"size:20;column:surface_1" json:"surface_1,omitempty"`
	Surface2 string `gorm:"size:20;column:surface_2" json:"surface_2,omitempty"`
	Surface3 string `gorm:"size:20;column:surface_3" json:"surface_3,omitempty"`
	Surface4 string `gorm:"size:20;column:surface_4" json:"surface_4,omitempty"`
	Surface5 string `gorm:"size:20;column:surface_5" json:"surface_5,omitempty"`

	PrimaryRamp   string `gorm:"size:20;column:primary_ramp" json:"primary_ramp,omitempty"`
	SecondaryRamp string `gorm:"size:20;column:secondary_ramp" json:"secondary_ramp,omitempty"`
	TertiaryRamp  string `gorm:"size:20;column:tertiary_ramp" json:"tertiary_ramp,omitempty"`
	ErrorColor    string `gorm:"size:20;column:error_color" json:"error_color,omitempty"`

	TextHigh  string `gorm:"size:20;column:text_high" json:"text_high,omitempty"`
	TextMuted string `gorm:"size:20;column:text_muted" json:"text_muted,omitempty"`
	TextHint  string `gorm:"size:20;column:text_hint" json:"text_hint,omitempty"`

	ChipColor     string `gorm:"size:30;column:chip_color" json:"chip_color,omitempty"`
	GhostBorder   string `gorm:"size:50;column:ghost_border" json:"ghost_border,omitempty"`
	AmbientShadow string `gorm:"size:100;column:ambient_shadow" json:"ambient_shadow,omitempty"`

	DarkMode bool      `gorm:"default:false" json:"dark_mode"`
	EsActivo bool      `gorm:"default:true" json:"es_activo"`
	State    string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (EstiloClinica) TableName() string { return "estilos_clinica" }

type Sucursal struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID     uint      `gorm:"not null;index" json:"clinica_id"`
	Clinica       *Clinica  `gorm:"foreignKey:ClinicaID" json:"clinica,omitempty"`
	Nombre        string    `gorm:"size:150;not null" json:"nombre"`
	Codigo        string    `gorm:"size:30" json:"codigo,omitempty"`
	Direccion     string    `gorm:"size:255" json:"direccion,omitempty"`
	Ciudad        string    `gorm:"size:120" json:"ciudad,omitempty"`
	Provincia     string    `gorm:"size:120" json:"provincia,omitempty"`
	Telefono      string    `gorm:"size:30" json:"telefono,omitempty"`
	Correo        string    `gorm:"size:150" json:"correo,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn time.Time `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

func (Sucursal) TableName() string { return "sucursales" }

type Consultorio struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID   uint      `gorm:"not null;index" json:"clinica_id"`
	Clinica     *Clinica  `gorm:"foreignKey:ClinicaID" json:"clinica,omitempty"`
	Nombre      string    `gorm:"size:120;not null" json:"nombre"`
	Codigo      string    `gorm:"size:30" json:"codigo,omitempty"`
	Piso        string    `gorm:"size:30" json:"piso,omitempty"`
	Descripcion string    `gorm:"size:255" json:"descripcion,omitempty"`
	State       string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn    time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

type UsuarioClinica struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UsuarioID uint      `gorm:"not null;index;uniqueIndex:udx_usuario_clinica" json:"usuario_id"`
	ClinicaID uint      `gorm:"not null;index;uniqueIndex:udx_usuario_clinica" json:"clinica_id"`
	RolID     *uint     `gorm:"index" json:"rol_id,omitempty"`
	State     string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn  time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (UsuarioClinica) TableName() string { return "usuarios_clinicas" }

type UsuarioConsultorio struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UsuarioID     uint      `gorm:"not null;index" json:"usuario_id"`
	ConsultorioID uint      `gorm:"not null;index" json:"consultorio_id"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (UsuarioConsultorio) TableName() string { return "usuarios_consultorios" }
