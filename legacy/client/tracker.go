package client

import (
	"sync"
	"time"
)

// taskTracker implements TaskTracker.
type taskTracker struct {
	tasks map[string]TaskStatus
	mutex sync.RWMutex
}

// TaskStatus represents the status of a tracked task.
type TaskStatus struct {
	Status    string
	Timestamp time.Time
}

// newTaskTracker creates a new task tracker (unexported).
func newTaskTracker() TaskTracker {
	return &taskTracker{
		tasks: make(map[string]TaskStatus),
	}
}

// TrackTask implements TaskTracker.
func (t *taskTracker) TrackTask(taskID string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.tasks[taskID] = TaskStatus{
		Status:    "pending",
		Timestamp: time.Now(),
	}
}

// UpdateTaskStatus implements TaskTracker.
func (t *taskTracker) UpdateTaskStatus(taskID, status string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.tasks[taskID] = TaskStatus{
		Status:    status,
		Timestamp: time.Now(),
	}
}

// GetTaskStatus implements TaskTracker.
func (t *taskTracker) GetTaskStatus(taskID string) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if status, exists := t.tasks[taskID]; exists {
		return status.Status
	}
	return "unknown"
}

// GetAllTasks returns all tracked tasks.
func (t *taskTracker) GetAllTasks() map[string]TaskStatus {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	result := make(map[string]TaskStatus)
	for id, status := range t.tasks {
		result[id] = status
	}
	return result
}

// RemoveTask removes a task from tracking.
func (t *taskTracker) RemoveTask(taskID string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	delete(t.tasks, taskID)
}
