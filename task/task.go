package task

import (
	"GoRunway/models"
	"errors"
	"fmt"
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
	startedAt   *time.Time
	completedAt *time.Time
	err         error
	rwmu        sync.RWMutex
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
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.id
}

func (t *Task) Title() string {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.title
}

func (t *Task) FlightID() string {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.flightID
}

func (t *Task) Priority() models.Priority {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.priority
}

func (t *Task) Status() models.TaskStatus {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.status
}

func (t *Task) Error() error {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.err
}

func (t *Task) Cancel() error {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	if t.status == models.StatusPending {
		t.status = models.StatusCancelled
		return nil
	}
	return fmt.Errorf("Task id: %s, title: %s, flightID: %s, status: %s, Err: %w", t.id, t.title, t.flightID, t.Status(), models.ErrTaskCannotCancel)
}

func (t *Task) Reset() error {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	if t.status == models.StatusRunning {
		return errors.New("cannot reset running task")
	}
	t.status = models.StatusPending
	t.err = nil
	t.startedAt = nil
	t.completedAt = nil
	return nil
}

func (t *Task) Execute() error {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	defer func() error {
		if r := recover(); r != nil {
			t.status = models.StatusFailed
			completedAt := time.Now()
			t.completedAt = &completedAt
			err := fmt.Errorf("task panicked: %v", r)
			t.err = err
			return err
		}
		return nil
	}()
	if t.status != models.StatusPending {
		return models.ErrTaskNotPending
	}

	t.status = models.StatusRunning
	startedAt := time.Now()
	t.startedAt = &startedAt

	err := t.action()
	if err != nil {
		t.status = models.StatusFailed
		t.err = err
		return err
	}

	t.status = models.StatusCompleted
	completedAt := time.Now()
	t.completedAt = &completedAt

	return nil
}

func (t *Task) GetInfo() models.TaskInfo {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return models.TaskInfo{
		ID:        t.id,
		Title:     t.title,
		FlightID:  t.flightID,
		Priority:  t.priority,
		Status:    t.status,
		CreatedAt: t.startedAt.Format(time.RFC3339),
	}
}
