package task

import "time"

type Status string
type RecurType string
type IntervalDays uint8
type ParityEnum string

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

const (
	Even ParityEnum = "even"
	Odd  ParityEnum = "odd"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RecurID     *int64    `json:"recur_id"`
	Status      Status    `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RecurTask struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	RecurType RecurType `json:"recure_type"`

	IntervalDays  *IntervalDays `json:"interval_days"`
	DayOfMonth    *uint8        `json:"day_of_month"`
	SpecificDates *[]time.Time  `json:"specific_dates"`
	Parity        *ParityEnum   `json:"parity"`

	StartDate         time.Time  `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	LastGeneratedDate time.Time  `json:"last_generated_date"`

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `jsom:"created_at"`
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

func (pe ParityEnum) Valid() bool {
	switch pe {
	case Odd, Even:
		return true
	default:
		return false
	}
}

func (pe ParityEnum) ValidateDate(date time.Time) bool {
	var dayParity ParityEnum
	dayNumber := date.Day()
	if dayNumber%2 == 0 {
		dayParity = Even
	} else {
		dayParity = Odd
	}
	switch pe {
	case Odd:
		if dayParity == Odd {
			return true
		}
	case Even:
		if dayParity == Even {
			return true
		}
	}

	return false
}
