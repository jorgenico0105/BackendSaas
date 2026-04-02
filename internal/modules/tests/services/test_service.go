package services

import (
	"errors"

	"saas-medico/internal/modules/tests/models"
	"saas-medico/internal/modules/tests/repositories"
)

var (
	ErrTestNotFound   = errors.New("test no encontrado")
	ErrReglaNotFound  = errors.New("regla no encontrada")
)

type TestService struct {
	repo *repositories.TestRepository
}

func NewTestService(repo *repositories.TestRepository) *TestService {
	return &TestService{repo: repo}
}

// ─── Reglas ───────────────────────────────────────────────────────────────────

func (s *TestService) ListReglas(formularioID uint) ([]models.TestRegla, error) {
	return s.repo.FindReglasByFormulario(formularioID)
}

func (s *TestService) GetReglaDetalles(reglaID uint) ([]models.TestReglaDetalle, error) {
	return s.repo.FindReglaDetalles(reglaID)
}

func (s *TestService) CreateRegla(r *models.TestRegla) error {
	return s.repo.CreateRegla(r)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func (s *TestService) ListTestsByPaciente(pacienteID uint) ([]models.Test, error) {
	return s.repo.FindTestsByPaciente(pacienteID)
}

func (s *TestService) GetTest(id uint) (*models.Test, []models.TestRespuesta, error) {
	t, err := s.repo.FindTestByID(id)
	if err != nil {
		return nil, nil, ErrTestNotFound
	}
	respuestas, err := s.repo.FindRespuestasByTest(id)
	if err != nil {
		return t, nil, err
	}
	return t, respuestas, nil
}

func (s *TestService) CreateTest(t *models.Test, respuestas []models.TestRespuesta) error {
	if err := s.repo.CreateTest(t); err != nil {
		return err
	}
	if len(respuestas) > 0 {
		for i := range respuestas {
			respuestas[i].TestID = t.ID
		}
		return s.repo.CreateRespuestas(respuestas)
	}
	return nil
}

func (s *TestService) LinkTestToSesion(sesionID, testID uint) error {
	return s.repo.LinkTestToSesion(&models.SesionTest{
		SesionID: sesionID,
		TestID:   testID,
	})
}

func (s *TestService) GetTestsBySesion(sesionID uint) ([]models.SesionTest, error) {
	return s.repo.FindTestsBySesion(sesionID)
}
