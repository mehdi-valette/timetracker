package repository

import (
	"errors"
	"testing"

	"github.com/mehdi-valette/timetracker/internal/test"
)

/*
The database must allow creating tasks and time ranges.
When a task is deleted, the associated time ranges must be deleted too.
*/
func TestConnectionInitializeDb(t *testing.T) {
	conn, connErr := CreateConnection(":memory:")

	if connErr != nil {
		t.Error(test.NoError(connErr))
	}

	if err := conn.InitializeDb(); err != nil {
		t.Error(test.NoError(err))
	}

	// create data (one task assigned to one time range)
	taskResult, taskErr := conn.Exec(`INSERT INTO "task" ("short_name", "name") VALUES ('abc', 'mytask')`)

	if taskErr != nil {
		t.Error(test.NoError(taskErr))
	}

	taskId, _ := taskResult.LastInsertId()

	_, timeRangeErr := conn.Exec(`INSERT INTO "time_range" ("task_fk") VALUES (?)`, taskId)

	if timeRangeErr != nil {
		t.Error(test.NoError(timeRangeErr))
	}

	// there should be one time range initially
	countBeforeDelete := conn.QueryOne(`SELECT COUNT(*) FROM "time_range"`)

	result := struct{ count int }{}

	if err := countBeforeDelete.Scan(&result.count); err != nil {
		t.Error(test.NoError(err))
	}

	if result.count != 1 {
		t.Errorf("should have one row, got %d", result.count)
	}

	// the time range must be deleted with the task
	if _, err := conn.Exec(`DELETE FROM "task" WHERE "id" = ?`, taskId); err != nil {
		t.Error(test.NoError(err))
	}

	countAfterDelete := conn.QueryOne(`SELECT COUNT(*) FROM "time_range"`)

	if err := countAfterDelete.Scan(&result.count); err != nil {
		t.Error(test.NoError(err))
	}

	if result.count != 0 {
		t.Errorf("should have 0 rows, got %d", result.count)
	}
}

func TestConnectionInitializeDbWithoutConnection(t *testing.T) {
	conn := DbConnection{}
	err := conn.InitializeDb()

	if !errors.Is(err, DbConnectorNoDbErr) {
		t.Error("should return an error")
	}
}
