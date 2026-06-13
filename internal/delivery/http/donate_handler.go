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
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse(http.StatusBadRequest, "missing id parameter"))
		return
	}

	donateDetails, err := h.donateUsecase.GetDonateDetailsByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	var amountTotal float64
	for _, d := range donateDetails {
		amountTotal += d.Amount
	}

	data := map[string]interface{}{
		"detail":       donateDetails,
		"amount_total": amountTotal,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse(http.StatusOK, "Donate details retrieved successfully", data))
}
