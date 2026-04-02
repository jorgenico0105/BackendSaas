// Seed de catálogos de nutrición: grupos de alimentos y alimentos.
// Ejecución: go run cmd/seed/main.go
// Es idempotente — usa FirstOrCreate, no duplica registros si se corre más de una vez.
package main

import (
	"log"

	"saas-medico/internal/config"
	"saas-medico/internal/database"
	nutricionModels "saas-medico/internal/modules/nutricion/models"
)

func main() {
	config.LoadConfig()
	database.Connect()
	db := database.GetDB()

	log.Println("Iniciando seed de nutrición...")

	// ── 0. Tipos de comida ────────────────────────────────────────────────────

	type tipoComidaSeed struct {
		Codigo, Nombre, HoraRef string
		Orden                   int
	}

	tiposComida := []tipoComidaSeed{
		{nutricionModels.TipoComidaDES, "Desayuno", "07:00:00", 1},
		{nutricionModels.TipoComidaMMA, "Media Mañana", "10:00:00", 2},
		{nutricionModels.TipoComidaALM, "Almuerzo", "13:00:00", 3},
		{nutricionModels.TipoComidaMTA, "Media Tarde", "16:00:00", 4},
		{nutricionModels.TipoComidaMER, "Merienda / Cena", "19:00:00", 5},
	}

	for _, t := range tiposComida {
		var count int64
		db.Model(&nutricionModels.NutricionTipoComida{}).Where("codigo = ?", t.Codigo).Count(&count)
		if count > 0 {
			continue // ya existe
		}
		// Usamos raw SQL con TIME(?) para que MySQL interprete correctamente la columna TIME
		err := db.Exec(
			"INSERT INTO nutricion_tipo_comida (codigo, nombre, hora_ref, orden, state) VALUES (?, ?, TIME(?), ?, ?)",
			t.Codigo, t.Nombre, t.HoraRef, t.Orden, "A",
		).Error
		if err != nil {
			log.Fatalf("Error insertando tipo comida %q: %v", t.Nombre, err)
		}
	}

	log.Printf("Tipos de comida insertados/verificados: %d", len(tiposComida))

	// ── 1. Grupos de alimentos ─────────────────────────────────────────────────

	type grupoSeed struct {
		Codigo, Nombre, Color string
		Orden                 int
	}

	gruposSeed := []grupoSeed{
		{nutricionModels.GrupoProteinaCod, "Proteína", "#E57373", 1},
		{nutricionModels.GrupoLacteoCod, "Lácteo", "#64B5F6", 2},
		{nutricionModels.GrupoLacteoVegCod, "Lácteo vegetal", "#81C784", 3},
		{nutricionModels.GrupoLegumbreCod, "Legumbre", "#FFB74D", 4},
		{nutricionModels.GrupoCarboCod, "Carbohidrato", "#FFF176", 5},
		{nutricionModels.GrupoFrutaCod, "Fruta", "#F06292", 6},
		{nutricionModels.GrupoFrutaSecaCod, "Fruta seca", "#A1887F", 7},
		{nutricionModels.GrupoGrasaSalCod, "Grasa saludable", "#4DB6AC", 8},
		{nutricionModels.GrupoFrutosSecosCod, "Frutos secos", "#8D6E63", 9},
		{nutricionModels.GrupoAzucarCod, "Azúcar", "#FFD54F", 10},
		{nutricionModels.GrupoVegetalCod, "Vegetal", "#AED581", 11},
		{nutricionModels.GrupoPreparacionCod, "Preparación", "#CE93D8", 12},
	}

	grupoMap := make(map[string]uint) // codigo → ID

	for _, g := range gruposSeed {
		var grupo nutricionModels.NutricionGrupoAlimento
		result := db.Where("codigo = ?", g.Codigo).FirstOrCreate(&grupo, nutricionModels.NutricionGrupoAlimento{
			Codigo: g.Codigo,
			Nombre: g.Nombre,
			Color:  g.Color,
			Orden:  g.Orden,
			State:  "A",
		})
		if result.Error != nil {
			log.Fatalf("Error insertando grupo %s: %v", g.Codigo, result.Error)
		}
		grupoMap[g.Codigo] = grupo.ID
	}

	log.Printf("Grupos insertados/verificados: %d", len(gruposSeed))

	// ── 2. Alimentos (100 g de porción cada uno) ───────────────────────────────

	type alimentoSeed struct {
		Nombre         string
		GrupoCodigo    string
		Calorias       float64
		ProteínaG      float64
		CarbohidratosG float64
		GrasasG        float64
	}

	alimentos := []alimentoSeed{
		// Proteínas
		{"Pechuga de pavo", nutricionModels.GrupoProteinaCod, 135, 30.0, 0, 1.0},
		{"Pollo pechuga", nutricionModels.GrupoProteinaCod, 140, 25.0, 0, 2.0},
		{"Carne de res", nutricionModels.GrupoProteinaCod, 250, 27.0, 0, 15.0},
		{"Carne de cerdo", nutricionModels.GrupoProteinaCod, 242, 25.0, 0, 14.0},
		{"Atún en aceite", nutricionModels.GrupoProteinaCod, 198, 26.0, 0, 10.0},
		{"Atún en agua", nutricionModels.GrupoProteinaCod, 132, 28.0, 0, 1.0},
		{"Huevo revuelto", nutricionModels.GrupoProteinaCod, 155, 13.0, 1.0, 11.0},
		{"Huevo duro", nutricionModels.GrupoProteinaCod, 155, 13.0, 1.1, 11.0},

		// Lácteos
		{"Queso sin grasa", nutricionModels.GrupoLacteoCod, 98, 11.0, 3.0, 0.5},
		{"Queso Mozzarella", nutricionModels.GrupoLacteoCod, 280, 22.0, 2.0, 21.0},
		{"Queso fresco", nutricionModels.GrupoLacteoCod, 264, 14.0, 4.0, 21.0},
		{"Yogur natural", nutricionModels.GrupoLacteoCod, 61, 3.5, 4.7, 3.0},
		{"Yogur griego", nutricionModels.GrupoLacteoCod, 59, 10.0, 3.0, 0.4},
		{"Leche semidescremada", nutricionModels.GrupoLacteoCod, 46, 3.4, 5.0, 1.5},
		{"Leche entera", nutricionModels.GrupoLacteoCod, 61, 3.2, 5.0, 3.3},
		{"Leche deslactosada", nutricionModels.GrupoLacteoCod, 50, 3.3, 5.0, 1.0},

		// Lácteo vegetal
		{"Leche de almendras", nutricionModels.GrupoLacteoVegCod, 17, 0.6, 0.3, 1.2},

		// Legumbres
		{"Frijoles", nutricionModels.GrupoLegumbreCod, 127, 9.0, 23.0, 0.5},
		{"Lentejas", nutricionModels.GrupoLegumbreCod, 116, 9.0, 20.0, 0.4},
		{"Garbanzos", nutricionModels.GrupoLegumbreCod, 164, 9.0, 27.0, 2.6},

		// Carbohidratos
		{"Arroz blanco cocido", nutricionModels.GrupoCarboCod, 130, 2.7, 28.0, 0.3},
		{"Arroz integral", nutricionModels.GrupoCarboCod, 123, 2.7, 26.0, 1.0},
		{"Quinoa cocida", nutricionModels.GrupoCarboCod, 120, 4.0, 21.0, 1.9},
		{"Pasta cocida", nutricionModels.GrupoCarboCod, 131, 5.0, 25.0, 1.1},
		{"Papa", nutricionModels.GrupoCarboCod, 77, 2.0, 17.0, 0.1},
		{"Camote", nutricionModels.GrupoCarboCod, 86, 1.6, 20.0, 0.1},
		{"Yuca", nutricionModels.GrupoCarboCod, 160, 1.4, 38.0, 0.3},
		{"Mote", nutricionModels.GrupoCarboCod, 96, 3.0, 21.0, 1.0},
		{"Cebada cocida", nutricionModels.GrupoCarboCod, 123, 2.3, 28.0, 0.4},
		{"Tortilla de maíz", nutricionModels.GrupoCarboCod, 218, 5.7, 45.0, 2.8},
		{"Pan integral", nutricionModels.GrupoCarboCod, 247, 13.0, 41.0, 4.2},
		{"Pan masa madre", nutricionModels.GrupoCarboCod, 230, 8.0, 47.0, 1.5},
		{"Galletas de arroz", nutricionModels.GrupoCarboCod, 387, 8.0, 81.0, 3.0},
		{"Galletas integrales", nutricionModels.GrupoCarboCod, 430, 7.0, 72.0, 12.0},
		{"Galletas de maíz", nutricionModels.GrupoCarboCod, 430, 7.0, 72.0, 12.0},
		{"Plátano verde", nutricionModels.GrupoCarboCod, 122, 1.3, 31.0, 0.4},
		{"Plátano maduro", nutricionModels.GrupoCarboCod, 116, 1.3, 31.0, 0.4},
		{"Granola", nutricionModels.GrupoCarboCod, 471, 10.0, 64.0, 20.0},

		// Frutas
		{"Guineo / banana", nutricionModels.GrupoFrutaCod, 89, 1.1, 23.0, 0.3},
		{"Manzana", nutricionModels.GrupoFrutaCod, 52, 0.3, 14.0, 0.2},
		{"Pera", nutricionModels.GrupoFrutaCod, 57, 0.4, 15.0, 0.1},
		{"Fresas", nutricionModels.GrupoFrutaCod, 32, 0.7, 8.0, 0.3},
		{"Kiwi", nutricionModels.GrupoFrutaCod, 61, 1.1, 15.0, 0.5},
		{"Durazno", nutricionModels.GrupoFrutaCod, 39, 0.9, 10.0, 0.3},
		{"Moras", nutricionModels.GrupoFrutaCod, 43, 1.4, 10.0, 0.5},
		{"Papaya", nutricionModels.GrupoFrutaCod, 43, 0.5, 11.0, 0.3},
		{"Piña", nutricionModels.GrupoFrutaCod, 50, 0.5, 13.0, 0.1},
		{"Melón", nutricionModels.GrupoFrutaCod, 34, 0.8, 8.0, 0.2},
		{"Naranja", nutricionModels.GrupoFrutaCod, 47, 0.9, 12.0, 0.1},
		{"Uvas", nutricionModels.GrupoFrutaCod, 69, 0.7, 18.0, 0.2},

		// Fruta seca
		{"Pasas", nutricionModels.GrupoFrutaSecaCod, 299, 3.0, 79.0, 0.5},

		// Grasa saludable
		{"Aguacate", nutricionModels.GrupoGrasaSalCod, 160, 2.0, 9.0, 15.0},

		// Frutos secos
		{"Almendras", nutricionModels.GrupoFrutosSecosCod, 579, 21.0, 22.0, 50.0},
		{"Nueces", nutricionModels.GrupoFrutosSecosCod, 654, 15.0, 14.0, 65.0},

		// Azúcar
		{"Miel", nutricionModels.GrupoAzucarCod, 304, 0.3, 82.0, 0.0},

		// Vegetales
		{"Lechuga", nutricionModels.GrupoVegetalCod, 15, 1.4, 3.0, 0.2},
		{"Tomate", nutricionModels.GrupoVegetalCod, 18, 0.9, 4.0, 0.2},
		{"Espinaca", nutricionModels.GrupoVegetalCod, 23, 2.9, 3.6, 0.4},
		{"Zanahoria", nutricionModels.GrupoVegetalCod, 41, 0.9, 10.0, 0.2},
		{"Remolacha", nutricionModels.GrupoVegetalCod, 43, 1.6, 10.0, 0.2},
		{"Pimiento", nutricionModels.GrupoVegetalCod, 31, 1.0, 6.0, 0.3},
		{"Champiñones", nutricionModels.GrupoVegetalCod, 22, 3.1, 3.3, 0.3},
		{"Maíz dulce", nutricionModels.GrupoVegetalCod, 96, 3.4, 21.0, 1.5},
		{"Brócoli", nutricionModels.GrupoVegetalCod, 34, 2.8, 7.0, 0.4},
		{"Cebolla", nutricionModels.GrupoVegetalCod, 40, 1.1, 9.0, 0.1},

		// Preparaciones típicas
		{"Muchín de yuca", nutricionModels.GrupoPreparacionCod, 280, 5.0, 45.0, 9.0},
		{"Muchín de papa", nutricionModels.GrupoPreparacionCod, 250, 6.0, 40.0, 8.0},
		{"Chifles", nutricionModels.GrupoPreparacionCod, 536, 2.0, 58.0, 34.0},
		{"Albóndigas de carne", nutricionModels.GrupoPreparacionCod, 250, 17.0, 10.0, 15.0},
		{"Nuggets de pollo", nutricionModels.GrupoPreparacionCod, 296, 15.0, 18.0, 18.0},
	}

	insertados := 0
	for _, a := range alimentos {
		grupoID, ok := grupoMap[a.GrupoCodigo]
		if !ok {
			log.Printf("WARN: grupo %q no encontrado para alimento %q", a.GrupoCodigo, a.Nombre)
			continue
		}

		var existing nutricionModels.NutricionAlimento
		result := db.Where("nombre = ?", a.Nombre).First(&existing)
		if result.Error == nil {
			// Ya existe — actualiza grupo_id si no está asignado
			if existing.GrupoID == nil {
				db.Model(&existing).Update("grupo_id", grupoID)
			}
			continue
		}

		nuevo := nutricionModels.NutricionAlimento{
			Nombre:        a.Nombre,
			GrupoID:       &grupoID,
			GramosPorcion: 100,
			Calorias:       a.Calorias,
			ProteínasG:     a.ProteínaG,
			CarbohidratosG: a.CarbohidratosG,
			GrasasG:        a.GrasasG,
			State:          "A",
		}
		if err := db.Create(&nuevo).Error; err != nil {
			log.Printf("Error insertando alimento %q: %v", a.Nombre, err)
			continue
		}
		insertados++
	}

	log.Printf("Alimentos insertados: %d / %d (los restantes ya existían)", insertados, len(alimentos))
	log.Println("Seed de nutrición completado.")
}
