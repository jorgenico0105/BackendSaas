package services

import (
	"errors"

	"saas-medico/internal/modules/historia/models"
	"saas-medico/internal/modules/historia/repositories"
)

var (
	ErrFormularioNotFound = errors.New("formulario no encontrado")
	ErrHistoriaNotFound   = errors.New("historia clínica no encontrada")
)

type HistoriaService struct {
	repo *repositories.HistoriaRepository
}

func NewHistoriaService(repo *repositories.HistoriaRepository) *HistoriaService {
	return &HistoriaService{repo: repo}
}

// ─── Formularios ──────────────────────────────────────────────────────────────

func (s *HistoriaService) ListFormularios(tipoID, usuarioID, clinicaID uint) ([]models.Formulario, error) {
	return s.repo.FindFormularios(tipoID, usuarioID, clinicaID)
}

func (s *HistoriaService) GetFormulario(id uint) (*models.Formulario, error) {
	return s.repo.FindFormularioByID(id)
}

func (s *HistoriaService) ListTiposFormulario() ([]models.TipoFormulario, error) {
	return s.repo.FindTiposFormulario()
}

func (s *HistoriaService) CreateFormularioCompleto(req models.CreateFormularioRequest, usuarioID, clinicaID uint) (*models.Formulario, error) {
	f := &models.Formulario{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		TipoFormularioID: req.TipoFormularioID,
		ProfesionID:      req.ProfesionID,
		UsuarioID:        usuarioID,
		ClinicaID:        &clinicaID,
		State:            "A",
	}
	if err := s.repo.CreateFormulario(f); err != nil {
		return nil, err
	}
	for i, pReq := range req.Preguntas {
		p := &models.FormularioPregunta{
			FormularioID:  f.ID,
			Pregunta:      pReq.Pregunta,
			TipoRespuesta: pReq.TipoRespuesta,
			Obligatorio:   pReq.Obligatorio,
			Orden:         i,
			State:         "A",
		}
		if err := s.repo.CreatePregunta(p); err != nil {
			return nil, err
		}
		for j, oReq := range pReq.Opciones {
			o := &models.FormularioOpcion{
				PreguntaID: p.ID,
				Valor:      oReq.Valor,
				Etiqueta:   oReq.Etiqueta,
				Orden:      j,
				Puntos:     oReq.Puntos,
				State:      "A",
			}
			if err := s.repo.CreateOpcion(o); err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}

func (s *HistoriaService) UpdateFormularioCompleto(id uint, req models.UpdateFormularioRequest) error {
	if err := s.repo.UpdateFormulario(id, req.Nombre, req.Descripcion); err != nil {
		return err
	}
	// Reemplazar preguntas: soft-delete todas y recrear
	if req.Preguntas != nil {
		if err := s.repo.DeleteOpcionesByFormulario(id); err != nil {
			return err
		}
		if err := s.repo.DeletePreguntasByFormulario(id); err != nil {
			return err
		}
		for i, pReq := range req.Preguntas {
			p := &models.FormularioPregunta{
				FormularioID:  id,
				Pregunta:      pReq.Pregunta,
				TipoRespuesta: pReq.TipoRespuesta,
				Obligatorio:   pReq.Obligatorio,
				Orden:         i,
				State:         "A",
			}
			if err := s.repo.CreatePregunta(p); err != nil {
				return err
			}
			for j, oReq := range pReq.Opciones {
				o := &models.FormularioOpcion{
					PreguntaID: p.ID,
					Valor:      oReq.Valor,
					Etiqueta:   oReq.Etiqueta,
					Orden:      j,
					Puntos:     oReq.Puntos,
					State:      "A",
				}
				if err := s.repo.CreateOpcion(o); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *HistoriaService) DeleteFormulario(id uint) error {
	return s.repo.DeleteFormulario(id)
}

func (s *HistoriaService) CreateFormulario(f *models.Formulario) error {
	return s.repo.CreateFormulario(f)
}

func (s *HistoriaService) GetPreguntasConOpciones(formularioID uint) ([]models.FormularioPregunta, map[uint][]models.FormularioOpcion, error) {
	preguntas, err := s.repo.FindPreguntasByFormulario(formularioID)
	if err != nil {
		return nil, nil, err
	}
	opciones := make(map[uint][]models.FormularioOpcion)
	for _, p := range preguntas {
		if p.TipoRespuesta == models.TipoRespuestaSelect || p.TipoRespuesta == models.TipoRespuestaMultiselect {
			opts, _ := s.repo.FindOpcionesByPregunta(p.ID)
			opciones[p.ID] = opts
		}
	}
	return preguntas, opciones, nil
}

// ─── Historia Clínica ─────────────────────────────────────────────────────────

func (s *HistoriaService) ListHistoriasByPaciente(pacienteID uint) ([]models.HistoriaClinica, error) {
	return s.repo.FindHistoriasByPaciente(pacienteID)
}
func (s *HistoriaService) FindHistoriasClinicasByPaciente(pacienteID uint) ([]models.ResultadoHistoria, error) {
	return s.repo.FindHistoriasClinicasByPaciente(pacienteID)
}

func (s *HistoriaService) GetHistoriaClinicaByUser(userId, clinicId, tipoFomr int) ([]models.Formulario, error) {
	return s.repo.GetHistoriaClinicaByUser(userId, clinicId, tipoFomr)
}

func (s *HistoriaService) GetHistoria(id uint) (*models.HistoriaClinica, error) {
	return s.repo.FindHistoriaByID(id)
}

func (s *HistoriaService) CreateHistoria(req *models.CreateHistoriaClinicaRequest, pacienteId int) (*models.HistoriaClinica, error) {

	hist, err := s.repo.CreateHistoria(req, pacienteId)
	if err != nil {
		return nil, err
	}
	return hist, nil
}

// ─── Alergias ─────────────────────────────────────────────────────────────────

func (s *HistoriaService) ListAlergias(pacienteID uint) ([]models.PacienteAlergia, error) {
	return s.repo.FindAlergiasByPaciente(pacienteID)
}

func (s *HistoriaService) AddAlergia(a *models.PacienteAlergia) error {
	return s.repo.CreateAlergia(a)
}

func (s *HistoriaService) RemoveAlergia(id, pacienteID uint) error {
	return s.repo.SoftDeleteAlergia(id, pacienteID)
}

// ─── Antecedentes ─────────────────────────────────────────────────────────────

func (s *HistoriaService) ListAntecedentes(pacienteID uint) ([]models.PacienteAntecedente, error) {
	return s.repo.FindAntecedentesByPaciente(pacienteID)
}

func (s *HistoriaService) AddAntecedente(a *models.PacienteAntecedente) error {
	return s.repo.CreateAntecedente(a)
}

// ─── Hábitos ──────────────────────────────────────────────────────────────────

func (s *HistoriaService) ListHabitos(pacienteID uint) ([]models.PacienteHabito, error) {
	return s.repo.FindHabitosByPaciente(pacienteID)
}

func (s *HistoriaService) SaveHabito(h *models.PacienteHabito) error {
	return s.repo.SaveHabito(h)
}

// ─── Diagnósticos ─────────────────────────────────────────────────────────────

func (s *HistoriaService) ListDiagnosticos(pacienteID uint) ([]models.PacienteDiagnostico, error) {
	return s.repo.FindDiagnosticosByPaciente(pacienteID)
}

func (s *HistoriaService) AddDiagnostico(d *models.PacienteDiagnostico) error {
	return s.repo.CreateDiagnostico(d)
}

func (s *HistoriaService) UpdateDiagnostico(d *models.PacienteDiagnostico) error {
	return s.repo.UpdateDiagnostico(d)
}

// ─── Imágenes del paciente ────────────────────────────────────────────────────

func (s *HistoriaService) AddPacienteImagen(img *models.PacienteImagen) error {
	return s.repo.CreatePacienteImagen(img)
}

func (s *HistoriaService) ListImagenesPaciente(pacienteID uint) ([]models.PacienteImagen, error) {
	return s.repo.FindImagenesByPaciente(pacienteID)
}
