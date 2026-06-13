package domain

import "context"

type Donate struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type DonateRepository interface {
	GetByID(ctx context.Context, id string) (*Donate, error)
}

type DonateUsecase interface {
	GetDonateByID(ctx context.Context, id string) (*Donate, error)
}
