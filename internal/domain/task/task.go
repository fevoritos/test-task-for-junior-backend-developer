package task

import "time"

type Status string
type RecurType string
type IntervalDays uint8

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

const (
	DailyType     RecurType = "daily"
	MonthlyType   RecurType = "monthly"
	SpecificDates RecurType = "specific_dates"
	Parity        RecurType = "parity"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RecureFields struct {
	Type     RecurType
	Interval IntervalDays
}

type RecurTask struct {
	Task

	IntervalDays IntervalDays `json:"interval_days"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (rt RecurType) Valid() bool {
	switch rt {
	case
		DailyType,
		MonthlyType,
		SpecificDates,
		Parity:
		return true
	default:
		return false
	}
}

func (i IntervalDays) Valid() bool {
	if i > 14 {
		return false
	}

	return true
}
