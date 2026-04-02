package services

import (
	"errors"
	"log"

	adminModels "saas-medico/internal/modules/admin/models"
	repositoriesAdmin "saas-medico/internal/modules/admin/repositories"
	"saas-medico/internal/modules/auth/models"
	"saas-medico/internal/modules/auth/repositories"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound          = errors.New("usuario no encontrado")
	ErrInvalidCredentials    = errors.New("credenciales inválidas")
	ErrUserInactive          = errors.New("usuario inactivo")
	ErrEmailAlreadyExists    = errors.New("el email ya está registrado")
	ErrInvalidRefreshToken   = errors.New("refresh token inválido")
	ErrRolNotFound           = errors.New("rol no encontrado")
	ErrClinicaAccesoDenegado = errors.New("no tiene acceso a esta clínica")
	ErrEstiloNoConfigurado   = errors.New("estilos no configurados para esta clínica")
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	tokenRepo *repositories.TokenRepository
	rolRepo   *repositories.RolRepository
	adminRepo *repositoriesAdmin.AdminRepository
	jwtSvc    *JWTService
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	tokenRepo *repositories.TokenRepository,
	rolRepo *repositories.RolRepository,
	adminRepo *repositoriesAdmin.AdminRepository,
	jwtSvc *JWTService,
) *AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo, rolRepo: rolRepo, adminRepo: adminRepo, jwtSvc: jwtSvc}
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.User
	var err error

	if req.Email != "" {
		user, err = s.userRepo.FindByEmail(req.Email)
	} else {
		user, err = s.userRepo.FindByUsername(req.UserName)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	clinicas, _, err := s.adminRepo.GetClinicasByUser(int(user.ID))
	if err != nil {
		return nil, err
	}

	// Determinar la clínica para el token
	var clinicaID uint
	if req.ClinicaID > 0 {
		for _, c := range clinicas {
			if c.ID == req.ClinicaID {
				clinicaID = c.ID
				break
			}
		}
		if clinicaID == 0 {
			return nil, ErrClinicaAccesoDenegado
		}
	} else if len(clinicas) > 0 {
		clinicaID = clinicas[0].ID
	}

	accessToken, expiresIn, err := s.jwtSvc.GenerateAccessToken(user, clinicaID)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, refreshExpiresAt, err := s.jwtSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: refreshExpiresAt,
	}
	if err := s.tokenRepo.Create(rt); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         user.ToResponse(),
		Clinicas:     clinicas,
	}, nil
}

func (s *AuthService) GetMenuPaleta(clinicaID, rolID uint) ([]*models.MenuItemResponse, error) {
	transacciones, err := s.adminRepo.GetMenuByRolAndClinica(rolID, clinicaID)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(transacciones), nil
}

func (s *AuthService) GetEstiloClinica(clinicaID uint) (*adminModels.EstiloClinica, error) {
	estilo, err := s.adminRepo.FindEstiloByClinica(clinicaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEstiloNoConfigurado
		}
		return nil, err
	}
	return estilo, nil
}

// buildMenuTree convierte una lista plana de transacciones en un árbol jerárquico.
// Elementos cuyo PadreID no está en la lista se tratan como raíces.
func buildMenuTree(items []adminModels.Transaccion) []*models.MenuItemResponse {
	if len(items) == 0 {
		return []*models.MenuItemResponse{}
	}

	nodeMap := make(map[uint]*models.MenuItemResponse, len(items))
	for _, t := range items {
		t := t
		nodeMap[t.ID] = &models.MenuItemResponse{
			ID:     t.ID,
			Nombre: t.Nombre,
			Ruta:   t.Ruta,
			Icono:  t.Icono,
			Tipo:   t.Tipo,
			Orden:  t.Orden,
		}
	}

	var roots []*models.MenuItemResponse
	for _, t := range items {
		node := nodeMap[t.ID]
		if t.PadreID == nil {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[*t.PadreID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			// padre no accesible para este rol → tratamos el hijo como raíz
			roots = append(roots, node)
		}
	}

	return roots
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	log.Println(req)
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, ErrEmailAlreadyExists
	}

	rol, err := s.rolRepo.FindByID(req.RolID)
	if err != nil {
		return nil, ErrRolNotFound
	}
	username := req.Cedula
	user := &models.User{
		Email:     req.Email,
		Nombre:    req.Nombre,
		Username:  username,
		Apellidos: req.Apellidos,
		Celular:   req.Celular,
		RolID:     req.RolID,
		Cedula:    req.Cedula,
		State:     "A",
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Rol = *rol
	return user, nil
}

func (s *AuthService) RefreshToken(refreshTokenStr string, clinicaID uint) (*models.LoginResponse, error) {
	rt, err := s.tokenRepo.FindByToken(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if !rt.IsValid() {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	if err := s.tokenRepo.Revoke(refreshTokenStr); err != nil {
		return nil, err
	}

	accessToken, expiresIn, err := s.jwtSvc.GenerateAccessToken(user, clinicaID)
	if err != nil {
		return nil, err
	}

	newRefreshStr, refreshExpiresAt, err := s.jwtSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRT := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshStr,
		ExpiresAt: refreshExpiresAt,
	}
	if err := s.tokenRepo.Create(newRT); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshStr,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         user.ToResponse(),
	}, nil
}

func (s *AuthService) Logout(userID uint) error {
	return s.tokenRepo.RevokeAllUserTokens(userID)
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) ChangePassword(userID uint, req models.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if !user.CheckPassword(req.CurrentPassword) {
		return ErrInvalidCredentials
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return s.tokenRepo.RevokeAllUserTokens(userID)
}

func (s *AuthService) UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Nombre != "" {
		user.Nombre = req.Nombre
	}
	if req.Apellidos != "" {
		user.Apellidos = req.Apellidos
	}
	if req.Celular != "" {
		user.Celular = req.Celular
	}
	if req.Foto != "" {
		user.Foto = req.Foto
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateFoto reemplaza la foto del usuario. Devuelve también la ruta anterior para que el caller la elimine.
func (s *AuthService) UpdateFoto(userID uint, rutaFoto string) (*models.User, string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, "", ErrUserNotFound
	}
	oldFoto := user.Foto
	user.Foto = rutaFoto
	if err := s.userRepo.Update(user); err != nil {
		return nil, "", err
	}
	return user, oldFoto, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	return s.jwtSvc.ValidateAccessToken(tokenString)
}
