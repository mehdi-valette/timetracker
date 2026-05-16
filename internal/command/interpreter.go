package command

import (
	"strings"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/repository"
)

type Interpreter interface {
	Interpret(rawLine string)
}

type TaskInterpreter struct {
	taskManager manager.TaskManager
}

var _ Interpreter = &TaskInterpreter{}

func CreateTaskInterpreter() Interpreter {
	conn := repository.DbConnection{}
	conn.Connect(":memory:")

	date := entity.CreateDate()
	taskRepo := repository.CreateTaskRepository(&conn, date)

	taskManager := manager.CreateTaskManager(taskRepo, nil, date)

	return &TaskInterpreter{taskManager: taskManager}
}

func (i *TaskInterpreter) Interpret(rawLine string) {
	line := strings.Trim(rawLine, " ")
	words := strings.Split(line, " ")

	if len(words) == 0 {
		return
	}

	cmd := words[0]
	params := words[1:]

	switch cmd {
	case "create":
		i.createTask(params)
	}
}

func (i *TaskInterpreter) createTask(params []string) {
	name := strings.Join(params, " ")

	i.taskManager.Create(name)
}
