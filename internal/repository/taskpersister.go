package repository

import (
	"database/sql"
	"errors"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
)

var TaskNotFoundErr = errors.New("task not found")

func CreateTaskRepository(conn DbConnector, date entity.Dater) manager.TaskPersister {
	return TaskRepository{
		conn: conn,
		date: date,
	}
}

type TaskRepository struct {
	conn DbConnector
	date entity.Dater
}

var _ manager.TaskPersister = TaskRepository{}

// Create implements [manager.TaskPersister].
func (t TaskRepository) Create() (entity.DbId, error) {
	execResult, execErr := t.conn.Exec(`INSERT INTO "task" ("name") VALUES(NULL);`)

	if execErr != nil {
		return entity.DbId(0), execErr
	}

	taskId, taskErr := execResult.LastInsertId()

	return entity.DbId(taskId), taskErr
}

// Delete implements [manager.TaskPersister].
func (t TaskRepository) Delete(taskId entity.DbId) error {
	deleteResult, deleteErr := t.conn.Exec(`DELETE FROM "task" WHERE "id" = ?`, taskId)

	if deleteErr != nil {
		return deleteErr
	}

	rowsAffected, rowsErr := deleteResult.RowsAffected()

	if rowsErr != nil {
		return deleteErr
	}

	if rowsAffected != 1 {
		return TaskNotFoundErr
	}

	return nil
}

// Get implements [manager.TaskPersister].
func (t TaskRepository) Get(taskId entity.DbId) (entity.Tasker, error) {
	queryResult := t.conn.QueryOne(`SELECT "id", "name" FROM "task" WHERE "id" = ?`, taskId)

	parsedResult := struct {
		id   int64
		name *string
	}{}

	getErr := queryResult.Scan(&parsedResult.id, &parsedResult.name)

	if errors.Is(getErr, sql.ErrNoRows) {
		return &entity.Task{}, TaskNotFoundErr
	} else if getErr != nil {
		return &entity.Task{}, getErr
	}

	name := ""
	if parsedResult.name != nil {
		name = *parsedResult.name
	}

	return entity.CreateTask(entity.DbId(parsedResult.id), name, t.date), nil
}

// Save implements [manager.TaskPersister].
func (t TaskRepository) Save(task entity.Tasker) error {
	updateResult, updateErr := t.conn.Exec(`UPDATE "task" SET "name" = ? WHERE "id" = ?`, task.GetName(), task.GetId())
	if updateErr != nil {
		return updateErr
	}

	rowsCount, rowsErr := updateResult.RowsAffected()

	if rowsErr != nil {
		return rowsErr
	}

	if rowsCount != 1 {
		return TaskNotFoundErr
	}

	return nil
}
