package models

import (
	"saas-medico/internal/modules/pacientes/models"
	"time"
)

// ─── Grupo de alimentos ───────────────────────────────────────────────────────

type NutricionGrupoAlimento struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo      string    `gorm:"size:30;uniqueIndex;not null" json:"codigo"`
	Nombre      string    `gorm:"size:80;not null" json:"nombre"`
	Descripcion string    `gorm:"size:255" json:"descripcion,omitempty"`
	Icono       string    `gorm:"size:100" json:"icono,omitempty"`
	Color       string    `gorm:"size:30" json:"color,omitempty"`
	Orden       int       `gorm:"default:0" json:"orden"`
	State       string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn    time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionGrupoAlimento) TableName() string { return "nutricion_grupos_alimento" }

// Códigos de grupo de alimento
const (
	GrupoProteinaCod    = "PROTEINA"
	GrupoLacteoCod      = "LACTEO"
	GrupoLacteoVegCod   = "LACTEO_VEG"
	GrupoLegumbreCod    = "LEGUMBRE"
	GrupoCarboCod       = "CARBOHIDRATO"
	GrupoFrutaCod       = "FRUTA"
	GrupoFrutaSecaCod   = "FRUTA_SECA"
	GrupoGrasaSalCod    = "GRASA_SAL"
	GrupoFrutosSecosCod = "FRUTOS_SECOS"
	GrupoAzucarCod      = "AZUCAR"
	GrupoVegetalCod     = "VEGETAL"
	GrupoPreparacionCod = "PREPARACION"
)

// ─── Catálogos ────────────────────────────────────────────────────────────────

type NutricionTipoComida struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo  string `gorm:"size:10;uniqueIndex;not null" json:"codigo"`
	Nombre  string `gorm:"size:80;not null" json:"nombre"`
	Orden   int    `gorm:"not null" json:"orden"`
	HoraRef string `gorm:"type:time" json:"hora_ref,omitempty"`
	State   string `gorm:"type:char(1);default:'A';not null" json:"state"`
	Main    bool   `gorm:"type:bool" json:"main,omitempty"`
}

func (NutricionTipoComida) TableName() string { return "nutricion_tipo_comida" }

type NutricionTipoComidaGrupo struct {
	ID              uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	TipoComidaID    uint   `gorm:"not null;index;uniqueIndex:ux_tipo_grupo" json:"tipo_comida_id"`
	GrupoAlimentoID uint   `gorm:"not null;index;uniqueIndex:ux_tipo_grupo" json:"grupo_alimento_id"`
	Obligatorio     bool   `gorm:"default:true" json:"obligatorio"`
	CantidadMin     *int   `gorm:"column:cantidad_min" json:"cantidad_min,omitempty"`
	CantidadMax     *int   `gorm:"column:cantidad_max" json:"cantidad_max,omitempty"`
	State           string `gorm:"type:char(1);default:'A';not null" json:"state"`

	TipoComida    NutricionTipoComida    `gorm:"foreignKey:TipoComidaID;references:ID" json:"tipo_comida"`
	GrupoAlimento NutricionGrupoAlimento `gorm:"foreignKey:GrupoAlimentoID;references:ID" json:"grupo_alimento"`
}

func (NutricionTipoComidaGrupo) TableName() string {
	return "nutricion_tipo_comida_grupos"
}

// Códigos tipo comida
const (
	TipoComidaDES = "DES" // Desayuno
	TipoComidaMMA = "MMA" // Media Mañana
	TipoComidaALM = "ALM" // Almuerzo
	TipoComidaMTA = "MTA" // Media Tarde
	TipoComidaMER = "MER" // Merienda/Cena
)

