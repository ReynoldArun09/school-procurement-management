package routes

import (
	"net/http"

	"github.com/ReynoldArun09/school-procurement-management/internal/handlers"
	"github.com/ReynoldArun09/school-procurement-management/internal/repositories"
	"github.com/ReynoldArun09/school-procurement-management/internal/services"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(mux *http.ServeMux, db *gorm.DB) {
	adminRepo := repositories.NewGormAdminAuthRepository(db)
	sessionRepo := repositories.NewGormAdminSessionRepository(db)

	adminService := services.NewAdminAuthService(adminRepo, sessionRepo)

	adminAuthHandler := handlers.NewAdminAuthHandler(adminService)

	mux.HandleFunc("POST /api/admin/login", adminAuthHandler.Login)
	mux.HandleFunc("GET /api/admin/me", adminAuthHandler.Me)
	mux.HandleFunc("GET /api/admin/logout", adminAuthHandler.Logout)

}
