package services

import (
	"errors"

	adminModels "saas-medico/internal/modules/admin/models"
	adminRepos "saas-medico/internal/modules/admin/repositories"
	authModels "saas-medico/internal/modules/auth/models"
	authRepos "saas-medico/internal/modules/auth/repositories"
)

var ErrRolNotFound = errors.New("rol no encontrado")

type RolService struct {
	rolRepo   *authRepos.RolRepository
	adminRepo *adminRepos.AdminRepository
}

func NewRolService(rolRepo *authRepos.RolRepository, adminRepo *adminRepos.AdminRepository) *RolService {
	return &RolService{rolRepo: rolRepo, adminRepo: adminRepo}
}

func (s *RolService) List() ([]authModels.Rol, error) {
	return s.rolRepo.FindAll()
}

func (s *RolService) GetByID(id uint) (*authModels.Rol, error) {
	rol, err := s.rolRepo.FindByID(id)
	if err != nil {
		return nil, ErrRolNotFound
	}
	return rol, nil
}

func (s *RolService) Create(req authModels.CreateRolRequest) (*authModels.Rol, error) {
	rol := &authModels.Rol{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		State:       "A",
	}
	return rol, s.rolRepo.Create(rol)
}

func (s *RolService) Update(id uint, req authModels.CreateRolRequest) (*authModels.Rol, error) {
	rol, err := s.rolRepo.FindByID(id)
	if err != nil {
		return nil, ErrRolNotFound
	}
	rol.Nombre = req.Nombre
	rol.Descripcion = req.Descripcion
	return rol, s.rolRepo.Update(rol)
}

// func (s *RolService) Delete(id uint) error {
// 	if _, err := s.rolRepo.FindByID(id); err != nil {
// 		return ErrRolNotFound
// 	}
// 	return s.rolRepo.SoftDelete(id)
// }

// ─── RolTransaccion ───────────────────────────────────────────────────────────

func (s *RolService) ListTransaccionesByRol(rolID uint) ([]adminModels.RolTransaccion, error) {
	return s.adminRepo.FindTransaccionesByRol(rolID)
}

func (s *RolService) AsignarTransacciones(rolID uint, transaccionIDs []uint) error {
	for _, tid := range transaccionIDs {
		rt := &adminModels.RolTransaccion{RolID: rolID, TransaccionID: tid, State: "A"}
		if err := s.adminRepo.AsignarTransaccionArol(rt); err != nil {
			return err
		}
	}
	return nil
}

func (s *RolService) RevocarTransaccion(rolID, transaccionID uint) error {
	return s.adminRepo.RevocarTransaccionDeRol(rolID, transaccionID)
}

func (s *RolService) ListTransacciones(clinicaID *uint) ([]adminModels.Transaccion, error) {
	return s.adminRepo.FindTransacciones(clinicaID)
}

// ─── UsuarioRol ───────────────────────────────────────────────────────────────

func (s *RolService) ListRolesByUsuario(usuarioID, clinicaID uint) ([]authModels.UsuarioRol, error) {
	return s.rolRepo.FindRolesByUser(usuarioID, clinicaID)
}

func (s *RolService) AsignarRolAUsuario(usuarioID, rolID, clinicaID uint) error {
	return s.rolRepo.AssignRolToUser(&authModels.UsuarioRol{
		UsuarioID: usuarioID,
		RolID:     rolID,
		ClinicaID: clinicaID,
		State:     "A",
	})
}

func (s *RolService) RevocarRolDeUsuario(usuarioID, rolID, clinicaID uint) error {
	return s.rolRepo.RemoveRolFromUser(usuarioID, rolID, clinicaID)
}
