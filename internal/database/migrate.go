package database

import (
	"log"

	//adminModels "saas-medico/internal/modules/admin/models"
	//agendaModels "saas-medico/internal/modules/agenda/models"
	historiaModels "saas-medico/internal/modules/historia/models"
	nutricionModels "saas-medico/internal/modules/nutricion/models"
	// pacientesModels "saas-medico/internal/modules/pacientes/models"
)

// RunMigrations ejecuta AutoMigrate en orden de dependencias FK.
// Descomenta los grupos que necesites migrar.
func RunMigrations() {
	db := GetDB()

	log.Println("Running database migrations...")

	// 1. Auth: roles, usuarios, tokens
	// if err := db.AutoMigrate(
	// 	&authModels.Rol{},
	// 	&authModels.User{},
	// 	&authModels.UsuarioRol{},
	// 	&authModels.RefreshToken{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (auth):", err)
	// }

	// // 2. Admin: menú y estilos (tablas necesarias para login/menu/estilos)
	// if err := db.AutoMigrate(
	// 	&adminModels.Clinica{},
	// 	&adminModels.EstiloClinica{},
	// 	&adminModels.Transaccion{},
	// 	&adminModels.RolTransaccion{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (admin - menu/estilos):", err)
	// }

	// 2b. Admin: resto (consultorios, profesiones, planes, suscripciones)
	// if err := db.AutoMigrate(
	// 	&adminModels.Sucursal{},
	// 	&adminModels.Consultorio{},
	// 	&adminModels.Profesion{},
	// 	&adminModels.UsuarioClinica{},
	// 	&adminModels.UsuarioConsultorio{},
	// 	&adminModels.PlanSaas{},
	// 	&adminModels.EstadoSuscripcion{},
	// 	&adminModels.Suscripcion{},
	// 	&adminModels.BloqueoAcceso{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (admin - resto):", err)
	// }

	// 3. Pacientes
	// if err := db.AutoMigrate(
	// 	&pacientesModels.Paciente{},
	// 	// &pacientesModels.PrePaciente{},
	// 	&pacientesModels.PacienteUsuario{},
	// 	&pacientesModels.Aplicacion{},
	// 	&pacientesModels.PacienteAplicacion{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (pacientes):", err)
	// }

	// 4. Agenda: tipos, estados, citas, sesiones, horarios, bloqueos
	//if err := db.AutoMigrate(
	//&agendaModels.TipoCita{},
	// &agendaModels.EstadoCita{},
	// &agendaModels.Cita{},
	// &agendaModels.Sesion{},
	// &agendaModels.HorarioMedico{},
	// &agendaModels.BloqueoAgenda{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (agenda):", err)
	// }

	// 5. Cobros
	// if err := db.AutoMigrate(
	// 	&cobrosModels.EstadoCobro{},
	// 	&cobrosModels.MedioPago{},
	// 	&cobrosModels.TipoEgreso{},
	// 	&cobrosModels.CobroSesion{},
	// 	&cobrosModels.Pago{},
	// 	&cobrosModels.Egreso{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (cobros):", err)
	// }

	// 6. Historia clínica — catálogos y formularios dinámicos
	if err := db.AutoMigrate(
		&historiaModels.FormularioCita{},
		// &historiaModels.AlergiaCatalogo{},
		// &historiaModels.TipoAntecedente{},
		// &historiaModels.HabitoCatalogo{},
		// &historiaModels.DiagnosticoCatalogo{},
		// &historiaModels.TipoExamen{},
		//&historiaModels.Formulario{},
		// &historiaModels.FormularioPregunta{},
		// &historiaModels.FormularioOpcion{},
	); err != nil {
		log.Fatal("Migration failed (historia - catalogos/formularios):", err)
	}

	// 7. Historia clínica — registros del paciente
	if err := db.AutoMigrate(
		&historiaModels.TipoImagenPaciente{},
	// 	&historiaModels.HistoriaRespuesta{},
	// 	&historiaModels.PacienteAlergia{},
	// 	&historiaModels.PacienteAntecedente{},
	// 	&historiaModels.PacienteHabito{},
	// 	&historiaModels.PacienteDiagnostico{},
	// 	&historiaModels.PacienteExamenResultado{},
	// 	&historiaModels.PacienteImagen{},
	// 	&historiaModels.PacienteCertificado{},
	); err != nil {
		log.Fatal("Migration failed (historia - paciente):", err)
	}

	// 8. Tests psicológicos
	// if err := db.AutoMigrate(
	// 	&testsModels.TestRegla{},
	// 	&testsModels.TestReglaDetalle{},
	// 	&testsModels.Test{},
	// 	&testsModels.TestRespuesta{},
	// 	&testsModels.TestArchivo{},
	// 	&testsModels.SesionTest{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (tests):", err)
	// }

	// 9. Nutrición — catálogos (sin FK a otras tablas)
	if err := db.AutoMigrate(
		//&nutricionModels.NutricionMenuPlantilla{},
		//&nutricionModels.NutricionMenuDetalle{},
		// 	&nutricionModels.NutricionTipoComida{},
		// 	&nutricionModels.NutricionGrupoAlimento{}, // debe ir antes de NutricionAlimento (FK)
		&nutricionModels.NutricionAlimento{},
		//&nutricionModels.NutricionTipoComidaGrupo{},
		//&nutricionModels.NutricionTipoRecurso{},
		//&nutricionModels.NutricionArchivoPDF{},
		//&nutricionModels.NutricionRegistroAlimento{},
		// 	&nutricionModels.NutricionDietaCatalogo{},
		// 	&nutricionModels.NutricionEjercicioCatalogo{},
		// 	&nutricionModels.NutricionLogroCatalogo{},
	); err != nil {
		log.Fatal("Migration failed (nutricion - catalogos):", err)
	}

	// 10. Nutrición — plan de dieta, menús y detalles de menú (requiere pacientes)
	// Jerarquía: DietaPaciente → Menu (semana) → MenuDetalle (día+comida) → MenuAlimento
	// if err := db.AutoMigrate(
	// 	&nutricionModels.NutricionDietaPaciente{},
	// 	&nutricionModels.NutricionMenu{},
	// 	&nutricionModels.NutricionMenuDetalle{},
	// 	&nutricionModels.NutricionMenuAlimento{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (nutricion - dieta/menu):", err)
	// }

	// 11. Nutrición — registros y seguimiento del paciente (menu_detalle_id en lugar de dieta_detalle_id)
	// if err := db.AutoMigrate(
	// 	&nutricionModels.NutricionR24H{},
	// 	&nutricionModels.NutricionR24HItem{},
	// 	&nutricionModels.NutricionPreferenciaAlimento{},
	// 	&nutricionModels.NutricionSintoma{},
	// 	&nutricionModels.NutricionArchivoPDF{},
	// 	&nutricionModels.NutricionEjercicioPaciente{},
	// 	&nutricionModels.NutricionRegistroComida{},
	// 	&nutricionModels.NutricionRegistroAlimento{},
	// 	&nutricionModels.NutricionRegistroEjercicio{},
	// 	&nutricionModels.NutricionProgresoPaciente{},
	// 	&nutricionModels.NutricionLogroPaciente{},
	// 	&nutricionModels.NutricionPacienteXP{},
	// ); err != nil {
	// 	log.Fatal("Migration failed (nutricion - registros):", err)
	// }

	// Nuevas tablas: PacienteAccesoApp (tracking de frecuencia de uso del app)
	// Descomentar una vez para agregar la columna estado a nutricion_registro_comidas
	// y crear paciente_acceso_app
	// if err := db.AutoMigrate(
	// 	&pacientesModels.PacienteAccesoApp{},
	// 	&pacientesModels.Aplicacion{},              // agrega columna "medico_id" si no existe
	// 	&nutricionModels.NutricionRegistroComida{}, // agrega columna "estado" si no existe
	// ); err != nil {
	// 	log.Fatal("Migration failed (nuevas tablas):", err)
	// }

	log.Println("Database migrations completed successfully")
}