type NutricionAlimento struct {
	ID               uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre           string                  `gorm:"size:150;not null" json:"nombre"`
	Descripcion      string                  `gorm:"size:255" json:"descripcion,omitempty"`
	GrupoID          *uint                   `gorm:"index" json:"grupo_id,omitempty"`
	Grupo            *NutricionGrupoAlimento `gorm:"foreignKey:GrupoID" json:"grupo,omitempty"`
	Categoria        string                  `gorm:"size:80" json:"categoria,omitempty"`
	GramosPorcion    float64                 `gorm:"type:decimal(8,2);default:100.00" json:"gramos_porcion"`
	Calorias         float64                 `gorm:"type:decimal(8,2);not null;default:0" json:"calorias"`
	ProteínasG       float64                 `gorm:"type:decimal(8,2);not null;default:0;column:proteinas_g" json:"proteinas_g"`
	CarbohidratosG   float64                 `gorm:"type:decimal(8,2);not null;default:0" json:"carbohidratos_g"`
	GrasasG          float64                 `gorm:"type:decimal(8,2);not null;default:0" json:"grasas_g"`
	FibraG           *float64                `gorm:"type:decimal(8,2)" json:"fibra_g,omitempty"`
	AzucaresG        *float64                `gorm:"type:decimal(8,2)" json:"azucares_g,omitempty"`
	SodioMg          *float64                `gorm:"type:decimal(8,2)" json:"sodio_mg,omitempty"`
	GrasasSaturadasG *float64                `gorm:"type:decimal(8,2)" json:"grasas_saturadas_g,omitempty"`
	GrasasTransG     *float64                `gorm:"type:decimal(8,2)" json:"grasas_trans_g,omitempty"`
	ColesterolMg     *float64                `gorm:"type:decimal(8,2)" json:"colesterol_mg,omitempty"`
	State            string                  `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoPor        *uint                   `gorm:"index" json:"creado_por,omitempty"`
	CreadoEn         time.Time               `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn    time.Time               `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
	Desayuno         bool                    `gorm:"type:bool" json:"desayuno,omitempty"`
	MediaTardeMana   bool                    `gorm:"type:bool" json:"media_tarde_mana,omitempty"`
	Almuerzo         bool                    `gorm:"type:bool" json:"almuerzo,omitempty"`
	Merienda         bool                    `gorm:"type:bool" json:"merienda,omitempty"`
}

func (NutricionAlimento) TableName() string { return "nutricion_alimentos" }

type NutricionDietaCatalogo struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre             string    `gorm:"size:150;not null" json:"nombre"`
	Descripcion        string    `gorm:"type:text" json:"descripcion,omitempty"`
	TipoPacientePerfil string    `gorm:"size:255" json:"tipo_paciente_perfil,omitempty"`
	Objetivo           string    `gorm:"size:100" json:"objetivo,omitempty"`
	CaloriasDia        *float64  `gorm:"type:decimal(8,2)" json:"calorias_dia,omitempty"`
	ProteínasGDia      *float64  `gorm:"type:decimal(8,2);column:proteinas_g_dia" json:"proteinas_g_dia,omitempty"`
	CarbohidratosGDia  *float64  `gorm:"type:decimal(8,2)" json:"carbohidratos_g_dia,omitempty"`
	GrasasGDia         *float64  `gorm:"type:decimal(8,2)" json:"grasas_g_dia,omitempty"`
	FibraGDia          *float64  `gorm:"type:decimal(8,2)" json:"fibra_g_dia,omitempty"`
	State              string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoPor          *uint     `gorm:"index" json:"creado_por,omitempty"`
	CreadoEn           time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn      time.Time `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

func (NutricionDietaCatalogo) TableName() string { return "nutricion_dietas_catalogo" }

// ─── Plan de dieta del paciente ───────────────────────────────────────────────

const (
	EstadoDietaActiva     = "ACTIVA"
	EstadoDietaCompletada = "COMPLETADA"
	EstadoDietaCancelada  = "CANCELADA"
	EstadoDietaPausada    = "PAUSADA"
)

type NutricionDietaPaciente struct {
	ID                  uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID          uint            `gorm:"not null;index" json:"paciente_id"`
	MedicoID            uint            `gorm:"not null;index" json:"medico_id"`
	DietaCatalogoID     *uint           `gorm:"index" json:"dieta_catalogo_id,omitempty"`
	Nombre              string          `gorm:"size:150;not null" json:"nombre"`
	Descripcion         string          `gorm:"type:text" json:"descripcion,omitempty"`
	Objetivo            string          `gorm:"size:150" json:"objetivo,omitempty"`
	PesoObjetivo        *float64        `gorm:"type:decimal(5,2);column:resultado_esperado" json:"resultado_esperado,omitempty"`
	FechaInicio         time.Time       `gorm:"type:date;not null" json:"fecha_inicio"`
	DuracionDias        int             `gorm:"not null;default:7" json:"duracion_dias"`
	NumComidas          int             `gorm:"not null;default:5" json:"num_comidas"`
	FechaFin            *time.Time      `gorm:"type:date" json:"fecha_fin,omitempty"`
	CaloriasDiaObjetivo *float64        `gorm:"type:decimal(8,2)" json:"calorias_dia_objetivo,omitempty"`
	ProteínasGDia       *float64        `gorm:"type:decimal(8,2);column:proteinas_g_dia" json:"proteinas_g_dia,omitempty"`
	CarbohidratosGDia   *float64        `gorm:"type:decimal(8,2)" json:"carbohidratos_g_dia,omitempty"`
	GrasasGDia          *float64        `gorm:"type:decimal(8,2)" json:"grasas_g_dia,omitempty"`
	FibraGDia           *float64        `gorm:"type:decimal(8,2)" json:"fibra_g_dia,omitempty"`
	Estado              string          `gorm:"size:20;default:'ACTIVA'" json:"estado"`
	State               string          `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn            time.Time       `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn       time.Time       `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
	Paciente            models.Paciente `gorm:"foreignKey:PacienteID"`
}

func (NutricionDietaPaciente) TableName() string { return "nutricion_dieta_paciente" }

// ─── Menú semanal ─────────────────────────────────────────────────────────────

const (
	EstadoMenuPendiente  = "PENDIENTE"
	EstadoMenuActivo     = "ACTIVO"
	EstadoMenuCompletado = "COMPLETADO"
)

// NutricionMenu — menús semanales; cada dieta tiene múltiples menús (1 por semana)
type NutricionMenu struct {
	ID              uint                   `gorm:"primaryKey;autoIncrement" json:"id"`
	DietaPacienteID uint                   `gorm:"not null;index" json:"dieta_paciente_id"`
	SemanaNumero    int                    `gorm:"not null" json:"semana_numero"`
	FechaInicio     time.Time              `gorm:"type:date;not null" json:"fecha_inicio"`
	FechaFin        time.Time              `gorm:"type:date;not null" json:"fecha_fin"`
	Nombre          string                 `gorm:"size:150" json:"nombre,omitempty"`
	Notas           string                 `gorm:"type:text" json:"notas,omitempty"`
	Estado          string                 `gorm:"size:20;default:'PENDIENTE'" json:"estado"` // PENDIENTE, ACTIVO, COMPLETADO
	State           string                 `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn        time.Time              `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn   time.Time              `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
	Detalles        []NutricionMenuDetalle `gorm:"foreignKey:MenuID;references:ID" json:"detalles"`
	Dieta           NutricionDietaPaciente `gorm:"foreignKey:DietaPacienteID;references:ID" json:"dieta"`
}

func (NutricionMenu) TableName() string { return "nutricion_menu" }

// ─── Detalle de menú ──────────────────────────────────────────────────────────
// Cada fila = una comida (tipo_comida) en un día de la semana del menú.
// Jerarquía: DietaPaciente → Menu (semana) → MenuDetalle (día+comida) → MenuAlimento

type NutricionMenuDetalle struct {
	ID                  uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	MenuID              uint                    `gorm:"not null;index;uniqueIndex:udx_menu_dia_comida" json:"menu_id"`
	TipoComidaID        uint                    `gorm:"not null;uniqueIndex:udx_menu_dia_comida" json:"tipo_comida_id"`
	DiaNúmero           int8                    `gorm:"not null;uniqueIndex:udx_menu_dia_comida;column:dia_numero" json:"dia_numero"` // 1=Lun … 7=Dom
	NombreComida        string                  `gorm:"size:150" json:"nombre_comida,omitempty"`
	Instrucciones       string                  `gorm:"type:text" json:"instrucciones,omitempty"`
	CaloriasTotal       *float64                `gorm:"type:decimal(8,2)" json:"calorias_total,omitempty"`
	ProteinasGTotal     *float64                `gorm:"type:decimal(8,2);column:proteinas_g_total" json:"proteinas_g_total,omitempty"`
	CarbohidratosGTotal *float64                `gorm:"type:decimal(8,2)" json:"carbohidratos_g_total,omitempty"`
	GrasasGTotal        *float64                `gorm:"type:decimal(8,2)" json:"grasas_g_total,omitempty"`
	State               string                  `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn            time.Time               `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn       time.Time               `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
	NombreReceta        string                  `gorm:"size:150" json:"nombre_receta,omitempty"`
	Alimentos           []NutricionMenuAlimento `gorm:"foreignKey:MenuDetalleID;references:ID" json:"alimentos"`
}

func (NutricionMenuDetalle) TableName() string { return "nutricion_menu_detalle" }

// Nutricion Plantilla Menu
type NutricionMenuPlantilla struct {
	ID              uint `gorm:"primaryKey;autoIncrement" json:"id"`
	MenuID          uint `gorm:"not null;index;uniqueIndex:udx_menu_dia_comida" json:"menu_id"`
	DietaPacienteID uint `gorm:"not null;index" json:"dieta_paciente_id"`
}

func (NutricionMenuPlantilla) TableName() string { return "nutricion_menu_plantilla" }

// NutricionMenuAlimento — alimentos asignados a un ítem del detalle del menú.
type NutricionMenuAlimento struct {
	ID                 uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	MenuDetalleID      uint              `gorm:"not null;index" json:"menu_detalle_id"`
	AlimentoID         uint              `gorm:"not null;index" json:"alimento_id"`
	GramosAsignados    float64           `gorm:"type:decimal(8,2);not null" json:"gramos_asignados"`
	CaloriasCalc       *float64          `gorm:"type:decimal(8,2)" json:"calorias_calc,omitempty"`
	ProteinasGCalc     *float64          `gorm:"type:decimal(8,2);column:proteinas_g_calc" json:"proteinas_g_calc,omitempty"`
	CarbohidratosGCalc *float64          `gorm:"type:decimal(8,2)" json:"carbohidratos_g_calc,omitempty"`
	GrasasGCalc        *float64          `gorm:"type:decimal(8,2)" json:"grasas_g_calc,omitempty"`
	Observacion        string            `gorm:"size:255" json:"observacion,omitempty"`
	State              string            `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn           time.Time         `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	Alimento           NutricionAlimento `gorm:"foreignKey:AlimentoID"`
}

func (NutricionMenuAlimento) TableName() string { return "nutricion_menu_alimentos" }

// ─── Recordatorio 24 horas ────────────────────────────────────────────────────

// NutricionR24H — encabezado del recordatorio 24 horas
type NutricionR24H struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint      `gorm:"not null;index" json:"paciente_id"`
	MedicoID      uint      `gorm:"not null;index" json:"medico_id"`
	Fecha         time.Time `gorm:"type:date;not null" json:"fecha"`
	Observaciones string    `gorm:"type:text" json:"observaciones,omitempty"`
	CaloriasTotal *float64  `gorm:"type:decimal(8,2)" json:"calorias_total,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionR24H) TableName() string { return "nutricion_r24h" }

// NutricionR24HItem — cada alimento registrado en un recordatorio 24 horas
type NutricionR24HItem struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	R24HID      uint      `gorm:"not null;index" json:"r24h_id"`
	HoraAprox   string    `gorm:"size:10" json:"hora_aprox,omitempty"`  // "07:30"
	TipoComida  string    `gorm:"size:50" json:"tipo_comida,omitempty"` // "Desayuno", "Almuerzo"...
	Alimento    string    `gorm:"size:200;not null" json:"alimento"`    // texto libre
	Cantidad    string    `gorm:"size:100" json:"cantidad,omitempty"`   // "1 taza", "200g"
	CaloriasEst *float64  `gorm:"type:decimal(8,2)" json:"calorias_est,omitempty"`
	Notas       string    `gorm:"size:255" json:"notas,omitempty"`
	State       string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn    time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionR24HItem) TableName() string { return "nutricion_r24h_items" }

// ─── Preferencias alimentarias ────────────────────────────────────────────────

const (
	PreferenciaGusto        = "GUSTO"
	PreferenciaDisgusto     = "DISGUSTO"
	PreferenciaIntolerancia = "INTOLERANCIA"
	PreferenciaAlergia      = "ALERGIA"
)

// NutricionPreferenciaAlimento — gustos y disgustos alimentarios del paciente
type NutricionPreferenciaAlimento struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID  uint              `gorm:"not null;index" json:"paciente_id"`
	AlimentoID  *uint             `gorm:"index" json:"alimento_id,omitempty"`
	NombreLibre string            `gorm:"size:150" json:"nombre_libre,omitempty"`
	Tipo        string            `gorm:"size:20;not null" json:"tipo"` // GUSTO, DISGUSTO, INTOLERANCIA, ALERGIA
	Notas       string            `gorm:"size:255" json:"notas,omitempty"`
	State       string            `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn    time.Time         `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	Alimento    NutricionAlimento `gorm:"foreignKey:AlimentoID"`
}

