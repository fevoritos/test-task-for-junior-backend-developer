package generator

import (
	"context"
	"fmt"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type TaskGenerator struct {
	taskRepo  TaskRepo
	recurRepo RecurRepo
	now       func() time.Time
}

func NewGenerator(taskRepo TaskRepo, recurRepo RecurRepo) *TaskGenerator {
	return &TaskGenerator{
		taskRepo:  taskRepo,
		recurRepo: recurRepo,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func (scheduler *TaskGenerator) Generate(ctx context.Context) {
	fmt.Println("generating...")
	now := scheduler.now().AddDate(0, 0, 1)

	list, err := scheduler.recurRepo.GetByDateRecurTask(ctx, now)
	if err != nil {
		fmt.Println(err.Error())
	}
	var NewScheduleAt time.Time

	for _, taskOfList := range list {
		interval := *taskOfList.IntervalDays
		NewScheduleAt = now.AddDate(0, 0, int(interval))
		fmt.Println(NewScheduleAt)

		task := taskdomain.Task{
			Title:       taskOfList.Title,
			Description: taskOfList.Description,
			RecurID:     &taskOfList.ID,
			Status:      taskdomain.StatusNew,
			ScheduledAt: NewScheduleAt,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		recurUpdate := taskOfList
		recurUpdate.ScheduledAt = NewScheduleAt
		recurUpdate.LastGeneratedDate = now

		updatedRecurTask, err := scheduler.recurRepo.UpdateRecurTask(ctx, taskOfList.ID, recurUpdate)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(updatedRecurTask)

		createdTask, err := scheduler.taskRepo.Create(ctx, &task)

		if err != nil {
			fmt.Print(err.Error())
		}

		fmt.Printf("Created a %s \n", createdTask.Title)
	}
}
