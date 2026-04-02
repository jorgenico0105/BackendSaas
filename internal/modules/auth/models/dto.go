package models

import modelsAdmin "saas-medico/internal/modules/admin/models"

// Request DTOs

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required,min=6"`
	UserName  string `json:"username" binding:"min=10"`
	ClinicaID uint   `json:"clinica_id"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Nombre    string `json:"nombre" binding:"required,min=2,max=100"`
	Apellidos string `json:"apellidos" binding:"required,min=2,max=100"`
	Celular   string `json:"celular" binding:"omitempty,max=20"`
	Cedula    string `json:"cedula" binding:"omitempty,max=13"`
	RolID     uint   `json:"rol_id" binding:"required"`
}

// ClinicasResult es un tipo interno para el canal del goroutine (no serializado).
type ClinicasResult struct {
	Clinicas []modelsAdmin.Clinica
	Err      error
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	ClinicaID    uint   `json:"clinica_id"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Nombre    string `json:"nombre" binding:"omitempty,min=2,max=100"`
	Apellidos string `json:"apellidos" binding:"omitempty,min=2,max=100"`
	Celular   string `json:"celular" binding:"omitempty,max=20"`
	Foto      string `json:"foto" binding:"omitempty"`
}

// Response DTOs

type LoginResponse struct {
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"`
	TokenType    string                `json:"token_type"`
	ExpiresIn    int64                 `json:"expires_in"`
	User         UserResponse          `json:"user"`
	Clinicas     []modelsAdmin.Clinica `json:"clinicas"`
}

type TokenClaims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	RolID   uint   `json:"rol_id"`
	RolName string `json:"rol_name"`
}

// MenuItemResponse representa un nodo del menú con sus hijos.
type MenuItemResponse struct {
	ID       uint               `json:"id"`
	Nombre   string             `json:"nombre"`
	Ruta     string             `json:"ruta,omitempty"`
	Icono    string             `json:"icono,omitempty"`
	Tipo     string             `json:"tipo"`
	Orden    int                `json:"orden"`
	Children []*MenuItemResponse `json:"children,omitempty"`
}

// Rol DTOs
type CreateRolRequest struct {
	Nombre      string `json:"nombre" binding:"required,min=2,max=50"`
	Descripcion string `json:"descripcion" binding:"omitempty,max=150"`
}
