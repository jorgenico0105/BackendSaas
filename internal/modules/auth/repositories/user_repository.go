package repositories

import (
	"saas-medico/internal/modules/auth/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Rol").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Preload("Rol").Preload("Clinicas").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Rol").Preload("Clinicas").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) SoftDelete(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("state", "I").Error
}

func (r *UserRepository) FindAll(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Where("state = ?", "A").Count(&total)

	offset := (page - 1) * pageSize
	if err := r.db.Preload("Rol").Where("state = ?", "A").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByRol(rolID uint) ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Rol").Where("rol_id = ? AND state = ?", rolID, "A").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
