package repositories

import (
	"saas-medico/internal/modules/auth/models"
	"time"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *TokenRepository) FindByUserID(userID uint) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	if err := r.db.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *TokenRepository) Revoke(token string) error {
	return r.db.Model(&models.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}

func (r *TokenRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

func (r *TokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ? OR revoked = ?", time.Now(), true).Delete(&models.RefreshToken{}).Error
}
