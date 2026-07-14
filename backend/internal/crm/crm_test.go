package crm_test

import (
	"testing"

	"github.com/crm-platform/backend/internal/crm"
)

func TestListParamsDefaults(t *testing.T) {
	p := crm.ListParams{}
	if p.Page != 0 {
		t.Errorf("expected default page 0, got %d", p.Page)
	}
	if p.PageSize != 0 {
		t.Errorf("expected default page_size 0, got %d", p.PageSize)
	}
}

func TestLeadModelJSON(t *testing.T) {
	l := crm.Lead{
		Title:    "Test Lead",
		Status:   "new",
		Currency: "UZS",
	}
	if l.Title != "Test Lead" {
		t.Errorf("expected title 'Test Lead', got %s", l.Title)
	}
	if l.Status != "new" {
		t.Errorf("expected status 'new', got %s", l.Status)
	}
}

func TestListResponsePagination(t *testing.T) {
	resp := crm.ListResponse{
		Total:      100,
		Page:       3,
		PageSize:   20,
		TotalPages: 5,
	}
	if resp.TotalPages != 5 {
		t.Errorf("expected 5 total pages, got %d", resp.TotalPages)
	}
	if resp.Page != 3 {
		t.Errorf("expected page 3, got %d", resp.Page)
	}
}

func TestDealDefaults(t *testing.T) {
	d := crm.Deal{
		Title:       "Big Deal",
		Currency:    "USD",
		Probability: 75,
	}
	if d.Probability != 75 {
		t.Errorf("expected probability 75, got %d", d.Probability)
	}
	if d.Won != nil {
		t.Error("expected Won to be nil by default")
	}
}
