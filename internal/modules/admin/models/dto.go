package models

// Clinica DTOs
type CreateClinicaRequest struct {
	Nombre             string `json:"nombre" binding:"required,min=2,max=150"`
	Ruc                string `json:"ruc" binding:"omitempty,max=13"`
	RazonSocial        string `json:"razon_social" binding:"omitempty,max=200"`
	Direccion          string `json:"direccion" binding:"omitempty,max=250"`
	Ciudad             string `json:"ciudad" binding:"omitempty,max=100"`
	Provincia          string `json:"provincia" binding:"omitempty,max=100"`
	Pais               string `json:"pais" binding:"omitempty,max=100"`
	Telefono           string `json:"telefono" binding:"omitempty,max=20"`
	Correo             string `json:"correo" binding:"omitempty,email"`
	SitioWeb           string `json:"sitio_web" binding:"omitempty,max=150"`
	RepresentanteLegal string `json:"representante_legal" binding:"omitempty,max=150"`
	TipoClinica        string `json:"tipo_clinica" binding:"omitempty,max=100"`
}

// Consultorio DTOs
type CreateConsultorioRequest struct {
	Nombre      string `json:"nombre" binding:"required,min=2,max=120"`
	Codigo      string `json:"codigo" binding:"omitempty,max=30"`
	Piso        string `json:"piso" binding:"omitempty,max=30"`
	Descripcion string `json:"descripcion" binding:"omitempty,max=255"`
}

// UsuarioConsultorio DTOs
type AsignarUsuarioConsultorioRequest struct {
	UsuarioID uint `json:"usuario_id" binding:"required"`
}

// Profesion DTOs
type CreateProfesionRequest struct {
	Nombre      string `json:"nombre" binding:"required,min=2,max=100"`
	Descripcion string `json:"descripcion" binding:"omitempty,max=255"`
}

// PlanSaas DTOs
type CreatePlanSaasRequest struct {
	Codigo        string  `json:"codigo" binding:"required,max=30"`
	Nombre        string  `json:"nombre" binding:"required,max=120"`
	Descripcion   string  `json:"descripcion" binding:"omitempty,max=255"`
	PrecioMensual float64 `json:"precio_mensual" binding:"required,min=0"`
	PrecioAnual   float64 `json:"precio_anual" binding:"omitempty,min=0"`
	MaxUsuarios   *int    `json:"max_usuarios"`
	MaxPacientes  *int    `json:"max_pacientes"`
}

// UsuarioClinica DTOs
type AsignarUsuarioClinicaRequest struct {
	UsuarioID uint  `json:"usuario_id" binding:"required"`
	RolID     *uint `json:"rol_id"`
}

// UsuarioRol DTOs
type AsignarRolUsuarioRequest struct {
	RolID     uint `json:"rol_id" binding:"required"`
	ClinicaID uint `json:"clinica_id" binding:"required"`
}

// RolTransaccion DTOs
type AsignarTransaccionesRolRequest struct {
	TransaccionIDs []uint `json:"transaccion_ids" binding:"required,min=1"`
}
