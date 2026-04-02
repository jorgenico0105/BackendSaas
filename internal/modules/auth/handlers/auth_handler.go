package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"saas-medico/internal/modules/auth/models"
	"saas-medico/internal/modules/auth/services"
	"saas-medico/internal/shared/responses"
	"saas-medico/internal/shared/uploads"
)

type AuthHandler struct {
	authSvc *services.AuthService
}

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos de login inválidos: "+err.Error())
		return
	}

	result, err := h.authSvc.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			responses.Unauthorized(c, "Credenciales inválidas")
		case errors.Is(err, services.ErrUserInactive):
			responses.Forbidden(c, "Usuario inactivo")
		default:
			responses.InternalError(c, "Error al iniciar sesión")
		}
		return
	}

	responses.Success(c, "Login exitoso", result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos de registro inválidos: "+err.Error())
		return
	}

	user, err := h.authSvc.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyExists):
			responses.BadRequest(c, "El email ya está registrado")
		case errors.Is(err, services.ErrRolNotFound):
			responses.BadRequest(c, "Rol no encontrado")
		default:
			responses.InternalError(c, "Error al registrar usuario")
		}
		return
	}

	responses.Created(c, "Usuario registrado exitosamente", user.ToResponse())
}
func (h *AuthHandler) BuildMenuPaleta(c *gin.Context) {
	clinicaID := c.GetUint("clinicaID")
	rolID := c.GetUint("rolID")

	menu, err := h.authSvc.GetMenuPaleta(clinicaID, rolID)
	if err != nil {
		responses.InternalError(c, "Error al obtener el menú")
		return
	}
	//log.Println("Aqui tengo erl menu ", menu)
	responses.Success(c, "Menú obtenido", menu)
}

func (h *AuthHandler) GetEstiloClinica(c *gin.Context) {
	clinicaID := c.GetUint("clinicaID")

	estilo, err := h.authSvc.GetEstiloClinica(clinicaID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEstiloNoConfigurado):
			responses.NotFound(c, "Estilos no configurados para esta clínica")
		default:
			responses.InternalError(c, "Error al obtener estilos")
		}
		return
	}

	responses.Success(c, "Estilos obtenidos", estilo)
}
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Refresh token requerido")
		return
	}

	result, err := h.authSvc.RefreshToken(req.RefreshToken, req.ClinicaID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidRefreshToken):
			responses.Unauthorized(c, "Refresh token inválido o expirado")
		case errors.Is(err, services.ErrUserInactive):
			responses.Forbidden(c, "Usuario inactivo")
		default:
			responses.InternalError(c, "Error al refrescar token")
		}
		return
	}

	responses.Success(c, "Token refrescado exitosamente", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint("userID")
	if err := h.authSvc.Logout(userID); err != nil {
		responses.InternalError(c, "Error al cerrar sesión")
		return
	}
	responses.Success(c, "Sesión cerrada exitosamente", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.authSvc.GetUserByID(userID)
	if err != nil {
		responses.NotFound(c, "Usuario no encontrado")
		return
	}
	responses.Success(c, "Perfil obtenido", user.ToResponse())
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}

	if err := h.authSvc.ChangePassword(userID, req); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			responses.BadRequest(c, "Contraseña actual incorrecta")
		default:
			responses.InternalError(c, "Error al cambiar contraseña")
		}
		return
	}

	responses.Success(c, "Contraseña actualizada exitosamente", nil)
}

func (h *AuthHandler) UploadFoto(c *gin.Context) {
	userID := c.GetUint("userID")

	fileHeader, err := c.FormFile("foto")
	if err != nil {
		responses.BadRequest(c, "Se requiere el campo 'foto'")
		return
	}

	result, err := uploads.SaveFile(c, fileHeader, "fotos", uploads.AllowedImageTypes)
	if err != nil {
		responses.BadRequest(c, err.Error())
		return
	}

	user, oldFoto, err := h.authSvc.UpdateFoto(userID, result.FilePath)
	if err != nil {
		uploads.DeleteFile(result.FilePath)
		responses.InternalError(c, "Error al actualizar la foto")
		return
	}

	// Eliminar foto anterior si existía
	if oldFoto != "" {
		uploads.DeleteFile(oldFoto)
	}

	responses.Success(c, "Foto actualizada", user.ToResponse())
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Datos inválidos: "+err.Error())
		return
	}

	user, err := h.authSvc.UpdateProfile(userID, req)
	if err != nil {
		responses.InternalError(c, "Error al actualizar perfil")
		return
	}

	responses.Success(c, "Perfil actualizado exitosamente", user.ToResponse())
}
