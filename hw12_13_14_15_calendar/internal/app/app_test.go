package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type mockLogger struct {
}

func (m *mockLogger) Info(msg string, fields ...zap.Field) {
}

func (m *mockLogger) Warn(msg string, fields ...zap.Field) {
}

func (m *mockLogger) Error(msg string, fields ...zap.Field) {
}

func (m *mockLogger) String(key string, val string) zap.Field {
	return zap.Field{}
}

func (m *mockLogger) Int64(key string, val int64) zap.Field {
	return zap.Field{}
}

func (m *mockLogger) Duration(key string, val time.Duration) zap.Field {
	return zap.Field{}
}

type AppSuite struct {
	suite.Suite
	mockCtl   *gomock.Controller
	mockStore *MockStorage
	app       *app.App
}

func (s *AppSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockStore = NewMockStorage(s.mockCtl)
	s.app = app.New(&mockLogger{}, s.mockStore)
}

func (s *AppSuite) TearDownTest() {
	s.mockCtl.Finish()
}

func (s *AppSuite) TestCreateEventSuccess() {
	event := app.Event{}
	ctx := context.Background()

	s.mockStore.EXPECT().NewEvent(ctx, event).Return(nil)
	err := s.app.CreateEvent(ctx, event)

	s.Require().NoError(err)
}

func (s *AppSuite) TestCreateEventFail() {
	event := app.Event{}
	sErr := errors.New("store_error")
	ctx := context.Background()

	s.mockStore.EXPECT().NewEvent(ctx, event).Return(sErr)
	err := s.app.CreateEvent(ctx, event)

	s.Require().Error(err)
	ok := errors.Is(err, sErr)
	s.Require().True(ok)
}

func (s *AppSuite) TestUpdateEventSuccess() {
	event := app.Event{}
	ctx := context.Background()

	s.mockStore.EXPECT().UpdateEvent(ctx, event).Return(nil)
	err := s.app.UpdateEvent(ctx, event)

	s.Require().NoError(err)
}

func (s *AppSuite) TestUpdateEventFail() {
	event := app.Event{}
	ctx := context.Background()
	sErr := errors.New("store_error")

	s.mockStore.EXPECT().UpdateEvent(ctx, event).Return(sErr)
	err := s.app.UpdateEvent(ctx, event)

	s.Require().Error(err)
	ok := errors.Is(err, sErr)
	s.Require().True(ok)
}

func (s *AppSuite) TestRemoveEventSuccess() {
	eventID := "unique_event_id"
	ctx := context.Background()

	s.mockStore.EXPECT().RemoveEvent(ctx, eventID).Return(nil)
	err := s.app.RemoveEvent(ctx, eventID)

	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveEventFail() {
	eventID := "unique_event_id"
	ctx := context.Background()
	sErr := errors.New("store_error")

	s.mockStore.EXPECT().RemoveEvent(ctx, eventID).Return(sErr)
	err := s.app.RemoveEvent(ctx, eventID)

	s.Require().Error(err)
	ok := errors.Is(err, sErr)
	s.Require().True(ok)
}

func (s *AppSuite) TestEventsQuerySuccess() {
	var from int64 = 0
	var to int64 = 1
	events := mockEvents()

	ctx := context.Background()

	s.mockStore.EXPECT().EventListFilterByStartDate(ctx, from, to).Return(events, nil)
	evs, err := s.app.Events(ctx, from, to)

	s.Require().NoError(err)
	s.Require().NotNil(evs)
	s.Require().Equal(events, evs)
}

func (s *AppSuite) TestEventsQueryFail() {
	var from int64 = 0
	var to int64 = 1
	sErr := errors.New("store_error")
	ctx := context.Background()

	s.mockStore.EXPECT().EventListFilterByStartDate(ctx, from, to).Return(nil, sErr)
	evs, err := s.app.Events(ctx, from, to)

	s.Require().Error(err)
	s.Require().Nil(evs)
}

func mockEvents() []app.Event {
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

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
