package handlers

import (
	"Mini-Quicko/internal/core/models"
	"Mini-Quicko/internal/core/ports"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	service ports.Service
}

func NewHTTPHandler(service ports.Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/products/{productId}/analyze", h.AnalyzeProduct).Methods("GET")
	router.HandleFunc("/products/{productId}/history", h.GetPriceHistory).Methods("GET")
	router.HandleFunc("/products/{productId}/info", h.GetProductInfo).Methods("GET")
	router.HandleFunc("/products/save-kaspi-data", h.SaveKaspiData).Methods("POST")
}

func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.service.HealthCheck(r.Context()); err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Service unavailable")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) AnalyzeProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	if productID == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	analysis, err := h.service.AnalyzeProduct(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, analysis)
}

func (h *HTTPHandler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	if productID == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	history, err := h.service.GetPriceHistory(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, history)
}

func (h *HTTPHandler) GetProductInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	if productID == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	info, err := h.service.GetProductInfo(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, info)
}

func (h *HTTPHandler) SaveKaspiData(w http.ResponseWriter, r *http.Request) {
	var request models.KaspiDataRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if request.ProductID == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	if len(request.Offers.Offers) == 0 {
		respondWithError(w, http.StatusBadRequest, "No offers provided")
		return
	}

	analysis, err := h.service.SaveKaspiData(r.Context(), &request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, analysis)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
