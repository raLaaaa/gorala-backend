package utilities

import "time"

type TimeUtil struct{}

func (t TimeUtil) RoundDate(date time.Time) time.Time {

	rounded := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return rounded
}