func (NutricionPreferenciaAlimento) TableName() string { return "nutricion_preferencias_alimento" }

// ─── Síntomas ─────────────────────────────────────────────────────────────────

// NutricionSintoma — síntomas reportados por el paciente relacionados a alimentación/nutrición
type NutricionSintoma struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID      uint      `gorm:"not null;index" json:"paciente_id"`
	Fecha           time.Time `gorm:"type:date;not null" json:"fecha"`
	Tipo            string    `gorm:"size:50" json:"tipo,omitempty"` // GASTROINTESTINAL, ENERGETICO, DIGESTIVO, OTRO
	Descripcion     string    `gorm:"type:text;not null" json:"descripcion"`
	Intensidad      *int8     `json:"intensidad,omitempty"` // 1-10
	AlimentoPosible string    `gorm:"size:255" json:"alimento_posible,omitempty"`
	State           string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn        time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionSintoma) TableName() string { return "nutricion_sintomas" }

// ─── Tipo de Recurso (catálogo) ───────────────────────────────────────────────

type NutricionTipoRecurso struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"size:255;not null" json:"nombre"`
	State  string `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (NutricionTipoRecurso) TableName() string { return "nutricion_tipo_recurso" }

// ─── Archivos PDF ─────────────────────────────────────────────────────────────

