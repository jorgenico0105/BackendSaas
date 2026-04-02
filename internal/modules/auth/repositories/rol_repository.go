package repositories

import (
	"saas-medico/internal/modules/auth/models"

	"gorm.io/gorm"
)

type RolRepository struct {
	db *gorm.DB
}

func NewRolRepository(db *gorm.DB) *RolRepository {
	return &RolRepository{db: db}
}

func (r *RolRepository) Create(rol *models.Rol) error {
	return r.db.Create(rol).Error
}

func (r *RolRepository) FindByID(id uint) (*models.Rol, error) {
	var rol models.Rol
	if err := r.db.First(&rol, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &rol, nil
}

func (r *RolRepository) FindByNombre(nombre string) (*models.Rol, error) {
	var rol models.Rol
	if err := r.db.Where("nombre = ? AND state = 'A'", nombre).First(&rol).Error; err != nil {
		return nil, err
	}
	return &rol, nil
}

func (r *RolRepository) FindAll() ([]models.Rol, error) {
	var roles []models.Rol
	if err := r.db.Where("state = 'A'").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RolRepository) Update(rol *models.Rol) error {
	return r.db.Save(rol).Error
}

func (r *RolRepository) SoftDelete(id uint) error {
	return r.db.Model(&models.Rol{}).Where("id = ?", id).Update("state", "I").Error
}

// UsuarioRol

func (r *RolRepository) AssignRolToUser(ur *models.UsuarioRol) error {
	return r.db.Create(ur).Error
}

func (r *RolRepository) FindRolesByUser(usuarioID, clinicaID uint) ([]models.UsuarioRol, error) {
	var list []models.UsuarioRol
	err := r.db.Where("usuario_id = ? AND clinica_id = ? AND state = 'A'", usuarioID, clinicaID).Find(&list).Error
	return list, err
}

func (r *RolRepository) RemoveRolFromUser(usuarioID, rolID, clinicaID uint) error {
	return r.db.Model(&models.UsuarioRol{}).
		Where("usuario_id = ? AND rol_id = ? AND clinica_id = ?", usuarioID, rolID, clinicaID).
		Update("state", "I").Error
}
