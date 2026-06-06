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
}

type TaskInterpreter struct {
	taskManager manager.TaskManager
}

var _ Interpreter = &TaskInterpreter{}

type ErrorMsg struct {
	error error
}

type TaskCreatedMsg struct {
	task entity.Tasker
}

type TaskListedMsg struct {
	taskList string
}

func CreateTaskInterpreter() (Interpreter, error) {
	conn, connErr := repository.CreateConnection(":memory:")
	conn.InitializeDb()

	if connErr != nil {
		return &TaskInterpreter{}, connErr
	}

	date := entity.CreateDate()
	taskRepo := repository.CreateTaskRepository(conn, date)

	taskManager := manager.CreateTaskManager(taskRepo, nil, date)

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
	}

	return func() tea.Msg { return ErrorMsg{error: errors.New("command unknown")} }
}

func (i *TaskInterpreter) createTask(params []string) tea.Cmd {
	name := strings.Join(params, " ")

	task, taskErr := i.taskManager.Create(name)

	if taskErr != nil {
		return func() tea.Msg {
			return ErrorMsg{error: fmt.Errorf("Error while creating the task: %w", taskErr)}
		}
	}

	return func() tea.Msg { return TaskCreatedMsg{task: task} }
}

func (i *TaskInterpreter) listTasks() tea.Cmd {
	tasks, listErr := i.taskManager.List()

	if listErr != nil {
		return func() tea.Msg { return ErrorMsg{fmt.Errorf("Error while listing the tasks: %w", listErr)} }
	}

	if len(tasks) == 0 {
		return func() tea.Msg { return TaskListedMsg{taskList: "no tasks"} }
	}

	info := ""

	for _, task := range tasks {
		info = info + strconv.FormatInt(int64(task.GetId()), 10) + " " + string(task.GetName()) + "\n"
	}

	return func() tea.Msg { return TaskListedMsg{taskList: info} }
}