// NutricionArchivoPDF — biblioteca de archivos clínicos por clínica/paciente
type NutricionArchivoPDF struct {
	ID            uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	ClinicaID     uint                  `gorm:"not null;index" json:"clinica_id"`
	MedicoID      uint                  `gorm:"not null;index" json:"medico_id"`
	PacienteID    *uint                 `gorm:"index" json:"paciente_id,omitempty"`
	TipoRecursoID uint                  `gorm:"column:tipo;not null;index" json:"tipo_recurso_id"`
	TipoRecurso   *NutricionTipoRecurso `gorm:"foreignKey:TipoRecursoID" json:"tipo_recurso,omitempty"`
	Titulo        string                `gorm:"size:255;not null" json:"titulo"`
	Descripcion   string                `gorm:"type:text" json:"descripcion,omitempty"`
	RutaArchivo   string                `gorm:"size:500;not null" json:"ruta_archivo"`
	State         string                `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn      time.Time             `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionArchivoPDF) TableName() string { return "nutricion_archivos_pdf" }

// ─── Ejercicios ───────────────────────────────────────────────────────────────

const (
	EstadoEjercicioPendiente  = "PENDIENTE"
	EstadoEjercicioCompletado = "COMPLETADO"
	EstadoEjercicioSaltado    = "SALTADO"
)

type NutricionEjercicioCatalogo struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre          string    `gorm:"size:150;not null" json:"nombre"`
	Descripcion     string    `gorm:"type:text" json:"descripcion,omitempty"`
	Categoria       string    `gorm:"size:80" json:"categoria,omitempty"`
	GrupoMuscular   string    `gorm:"size:120" json:"grupo_muscular,omitempty"`
	CaloriasPorHora *float64  `gorm:"type:decimal(8,2)" json:"calorias_por_hora,omitempty"`
	UnidadMedida    string    `gorm:"size:30;default:'minutos'" json:"unidad_medida"`
	Nivel           string    `gorm:"size:20" json:"nivel,omitempty"`
	UrlReferencia   string    `gorm:"size:500" json:"url_referencia,omitempty"`
	State           string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoPor       *uint     `gorm:"index" json:"creado_por,omitempty"`
	CreadoEn        time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionEjercicioCatalogo) TableName() string { return "nutricion_ejercicios_catalogo" }

