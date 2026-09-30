package domain

import "context"

type Period struct {
	Year  int    `json:"year"`
	Month string `json:"month"`
	Key   string `json:"key"`
}

type MonthlyPayment struct {
	Year   int     `json:"year"`
	Month  string  `json:"month"`
	Period string  `json:"period"`
	Amount float64 `json:"amount"`
	Paid   bool    `json:"paid"`
}

type ResidentContribution struct {
	No                int              `json:"no"`
	Name              string           `json:"name"`
	Cluster           string           `json:"cluster"`
	Block             string           `json:"block"`
	Number            string           `json:"number"`
	HasAccessCard     bool             `json:"has_access_card"`
	HasAdditionalCard bool             `json:"has_additional_card"`
	CardNumber1       string           `json:"card_number_1"`
	CardNumber2       string           `json:"card_number_2"`
	CardExpired       string           `json:"card_expired"`
	Notes             string           `json:"notes"`
	TotalPaid         float64          `json:"total_paid"`
	MonthsPaidCount   int              `json:"months_paid_count"`
	Payments          []MonthlyPayment `json:"payments"`
}

type ClusterSummary struct {
	Name             string  `json:"name"`
	TotalResidents   int     `json:"total_residents"`
	TotalPaidAmount  float64 `json:"total_paid_amount"`
	TotalAccessCards int     `json:"total_access_cards"`
}

type ContributionFeeSummary struct {
	TotalResidents     int              `json:"total_residents"`
	TotalCollected     float64          `json:"total_collected"`
	TotalAccessCards   int              `json:"total_access_cards"`
	TotalExtraCards    int              `json:"total_extra_cards"`
	TotalMonthsCovered int              `json:"total_months_covered"`
	ClusterSummaries   []ClusterSummary `json:"cluster_summaries"`
}

type ContributionFeeData struct {
	Source    string                 `json:"source"`
	Summary   ContributionFeeSummary `json:"summary"`
	Periods   []Period               `json:"periods"`
	Residents []ResidentContribution `json:"residents"`
}

type ContributionFeeUsecase interface {
	GetContributionFees(ctx context.Context, sourceURL string) (*ContributionFeeData, error)
}
