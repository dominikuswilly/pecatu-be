package http

import (
	"encoding/json"
	"net/http"
	"pecatu-be/internal/domain"
)

type ContributionFeeHandler struct {
	usecase domain.ContributionFeeUsecase
}

func NewContributionFeeHandler(usecase domain.ContributionFeeUsecase) *ContributionFeeHandler {
	return &ContributionFeeHandler{
		usecase: usecase,
	}
}

func (h *ContributionFeeHandler) GetContributionFees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sourceURL := r.URL.Query().Get("url")

	data, err := h.usecase.GetContributionFees(r.Context(), sourceURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse(http.StatusOK, "Contribution fee data retrieved successfully", data))
}
