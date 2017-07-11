package eats

import (
	"errors"
	"path"
	"strconv"

	"go.uber.org/cadence"
	s "go.uber.org/cadence/.gen/go/shared"
)

type (
	// TransformFunc type defining the signature of transform function.
	transformFunc func(event *s.HistoryEvent, tasks *TaskGroup) error

	// TaskGroupExecution implements object to transform a workflow history into a TaskGroup.
	TaskGroupExecution struct {
		client       cadence.Client
		transformers map[s.EventType]transformFunc
	}
)

// NewTaskGroupExecution returns a new instanc of TaskGroupExecution.
func NewTaskGroupExecution(c cadence.Client) *TaskGroupExecution {
	obj := &TaskGroupExecution{
		client:       c,
		transformers: make(map[s.EventType]transformFunc),
	}

	obj.transformers[s.EventType_WorkflowExecutionStarted] = obj.tfWorkflowExecutionStarted
	obj.transformers[s.EventType_WorkflowExecutionCompleted] = obj.tfWorkflowExecutionCompleted
	obj.transformers[s.EventType_WorkflowExecutionFailed] = obj.tfWorkflowExecutionFailed

	obj.transformers[s.EventType_ActivityTaskScheduled] = obj.tfActivityTaskScheduled
	obj.transformers[s.EventType_ActivityTaskStarted] = obj.tfActivityTaskStarted
	obj.transformers[s.EventType_ActivityTaskCompleted] = obj.tfActivityTaskCompleted
	obj.transformers[s.EventType_ActivityTaskFailed] = obj.tfActivityTaskFailed
	obj.transformers[s.EventType_ActivityTaskTimedOut] = obj.tfActivityTaskTimedOut

	obj.transformers[s.EventType_StartChildWorkflowExecutionInitiated] = obj.tfStartChildWorkflowExecutionInitiated
	obj.transformers[s.EventType_ChildWorkflowExecutionStarted] = obj.tfChildWorkflowExecutionStarted
	obj.transformers[s.EventType_ChildWorkflowExecutionCompleted] = obj.tfChildWorkflowExecutionCompleted
	obj.transformers[s.EventType_ChildWorkflowExecutionFailed] = obj.tfChildWorkflowExecutionFailed
	obj.transformers[s.EventType_ChildWorkflowExecutionTimedOut] = obj.tfChildWorkflowExecutionTimedOut

	obj.transformers[s.EventType_TimerStarted] = obj.tfTimerStarted
	obj.transformers[s.EventType_TimerFired] = obj.tfTimerFired
	obj.transformers[s.EventType_TimerCanceled] = obj.tfTimerCanceled
	return obj
}

// Transform converts a workflow execution history into a TaskGroup structure.
func (h *TaskGroupExecution) Transform(workflowID string, runID string) (*TaskGroup, error) {
	tasks := &TaskGroup{
		ID:      workflowID,
		RunID:   runID,
		Tasks:   make([]*Task, 0),
		TaskMap: make(map[int64]*Task),
	}

	history, err := h.client.GetWorkflowHistory(workflowID, runID)
	if err != nil {
		return nil, err
	}
	tasks.History = history

	for _, event := range history.Events {
		transFunc, found := h.transformers[*event.EventType]
		if !found {
			continue
		}

		err := transFunc(event, tasks)
		if err != nil {
			return nil, err
		}
	}

	tasks.TaskMap = nil
	return tasks, nil
}

func (h *TaskGroupExecution) tfActivityTaskScheduled(event *s.HistoryEvent, tasks *TaskGroup) error {
	name := event.ActivityTaskScheduledEventAttributes.ActivityType.Name
	return h.createTask(event, name, tasks)
}

func (h *TaskGroupExecution) tfActivityTaskStarted(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ActivityTaskStartedEventAttributes.ScheduledEventId
	return h.setTaskStatus(tasks, id, "r")
}

func (h *TaskGroupExecution) tfActivityTaskCompleted(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ActivityTaskCompletedEventAttributes.ScheduledEventId
	return h.setTaskStatus(tasks, id, "c")
}

