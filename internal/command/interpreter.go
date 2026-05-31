package command

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/repository"
)

type Interpreter interface {
	Interpret(rawLine string) (textinput.Model, string)
}

type TaskInterpreter struct {
	taskManager manager.TaskManager
}

var _ Interpreter = &TaskInterpreter{}

type TaskCreated struct{}

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

func (i *TaskInterpreter) Interpret(rawLine string) (textinput.Model, string) {
	line := strings.Trim(rawLine, " ")
	words := strings.Split(line, " ")

	if len(words) == 0 {
		return textinput.New(), ""
	}

	cmd := words[0]
	params := words[1:]

	info := "command not found"

	switch cmd {
	case "create":
		info = i.createTask(params)
	case "list":
		info = i.listTasks()
	}

	return textinput.New(), info
}

func (i *TaskInterpreter) createTask(params []string) string {
	name := strings.Join(params, " ")

	task, err := i.taskManager.Create(name)

	if err != nil {
		return "error while creating the task: " + err.Error()
	}

	return "created task \"" + string(task.GetName()) + "\""
}

func (i *TaskInterpreter) listTasks() string {
	tasks, err := i.taskManager.List()

	if err != nil {
		return "Error while listing the tasks: " + err.Error()
	}

	if len(tasks) == 0 {
		return "no tasks"
	}

	info := ""

	for _, task := range tasks {
		info = info + string(task.GetName()) + "\n"
	}

	return info
}
