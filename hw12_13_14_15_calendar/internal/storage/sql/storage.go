package sqlstorage

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // nolint: gci
	"github.com/jmoiron/sqlx"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type sqlError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *sqlError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *sqlError) Unwrap() error {
	return e.Err
}

type Storage struct {
	ctx context.Context
	db  *sqlx.DB
}

func New(ctx context.Context, user, pass, addr, dbName string) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?ssmode=disable", user, pass, addr, dbName)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, &sqlError{Message: "can't create db store", Err: err}
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, &sqlError{Message: "ping error", Err: err}
	}

	return &Storage{ctx: ctx, db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) NewEvent(e storage.Event) error {
	_, err := s.db.ExecContext(
		s.ctx,
		`INSERT event 
    		SET id=?, 
    		    title=?,
    		    start_date=FROM_UNIXTIME(?), 
    		    end_date=FROM_UNIXTIME(?), 
    		    description=?, 
    		    owner_id=?, 
    		    remind_in=?`,
		e.ID,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		e.OwnerID,
		e.RemindIn,
	)
	if err != nil {
		return &sqlError{Message: "can't add event to db", Err: err}
	}
	return nil
}

func (s *Storage) UpdateEvent(e storage.Event) error {
	_, err := s.db.ExecContext(
		s.ctx,
		`UPDATE event
			SET title=?,
    		    start_date=FROM_UNIXTIME(?), 
    		    end_date=FROM_UNIXTIME(?), 
    		    description=?, 
    		    owner_id=?, 
    		    remind_in=?
			WHERE id=?`,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		e.OwnerID,
		e.RemindIn,
		e.ID,
	)

	if err != nil {
		return &sqlError{Message: "can't update event", Err: err}
	}
	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	_, err := s.db.ExecContext(s.ctx, "DELETE FROM event WHERE id=$1", id)
	if err != nil {
		return &sqlError{Message: "can't delete event from db", Err: err}
	}
	return nil
}

func (s *Storage) EventList(from int64, to int64) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.SelectContext(
		s.ctx,
		&events,
		`SELECT id, 
       			title, 
       			UNIX_TIMESTAMP(start_date), 
    		    UNIX_TIMESTAMP(end_date), 
    		    description, 
    		    owner_id, 
    		    remind_in
			FROM event
			WHERE UNIX_TIMESTAMP(start_date) >=$1 AND UNIX_TIMESTAMP(start_date) <=$2`,
		from, to,
	)
	if err != nil {
		return nil, &sqlError{Message: "can't select events from db", Err: err}
	}
	return events, nil
}
