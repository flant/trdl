package tasks_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/werf/vault-plugin-secrets-trdl/pkg/tasks_manager/worker"
)

const taskReasonInvalidatedTask = "the task failed due to restart of the plugin"

func (m *Manager) RunTask(ctx context.Context, reqStorage logical.Storage, taskFunc func(context.Context, logical.Storage) error) (string, error) {
	var taskUUID string
	err := m.doTaskWrap(ctx, reqStorage, taskFunc, func(newTaskFunc func(ctx context.Context) error) error {
		busy, err := m.isBusy(ctx, reqStorage)
		if err != nil {
			return err
		}

		if busy {
			return BusyError
		}

		taskUUID, err = m.queueTask(ctx, newTaskFunc)
		return err
	})

	return taskUUID, err
}

func (m *Manager) AddOptionalTask(ctx context.Context, reqStorage logical.Storage, taskFunc func(context.Context, logical.Storage) error) (string, bool, error) {
	taskUUID, err := m.RunTask(ctx, reqStorage, taskFunc)
	if err != nil {
		if err == BusyError {
			return taskUUID, false, nil
		}

		return "", false, err
	}

	return taskUUID, true, nil
}

func (m *Manager) AddTask(ctx context.Context, reqStorage logical.Storage, taskFunc func(context.Context, logical.Storage) error) (string, error) {
	var taskUUID string
	err := m.doTaskWrap(ctx, reqStorage, taskFunc, func(newTaskFunc func(ctx context.Context) error) error {
		var err error
		taskUUID, err = m.queueTask(ctx, newTaskFunc)

		return err
	})

	return taskUUID, err
}

func (m *Manager) doTaskWrap(ctx context.Context, reqStorage logical.Storage, taskFunc func(context.Context, logical.Storage) error, f func(func(ctx context.Context) error) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// initialize on first task
	if m.Storage == nil {
		m.Storage = reqStorage
		if err := m.invalidateStorage(ctx, reqStorage); err != nil {
			return fmt.Errorf("unable to invalidate storage: %s", err)
		}
	}

	config, err := getConfiguration(ctx, reqStorage)
	if err != nil {
		return fmt.Errorf("unable to get tasks manager configuration: %s", err)
	}

	var taskTimeoutDuration time.Duration
	if config != nil {
		taskTimeoutDuration = config.TaskTimeout
	} else {
		taskTimeoutDuration = defaultTaskTimeoutDuration
	}

	workerTaskFunc := func(ctx context.Context) error {
		ctxWithTimeout, ctxCancelFunc := context.WithTimeout(ctx, taskTimeoutDuration)
		defer ctxCancelFunc()

		if err := taskFunc(ctxWithTimeout, m.Storage); err != nil {
			hclog.L().Debug(fmt.Sprintf("task failed: %s", err))
			return err
		}

		hclog.L().Debug(fmt.Sprintf("task succeeded"))
		return nil
	}

	return f(workerTaskFunc)
}

func (m *Manager) invalidateStorage(ctx context.Context, reqStorage logical.Storage) error {
	invalidateTask := func(task *Task) error {
		task.Status = taskStatusFailed
		task.Reason = taskReasonInvalidatedTask
		task.Modified = time.Now()

		if err := putTaskIntoStorage(ctx, reqStorage, task); err != nil {
			return fmt.Errorf("unable to put task %q into the storage: %q", task.UUID, err)
		}

		return nil
	}

	// invalidate current running task
	{
		currentTaskUUID, err := getCurrentTaskUUIDFromStorage(ctx, reqStorage)
		if err != nil {
			return fmt.Errorf("unable to get current running task uuid from storage: %s", err)
		}

		if currentTaskUUID != "" {
			runningTask, err := getTaskFromStorage(ctx, reqStorage, currentTaskUUID)
			if err != nil {
				return fmt.Errorf("unable to get task %q from storage: %s", currentTaskUUID, err)
			}

			if runningTask != nil {
				if err := invalidateTask(runningTask); err != nil {
					return fmt.Errorf("unable to invalidate task %q: %s", currentTaskUUID, err)
				}
			}

			if err := m.Storage.Delete(ctx, storageKeyCurrentRunningTask); err != nil {
				return fmt.Errorf("unable to delete %q from storage: %q", storageKeyCurrentRunningTask, err)
			}
		}
	}

	// invalidate queued tasks
	{
		queuedTasksUUID, err := reqStorage.List(ctx, storageKeyPrefixQueuedTask)
		if err != nil {
			return fmt.Errorf("unable to get queued tasks from storage: %s", err)
		}

		for _, uuid := range queuedTasksUUID {
			queuedTask, err := getQueuedTaskFromStorage(ctx, m.Storage, uuid)
			if err != nil {
				return fmt.Errorf("unable to get queued task %q from storage: %s", uuid, err)
			}

			if err := invalidateTask(queuedTask); err != nil {
				return fmt.Errorf("unable to invalidate task %q: %s", uuid, err)
			}

			if err := m.Storage.Delete(ctx, queuedTaskStorageKey(uuid)); err != nil {
				return fmt.Errorf("unable to delete %q from storage: %q", storageKeyCurrentRunningTask, err)
			}
		}
	}

	return nil
}

func (m *Manager) queueTask(ctx context.Context, workerTaskFunc func(context.Context) error) (string, error) {
	task := newTask()
	if err := putQueuedTaskIntoStorage(ctx, m.Storage, task); err != nil {
		return "", fmt.Errorf("unable to put queued task %q into storage: %s", task.UUID, err)
	}

	m.taskChan <- &worker.Task{Context: ctx, UUID: task.UUID, Action: workerTaskFunc}

	return task.UUID, nil
}

func (m *Manager) isBusy(ctx context.Context, reqStorage logical.Storage) (bool, error) {
	// busy if there is task in progress
	{
		currentTaskUUID, err := getCurrentTaskUUIDFromStorage(ctx, reqStorage)
		if err != nil {
			return false, fmt.Errorf("unable to get current task uuid from storage: %s", err)
		}

		if currentTaskUUID != "" {
			return true, nil
		}
	}

	// busy if there are queued tasks
	{
		queuedTasksUUID, err := reqStorage.List(ctx, storageKeyPrefixQueuedTask)
		if err != nil {
			return false, fmt.Errorf("unable to get queued tasks from storage: %s", err)
		}

		if len(queuedTasksUUID) != 0 {
			return true, nil
		}
	}

	return false, nil
}
