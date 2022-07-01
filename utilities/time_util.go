package utilities

import "time"

type TimeUtil struct{}

func (t TimeUtil) RoundDate(date time.Time) time.Time {

	rounded := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return rounded
}

func (t TimeUtil) GetTodayWithoutTime(date time.Time) time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	return today
}
