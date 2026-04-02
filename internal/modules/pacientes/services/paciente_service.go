package services

import (
	"errors"

	"log"
	authServices "saas-medico/internal/modules/auth/services"
	"saas-medico/internal/modules/pacientes/models"
	"saas-medico/internal/modules/pacientes/repositories"
)

var (
	ErrPacienteNotFound      = errors.New("paciente no encontrado")
	ErrPrePacienteNotFound   = errors.New("pre-paciente no encontrado")
	ErrPacienteUsuarioExists = errors.New("el paciente ya tiene usuario en esta clínica")
	ErrPacienteCredenciales  = errors.New("credenciales inválidas")
	ErrSinAccesoAplicacion   = errors.New("el paciente no tiene acceso a esta aplicación")
)

type PacienteService struct {
	repo   *repositories.PacienteRepository
	jwtSvc *authServices.JWTService
}

func NewPacienteService(repo *repositories.PacienteRepository, jwtSvc *authServices.JWTService) *PacienteService {
	return &PacienteService{repo: repo, jwtSvc: jwtSvc}
}

func (s *PacienteService) Create(req models.CreatePacienteRequest, createdBy, clinicaID uint) (*models.Paciente, error) {
	fechaNac, err := models.ParseFechaNacimiento(req.FechaNacimiento)
	if err != nil {
		return nil, errors.New("formato de fecha_nacimiento inválido, use YYYY-MM-DD")
	}
	p := &models.Paciente{
		ClinicaID:          clinicaID,
		Nombres:            req.Nombres,
		Apellidos:          req.Apellidos,
		Sexo:               req.Sexo,
		FechaNacimiento:    fechaNac,
		LugarNacimiento:    req.LugarNacimiento,
		Direccion:          req.Direccion,
		Telefono:           req.Telefono,
		Correo:             req.Correo,
		ContactoEmergencia: req.ContactoEmergencia,
		TelefonoEmergencia: req.TelefonoEmergencia,
		NumeroDocumento:    req.NumeroDocumento,
		TipoSangre:         req.TipoSangre,
		TipoPaciente:       req.TipoPaciente,
		State:              "A",
		CreatedBy:          createdBy,
	}
	if p.TipoPaciente == 0 {
		p.TipoPaciente = 1
	}
	return p, s.repo.Create(p)
}

func (s *PacienteService) GetByID(id uint) (*models.Paciente, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrPacienteNotFound
	}
	return p, nil
}

func (s *PacienteService) List(search string, page, pageSize, clinicaId, usuarioId int) ([]models.Paciente, int64, error) {
	return s.repo.FindAll(search, page, pageSize, clinicaId, usuarioId)
}

func (s *PacienteService) Update(id uint, req models.UpdatePacienteRequest) (*models.Paciente, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrPacienteNotFound
	}
	if req.Nombres != "" {
		p.Nombres = req.Nombres
	}
	if req.Apellidos != "" {
		p.Apellidos = req.Apellidos
	}
	if req.Sexo != "" {
		p.Sexo = req.Sexo
	}
	if req.FechaNacimiento != "" {
		fechaNac, err := models.ParseFechaNacimiento(req.FechaNacimiento)
		if err != nil {
			return nil, errors.New("formato de fecha_nacimiento inválido, use YYYY-MM-DD")
		}
		p.FechaNacimiento = fechaNac
	}
	if req.LugarNacimiento != "" {
		p.LugarNacimiento = req.LugarNacimiento
	}
	if req.Direccion != "" {
		p.Direccion = req.Direccion
	}
	if req.Telefono != "" {
		p.Telefono = req.Telefono
	}
	if req.Correo != "" {
		p.Correo = req.Correo
	}
	if req.ContactoEmergencia != "" {
		p.ContactoEmergencia = req.ContactoEmergencia
	}
	if req.TelefonoEmergencia != "" {
		p.TelefonoEmergencia = req.TelefonoEmergencia
	}
	if req.NumeroDocumento != "" {
		p.NumeroDocumento = req.NumeroDocumento
	}
	if req.TipoSangre != "" {
		p.TipoSangre = req.TipoSangre
	}
	if req.Foto != "" {
		p.Foto = req.Foto
	}
	return p, s.repo.Update(p)
}

