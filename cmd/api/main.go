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
	// database.RunMigrations()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "nico1234.",
		DB:       0,
		Protocol: 2,
	})

	auth.Setup()

	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	openiaSvc := openiaService.NewOpenIaService(rdb)

	router := gin.Default()
	router.Use(middleware.CORS())

	// Confiar solo en Nginx/local
	err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	if err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong", "status": "healthy"})
	})

	router.Static("/storage", "./storage")

	api := router.Group("/api/v1")
	authMiddleware := auth.GetAuthMiddleware()

	auth.RegisterRoutes(api)
	admin.RegisterRoutes(api, authMiddleware)
	pacientes.RegisterRoutes(api, authMiddleware)
	agenda.RegisterRoutes(api, authMiddleware)
	cobros.RegisterRoutes(api, authMiddleware)
	psicologia.RegisterRoutes(api, authMiddleware)
	nutricion.RegisterRoutes(api, authMiddleware, rdb, openiaSvc)
	odontologia.RegisterRoutes(api, authMiddleware)
	historia.RegisterRoutes(api, authMiddleware)
	tests.RegisterRoutes(api, authMiddleware)

	port := config.AppConfig.ServerPort
	addr := "127.0.0.1:" + port

	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", config.AppConfig.Environment)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