type NutricionEjercicioPaciente struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID        uint      `gorm:"not null;index" json:"paciente_id"`
	MedicoID          uint      `gorm:"not null;index" json:"medico_id"`
	DietaPacienteID   *uint     `gorm:"index" json:"dieta_paciente_id,omitempty"`
	EjercicioID       uint      `gorm:"not null;index" json:"ejercicio_id"`
	DiaNúmero         *int8     `gorm:"column:dia_numero" json:"dia_numero,omitempty"`
	DiaSemana         string    `gorm:"size:15" json:"dia_semana,omitempty"`
	DuracionMin       *int      `json:"duracion_min,omitempty"`
	Series            *int      `json:"series,omitempty"`
	Repeticiones      *int      `json:"repeticiones,omitempty"`
	PesoKg            *float64  `gorm:"type:decimal(6,2)" json:"peso_kg,omitempty"`
	DescansoSeg       *int      `json:"descanso_seg,omitempty"`
	CaloriasEstimadas *float64  `gorm:"type:decimal(8,2)" json:"calorias_estimadas,omitempty"`
	Instrucciones     string    `gorm:"type:text" json:"instrucciones,omitempty"`
	Estado            string    `gorm:"size:20;default:'PENDIENTE'" json:"estado"`
	State             string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn          time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
	ActualizadoEn     time.Time `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

func (NutricionEjercicioPaciente) TableName() string { return "nutricion_ejercicios_paciente" }

// ─── Registros del paciente ───────────────────────────────────────────────────

// Estados de registro de comida
const (
	EstadoRegistroComidaPendiente = "P" // P = pendiente (en plan, no consumida aún)
	EstadoRegistroComidaConsumida = "C" // C = consumida (paciente confirmó que la comió)
)

type NutricionRegistroComida struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID         uint      `gorm:"not null;index:idx_pac_fecha" json:"paciente_id"`
	Fecha              time.Time `gorm:"type:date;not null;index:idx_pac_fecha" json:"fecha"`
	TipoComidaID       uint      `gorm:"not null;index" json:"tipo_comida_id"`
	MenuDetalleID      *uint     `gorm:"index" json:"menu_detalle_id,omitempty"`
	FueraDePlan        bool      `gorm:"default:false" json:"fuera_de_plan"`
	DescripcionLibre   string    `gorm:"size:255" json:"descripcion_libre,omitempty"`
	CaloriasConsumidas *float64  `gorm:"type:decimal(8,2)" json:"calorias_consumidas,omitempty"`
	ProteínasG         *float64  `gorm:"type:decimal(8,2);column:proteinas_g" json:"proteinas_g,omitempty"`
	CarbohidratosG     *float64  `gorm:"type:decimal(8,2)" json:"carbohidratos_g,omitempty"`
	GrasasG            *float64  `gorm:"type:decimal(8,2)" json:"grasas_g,omitempty"`
	PorcentajeCumplido *int      `json:"porcentaje_cumplido,omitempty"`
	FotoComida         string    `gorm:"size:500" json:"foto_comida,omitempty"`
	Notas              string    `gorm:"size:255" json:"notas,omitempty"`
	Estado             string    `gorm:"type:char(1);default:'C';not null" json:"estado"` // P=pendiente, C=consumida
	State              string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn           time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionRegistroComida) TableName() string { return "nutricion_registro_comidas" }

type NutricionRegistroAlimento struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RegistroComidaID   uint      `gorm:"not null;index" json:"registro_comida_id"`
	AlimentoID         *uint     `gorm:"index" json:"alimento_id,omitempty"`
	NombreLibre        string    `gorm:"size:150" json:"nombre_libre,omitempty"`
	GramosConsumidos   float64   `gorm:"type:decimal(8,2);not null" json:"gramos_consumidos"`
	CaloriasCalc       *float64  `gorm:"type:decimal(8,2)" json:"calorias_calc,omitempty"`
	ProteínasGCalc     *float64  `gorm:"type:decimal(8,2);column:proteinas_g_calc" json:"proteinas_g_calc,omitempty"`
	CarbohidratosGCalc *float64  `gorm:"type:decimal(8,2)" json:"carbohidratos_g_calc,omitempty"`
	GrasasGCalc        *float64  `gorm:"type:decimal(8,2)" json:"grasas_g_calc,omitempty"`
	State              string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	FueraDePlan        bool      `gorm:"default:false" json:"fuera_de_plan"`
	CreadoEn           time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionRegistroAlimento) TableName() string { return "nutricion_registro_alimentos" }

type NutricionRegistroEjercicio struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID            uint      `gorm:"not null;index:idx_pac_fecha_ej" json:"paciente_id"`
	Fecha                 time.Time `gorm:"type:date;not null;index:idx_pac_fecha_ej" json:"fecha"`
	EjercicioPacienteID   *uint     `gorm:"index" json:"ejercicio_paciente_id,omitempty"`
	EjercicioID           *uint     `gorm:"index" json:"ejercicio_id,omitempty"`
	NombreLibre           string    `gorm:"size:150" json:"nombre_libre,omitempty"`
	DuracionMinReal       *int      `json:"duracion_min_real,omitempty"`
	SeriesReal            *int      `json:"series_real,omitempty"`
	RepeticionesReal      *int      `json:"repeticiones_real,omitempty"`
	PesoKgReal            *float64  `gorm:"type:decimal(6,2)" json:"peso_kg_real,omitempty"`
	CaloriasQuemadas      *float64  `gorm:"type:decimal(8,2)" json:"calorias_quemadas,omitempty"`
	FrecuenciaCardiacaMax *int      `json:"frecuencia_cardiaca_max,omitempty"`
	NivelEsfuerzo         *int8     `json:"nivel_esfuerzo,omitempty"`
	Notas                 string    `gorm:"size:255" json:"notas,omitempty"`
	State                 string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn              time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionRegistroEjercicio) TableName() string { return "nutricion_registro_ejercicios" }

// ─── Progreso ─────────────────────────────────────────────────────────────────

type NutricionProgresoPaciente struct {
	ID                       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID               uint      `gorm:"not null;index:idx_progreso_pac_fecha" json:"paciente_id"`
	MedicoID                 *uint     `gorm:"index" json:"medico_id,omitempty"`
	DietaPacienteID          *uint     `gorm:"index" json:"dieta_paciente_id,omitempty"`
	Fecha                    time.Time `gorm:"type:date;not null;index:idx_progreso_pac_fecha" json:"fecha"`
	PesoKg                   *float64  `gorm:"type:decimal(6,2)" json:"peso_kg,omitempty"`
	AlturaCm                 *float64  `gorm:"type:decimal(6,2)" json:"altura_cm,omitempty"`
	IMC                      *float64  `gorm:"type:decimal(5,2)" json:"imc,omitempty"`
	GrasaCorporalPct         *float64  `gorm:"type:decimal(5,2)" json:"grasa_corporal_pct,omitempty"`
	MasaMuscularKg           *float64  `gorm:"type:decimal(6,2)" json:"masa_muscular_kg,omitempty"`
	CinturaCm                *float64  `gorm:"type:decimal(6,2)" json:"cintura_cm,omitempty"`
	CaderaCm                 *float64  `gorm:"type:decimal(6,2)" json:"cadera_cm,omitempty"`
	PechoCm                  *float64  `gorm:"type:decimal(6,2)" json:"pecho_cm,omitempty"`
	BrazoCm                  *float64  `gorm:"type:decimal(6,2)" json:"brazo_cm,omitempty"`
	MusloCm                  *float64  `gorm:"type:decimal(6,2)" json:"muslo_cm,omitempty"`
	CaloriasConsumidasDia    *float64  `gorm:"type:decimal(8,2)" json:"calorias_consumidas_dia,omitempty"`
	PctCumplimientoDieta     *int      `json:"pct_cumplimiento_dieta,omitempty"`
	PctCumplimientoEjercicio *int      `json:"pct_cumplimiento_ejercicio,omitempty"`
	EnergiaNivel             *int8     `json:"energia_nivel,omitempty"`
	SuenoHoras               *float64  `gorm:"type:decimal(4,2)" json:"sueno_horas,omitempty"`
	HidratacionLitros        *float64  `gorm:"type:decimal(4,2)" json:"hidratacion_litros,omitempty"`
	Notas                    string    `gorm:"type:text" json:"notas,omitempty"`
	FotoProgreso             string    `gorm:"size:500" json:"foto_progreso,omitempty"`
	State                    string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn                 time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionProgresoPaciente) TableName() string { return "nutricion_progreso_paciente" }

// ─── Logros / Gamificación ────────────────────────────────────────────────────

type NutricionLogroCatalogo struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo         string    `gorm:"size:30;uniqueIndex;not null" json:"codigo"`
	Nombre         string    `gorm:"size:120;not null" json:"nombre"`
	Descripcion    string    `gorm:"size:255" json:"descripcion,omitempty"`
	Icono          string    `gorm:"size:100" json:"icono,omitempty"`
	Categoria      string    `gorm:"size:50" json:"categoria,omitempty"`
	CondicionTipo  string    `gorm:"size:50" json:"condicion_tipo,omitempty"`
	CondicionValor *int      `json:"condicion_valor,omitempty"`
	PuntosXP       int       `gorm:"default:0" json:"puntos_xp"`
	State          string    `gorm:"type:char(1);default:'A';not null" json:"state"`
	CreadoEn       time.Time `gorm:"autoCreateTime;column:creado_en" json:"creado_en"`
}

func (NutricionLogroCatalogo) TableName() string { return "nutricion_logros_catalogo" }

type NutricionLogroPaciente struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID    uint      `gorm:"not null;uniqueIndex:udx_pac_logro" json:"paciente_id"`
	LogroID       uint      `gorm:"not null;uniqueIndex:udx_pac_logro" json:"logro_id"`
	FechaObtenido time.Time `gorm:"not null" json:"fecha_obtenido"`
	PuntosXP      int       `gorm:"default:0" json:"puntos_xp"`
	Notas         string    `gorm:"size:255" json:"notas,omitempty"`
	State         string    `gorm:"type:char(1);default:'A';not null" json:"state"`
}

func (NutricionLogroPaciente) TableName() string { return "nutricion_logros_paciente" }

type NutricionPacienteXP struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	PacienteID     uint       `gorm:"not null;uniqueIndex" json:"paciente_id"`
	XPTotal        int        `gorm:"default:0" json:"xp_total"`
	Nivel          int        `gorm:"default:1" json:"nivel"`
	RachaActual    int        `gorm:"default:0" json:"racha_actual"`
	RachaMaxima    int        `gorm:"default:0" json:"racha_maxima"`
	UltimoRegistro *time.Time `gorm:"type:date" json:"ultimo_registro,omitempty"`
	State          string     `gorm:"type:char(1);default:'A';not null" json:"state"`
	ActualizadoEn  time.Time  `gorm:"autoUpdateTime;column:actualizado_en" json:"actualizado_en"`
}

func (NutricionPacienteXP) TableName() string { return "nutricion_paciente_xp" }

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAlimentoRequest struct {
	Nombre         string   `json:"nombre" binding:"required,min=2,max=150"`
	Descripcion    string   `json:"descripcion"`
	Categoria      string   `json:"categoria"`
	GramosPorcion  float64  `json:"gramos_porcion"`
	Calorias       float64  `json:"calorias" binding:"required"`
	ProteínasG     float64  `json:"proteinas_g"`
	CarbohidratosG float64  `json:"carbohidratos_g"`
	GrasasG        float64  `json:"grasas_g"`
	FibraG         *float64 `json:"fibra_g"`
	AzucaresG      *float64 `json:"azucares_g"`
	SodioMg        *float64 `json:"sodio_mg"`
	Desayuno       bool     `json:"desayuno"`
	MediaTardeMana bool     `json:"media_tarde_mana"`
	Almuerzo       bool     `json:"almuerzo"`
	Merienda       bool     `json:"merienda"`
}

type UpdateAlimentoRequest struct {
	Nombre         string   `json:"nombre" binding:"required,min=2,max=150"`
	Descripcion    string   `json:"descripcion"`
	Categoria      string   `json:"categoria"`
	GramosPorcion  float64  `json:"gramos_porcion"`
	Calorias       float64  `json:"calorias" binding:"required"`
	ProteínasG     float64  `json:"proteinas_g"`
	CarbohidratosG float64  `json:"carbohidratos_g"`
	GrasasG        float64  `json:"grasas_g"`
	FibraG         *float64 `json:"fibra_g"`
	AzucaresG      *float64 `json:"azucares_g"`
	SodioMg        *float64 `json:"sodio_mg"`
	Desayuno       bool     `json:"desayuno"`
	MediaTardeMana bool     `json:"media_tarde_mana"`
	Almuerzo       bool     `json:"almuerzo"`
	Merienda       bool     `json:"merienda"`
}

type CreateDietaRequest struct {
	DietaCatalogoID     *uint    `json:"dieta_catalogo_id"`
	Nombre              string   `json:"nombre" binding:"required,min=2,max=150"`
	Descripcion         string   `json:"descripcion"`
	Objetivo            string   `json:"objetivo"`
	ResultadoEsperado   *float64 `json:"resultado_esperado"`
	FechaInicio         string   `json:"fecha_inicio" binding:"required"`
	DuracionDias        int      `json:"duracion_dias"`
	NumComidas          int      `json:"num_comidas"`
	CaloriasDiaObjetivo *float64 `json:"calorias_dia_objetivo"`
	ProteínasGDia       *float64 `json:"proteinas_g_dia"`
	CarbohidratosGDia   *float64 `json:"carbohidratos_g_dia"`
	GrasasGDia          *float64 `json:"grasas_g_dia"`
	FibraGDia           *float64 `json:"fibra_g_dia"`
}

type UpdateDietaRequest struct {
	Nombre              string   `json:"nombre"`
	Descripcion         string   `json:"descripcion"`
	Objetivo            string   `json:"objetivo"`
	ResultadoEsperado   *float64 `json:"resultado_esperado"`
	NumComidas          int      `json:"num_comidas"`
	CaloriasDiaObjetivo *float64 `json:"calorias_dia_objetivo"`
	ProteínasGDia       *float64 `json:"proteinas_g_dia"`
	CarbohidratosGDia   *float64 `json:"carbohidratos_g_dia"`
	GrasasGDia          *float64 `json:"grasas_g_dia"`
	FibraGDia           *float64 `json:"fibra_g_dia"`
	Estado              string   `json:"estado"`
}

type CreateMenuRequest struct {
	SemanaNumero int    `json:"semana_numero" binding:"required,min=1"`
	FechaInicio  string `json:"fecha_inicio" binding:"required"` // "2026-03-01"
	Nombre       string `json:"nombre"`
	Notas        string `json:"notas"`
}

type AddDetalleMenuRequest struct {
	TipoComidaID  uint   `json:"tipo_comida_id" binding:"required"`
	DiaNúmero     int8   `json:"dia_numero" binding:"required,min=1,max=7"`
	NombreComida  string `json:"nombre_comida"`
	Instrucciones string `json:"instrucciones"`
}

type UpdateDetalleMenuRequest struct {
	NombreReceta string `json:"nombre_receta"`
}

type AddAlimentoMenuRequest struct {
	AlimentoID      uint    `json:"alimento_id" binding:"required"`
	GramosAsignados float64 `json:"gramos_asignados" binding:"required,gt=0"`
	Observacion     string  `json:"observacion"`
}

type CalcularFormulasRequest struct {
	Sexo            string   `json:"sexo" binding:"required,oneof=M F"`
	EdadAnos        *int     `json:"edad_anos"`
	AlturaCm        *float64 `json:"altura_cm"`
	PesoKg          *float64 `json:"peso_kg"`
	CinturaCm       *float64 `json:"cintura_cm"`
	CaderaCm        *float64 `json:"cadera_cm"`
	FactorActividad *float64 `json:"factor_actividad"` // 1.2, 1.375, 1.55, 1.725, 1.9
}

// NutricionFormulasResult — valores nutricionales calculados devueltos en la respuesta
type NutricionFormulasResult struct {
	IMC              *float64 `json:"imc,omitempty"`
	ClasificacionIMC string   `json:"clasificacion_imc,omitempty"`
	ICC              *float64 `json:"icc,omitempty"`
	RiesgoMetabolico string   `json:"riesgo_metabolico,omitempty"`
	TMB              *float64 `json:"tmb,omitempty"` // Harris Benedict
	GEB              *float64 `json:"geb,omitempty"` // = TMB
	GET              *float64 `json:"get,omitempty"` // TMB x factor actividad
}

type CreateR24HRequest struct {
	Fecha         string `json:"fecha" binding:"required"` // "2026-03-16"
	Observaciones string `json:"observaciones"`
}

type AddR24HItemRequest struct {
	HoraAprox   string   `json:"hora_aprox"`
	TipoComida  string   `json:"tipo_comida" binding:"required"`
	Alimento    string   `json:"alimento" binding:"required"`
	Cantidad    string   `json:"cantidad"`
	CaloriasEst *float64 `json:"calorias_est"`
	Notas       string   `json:"notas"`
}

type CreatePreferenciaRequest struct {
	AlimentoID  *uint  `json:"alimento_id"`
	NombreLibre string `json:"nombre_libre"`
	Tipo        string `json:"tipo" binding:"required,oneof=GUSTO DISGUSTO INTOLERANCIA ALERGIA"`
	Notas       string `json:"notas"`
}

type CreateSintomaRequest struct {
	Fecha           string `json:"fecha" binding:"required"`
	Tipo            string `json:"tipo"`
	Descripcion     string `json:"descripcion" binding:"required"`
	Intensidad      *int8  `json:"intensidad"`
	AlimentoPosible string `json:"alimento_posible"`
}

type CreateRegistroComidaRequest struct {
	Fecha               string   `json:"fecha" binding:"required"`
	TipoComidaID        uint     `json:"tipo_comida_id" binding:"required"`
	MenuDetalleID       *uint    `json:"menu_detalle_id"`
	FueraDePlan         bool     `json:"fuera_de_plan"`
	DescripcionLibre    string   `json:"descripcion_libre"`
	CaloriasConsumidas  *float64 `json:"calorias_consumidas"`
	FotoComida          string   `json:"foto_comida"`
	Notas               string   `json:"notas"`
	ProteinasConsumidas *float64 `json:"proteinas_g"`
	GrasasConsumidas    *float64 `json:"grasas_g"`
	CarbosConsumidos    *float64 `json:"carbohidratos_g"`
}

type AddRegistroAlimentoRequest struct {
	AlimentoID       *uint   `json:"alimento_id"`
	NombreLibre      string  `json:"nombre_libre"`
	GramosConsumidos float64 `json:"gramos_consumidos" binding:"required,gt=0"`
}

type CreateRegistroEjercicioRequest struct {
	Fecha                 string   `json:"fecha" binding:"required"`
	EjercicioPacienteID   *uint    `json:"ejercicio_paciente_id"`
	EjercicioID           *uint    `json:"ejercicio_id"`
	NombreLibre           string   `json:"nombre_libre"`
	DuracionMinReal       *int     `json:"duracion_min_real"`
	SeriesReal            *int     `json:"series_real"`
	RepeticionesReal      *int     `json:"repeticiones_real"`
	PesoKgReal            *float64 `json:"peso_kg_real"`
	CaloriasQuemadas      *float64 `json:"calorias_quemadas"`
	FrecuenciaCardiacaMax *int     `json:"frecuencia_cardiaca_max"`
	NivelEsfuerzo         *int8    `json:"nivel_esfuerzo"`
	Notas                 string   `json:"notas"`
}

type CreateProgresoRequest struct {
	Fecha                string   `json:"fecha"                  form:"fecha"                  binding:"required"`
	DietaPacienteID      *uint    `json:"dieta_paciente_id"      form:"dieta_paciente_id"`
	PesoKg               *float64 `json:"peso_kg"                form:"peso_kg"`
	AlturaCm             *float64 `json:"altura_cm"              form:"altura_cm"`
	GrasaCorporalPct     *float64 `json:"grasa_corporal_pct"     form:"grasa_corporal_pct"`
	MasaMuscularKg       *float64 `json:"masa_muscular_kg"       form:"masa_muscular_kg"`
	CinturaCm            *float64 `json:"cintura_cm"             form:"cintura_cm"`
	CaderaCm             *float64 `json:"cadera_cm"              form:"cadera_cm"`
	PechoCm              *float64 `json:"pecho_cm"               form:"pecho_cm"`
	BrazoCm              *float64 `json:"brazo_cm"               form:"brazo_cm"`
	MusloCm              *float64 `json:"muslo_cm"               form:"muslo_cm"`
	HidratacionLitros    *float64 `json:"hidratacion_litros"     form:"hidratacion_litros"`
	SuenoHoras           *float64 `json:"sueno_horas"            form:"sueno_horas"`
	EnergiaNivel         *int8    `json:"energia_nivel"          form:"energia_nivel"`
	PctCumplimientoDieta *int     `json:"pct_cumplimiento_dieta" form:"pct_cumplimiento_dieta"`
	FotoProgreso         string   `json:"foto_progreso"          form:"foto_progreso"`
	Notas                string   `json:"notas"                  form:"notas"`
}
type AskIaNutricionQuestion struct {
	Prompt string `json:"prompt" binding:"required"`
}
type CreateEjercicioPacienteRequest struct {
	EjercicioID     uint   `json:"ejercicio_id" binding:"required"`
	DietaPacienteID *uint  `json:"dieta_paciente_id"`
	DiaNúmero       *int8  `json:"dia_numero"`
	DiaSemana       string `json:"dia_semana"`
	DuracionMin     *int   `json:"duracion_min"`
	Series          *int   `json:"series"`
	Repeticiones    *int   `json:"repeticiones"`
	Instrucciones   string `json:"instrucciones"`
}

type CreateArchivoPDFRequest struct {
	PacienteID    *uint  `json:"paciente_id"`
	TipoRecursoID uint   `json:"tipo_recurso_id" binding:"required"`
	Titulo        string `json:"titulo" binding:"required"`
	Descripcion   string `json:"descripcion"`
	RutaArchivo   string `json:"ruta_archivo" binding:"required"`
}

type CreateTipoRecursoRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateTipoRecursoRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type CreateEjercicioCatalogoRequest struct {
	Nombre          string   `json:"nombre" binding:"required"`
	Descripcion     string   `json:"descripcion"`
	Categoria       string   `json:"categoria"`
	GrupoMuscular   string   `json:"grupo_muscular"`
	CaloriasPorHora *float64 `json:"calorias_por_hora"`
	UnidadMedida    string   `json:"unidad_medida"`
	Nivel           string   `json:"nivel"`
}

// ─── Resumen diario (endpoint unificado para móvil y web) ─────────────────────

// ResumenDiarioResponse — snapshot del día de un paciente para el dashboard móvil
type ResumenDiarioResponse struct {
	Fecha                  string                       `json:"fecha"`
	CaloriasObjetivo       float64                      `json:"calorias_objetivo"`
	CaloriasConsumidas     float64                      `json:"calorias_consumidas"`
	CaloriasQuemadas       float64                      `json:"calorias_quemadas"`
	ProteinasG             float64                      `json:"proteinas_g"`
	CarbohidratosG         float64                      `json:"carbohidratos_g"`
	GrasasG                float64                      `json:"grasas_g"`
	PorcentajeCumplimiento int                          `json:"porcentaje_cumplimiento"`
	RegistrosComida        []NutricionRegistroComida    `json:"registros_comida"`
	RegistrosEjercicio     []NutricionRegistroEjercicio `json:"registros_ejercicio"`
	Progreso               *NutricionProgresoPaciente   `json:"progreso,omitempty"`
}

// PacienteRegistrosStats — estadísticas de un paciente para el dashboard web del nutriólogo
type PacienteRegistrosStats struct {
	PacienteID         uint     `json:"paciente_id"`
	TotalComidas       int      `json:"total_comidas"`
	TotalEjercicios    int      `json:"total_ejercicios"`
	CaloriasConsumidas float64  `json:"calorias_consumidas"`
	CaloriasQuemadas   float64  `json:"calorias_quemadas"`
	DiasActivos        int      `json:"dias_activos"`
	UltimoAcceso       string   `json:"ultimo_acceso,omitempty"`
	UltimoPeso         *float64 `json:"ultimo_peso,omitempty"`
	PesoInicial        *float64 `json:"peso_inicial,omitempty"`
	TotalAccesos       int      `json:"total_accesos"`
}
