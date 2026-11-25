package tool

import "time"

const calendar_date_layout = "2006-01-02"
const calendar_month_layout = "2006-01"

func Calendar_date_normalize(data string) (string, bool) {
	parsed, err := time.Parse(calendar_date_layout, data)
	if err != nil {
		return "", false
	}

	return parsed.Format(calendar_date_layout), true
}

func Calendar_month_normalize(data string) (string, bool) {
	parsed, err := time.Parse(calendar_month_layout, data)
	if err != nil {
		return "", false
	}

	return parsed.Format(calendar_month_layout), true
}

func Calendar_month_range(month string) (string, string) {
	parsed, err := time.Parse(calendar_month_layout, month)
	if err != nil {
		return "", ""
	}

	return parsed.Format(calendar_date_layout), parsed.AddDate(0, 1, 0).Format(calendar_date_layout)
}

func Calendar_month_shift(month string, amount int) string {
	parsed, err := time.Parse(calendar_month_layout, month)
	if err != nil {
		return ""
	}

	return parsed.AddDate(0, amount, 0).Format(calendar_month_layout)
}

func Calendar_month_grid(month string) (int, int) {
	parsed, err := time.Parse(calendar_month_layout, month)
	if err != nil {
		return 0, 0
	}

	leading := (int(parsed.Weekday()) + 6) % 7
	days := parsed.AddDate(0, 1, -1).Day()

	return leading, days
}
