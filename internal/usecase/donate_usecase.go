package usecase

import (
	"context"
	"pecatu-be/internal/domain"
)

type donateUsecase struct {
	donateRepo domain.DonateRepository
}

func NewDonateUsecase(donateRepo domain.DonateRepository) domain.DonateUsecase {
	return &donateUsecase{
		donateRepo: donateRepo,
	}
}

func (u *donateUsecase) GetDonateByID(ctx context.Context, id string) (*domain.Donate, error) {
	return u.donateRepo.GetByID(ctx, id)
}
