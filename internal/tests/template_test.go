package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/rearrange/opreiadrcligo/internal/template"
)

// --- RenderIndex ---

func TestRenderIndex_ReturnsContent(t *testing.T) {
	content, err := template.RenderIndex()
	if err != nil {
		t.Fatalf("RenderIndex() error = %v", err)
	}
	if content == "" {
		t.Error("RenderIndex() returned empty string")
	}
}

func TestRenderIndex_ContainsExpectedSections(t *testing.T) {
	content, err := template.RenderIndex()
	if err != nil {
		t.Fatalf("RenderIndex() error = %v", err)
	}
	for _, want := range []string{
		"Architecture Decision Records",
		"## What is an ADR?",
		"## Index",
		"| # | Title | Created Date | Author |",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("RenderIndex() missing %q", want)
		}
	}
}

func TestRenderIndex_ContainsToday(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	content, err := template.RenderIndex()
	if err != nil {
		t.Fatalf("RenderIndex() error = %v", err)
	}
	if !strings.Contains(content, today) {
		t.Errorf("RenderIndex() missing today's date %q", today)
	}
}

// --- RenderADR ---

func TestRenderADR_FormatsNumberWithLeadingZeros(t *testing.T) {
	data := template.ADRData{Number: 7, Title: "Test Decision", Date: "2 Apr 2026", Author: "Alice", Status: "Draft"}
	content, err := template.RenderADR(data)
	if err != nil {
		t.Fatalf("RenderADR() error = %v", err)
	}
	if !strings.Contains(content, "0007") {
		t.Errorf("RenderADR() number not zero-padded to 4 digits; output:\n%s", content)
	}
}

func TestRenderADR_ContainsTitleAuthorDateStatus(t *testing.T) {
	data := template.ADRData{
		Number: 1,
		Title:  "Use Go for CLI",
		Date:   "2 Apr 2026",
		Author: "Bob",
		Status: "Draft",
	}
	content, err := template.RenderADR(data)
	if err != nil {
		t.Fatalf("RenderADR() error = %v", err)
	}
	for _, want := range []string{"Use Go for CLI", "Bob", "2 Apr 2026", "Draft"} {
		if !strings.Contains(content, want) {
			t.Errorf("RenderADR() missing %q in output", want)
		}
	}
}

func TestRenderADR_DefaultStatusIsDraft(t *testing.T) {
	data := template.ADRData{Number: 1, Title: "Test", Date: "2 Apr 2026", Author: "Alice", Status: "Draft"}
	content, err := template.RenderADR(data)
	if err != nil {
		t.Fatalf("RenderADR() error = %v", err)
	}
	if !strings.Contains(content, "Draft") {
		t.Errorf("RenderADR() expected 'Draft' status in output")
	}
}

func TestRenderADR_DateAndAuthorInCorrectColumns(t *testing.T) {
	data := template.ADRData{Number: 1, Title: "Test", Date: "2 Apr 2026", Author: "Alice", Status: "Draft"}
	content, err := template.RenderADR(data)
	if err != nil {
		t.Fatalf("RenderADR() error = %v", err)
	}
	// The row must appear as "| <date> | <author> | <status> |" in that order.
	if !strings.Contains(content, "| 2 Apr 2026 | Alice | Draft |") {
		t.Errorf("RenderADR() columns out of order; output:\n%s", content)
	}
}