func (s *PacienteService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrPacienteNotFound
	}
	return s.repo.SoftDelete(id)
}

// ── PrePaciente ───────────────────────────────────────────────────────────────

func (s *PacienteService) CreatePrePaciente(req models.CreatePrePacienteRequest) (*models.PrePaciente, error) {
	if req.Origen == "" {
		req.Origen = models.OrigenManual
	}
	fechaNacPre, err := models.ParseFechaNacimiento(req.FechaNacimiento)
	if err != nil {
		return nil, errors.New("formato de fecha_nacimiento inválido, use YYYY-MM-DD")
	}
	pp := &models.PrePaciente{
		ClinicaID:       req.ClinicaID,
		Nombres:         req.Nombres,
		Apellidos:       req.Apellidos,
		Telefono:        req.Telefono,
		Correo:          req.Correo,
		Identificacion:  req.Identificacion,
		FechaNacimiento: fechaNacPre,
		Sexo:            req.Sexo,
		Origen:          req.Origen,
		Notas:           req.Notas,
		State:           "A",
	}
	return pp, s.repo.CreatePrePaciente(pp)
}

func (s *PacienteService) ListPrePacientes(clinicaID uint, page, pageSize int) ([]models.PrePaciente, int64, error) {
	return s.repo.FindPrePacientesByClinica(clinicaID, page, pageSize)
}

func (s *PacienteService) DeletePrePaciente(id uint) error {
	if _, err := s.repo.FindPrePacienteByID(id); err != nil {
		return ErrPrePacienteNotFound
	}
	return s.repo.DeletePrePaciente(id)
}

// ── Aplicaciones ──────────────────────────────────────────────────────────────

var ErrAplicacionYaAsignada = errors.New("el paciente ya tiene acceso a esta aplicación en esta clínica")

func (s *PacienteService) ListAplicaciones(clinicaID uint) ([]models.Aplicacion, error) {
	return s.repo.FindAplicaciones(clinicaID)
}

func (s *PacienteService) ListAplicacionesPaciente(pacienteID, clinicaID uint) ([]models.PacienteAplicacion, error) {
	return s.repo.FindAplicacionesByPaciente(pacienteID, clinicaID)
}

func (s *PacienteService) CreateAplicacion(req models.CreateAplicacionRequest, clinicaID uint) (*models.Aplicacion, error) {
	a := &models.Aplicacion{
		ClinicaID:   clinicaID,
		MedicoID:    req.MedicoID,
		Codigo:      req.Codigo,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		State:       "A",
	}
	return a, s.repo.CreateAplicacion(a)
}

func (s *PacienteService) AsignarAplicacion(pacienteID uint, req models.AsignarAplicacionRequest, clinicaID, creadoPor uint) (*models.PacienteAplicacion, error) {
	// Si ya existe, reactivar si estaba inactiva
	existing, err := s.repo.FindPacienteAplicacion(pacienteID, req.AplicacionID, clinicaID)
	if err == nil {
		if existing.State == "A" {
			return nil, ErrAplicacionYaAsignada
		}
		existing.State = "A"
		return existing, s.repo.UpdatePacienteAplicacion(existing)
	}

	pa := &models.PacienteAplicacion{
		PacienteID:   pacienteID,
		AplicacionID: req.AplicacionID,
		ClinicaID:    clinicaID,
		CreadoPor:    creadoPor,
		State:        "A",
	}
	if err := s.repo.CreatePacienteAplicacion(pa); err != nil {
		return nil, err
	}

	// Crear PacienteUsuario si no existe aún en esta clínica
	if !s.repo.ExistsPacienteUsuario(pacienteID, clinicaID) {
		paciente, err := s.repo.FindByID(pacienteID)
		if err != nil {
			return pa, nil // acceso creado; fallo silencioso en usuario (cedula no encontrada)
		}
		username := paciente.NumeroDocumento
		if username == "" {
			username = paciente.Telefono // fallback si no tiene cédula
		}
		if username != "" {
			pu := &models.PacienteUsuario{
				PacienteID: pacienteID,
				ClinicaID:  clinicaID,
				Username:   username,
				State:      "A",
			}
			_ = pu.SetPassword("Usuario123")
			_ = s.repo.CreatePacienteUsuario(pu)
		}
	}

	return pa, nil
}

