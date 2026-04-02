package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"saas-medico/internal/config"
	"saas-medico/internal/modules/auth/models"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
)

type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	RolID        uint   `json:"rol_id"`
	RolName      string `json:"rol_name"`
	ClinicaID    uint   `json:"clinica_id"`
	AplicacionID uint   `json:"aplicacion_id,omitempty"` // solo en tokens de paciente
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTService() *JWTService {
	cfg := config.AppConfig
	return &JWTService{
		secretKey:            []byte(cfg.JWTSecret),
		accessTokenDuration:  time.Duration(cfg.JWTExpirationHours) * time.Hour,
		refreshTokenDuration: time.Duration(cfg.JWTRefreshDays) * 24 * time.Hour,
	}
}

func (s *JWTService) GenerateAccessToken(user *models.User, clinicaID uint) (string, int64, error) {
	expiresAt := time.Now().Add(s.accessTokenDuration)

	claims := JWTClaims{
		UserID:    user.ID,
		Email:     user.Email,
		RolID:     user.RolID,
		RolName:   user.Rol.Nombre,
		ClinicaID: clinicaID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "saas-medico",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", 0, err
	}

	return signed, int64(s.accessTokenDuration.Seconds()), nil
}

func (s *JWTService) GeneratePacienteToken(pacienteID, clinicaID, aplicacionID uint) (string, int64, error) {
	expiresAt := time.Now().Add(s.accessTokenDuration)
	claims := JWTClaims{
		UserID:       pacienteID,
		RolName:      "paciente",
		ClinicaID:    clinicaID,
		AplicacionID: aplicacionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "saas-medico",
			Subject:   fmt.Sprintf("%d", pacienteID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", 0, err
	}
	return signed, int64(s.accessTokenDuration.Seconds()), nil
}

func (s *JWTService) GenerateRefreshToken() (string, time.Time, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, err
	}
	token := base64.URLEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(s.refreshTokenDuration)
	return token, expiresAt, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *JWTService) GetAccessTokenDuration() time.Duration  { return s.accessTokenDuration }
func (s *JWTService) GetRefreshTokenDuration() time.Duration { return s.refreshTokenDuration }
