package http

import (
	"encoding/json"
	"net/http"
	"pecatu-be/internal/domain"
)

type DonateHandler struct {
	donateUsecase domain.DonateUsecase
}

func NewDonateHandler(donateUsecase domain.DonateUsecase) *DonateHandler {
	return &DonateHandler{
		donateUsecase: donateUsecase,
	}
}

func (h *DonateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	donate, err := h.donateUsecase.GetDonateByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(donate)
}
