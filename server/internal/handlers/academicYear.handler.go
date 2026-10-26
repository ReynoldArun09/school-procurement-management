package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ReynoldArun09/school-procurement-management/internal/services"
	"github.com/ReynoldArun09/school-procurement-management/internal/utils"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input services.CreateInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	year, err := h.service.Create(r.Context(), input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.ApiResponse{
		Message:    "Academic Year created",
		StatusCode: http.StatusCreated,
		Data:       year,
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	rawID := strings.TrimSpace(r.PathValue("academicYear"))
	id, err := strconv.ParseUint(rawID, 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	year, err := h.service.GetByID(r.Context(), uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ApiResponse{
		Message:    "Academic Year created",
		StatusCode: http.StatusOK,
		Data:       year,
	})
}
