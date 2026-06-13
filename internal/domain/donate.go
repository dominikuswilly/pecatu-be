package domain

import (
	"context"
	"time"
)

type DonateDetail struct {
	ID        string    `json:"id"`
	DonateID  string    `json:"donate_id"`
	Cluster   *string   `json:"cluster"`
	Block     *string   `json:"block"`
	Number    *string   `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	Amount    float64   `json:"amount"`
	Name      string    `json:"name"`
}

type DonateRepository interface {
	GetByDonateID(ctx context.Context, donateID string) ([]*DonateDetail, error)
}

type DonateUsecase interface {
	GetDonateDetailsByID(ctx context.Context, donateID string) ([]*DonateDetail, error)
}
