package routes

import (
	"net/http"

	"github.com/ReynoldArun09/school-procurement-management/internal/handlers"
	"github.com/ReynoldArun09/school-procurement-management/internal/repositories"
	"github.com/ReynoldArun09/school-procurement-management/internal/services"
	"gorm.io/gorm"
)

func RegisterAcademicYearRoutes(mux *http.ServeMux, db *gorm.DB) {

	academicYearRepo := repositories.NewGormRepository(db)
	academicYearService := services.NewService(academicYearRepo)
	academicYearHandler := handlers.NewHandler(academicYearService)

	mux.HandleFunc("POST /create-academic-year", academicYearHandler.Create)
}
