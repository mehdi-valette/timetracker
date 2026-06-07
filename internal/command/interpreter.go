package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/repository"
)

type Interpreter interface {
	Interpret(rawLine string) tea.Cmd
	GetCurrentTask() entity.Tasker
}

type TaskInterpreter struct {
	taskManager manager.TaskManager
	currentTask entity.Tasker
}

var _ Interpreter = &TaskInterpreter{}

type TaskStartedMsg struct {
	task entity.Tasker
}

type TaskStoppedMsg struct {
	task entity.Tasker
}

type TaskCreatedMsg struct {
	task entity.Tasker
}

type TaskListedMsg struct {
	taskList string
}

type TaskDeletedMsg struct {
	task entity.Tasker
}

func CreateTaskInterpreter() (Interpreter, error) {
	conn, connErr := repository.CreateConnection(":memory:")
	conn.InitializeDb()

	if connErr != nil {
		return &TaskInterpreter{}, connErr
	}

	date := entity.CreateDate()

	timeRangeRepo := repository.CreateTimeRangeRepository(conn, date)
	timeRangeManager := manager.CreateTimeRangeManager(timeRangeRepo, date)

	taskRepo := repository.CreateTaskRepository(conn, date)
	taskManager := manager.CreateTaskManager(taskRepo, timeRangeManager, date)

	return &TaskInterpreter{taskManager: taskManager}, nil
}

func (i *TaskInterpreter) Interpret(rawLine string) tea.Cmd {
	line := strings.Trim(rawLine, " ")
	words := strings.Split(line, " ")

	if len(words) == 0 {
		return nil
	}

	cmd := words[0]
	params := words[1:]

	switch cmd {
	case "quit", "exit":
		return tea.Quit
	case "create":
		return i.createTask(params)
	case "list":
		return i.listTasks()
	case "start":
		return i.startTask(params)
	case "stop":
		return i.stopTask(params)
	case "delete":
		return i.deleteTask(params)
	}

	return ErrorCmd(errors.New("command unknown"))
}

func (i *TaskInterpreter) createTask(params []string) tea.Cmd {
	name := strings.Join(params, " ")

	task, taskErr := i.taskManager.Create(name)

	if taskErr != nil {
		return ErrorCmd(fmt.Errorf("Error while creating the task: %w", taskErr))
	}

	return func() tea.Msg { return TaskCreatedMsg{task: task} }
}

func (i *TaskInterpreter) listTasks() tea.Cmd {
	tasks, listErr := i.taskManager.List()

	if listErr != nil {
		return ErrorCmd(fmt.Errorf("Error while listing the tasks: %w", listErr))
	}

	if len(tasks) == 0 {
		return func() tea.Msg { return TaskListedMsg{taskList: "no tasks"} }
	}

	info := ""

	for _, task := range tasks {
		duration, _ := task.Duration()

		info += fmt.Sprintf("(%d | %s) %s \n", task.GetId(), duration.ToString(), task.GetName())
	}

	return func() tea.Msg { return TaskListedMsg{taskList: info} }
}

func (i *TaskInterpreter) GetCurrentTask() entity.Tasker {
	return i.currentTask
}

func (i *TaskInterpreter) startTask(params []string) tea.Cmd {
	if len(params) > 1 {
		return ErrorCmd(errors.New("expected a single parameter"))
	}

	if len(params) == 1 {
		id, convErr := strconv.ParseInt(params[0], 10, 8)
		taskId := entity.DbId(id)

		if convErr != nil {
			return ErrorCmd(convErr)
		}

		task, getErr := i.taskManager.Get(taskId)

		if getErr != nil {
			return ErrorCmd(getErr)
		}

		i.currentTask = task
	}

	if i.currentTask == nil {
		return ErrorCmd(errors.New("no current task, please give an ID"))
	}

	task, startErr := i.taskManager.Start(i.currentTask.GetId())

	if startErr != nil {
		return ErrorCmd(startErr)
	}

	i.currentTask = task

	return func() tea.Msg { return TaskStartedMsg{task: i.currentTask} }
}

func (i *TaskInterpreter) stopTask(params []string) tea.Cmd {
	if len(params) > 1 {
		return ErrorCmd(errors.New("expected a single parameter"))
	}

	if len(params) == 1 {
		id, convErr := strconv.ParseInt(params[0], 10, 8)
		taskId := entity.DbId(id)

		if convErr != nil {
			return ErrorCmd(convErr)
		}

		task, getErr := i.taskManager.Get(taskId)

		if getErr != nil {
			return ErrorCmd(getErr)
		}

		i.currentTask = task
	}

	if i.currentTask == nil {
		return ErrorCmd(errors.New("no current task, please give an ID"))
	}

	task, stopErr := i.taskManager.Stop(i.currentTask.GetId())

	if stopErr != nil {
		return ErrorCmd(stopErr)
	}

	i.currentTask = task

	return func() tea.Msg { return TaskStoppedMsg{task: task} }
}

func (i *TaskInterpreter) deleteTask(params []string) tea.Cmd {
	if len(params) > 1 {
		return ErrorCmd(errors.New("expected a single parameter"))
	}

	var taskId entity.DbId
	var task entity.Tasker

	if len(params) == 1 {
		id, convErr := strconv.ParseInt(params[0], 10, 8)
		taskId = entity.DbId(id)

		if convErr != nil {
			return ErrorCmd(convErr)
		}

		var getErr error
		task, getErr = i.taskManager.Get(taskId)

		if getErr != nil {
			return ErrorCmd(getErr)
		}
	} else if i.currentTask != nil {
		taskId = i.currentTask.GetId()
		task = i.currentTask
	} else {
		return ErrorCmd(errors.New("no current task, please give an ID"))
	}

	deleteErr := i.taskManager.Delete(taskId)

	if deleteErr != nil {
		return ErrorCmd(deleteErr)
	}

	i.currentTask = nil

	return func() tea.Msg { return TaskDeletedMsg{task: task} }
}
