// +build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
	sqlstorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/suite"
)

var (
	psUsr  = "postgres"
	psPass = "password"
	psAddr = "db:5432"
	psDB   = "postgres"

	restURL = "http://calendar:8888"
)

type IntegrationSuite struct {
	suite.Suite
	db      *sqlx.DB
	storage *sqlstorage.EventDataStore
	events  []app.Event
}

func (s *IntegrationSuite) SetupTest() {
	storage, err := sqlstorage.New(context.Background(), psUsr, psPass, psAddr, psDB)
	if err != nil {
		log.Fatal(err.Error())
	}
	s.storage = storage
	s.db = s.initDB()
	s.events = s.defaultEvents()
	s.saveDefaultEvents()
}

func (s *IntegrationSuite) TearDownTest() {
	s.removeDefaultEvents()
	s.removeNotifications()
	_ = s.db.Close()
	_ = s.storage.Close()
}

func (s *IntegrationSuite) initDB() *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", psUsr, psPass, psAddr, psDB)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func (s *IntegrationSuite) saveDefaultEvents() {
	for _, e := range s.events {
		_ = s.storage.NewEvent(context.Background(), e)
	}
}

func (s *IntegrationSuite) removeDefaultEvents() {
	for _, e := range s.events {
		_ = s.storage.RemoveEvent(context.Background(), e.ID)
	}
}

func (s *IntegrationSuite) removeNotifications() {
	_, _ = s.db.Exec("DELETE FROM notification")
}

func (s *IntegrationSuite) notificationIsExist(eventID string) bool {
	var count int

	_ = s.db.Get(
		&count,
		`SELECT COUNT(*)
			FROM notification
			WHERE id=$1`,
		eventID,
	)

	return count > 0
}

func (s *IntegrationSuite) getEvent(id string) (app.Event, error) {
	var event app.Event
	err := s.db.Get(&event, "SELECT * FROM event WHERE id=$1", id)
	return event, err
}

func (s *IntegrationSuite) defaultEvents() []app.Event {
	return []app.Event{
		{
			ID:          "unique_event_id_1",
			Title:       "Event_Title_1",
			StartDate:   100500,
			EndDate:     300800,
			Description: "Event_Description_1",
			OwnerID:     "unique_owner_uid",
			RemindIn:    0,
		},
		{
			ID:          "unique_event_id_2",
			Title:       "Event_Title_2",
			StartDate:   300800,
			EndDate:     300900,
			Description: "Event_Description_2",
			OwnerID:     "unique_owner_uid",
			RemindIn:    0,
		},
	}
}

func (s *IntegrationSuite) TestCreateEventSuccess() {
	newEvent := app.Event{
		ID:          "unique_event_id_test",
		Title:       "Event_Title_test",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_test",
		OwnerID:     "unique_owner_uid_test",
		RemindIn:    0,
	}
	s.events = append(s.events, newEvent)

	data, err := json.Marshal(&newEvent)

	s.Require().NoError(err)

	resp, err := http.Post(restURL+"/event/create", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	event, err := s.getEvent(newEvent.ID)

	s.Require().NoError(err)
	s.Require().Equal(newEvent, event)
}

func (s *IntegrationSuite) TestCreateEventFail() {
	data, err := json.Marshal(&s.events[0])

	s.Require().NoError(err)

	resp, err := http.Post(restURL+"/event/create", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *IntegrationSuite) TestUpdateEventSuccess() {
	eventToUpdate := s.events[0]
	eventToUpdate.Title = "Updated title"

	data, err := json.Marshal(&eventToUpdate)

	s.Require().NoError(err)

	resp, err := http.Post(restURL+"/event/update", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	event, err := s.getEvent(eventToUpdate.ID)

	s.Require().NoError(err)
	s.Require().Equal(eventToUpdate, event)
}

func (s *IntegrationSuite) TestUpdateEventFail() {
	newEvent := app.Event{
		ID:          "unique_event_id_test",
		Title:       "Event_Title_test",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_test",
		OwnerID:     "unique_owner_uid_test",
		RemindIn:    0,
	}

	data, err := json.Marshal(&newEvent)
	s.Require().NoError(err)

	resp, err := http.Post(restURL+"/event/update", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *IntegrationSuite) TestRemoveEventSuccess() {
	eventID := s.events[0].ID
	reqForm := rest.EventRemoveForm{EventID: eventID}

	data, err := json.Marshal(&reqForm)
	s.Require().NoError(err)

	resp, err := http.Post(restURL+"/event/remove", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	_, err = s.getEvent(eventID)

	ok := errors.Is(err, sql.ErrNoRows)
	s.Require().True(ok)
	s.Require().Error(err)
}

func (s *IntegrationSuite) TestRemoveEventFail() {
	reqForm := rest.EventRemoveForm{EventID: "invalid_id"}

	data, err := json.Marshal(&reqForm)
	s.Require().NoError(err)

	resp, err := http.Post(restURL+"/event/remove", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *IntegrationSuite) TestEventsQuerySuccess() {
	resp, err := http.Get(restURL + "/events?from=100500&to=300700")

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		Error  interface{} `json:"error"`
		Events []app.Event `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	s.Require().NoError(err)
	s.Require().Len(response.Events, 1)
}

func (s *IntegrationSuite) TestEventsQueryFail() {
	resp, err := http.Get(restURL + "/events?from=900500&to=900700")

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *IntegrationSuite) TestNotification() {
	newEvent := app.Event{
		ID:          "unique_event_id_notification",
		Title:       "Event_Title_test",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_test",
		OwnerID:     "unique_owner_uid_test",
		RemindIn:    time.Now().Add(5 * time.Second).Unix(),
	}
	err := s.storage.NewEvent(context.Background(), newEvent)

	s.Require().NoError(err)

	s.events = append(s.events, newEvent)

	time.Sleep(15 * time.Second)

	s.Require().True(s.notificationIsExist(newEvent.ID))
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
