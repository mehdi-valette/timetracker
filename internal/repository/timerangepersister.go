package repository

import (
	"errors"

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

	parse := struct {
		id     entity.DbId
		taskId entity.DbId
		start  *entity.Timestamp
		end    *entity.Timestamp
	}{}

	if err := result.Scan(&parse.id, &parse.taskId, &parse.start, &parse.end); err != nil {
		return &entity.TimeRange{}, err
	}

	return entity.CreateTimeRange(parse.id, parse.taskId, t.date), nil
}

// ListByTaskId implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) ListByTaskId(taskId entity.DbId) []entity.TimeRanger {
	panic("unimplemented")
}

// Save implements [manager.TimeRangePersister].
func (t *TimeRangeRepository) Save(timeRange entity.TimeRanger) error {
	result, execErr := t.conn.Exec(`UPDATE "time_range" SET "start" = ?, "end" = ? WHERE "id" = ?`, timeRange.GetStart().GetSeconds(), timeRange.GetEnd().GetSeconds(), timeRange.GetId())

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
