package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // nolint: gci
	"github.com/jmoiron/sqlx"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type SQLError struct {
	app.BaseError
}

func NewError(msg string, err error) *SQLError {
	return &SQLError{BaseError: app.BaseError{Message: msg, Err: err}}
}

type EventDataStore struct {
	db *sqlx.DB
}

func New(ctx context.Context, user, pass, addr, dbName string) (*EventDataStore, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, addr, dbName)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, NewError("can't create db store", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, NewError("ping error", err)
	}

	return &EventDataStore{db: db}, nil
}

func (s *EventDataStore) Close() error {
	return s.db.Close()
}

func (s *EventDataStore) NewEvent(ctx context.Context, e app.Event) error {
	isExist, err := s.eventIsExist(ctx, e.ID)
	if err != nil {
		return err
	}

	if isExist {
		return storage.ErrEventAlreadyExist
	}

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO event (id, title, start_date, end_date, description,  owner_id,  remind_in) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		e.ID,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		e.OwnerID,
		e.RemindIn,
	)
	if err != nil {
		return NewError("can't add event to db", err)
	}
	return nil
}

func (s *EventDataStore) UpdateEvent(ctx context.Context, e app.Event) error {
	isExist, err := s.eventIsExist(ctx, e.ID)
	if err != nil {
		return err
	}

	if !isExist {
		return storage.ErrEventDoesNotExist
	}

	_, err = s.db.ExecContext(
		ctx,
		`UPDATE event
			SET title=$1,
    		    start_date=$2, 
    		    end_date=$3, 
    		    description=$4, 
    		    owner_id=$5, 
    		    remind_in=$6
			WHERE id=$7`,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		e.OwnerID,
		e.RemindIn,
		e.ID,
	)
	if err != nil {
		return NewError("can't update event", err)
	}
	return nil
}

func (s *EventDataStore) RemoveEvent(ctx context.Context, id string) error {
	isExist, err := s.eventIsExist(ctx, id)
	if err != nil {
		return err
	}

	if !isExist {
		return storage.ErrEventDoesNotExist
	}

	_, err = s.db.ExecContext(ctx, "DELETE FROM event WHERE id=$1", id)
	if err != nil {
		return NewError("can't delete event from db", err)
	}
	return nil
}

func (s *EventDataStore) EventListFilterByStartDate(ctx context.Context, from int64, to int64) ([]app.Event, error) {
	var events []app.Event
	err := s.db.SelectContext(
		ctx,
		&events,
		`SELECT id, 
       			title, 
       			start_date, 
    		    end_date, 
    		    description, 
    		    owner_id, 
    		    remind_in
			FROM event
			WHERE start_date >=$1 AND start_date <=$2`,
		from, to,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoEvents
		}
		return nil, NewError("can't select events from db", err)
	}
	return events, nil
}

func (s *EventDataStore) EventListFilterByReminderIn(ctx context.Context, from int64, to int64) ([]app.Event, error) {
	var events []app.Event
	err := s.db.SelectContext(
		ctx,
		&events,
		`SELECT id, 
       			title, 
       			start_date, 
    		    end_date, 
    		    description, 
    		    owner_id, 
    		    remind_in
			FROM event
			WHERE remind_in >=$1 AND remind_in <=$2`,
		from, to,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoEvents
		}
		return nil, NewError("can't select events from db", err)
	}
	return events, nil
}

func (s *EventDataStore) eventIsExist(ctx context.Context, id string) (bool, error) {
	var count int

	err := s.db.GetContext(
		ctx,
		&count,
		`SELECT COUNT(*)
			FROM event
			WHERE id=$1`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, NewError("can't get event", err)
	}

	return count > 0, nil
}
