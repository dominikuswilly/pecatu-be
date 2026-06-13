package postgres

import (
	"context"
	"database/sql"
	"pecatu-be/internal/domain"
)

type donateRepository struct {
	db *sql.DB
}

func NewDonateRepository(db *sql.DB) domain.DonateRepository {
	return &donateRepository{
		db: db,
	}
}

func (r *donateRepository) GetByID(ctx context.Context, id string) (*domain.Donate, error) {
	// Dummy data implementation as requested
	return &domain.Donate{
		ID:     id,
		Amount: 10000,
		Status: "success",
	}, nil
}
