package task

import (
	"GoRunway/models"
	"sync"
	"time"
)

type Task struct {
	id          string
	title       string
	flightID    string
	priority    models.Priority
	status      models.TaskStatus
	action      models.Action
	createdAt   time.Time
	startedAt   time.Time
	completedAt time.Time
	err         error
	mu          sync.Mutex
}

func NewTask(id, title, flightID string, priority models.Priority, action models.Action) (*Task, error) {
	if id == "" {
		return nil, &models.ValidationError{
			Field:  "id",
			Reason: models.ErrInvalidTaskID.Error(),
		}
	}
	if title == "" {
		return nil, &models.ValidationError{
			Field:  "title",
			Reason: models.ErrEmptyTitle.Error(),
		}
	}
	if flightID == "" {
		return nil, &models.ValidationError{
			Field:  "flightID",
			Reason: models.ErrEmptyFlightID.Error(),
		}
	}
	if priority < models.PriorityLow || models.PriorityCritical < priority {
		return nil, &models.ValidationError{
			Field:  "priority",
			Reason: models.ErrInvalidPriority.Error(),
		}
	}
	if action == nil {
		return nil, &models.ValidationError{
			Field:  "action",
			Reason: models.ErrNilAction.Error(),
		}
	}

	return new(Task{
		id:        id,
		title:     title,
		flightID:  flightID,
		priority:  priority,
		status:    models.StatusPending,
		action:    action,
		createdAt: time.Now(),
	}), nil
}

func (t *Task) ID() string {
	// TODO: Safely return task ID
	return ""
}

func (t *Task) Title() string {
	// TODO: Safely return task title
	return ""
}

func (t *Task) FlightID() string {
	// TODO: Safely return flight ID
	return ""
}

func (t *Task) Priority() models.Priority {
	// TODO: Safely return priority
	return 0
}

func (t *Task) Status() models.TaskStatus {
	// TODO: Safely return status
	return ""
}

func (t *Task) Error() error {
	// TODO: Safely return recorded error
	return nil
}

func (t *Task) Cancel() error {
	// TODO: Cancel pending task or return error
	return nil
}

func (t *Task) Reset() error {
	// TODO: Reset completed/failed/cancelled task to PENDING or return error
	return nil
}

func (t *Task) Execute() error {
	// TODO: Execute action with panic recovery, update status and timestamps
	return nil
}

func (t *Task) GetInfo() models.TaskInfo {
	// TODO: Return immutable snapshot of task info
	return models.TaskInfo{}
}
