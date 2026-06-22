package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
)

var TaskNotFoundErr = errors.New("task not found")

type taskRecord struct {
	id        int64
	shortName *string
	name      *string
}

func (taskRecord taskRecord) recordToTasker(date entity.Dater) entity.Tasker {
	name := ""
	if taskRecord.name != nil {
		name = *taskRecord.name
	}

	shortName := ""
	if taskRecord.shortName != nil {
		shortName = *taskRecord.shortName
	}

	return entity.CreateTask(entity.DbId(taskRecord.id), shortName, name, date)
}

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

// GetByShortName implements [manager.TaskPersister].
func (t TaskRepository) GetByShortName(shortname string) (entity.Tasker, error) {
	queryResult := t.conn.QueryOne(`SELECT "id", "short_name", "name" FROM "task" WHERE "short_name" = ?`, entity.CreateShortName(shortname))

	parsedResult := taskRecord{}

	getErr := queryResult.Scan(&parsedResult.id, &parsedResult.shortName, &parsedResult.name)

	if getErr != nil {
		if strings.Contains(getErr.Error(), "sql: no rows in result set") {
			return nil, TaskNotFoundErr
		}

		return nil, getErr
	}

	return parsedResult.recordToTasker(t.date), nil
}

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
	queryResult := t.conn.QueryOne(`SELECT "id", "short_name", "name" FROM "task" WHERE "id" = ?`, taskId)

	parsedResult := taskRecord{}

	getErr := queryResult.Scan(&parsedResult.id, &parsedResult.shortName, &parsedResult.name)

	if errors.Is(getErr, sql.ErrNoRows) {
		return &entity.Task{}, TaskNotFoundErr
	} else if getErr != nil {
		return &entity.Task{}, getErr
	}

	return parsedResult.recordToTasker(t.date), nil
}

// Save implements [manager.TaskPersister].
func (t TaskRepository) Save(task entity.Tasker) error {
	updateResult, updateErr := t.conn.Exec(`UPDATE "task" SET "short_name" = ?, "name" = ? WHERE "id" = ?`, task.GetShortName(), task.GetName(), task.GetId())
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

func (t TaskRepository) List() ([]entity.Tasker, error) {
	// get the number of entries for future insertion performance
	countResult := t.conn.QueryOne(`SELECT COUNT(*) FROM "task"`)

	count := 0
	if err := countResult.Scan(&count); err != nil {
		return []entity.Tasker{}, err
	}

	taskList := make([]entity.Tasker, 0, count)

	// query all the tasks from the database
	selectResult, selectErr := t.conn.QueryMany(`SELECT "id", "short_name", "name" FROM "task"`)

	if selectErr != nil {
		return []entity.Tasker{}, selectErr
	}

	// transform each row into a task
	parsedResult := taskRecord{}

	for selectResult.Next() {
		if err := selectResult.Scan(&parsedResult.id, &parsedResult.shortName, &parsedResult.name); err != nil {
			return []entity.Tasker{}, err
		}

		taskList = append(
			taskList,
			parsedResult.recordToTasker(t.date),
		)
	}

	return taskList, nil
}
