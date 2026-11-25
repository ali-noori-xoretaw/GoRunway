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

type validationError struct {
	ve  *models.ValidationError
	err error
}

func (e *validationError) Error() string {
	return e.ve.Error()
}

func (e *validationError) Unwrap() []error {
	return []error{e.ve, e.err}

}

func NewTask(id, title, flightID string, priority models.Priority, action models.Action) (*Task, error) {
	if id == "" {
		return nil, &validationError{
			ve: &models.ValidationError{
				Field:  "id",
				Reason: models.ErrInvalidTaskID.Error(),
			},
			err: models.ErrInvalidTaskID,
		}
	}
	if title == "" {
		return nil, &validationError{
			ve: &models.ValidationError{
				Field:  "title",
				Reason: models.ErrEmptyTitle.Error(),
			},
			err: models.ErrEmptyTitle,
		}
	}
	if flightID == "" {
		return nil, &validationError{
			ve: &models.ValidationError{
				Field:  "flightID",
				Reason: models.ErrEmptyFlightID.Error(),
			},
			err: models.ErrEmptyFlightID,
		}
	}
	if priority < models.PriorityLow || models.PriorityCritical < priority {
		return nil, &validationError{
			ve: &models.ValidationError{
				Field:  "priority",
				Reason: models.ErrInvalidPriority.Error(),
			},
			err: models.ErrInvalidPriority,
		}
	}
	if action == nil {
		return nil, &validationError{
			ve: &models.ValidationError{
				Field:  "action",
				Reason: models.ErrNilAction.Error(),
			},
			err: models.ErrNilAction,
		}
	}

	return &Task{
		id:        id,
		title:     title,
		flightID:  flightID,
		priority:  priority,
		status:    models.StatusPending,
		action:    action,
		createdAt: time.Now(),
	}, nil
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
	return fmt.Errorf("Task id: %s, title: %s, flightID: %s, status: %s, Err: %w", t.id, t.title, t.flightID, t.status, models.ErrTaskCannotCancel)
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

func (t *Task) Execute() (retErr error) {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			t.status = models.StatusFailed
			completedAt := time.Now()
			t.completedAt = &completedAt
			retErr = fmt.Errorf("task panicked: %v", r)
			t.err = retErr
		}
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
		CreatedAt: t.createdAt.Format(time.RFC3339),
	}
}
