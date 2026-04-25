package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	// если есть recureType то должны создать строку в recurrence_tasks и создать первую таску
	// crate будет возвращать первую задачу из правил периодичности
	// на каждый тип периодичноти - своя логика генерации задач

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		ScheduledAt: normalized.ScheduledAt,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	if normalized.RecurType != nil {

		recModel := &taskdomain.RecurTask{
			Title:             normalized.Title,
			Description:       normalized.Description,
			RecurType:         *normalized.RecurType,
			IntervalDays:      normalized.IntervalDays,
			DayOfMonth:        normalized.DayOfMonth,
			SpecificDates:     normalized.SpecificDates,
			Parity:            normalized.Parity,
			ScheduledAt:       normalized.ScheduledAt,
			CreatedAt:         now,
			LastGeneratedDate: now,
		}

		recurTask, err := s.repo.CreateRecurTask(ctx, recModel)
		if err != nil {
			return nil, err
		}
		model.RecurID = &recurTask.ID
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	nulOtherRecurRows(*input.RecurType, &input)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurType != nil && !input.RecurType.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid recur type", ErrInvalidInput)
	}

	switch *input.RecurType {
	case taskdomain.DailyType:
		if !input.IntervalDays.Valid() {
			return CreateInput{}, fmt.Errorf("%w: invalid interval range", ErrInvalidInput)
		}
	case taskdomain.Parity:
		if !input.Parity.Valid() {
			return CreateInput{}, fmt.Errorf("%w: invalid parity", ErrInvalidInput)
		}
		if !input.Parity.ValidateDate(input.ScheduledAt) {
			return CreateInput{}, fmt.Errorf("%w: the date doesn't match with parity", ErrInvalidInput)
		}
	case taskdomain.MonthlyType:
		if !input.DayOfMonth.Valid() {
			return CreateInput{}, fmt.Errorf("%w: invalid day of month", ErrInvalidInput)
		}
	case taskdomain.SpecificDates:
		for _, date := range *input.SpecificDates {
			if !date.Valid() {
				return CreateInput{}, fmt.Errorf("%w: ivalid date format. Valid is DD-MM", ErrInvalidInput)
			}
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func nulOtherRecurRows(rt taskdomain.RecurType, input *CreateInput) {
	switch rt {
	case taskdomain.DailyType:
		input.DayOfMonth = nil
		input.SpecificDates = nil
		input.Parity = nil
	case taskdomain.MonthlyType:
		input.SpecificDates = nil
		input.IntervalDays = nil
		input.Parity = nil
	case taskdomain.SpecificDates:
		input.DayOfMonth = nil
		input.IntervalDays = nil
		input.Parity = nil
	case taskdomain.Parity:
		input.DayOfMonth = nil
		input.IntervalDays = nil
		input.SpecificDates = nil
	}
}
