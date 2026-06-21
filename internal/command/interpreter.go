package command

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/repository"
)

var explanation = `## Commands

- *create [short_name] [name]*: create the task *name* with the short name *short_name*
-    *start [short_name]*: start a new time range entry for *short_name*
-     *stop [short_name]*: stop the last time range entry for *short_name*
-    *begin [short_name]*: shortcut for *create* and *start*
-   *delete [short_name]*: delete the task *short_name*
-          *list*: list all tasks and their duration
-  *exit*, *quit*: stop the current task and exit the application

> note: the [short_name] can be omitted when there's a current task

## Shortcuts

- *ctrl+v*, *ctrl+shift+v*: paste from clipboard
-                 *ctrl+l*: clear the line
-                 *ctrl+c*: stop the current task and exit the application
`

type Interpreter interface {
	Interpret(rawLine string) tea.Cmd
	GetCurrentTask() entity.Tasker
	ListTasks() tea.Cmd
}

type TaskInterpreter struct {
	taskManager manager.TaskManager
	currentTask entity.Tasker
}

var _ Interpreter = &TaskInterpreter{}

func CreateTaskInterpreter(databasePath string) (Interpreter, error) {
	conn, connErr := repository.CreateConnection(databasePath)
	conn.InitializeDb()

	if connErr != nil {
		return &TaskInterpreter{}, connErr
	}

	date := entity.CreateDate()

	timeRangeRepo := repository.CreateTimeRangeRepository(conn, date)
	timeRangeManager := manager.CreateTimeRangeManager(timeRangeRepo, date)

	taskRepo := repository.CreateTaskRepository(conn, date)
	taskManager := manager.CreateTaskManager(taskRepo, timeRangeManager, date)

	lastTimeRange, getLastRangeErr := timeRangeManager.GetLastTimeRange()

	taskInterpreter := &TaskInterpreter{taskManager: taskManager}

	if getLastRangeErr == nil {
		task, taskGetErr := taskManager.Get(lastTimeRange.GetTaskId())

		if taskGetErr != nil {
			panic(taskGetErr)
		}

		taskInterpreter.currentTask = task
	}

	return taskInterpreter, nil
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
		i.stopTask([]string{})
		return tea.Quit
	case "create":
		return i.createTask(params)
	case "begin":
		return i.beginTask(params)
	case "list":
		return i.ListTasks()
	case "start":
		return i.startTask(params)
	case "stop":
		return i.stopTask(params)
	case "delete":
		return i.deleteTask(params)
	}

	return func() tea.Msg { return HelpMsg{help: explanation} }
}

func (i *TaskInterpreter) createTask(params []string) tea.Cmd {
	if len(params) < 2 {
		return ErrorCmd(errors.New("give a short-name followed by the task's name"))
	}

	shortName := params[0]
	name := strings.Join(params[1:], " ")

	task, taskErr := i.taskManager.Create(shortName, name)

	if taskErr != nil {
		return ErrorCmd(fmt.Errorf("Error while creating the task: %w", taskErr))
	}

	return func() tea.Msg { return TaskCreatedMsg{task: task} }
}

func (i *TaskInterpreter) ListTasks() tea.Cmd {
	tasks, listErr := i.taskManager.List()

	if listErr != nil {
		return ErrorCmd(fmt.Errorf("Error while listing the tasks: %w", listErr))
	}

	if len(tasks) == 0 {
		return func() tea.Msg { return TaskListedMsg{taskList: "no tasks"} }
	}

	info := new(bytes.Buffer)

	for _, task := range tasks {
		fmt.Fprintf(info, "%s\n", task.String())
	}

	return func() tea.Msg { return TaskListedMsg{taskList: info.String()} }
}

func (i *TaskInterpreter) GetCurrentTask() entity.Tasker {
	return i.currentTask
}

func (i *TaskInterpreter) beginTask(params []string) tea.Cmd {
	cmd := i.createTask(params)

	msg, ok := cmd().(TaskCreatedMsg)

	if !ok {
		return ErrorCmd(errors.New("error beginning the task"))
	}

	return i.startTask([]string{msg.task.GetShortName().String()})
}

func (i *TaskInterpreter) startTask(params []string) tea.Cmd {
	if len(params) > 1 {
		return ErrorCmd(errors.New("expected a single parameter"))
	}

	if len(params) == 1 {
		task, getErr := i.taskManager.GetByShortName(params[0])

		if getErr != nil {
			return ErrorCmd(getErr)
		}

		if i.currentTask != nil {
			if _, stopErr := i.taskManager.Stop(i.currentTask.GetId()); stopErr != nil {
				return ErrorCmd(stopErr)
			}
		}

		i.currentTask = task
	}

	if i.currentTask == nil {
		return ErrorCmd(errors.New("no current task, please give a short name"))
	}

	task, startErr := i.taskManager.Start(i.currentTask.GetId())

	if startErr != nil && startErr != manager.TaskManagerTaskRunningErr {
		return ErrorCmd(startErr)
	}

	if startErr == nil {
		i.currentTask = task
	}

	return func() tea.Msg { return TaskStartedMsg{task: i.currentTask} }
}

func (i *TaskInterpreter) stopTask(params []string) tea.Cmd {
	if len(params) > 1 {
		return ErrorCmd(errors.New("expected a single parameter"))
	}

	currentTask := i.currentTask

	if len(params) == 1 {
		task, getErr := i.taskManager.GetByShortName(params[0])

		if getErr != nil {
			return ErrorCmd(getErr)
		}

		currentTask = task
	}

	if currentTask == nil {
		return ErrorCmd(errors.New("no current task, please give a short name"))
	}

	task, stopErr := i.taskManager.Stop(currentTask.GetId())

	if stopErr != nil {
		return ErrorCmd(stopErr)
	}

	if task.GetId() == i.currentTask.GetId() {
		i.currentTask = task
	}

	return func() tea.Msg { return TaskStoppedMsg{task: task} }
}

func (i *TaskInterpreter) deleteTask(params []string) tea.Cmd {
	if len(params) > 1 {
		return ErrorCmd(errors.New("expected a single parameter"))
	}

	var taskId entity.DbId
	var task entity.Tasker

	if len(params) == 1 {
		var getErr error
		task, getErr = i.taskManager.GetByShortName(params[0])

		if getErr != nil {
			return ErrorCmd(getErr)
		}

		taskId = task.GetId()
	} else if i.currentTask != nil {
		taskId = i.currentTask.GetId()
		task = i.currentTask
	} else {
		return ErrorCmd(errors.New("no current task, please give a short name"))
	}

	deleteErr := i.taskManager.Delete(taskId)

	if deleteErr != nil {
		return ErrorCmd(deleteErr)
	}

	i.currentTask = nil

	return func() tea.Msg { return TaskDeletedMsg{task: task} }
}
