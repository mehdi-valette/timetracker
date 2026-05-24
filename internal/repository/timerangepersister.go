package repository

import (
	"errors"
	"strings"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
)

var TimeRangeNotFoundErr = errors.New("time range not found")

func CreateTimeRangeRepository(conn DbConnector, date entity.Dater) manager.TimeRangePersister {
	return &TimeRangeRepository{
		conn: conn,
		date: date,
	}
}

type TimeRangeRepository struct {
	conn DbConnector
	date entity.Dater
}

// Create implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) Create(taskId entity.DbId) (entity.DbId, error) {
	result, execErr := t.conn.Exec(`INSERT INTO "time_range" ("task_fk") VALUES(?)`, taskId)

	if execErr != nil {
		return entity.DbId(0), execErr
	}

	lastId, lastErr := result.LastInsertId()

	if lastErr != nil {
		return entity.DbId(0), lastErr
	}

	return entity.DbId(lastId), nil
}

// Delete implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) Delete(id entity.DbId) error {
	panic("unimplemented")
}

// Get implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) Get(timeRangeId entity.DbId) (entity.TimeRanger, error) {
	result := t.conn.QueryOne(`SELECT id, task_fk, start, end FROM "time_range" WHERE "id" = ?`, timeRangeId)

	record := entity.TimeRangeRecord{}

	if err := result.Scan(&record.Id, &record.TaskId, &record.Start, &record.End); err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return &entity.TimeRange{}, TimeRangeNotFoundErr
		}

		return &entity.TimeRange{}, err
	}

	return entity.CreateTimeRangeFromRecord(record, t.date), nil
}

// ListByTaskId implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) ListByTaskId(taskId entity.DbId) []entity.TimeRanger {
	panic("unimplemented")
}

// Save implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) Save(timeRange entity.TimeRanger) error {
	var start *int64 = nil
	var end *int64 = nil

	if timeRange.HasStarted() {
		start = new(timeRange.GetStart().GetSeconds())
	}

	if timeRange.HasEnded() {
		end = new(timeRange.GetEnd().GetSeconds())
	}

	result, execErr := t.conn.Exec(`UPDATE "time_range" SET "start" = ?, "end" = ? WHERE "id" = ?`, start, end, timeRange.GetId())

	if execErr != nil {
		return execErr
	}

	countRows, countErr := result.RowsAffected()

	if countErr != nil {
		return countErr
	}

	if countRows != 1 {
		return TimeRangeNotFoundErr
	}

	return nil
}

var _ manager.TimeRangePersister = &TimeRangeRepository{}
