package main

import (
	"log"

	"saas-medico/internal/config"
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/admin"
	"saas-medico/internal/modules/agenda"
	"saas-medico/internal/modules/auth"
	"saas-medico/internal/modules/cobros"
	"saas-medico/internal/modules/historia"
	"saas-medico/internal/modules/nutricion"
	"saas-medico/internal/modules/odontologia"
	"saas-medico/internal/modules/pacientes"
	"saas-medico/internal/modules/psicologia"
	"saas-medico/internal/modules/tests"
	openiaService "saas-medico/internal/shared/openia"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	config.LoadConfig()

	database.Connect()
	//database.RunMigrations()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "162.243.161.156:6379",
		Password: "nico1234.",
		DB:       0, // use default DB
		Protocol: 2,
	})

	auth.Setup()

	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	openiaService := openiaService.NewOpenIaService(rdb)

	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong", "status": "healthy"})
	})

	// Serve uploaded files (photos, PDFs, resources)
	router.Static("/storage", "./storage")

	api := router.Group("/api/v1")
	authMiddleware := auth.GetAuthMiddleware()

	// Auth (públicas + protegidas)
	auth.RegisterRoutes(api)

	// Admin: clínicas, sucursales, consultorios, profesiones, planes
	admin.RegisterRoutes(api, authMiddleware)

	// Pacientes
	pacientes.RegisterRoutes(api, authMiddleware)

	// Agenda: citas, sesiones, horarios, bloqueos
	agenda.RegisterRoutes(api, authMiddleware)

	// Cobros y pagos
	cobros.RegisterRoutes(api, authMiddleware)

	// Especialidades
	psicologia.RegisterRoutes(api, authMiddleware)
	nutricion.RegisterRoutes(api, authMiddleware, rdb, openiaService)
	odontologia.RegisterRoutes(api, authMiddleware)

	// Historia clínica, formularios, tests
	historia.RegisterRoutes(api, authMiddleware)
	tests.RegisterRoutes(api, authMiddleware)

	port := config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", config.AppConfig.Environment)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
