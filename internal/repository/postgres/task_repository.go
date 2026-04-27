package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (recur_id, title, description, status, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, recur_id, title, description, status, scheduled_at, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.RecurID, task.Title, task.Description, task.Status, task.ScheduledAt, task.CreatedAt, task.UpdatedAt)
	created, err := scanTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, recur_id, title, description, status, scheduled_at, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, recur_id, title, description, status, scheduled_at, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID)
	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, recur_id, title, description, status, scheduled_at, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) CreateRecurTask(ctx context.Context, recTask *taskdomain.RecurTask) (*taskdomain.RecurTask, error) {
	const query = `
	INSERT INTO recurrence_tasks (
		title, 
		description, 
		recur_type, 
		interval_days,
		day_of_month,
		specific_dates,
		parity,
		scheduled_at,
		last_generated_date, 
		is_active,
		created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, title, description, recur_type, interval_days, day_of_month, specific_dates, parity, scheduled_at, last_generated_date, is_active, created_at
	`

	row := r.pool.QueryRow(ctx, query,
		recTask.Title,
		recTask.Description,
		recTask.RecurType,
		recTask.IntervalDays,
		recTask.DayOfMonth,
		recTask.SpecificDates,
		recTask.Parity,
		recTask.ScheduledAt,
		recTask.LastGeneratedDate,
		true,
		recTask.CreatedAt,
	)
	created, err := scanRecurTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByScheduledDateRecurTask(ctx context.Context, date time.Time) ([]taskdomain.RecurTask, error) {
	const query = `
	SELECT * FROM recurrence_tasks
	WHERE scheduled_at = $1;
	`

	rows, err := r.pool.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}

	var list []taskdomain.RecurTask
	for rows.Next() {
		taskPtr, err := scanRecurTask(rows)
		if err != nil {
			return nil, fmt.Errorf("Scan failed: %w", err)
		}

		list = append(list, *taskPtr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error %w", err)
	}

	return list, err
}

func (r *Repository) UpdateRecurTask(ctx context.Context, id int64, task taskdomain.RecurTask) (*taskdomain.RecurTask, error) {
	fmt.Println(task.ScheduledAt)
	const query = `
	UPDATE recurrence_tasks
	SET 
		title = $1,
		description = $2,
		recur_type = $3,
		interval_days = $4,
		day_of_month = $5,
		specific_dates = $6,
		parity = $7,
		scheduled_at = $8,
		last_generated_date = $9,
		is_active = $10
	WHERE id = $11
	RETURNING 
		id, title, description, recur_type, interval_days, 
		day_of_month, specific_dates, parity, scheduled_at, 
		last_generated_date, is_active, created_at;
	`

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.RecurType,
		task.IntervalDays,
		task.DayOfMonth,
		task.SpecificDates,
		task.Parity,
		task.ScheduledAt,
		task.LastGeneratedDate,
		task.IsActive,
		id,
	)

	returnedTask, err := scanRecurTask(row)
	if err != nil {
		return nil, fmt.Errorf("failed to update and scan recur task: %w", err)
	}

	return returnedTask, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.RecurID,
		&task.Title,
		&task.Description,
		&status,
		&task.ScheduledAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	return &task, nil
}

type recurTaskScanner interface {
	Scan(dest ...any) error
}

func scanRecurTask(scanner recurTaskScanner) (*taskdomain.RecurTask, error) {
	var recurTask taskdomain.RecurTask

	if err := scanner.Scan(
		&recurTask.ID,
		&recurTask.Title,
		&recurTask.Description,
		&recurTask.RecurType,
		&recurTask.IntervalDays,
		&recurTask.DayOfMonth,
		&recurTask.SpecificDates,
		&recurTask.Parity,
		&recurTask.ScheduledAt,
		&recurTask.LastGeneratedDate,
		&recurTask.IsActive,
		&recurTask.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &recurTask, nil
}
