package command

import (
	tea "charm.land/bubbletea/v2"
	"github.com/mehdi-valette/timetracker/internal/command/alltasks"
	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/repository"
)

func Run(databasePath string) (returnModel tea.Model, returnErr error) {
	taskManager, timeRangeManager, managerErr := createManagers((databasePath))

	if managerErr != nil {
		return alltasks.AllTasksModel{}, managerErr
	}

	program := tea.NewProgram(
		alltasks.CreateAllTasksModel(taskManager, timeRangeManager),
	)

	return program.Run()
}

func createManagers(databasePath string) (manager.TaskManager, manager.TimeRangeManager, error) {
	conn, connErr := repository.CreateConnection(databasePath)
	conn.InitializeDb()

	if connErr != nil {
		return &manager.TaskManagement{}, &manager.TimeRangeManagement{}, connErr
	}

	date := entity.CreateDate()

	timeRangeRepo := repository.CreateTimeRangeRepository(conn, date)
	timeRangeManager := manager.CreateTimeRangeManager(timeRangeRepo, date)

	taskRepo := repository.CreateTaskRepository(conn, date)
	taskManager := manager.CreateTaskManager(taskRepo, timeRangeManager, date)

	return taskManager, timeRangeManager, nil
}
