package generator

import (
	"context"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type TaskRepo interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
}

type RecurRepo interface {
	GetByScheduledDateRecurTask(ctx context.Context, date time.Time) ([]taskdomain.RecurTask, error)
	CreateRecurTask(ctx context.Context, recTask *taskdomain.RecurTask) (*taskdomain.RecurTask, error)
	UpdateRecurTask(ctx context.Context, id int64, task taskdomain.RecurTask) (*taskdomain.RecurTask, error)
}
