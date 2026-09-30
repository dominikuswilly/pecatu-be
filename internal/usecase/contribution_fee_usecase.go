package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"pecatu-be/internal/domain"
)

type contributionFeeUsecase struct {
	httpClient *http.Client
}

func NewContributionFeeUsecase() domain.ContributionFeeUsecase {
	return &contributionFeeUsecase{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (u *contributionFeeUsecase) GetContributionFees(ctx context.Context, sourceURL string) (*domain.ContributionFeeData, error) {
	// Step 1: Resolve content from URL or local fallback
	content, sourceUsed, err := u.fetchContent(ctx, sourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contribution fee content: %w", err)
	}

	// Step 2 & 3: Process CSV content and extract into structured domain model
	data, err := u.parseCSV(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse contribution fee CSV: %w", err)
	}

	data.Source = sourceUsed
	return data, nil
}

func (u *contributionFeeUsecase) fetchContent(ctx context.Context, sourceURL string) ([]byte, string, error) {
	cleanURL := strings.TrimSpace(sourceURL)
	if cleanURL == "" {
		cleanURL = os.Getenv("CONTRIBUTION_FEE_URL")
		if cleanURL == "" {
			cleanURL = os.Getenv("IURAN_CSV_URL")
		}
	}

	// If HTTP/HTTPS URL, download it
	if strings.HasPrefix(cleanURL, "http://") || strings.HasPrefix(cleanURL, "https://") {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, cleanURL, nil)
		if err != nil {
			log.Printf("[ContributionFee] Failed to create request for %s: %v", cleanURL, err)
		} else {
			resp, err := u.httpClient.Do(req)
			if err != nil {
				log.Printf("[ContributionFee] Error downloading from %s: %v, falling back to local file", cleanURL, err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					body, err := io.ReadAll(resp.Body)
					if err == nil && len(body) > 0 {
						return body, cleanURL, nil
					}
					log.Printf("[ContributionFee] Error reading body from %s: %v", cleanURL, err)
				} else {
					log.Printf("[ContributionFee] Non-200 status code %d from %s", resp.StatusCode, cleanURL)
				}
			}
		}
	}

	// If explicit local file path provided
	if cleanURL != "" && !strings.HasPrefix(cleanURL, "http://") && !strings.HasPrefix(cleanURL, "https://") {
		cleanURL = strings.TrimPrefix(cleanURL, "file://")
		if data, err := os.ReadFile(cleanURL); err == nil {
			return data, cleanURL, nil
		}
	}

	// Default fallback paths
	candidatePaths := []string{
		"data/01.csv",
		filepath.Join("..", "agent", "plan", "01.csv"),
		filepath.Join("..", "agent", "temp", "01.csv"),
		filepath.Join("agent", "plan", "01.csv"),
		filepath.Join("agent", "temp", "01.csv"),
	}

	for _, path := range candidatePaths {
		if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
			return data, path, nil
		}
	}

	return nil, "", fmt.Errorf("no valid CSV data source found (tried URL: %q and fallback local files)", cleanURL)
}

