package dto

import "time"

type ChargeData struct {
	Id           string    `json:"id"`
	JournalId    string    `json:"journal_id"`
	DaysLate     int       `json:"days_late"`
	DailyLateFee int       `json:"daily_late_fee"`
	Total        int       `json:"total"`
	UserId       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
}
