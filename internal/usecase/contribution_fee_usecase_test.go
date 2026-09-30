package usecase

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseCSV_SampleData(t *testing.T) {
	// Read sample data
	content, err := os.ReadFile("../../data/01.csv")
	if err != nil {
		content, err = os.ReadFile("../../../agent/plan/01.csv")
		if err != nil {
			t.Fatalf("Failed to read sample csv: %v", err)
		}
	}

	uc := &contributionFeeUsecase{}
	data, err := uc.parseCSV(content)
	if err != nil {
		t.Fatalf("parseCSV failed: %v", err)
	}

	if data.Summary.TotalResidents == 0 {
		t.Errorf("Expected total residents > 0, got %d", data.Summary.TotalResidents)
	}

	if data.Summary.TotalCollected == 0 {
		t.Errorf("Expected total collected > 0, got %f", data.Summary.TotalCollected)
	}

	if len(data.Periods) == 0 {
		t.Errorf("Expected periods > 0, got %d", len(data.Periods))
	}

	t.Logf("Total residents: %d", data.Summary.TotalResidents)
	t.Logf("Total collected: Rp %.2f", data.Summary.TotalCollected)
	t.Logf("Total access cards: %d", data.Summary.TotalAccessCards)
	t.Logf("Total extra cards: %d", data.Summary.TotalExtraCards)
	t.Logf("Periods count: %d (from %s %d to %s %d)",
		len(data.Periods),
		data.Periods[0].Month, data.Periods[0].Year,
		data.Periods[len(data.Periods)-1].Month, data.Periods[len(data.Periods)-1].Year,
	)

	// Validate first resident
	if len(data.Residents) > 0 {
		r0 := data.Residents[0]
		t.Logf("Resident 0: #%d %s (%s %s-%s), Total: %.0f, Months: %d",
			r0.No, r0.Name, r0.Cluster, r0.Block, r0.Number, r0.TotalPaid, r0.MonthsPaidCount)
		if r0.TotalPaid <= 0 {
			t.Errorf("Expected resident 0 to have paid > 0, got %f", r0.TotalPaid)
		}
	}
}

func TestGetContributionFees_HTTPURL(t *testing.T) {
	// Create mock HTTP server serving sample CSV
	sampleCSV := `No,Nama,Cluster,Blok,Nomor,Kartu Akses,Penambahan Kartu Akses Rp 25.000,Nomor Kartu 1,Nomor Kartu 2,Expired Kartu,Notes,KAS,,,,
,,,,,,,,,,,2026,,,,
,,,,,,,,,,, September, Oktober, November , Desember
1,Warga Test,Akasha,A1,01,TRUE,FALSE,101,,,,"Rp15,000","Rp15,000",,`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sampleCSV))
	}))
	defer ts.Close()

	uc := NewContributionFeeUsecase()
	data, err := uc.GetContributionFees(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("GetContributionFees failed: %v", err)
	}

	if data.Summary.TotalResidents != 1 {
		t.Errorf("Expected 1 resident, got %d", data.Summary.TotalResidents)
	}
	if data.Residents[0].TotalPaid != 30000 {
		t.Errorf("Expected total paid 30000, got %f", data.Residents[0].TotalPaid)
	}
	if data.Residents[0].MonthsPaidCount != 2 {
		t.Errorf("Expected 2 months paid, got %d", data.Residents[0].MonthsPaidCount)
	}
}

func TestNormalizeSpreadsheetURL(t *testing.T) {
	raw := "https://docs.google.com/spreadsheets/d/1dinqvqrq8g7drc12yM5kpH1e1rv4-XXw/edit?gid=353831481#gid=353831481"
	expected := "https://docs.google.com/spreadsheets/d/1dinqvqrq8g7drc12yM5kpH1e1rv4-XXw/export?format=csv&gid=353831481"
	got := normalizeSpreadsheetURL(raw)
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestGetContributionFees_RealGoogleSheets(t *testing.T) {
	sheetWebURL := "https://docs.google.com/spreadsheets/d/1dinqvqrq8g7drc12yM5kpH1e1rv4-XXw/edit?gid=353831481#gid=353831481"
	uc := NewContributionFeeUsecase()
	data, err := uc.GetContributionFees(context.Background(), sheetWebURL)
	if err != nil {
		t.Fatalf("Failed to fetch real Google Sheets data: %v", err)
	}

	if data.Summary.TotalResidents == 0 {
		t.Errorf("Expected residents > 0 from real sheets")
	}
	t.Logf("Fetched real sheets data successfully! Residents: %d, Collected: Rp %.0f, Source: %s",
		data.Summary.TotalResidents, data.Summary.TotalCollected, data.Source)
}