func (u *contributionFeeUsecase) parseCSV(content []byte) (*domain.ContributionFeeData, error) {
	// Strip UTF-8 BOM if present
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))

	reader := csv.NewReader(bytes.NewReader(content))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv read error: %w", err)
	}

	if len(records) < 3 {
		return nil, fmt.Errorf("csv content has insufficient rows (got %d, need at least 3)", len(records))
	}

	// Row 1: Years row (index 1)
	// Row 2: Months row (index 2)
	yearsRow := records[1]
	monthsRow := records[2]

	var periods []domain.Period
	var periodColIndices []int

	currentYear := 2026 // default baseline if first column doesn't specify
	maxCol := len(monthsRow)
	if len(yearsRow) > maxCol {
		maxCol = len(yearsRow)
	}

	// Contribution months start at column index 11
	startIndex := 11
	for col := startIndex; col < maxCol; col++ {
		// Update currentYear if specified in yearsRow
		if col < len(yearsRow) {
			yrStr := strings.TrimSpace(yearsRow[col])
			if yrStr != "" {
				if parsedYr, err := strconv.Atoi(yrStr); err == nil && parsedYr > 2000 {
					currentYear = parsedYr
				}
			}
		}

		if col < len(monthsRow) {
			monthStr := strings.TrimSpace(monthsRow[col])
			if monthStr != "" {
				periodKey := fmt.Sprintf("%d-%s", currentYear, strings.ToLower(monthStr))
				periods = append(periods, domain.Period{
					Year:  currentYear,
					Month: monthStr,
					Key:   periodKey,
				})
				periodColIndices = append(periodColIndices, col)
			}
		}
	}

	var residents []domain.ResidentContribution
	var totalCollected float64
	var totalAccessCards int
	var totalExtraCards int
	clusterMap := make(map[string]*domain.ClusterSummary)

	// Resident data rows start from row index 3
	for i := 3; i < len(records); i++ {
		row := records[i]
		if len(row) <= 1 {
			continue
		}

		name := ""
		if len(row) > 1 {
			name = strings.TrimSpace(row[1])
		}
		if name == "" && strings.TrimSpace(row[0]) == "" {
			continue
		}

		noVal, _ := strconv.Atoi(strings.TrimSpace(row[0]))

		cluster := ""
		if len(row) > 2 {
			cluster = strings.TrimSpace(row[2])
		}

		block := ""
		if len(row) > 3 {
			block = strings.TrimSpace(row[3])
		}

		number := ""
		if len(row) > 4 {
			number = strings.TrimSpace(row[4])
		}

		hasCard := false
		if len(row) > 5 {
			hasCard = strings.EqualFold(strings.TrimSpace(row[5]), "TRUE")
		}

		hasExtraCard := false
		if len(row) > 6 {
			hasExtraCard = strings.EqualFold(strings.TrimSpace(row[6]), "TRUE")
		}

		card1 := ""
		if len(row) > 7 {
			card1 = strings.TrimSpace(row[7])
		}

		card2 := ""
		if len(row) > 8 {
			card2 = strings.TrimSpace(row[8])
		}

		expired := ""
		if len(row) > 9 {
			expired = strings.TrimSpace(row[9])
		}

		notes := ""
		if len(row) > 10 {
			notes = strings.TrimSpace(row[10])
		}

		var payments []domain.MonthlyPayment
		var residentTotalPaid float64
		var monthsPaidCount int

		for pIdx, colIdx := range periodColIndices {
			cellVal := ""
			if colIdx < len(row) {
				cellVal = row[colIdx]
			}
			amount := parseCurrencyAmount(cellVal)
			paid := amount > 0
			if paid {
				residentTotalPaid += amount
				monthsPaidCount++
			}

			pInfo := periods[pIdx]
			payments = append(payments, domain.MonthlyPayment{
				Year:   pInfo.Year,
				Month:  pInfo.Month,
				Period: fmt.Sprintf("%s %d", pInfo.Month, pInfo.Year),
				Amount: amount,
				Paid:   paid,
			})
		}

		if hasCard {
			totalAccessCards++
		}
		if hasExtraCard {
			totalExtraCards++
		}
		totalCollected += residentTotalPaid

		// Cluster aggregation
		cKey := cluster
		if cKey == "" {
			cKey = "Lainnya"
		}
		if _, exists := clusterMap[cKey]; !exists {
			clusterMap[cKey] = &domain.ClusterSummary{
				Name: cKey,
			}
		}
		clusterMap[cKey].TotalResidents++
		clusterMap[cKey].TotalPaidAmount += residentTotalPaid
		if hasCard {
			clusterMap[cKey].TotalAccessCards++
		}

		residents = append(residents, domain.ResidentContribution{
			No:                noVal,
			Name:              name,
			Cluster:           cluster,
			Block:             block,
			Number:            number,
			HasAccessCard:     hasCard,
			HasAdditionalCard: hasExtraCard,
			CardNumber1:       card1,
			CardNumber2:       card2,
			CardExpired:       expired,
			Notes:             notes,
			TotalPaid:         residentTotalPaid,
			MonthsPaidCount:   monthsPaidCount,
			Payments:          payments,
		})
	}

	var clusterSummaries []domain.ClusterSummary
	for _, cs := range clusterMap {
		clusterSummaries = append(clusterSummaries, *cs)
	}

	summary := domain.ContributionFeeSummary{
		TotalResidents:     len(residents),
		TotalCollected:     totalCollected,
		TotalAccessCards:   totalAccessCards,
		TotalExtraCards:    totalExtraCards,
		TotalMonthsCovered: len(periods),
		ClusterSummaries:   clusterSummaries,
	}

	return &domain.ContributionFeeData{
		Summary:   summary,
		Periods:   periods,
		Residents: residents,
	}, nil
}

func parseCurrencyAmount(raw string) float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}

	// Filter digits and dots
	var b strings.Builder
	hasDot := false
	for _, r := range raw {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if r == '.' && !hasDot {
			b.WriteRune(r)
			hasDot = true
		}
	}

	s := b.String()
	if s == "" {
		return 0
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}