func (h *TaskGroupExecution) tfActivityTaskFailed(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ActivityTaskFailedEventAttributes.ScheduledEventId
	return h.setTaskStatus(tasks, id, "f")
}

func (h *TaskGroupExecution) tfActivityTaskTimedOut(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ActivityTaskTimedOutEventAttributes.ScheduledEventId
	return h.setTaskStatus(tasks, id, "t")
}

func (h *TaskGroupExecution) tfStartChildWorkflowExecutionInitiated(event *s.HistoryEvent, tasks *TaskGroup) error {
	name := event.StartChildWorkflowExecutionInitiatedEventAttributes.WorkflowType.Name
	return h.createTask(event, name, tasks)
}

func (h *TaskGroupExecution) tfChildWorkflowExecutionStarted(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ChildWorkflowExecutionStartedEventAttributes.InitiatedEventId
	task, found := tasks.TaskMap[id]
	if !found {
		return errors.New("Could not find ActivityTaskScheduled event: " + strconv.FormatInt(id, 10))
	}

	execution := event.ChildWorkflowExecutionStartedEventAttributes.WorkflowExecution
	taskGroup, err := h.Transform(*execution.WorkflowId, *execution.RunId)
	if err != nil {
		return err
	}

	task.SubTasks = taskGroup.Tasks
	task.Status = "r"
	return nil
}

func (h *TaskGroupExecution) tfChildWorkflowExecutionCompleted(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ChildWorkflowExecutionCompletedEventAttributes.InitiatedEventId
	return h.setTaskStatus(tasks, id, "c")
}

func (h *TaskGroupExecution) tfChildWorkflowExecutionFailed(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ChildWorkflowExecutionFailedEventAttributes.InitiatedEventId
	return h.setTaskStatus(tasks, id, "f")
}

func (h *TaskGroupExecution) tfChildWorkflowExecutionTimedOut(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.ChildWorkflowExecutionTimedOutEventAttributes.InitiatedEventId
	return h.setTaskStatus(tasks, id, "t")
}

func (h *TaskGroupExecution) tfTimerStarted(event *s.HistoryEvent, tasks *TaskGroup) error {
	name := "timer.WaitForDeadline"
	h.createTask(event, &name, tasks)
	return h.setTaskStatus(tasks, *event.EventId, "r")
}

func (h *TaskGroupExecution) tfTimerCanceled(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.TimerCanceledEventAttributes.StartedEventId
	return h.setTaskStatus(tasks, id, "ca")
}

func (h *TaskGroupExecution) tfTimerFired(event *s.HistoryEvent, tasks *TaskGroup) error {
	id := *event.TimerFiredEventAttributes.StartedEventId
	return h.setTaskStatus(tasks, id, "c")
}

func (h *TaskGroupExecution) tfWorkflowExecutionStarted(event *s.HistoryEvent, tasks *TaskGroup) error {
	tasks.Status = "r"
	return nil
}

func (h *TaskGroupExecution) tfWorkflowExecutionCompleted(event *s.HistoryEvent, tasks *TaskGroup) error {
	tasks.Status = "c"
	return nil
}

func (h *TaskGroupExecution) tfWorkflowExecutionFailed(event *s.HistoryEvent, tasks *TaskGroup) error {
	tasks.Status = "f"
	return nil
}

func (h *TaskGroupExecution) createTask(event *s.HistoryEvent, name *string, tasks *TaskGroup) error {
	task := &Task{
		ID:     *event.EventId,
		Name:   path.Ext(*name)[1:],
		Status: "s",
	}

	tasks.TaskMap[task.ID] = task
	tasks.Tasks = append(tasks.Tasks, task)

	return nil
}

func (h *TaskGroupExecution) setTaskStatus(tasks *TaskGroup, id int64, status TaskStatus) error {
	task, found := tasks.TaskMap[id]
	if !found {
		return errors.New("Could not find ActivityTaskScheduled event: " + strconv.FormatInt(id, 10))
	}

	task.Status = status
	return nil
}