func (s *PacienteService) RevocarAplicacion(pacienteID, aplicacionID, clinicaID uint) error {
	pa, err := s.repo.FindPacienteAplicacion(pacienteID, aplicacionID, clinicaID)
	if err != nil {
		return errors.New("acceso no encontrado")
	}
	pa.State = "I"
	return s.repo.UpdatePacienteAplicacion(pa)
}

// ── PacienteUsuario ───────────────────────────────────────────────────────────

func (s *PacienteService) CreatePacienteUsuario(req models.CreatePacienteUsuarioRequest, clinicaID uint) (*models.PacienteUsuario, error) {
	if s.repo.ExistsPacienteUsuario(req.PacienteID, clinicaID) {
		return nil, ErrPacienteUsuarioExists
	}
	pu := &models.PacienteUsuario{
		PacienteID: req.PacienteID,
		ClinicaID:  clinicaID,
		Username:   req.Username,
		State:      "A",
	}
	if err := pu.SetPassword(req.Password); err != nil {
		return nil, err
	}
	return pu, s.repo.CreatePacienteUsuario(pu)
}

func (s *PacienteService) CountAccesos(pacienteID uint, desde, hasta string) int64 {
	return s.repo.CountAccesosByPaciente(pacienteID, desde, hasta)
}

func (s *PacienteService) UltimoAcceso(pacienteID uint) string {
	return s.repo.FindUltimoAccesoPaciente(pacienteID)
}

func (s *PacienteService) LoginPaciente(req models.PacienteLoginRequest) (*models.PacienteLoginResponse, error) {
	pu, err := s.repo.FindPacienteUsuario(req.Username, req.ClinicaID)
	if err != nil {
		return nil, ErrPacienteCredenciales
	}
	if !pu.CheckPassword(req.Password) {
		return nil, ErrPacienteCredenciales
	}
	// Verificar acceso a la aplicación solicitada
	pa, err := s.repo.FindPacienteAplicacion(pu.PacienteID, req.AplicacionID, req.ClinicaID)
	if err != nil || pa.State != "A" {
		return nil, ErrSinAccesoAplicacion
	}
	log.Printf("[INFO PACIENTE APP]%v", pa.Aplicacion.Medico)

	token, expiresIn, err := s.jwtSvc.GeneratePacienteToken(pu.PacienteID, pu.ClinicaID, req.AplicacionID)
	if err != nil {
		return nil, err
	}

	resp := &models.PacienteLoginResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		PacienteID:   pu.PacienteID,
		ClinicaID:    pu.ClinicaID,
		AplicacionID: req.AplicacionID,
		UserName:     pu.Paciente.Nombres + " " + pu.Paciente.Apellidos,
		Medico:       pa.Aplicacion.Medico,
	}

	// Enriquecer respuesta con info del doctor (médico de la aplicación)
	if app, err := s.repo.FindAplicacionByID(req.AplicacionID); err == nil && app.MedicoID != nil {
		if doctor, err := s.repo.FindDoctorByID(*app.MedicoID); err == nil {
			resp.DoctorID = app.MedicoID
			resp.DoctorNombre = doctor.Nombre
			resp.DoctorApellidos = doctor.Apellidos
			resp.DoctorEspecialidad = doctor.Especialidad
		}
	}

	// Registrar acceso (para métricas de frecuencia de uso)
	if err := s.repo.RegistrarAccesoApp(pu.PacienteID, pu.ClinicaID, req.AplicacionID, "LOGIN"); err != nil {
		log.Printf("[WARN] No se pudo registrar acceso de paciente %d: %v", pu.PacienteID, err)
	}

	return resp, nil
}
