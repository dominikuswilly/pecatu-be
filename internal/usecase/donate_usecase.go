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

func (u *donateUsecase) GetDonateDetailsByID(ctx context.Context, donateID string) ([]*domain.DonateDetail, error) {
	return u.donateRepo.GetByDonateID(ctx, donateID)
}
