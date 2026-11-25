package test

import (
	"testing"

	"opennamu/route/tool"
)

func Test_calendar_date_normalize(t *testing.T) {
	valid := []string{"2024-02-29", "2024-01-01", "2000-12-31"}
	for _, date := range valid {
		got, ok := tool.Calendar_date_normalize(date)
		if !ok || got != date {
			t.Errorf("Calendar_date_normalize(%q) = %q, %v, want %q, true", date, got, ok, date)
		}
	}

	invalid := []string{"", "2023-02-29", "2024-02-30", "2024-13-01", "2024-00-10", "2024-2-9", "2024-02-29 ", "abcd-ef-gh"}
	for _, date := range invalid {
		if got, ok := tool.Calendar_date_normalize(date); ok {
			t.Errorf("Calendar_date_normalize(%q) = %q, true, want false", date, got)
		}
	}
}

func Test_calendar_month_normalize(t *testing.T) {
	valid := []string{"2024-01", "2024-12", "2000-02"}
	for _, month := range valid {
		got, ok := tool.Calendar_month_normalize(month)
		if !ok || got != month {
			t.Errorf("Calendar_month_normalize(%q) = %q, %v, want %q, true", month, got, ok, month)
		}
	}

	invalid := []string{"", "2024-13", "2024-00", "2024-1", "2024-01-01", "abcd-ef"}
	for _, month := range invalid {
		if got, ok := tool.Calendar_month_normalize(month); ok {
			t.Errorf("Calendar_month_normalize(%q) = %q, true, want false", month, got)
		}
	}
}

func Test_calendar_month_range_and_shift(t *testing.T) {
	start, end := tool.Calendar_month_range("2024-12")
	if start != "2024-12-01" || end != "2025-01-01" {
		t.Errorf("Calendar_month_range(2024-12) = %q, %q", start, end)
	}

	start, end = tool.Calendar_month_range("2024-02")
	if start != "2024-02-01" || end != "2024-03-01" {
		t.Errorf("Calendar_month_range(2024-02) = %q, %q", start, end)
	}

	cases := []struct {
		month  string
		amount int
		want   string
	}{
		{"2024-01", -1, "2023-12"},
		{"2024-12", 1, "2025-01"},
		{"2024-01", 12, "2025-01"},
		{"2024-01", -12, "2023-01"},
	}
	for _, test_case := range cases {
		if got := tool.Calendar_month_shift(test_case.month, test_case.amount); got != test_case.want {
			t.Errorf("Calendar_month_shift(%q, %d) = %q, want %q", test_case.month, test_case.amount, got, test_case.want)
		}
	}
}

func Test_calendar_month_grid(t *testing.T) {
	cases := []struct {
		month   string
		leading int
		days    int
	}{
		{"2024-01", 0, 31},
		{"2024-02", 3, 29},
		{"2023-02", 2, 28},
		{"2024-09", 6, 30},
	}
	for _, test_case := range cases {
		leading, days := tool.Calendar_month_grid(test_case.month)
		if leading != test_case.leading || days != test_case.days {
			t.Errorf("Calendar_month_grid(%q) = %d, %d, want %d, %d", test_case.month, leading, days, test_case.leading, test_case.days)
		}
	}

	if leading, days := tool.Calendar_month_grid("2024-13"); leading != 0 || days != 0 {
		t.Errorf("Calendar_month_grid(invalid) = %d, %d, want 0, 0", leading, days)
	}
}
