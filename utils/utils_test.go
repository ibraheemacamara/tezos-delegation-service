package utils

import "testing"

func TestGetYearFromTimestamp(t *testing.T) {
	timestamp := "2025-08-17T12:00:00Z"
	expectedYear := "2025"
	year, err := GetYearFromTimestamp(timestamp)
	if err != nil {
		t.Fatalf("GetYearFromTimestamp failed: %v", err)
	}
	if year != expectedYear {
		t.Errorf("expected year %s, got %s", expectedYear, year)
	}
}

func TestGetYearFromTimestampWrong(t *testing.T) {
	timestamp := "2025-08-17"
	_, err := GetYearFromTimestamp(timestamp)
	if err == nil {
		t.Fatalf("GetYearFromTimestamp should have failed")
	}
}
