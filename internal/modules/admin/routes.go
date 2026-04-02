package admin

import (
	"saas-medico/internal/database"
	"saas-medico/internal/middleware"
	"saas-medico/internal/modules/admin/handlers"
	"saas-medico/internal/modules/admin/repositories"
	"saas-medico/internal/modules/admin/services"
	authModels "saas-medico/internal/modules/auth/models"
	authRepos "saas-medico/internal/modules/auth/repositories"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewAdminRepository(database.GetDB())
	svc := services.NewAdminService(repo)
	h := handlers.NewAdminHandler(svc)

	db := database.GetDB()
	rolSvc := services.NewRolService(authRepos.NewRolRepository(db), repo)
	rh := handlers.NewRolHandler(rolSvc)

	admin := router.Group("/admin")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireRoles(authModels.RolSuperAdmin, authModels.RolAdmin))
	{
		// Clínicas
		admin.GET("/clinicas", h.ListClinicas)
		admin.POST("/clinicas", h.CreateClinica)
		//admin.GET("/clinicaByUser", h.get)

		clinicas := admin.Group("/clinicas/:id")
		{
			clinicas.GET("", h.GetClinica)
			clinicas.PUT("", h.UpdateClinica)
			clinicas.DELETE("", h.DeleteClinica)

			// Usuarios por clínica
			clinicas.GET("/usuarios", h.ListUsuariosClinica)
			clinicas.POST("/usuarios", h.AsignarUsuario)
			clinicas.DELETE("/usuarios/:usuarioID", h.RemoverUsuario)
		}

		// Profesiones
		admin.GET("/profesiones", h.ListProfesiones)
		admin.POST("/profesiones", h.CreateProfesion)
		admin.PUT("/profesiones/:id", h.UpdateProfesion)
		admin.DELETE("/profesiones/:id", h.DeleteProfesion)

		// Planes SaaS
		admin.GET("/planes-saas", h.ListPlanes)
		admin.POST("/planes-saas", h.CreatePlan)

		// Transacciones (menú del sistema)
		admin.GET("/transacciones", rh.ListTransacciones)

		// Roles
		admin.GET("/roles", rh.List)
		admin.POST("/roles", rh.Create)
		admin.GET("/roles/:id", rh.Get)
		admin.PUT("/roles/:id", rh.Update)
		admin.GET("/roles/:id/transacciones", rh.ListTransaccionesByRol)
		admin.POST("/roles/:id/transacciones", rh.AsignarTransacciones)
		admin.DELETE("/roles/:id/transacciones/:transaccionId", rh.RevocarTransaccion)

		// Usuarios → Roles
		admin.GET("/usuarios/:usuarioId/roles", rh.ListRolesByUsuario)
		admin.POST("/usuarios/:usuarioId/roles", rh.AsignarRolAUsuario)
		admin.DELETE("/usuarios/:usuarioId/roles/:rolId", rh.RevocarRolDeUsuario)
	}
}
