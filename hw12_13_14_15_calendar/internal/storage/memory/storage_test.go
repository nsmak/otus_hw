package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/suite"
)

type MemStoreSuite struct {
	suite.Suite
	store *Storage
}

func (m *MemStoreSuite) SetupTest() {
	m.store = New()
	m.store.events = map[string]*storage.Event{
		"1": &storage.Event{
			ID:          "1",
			Title:       "Title1",
			StartDate:   1,
			EndDate:     2,
			Description: "Description1",
			OwnerID:     "",
			RemindIn:    5,
		},
		"2": &storage.Event{
			ID:          "2",
			Title:       "Title2",
			StartDate:   5,
			EndDate:     25,
			Description: "Description2",
			OwnerID:     "",
			RemindIn:    0,
		},
		"3": &storage.Event{
			ID:          "3",
			Title:       "Title3",
			StartDate:   6,
			EndDate:     18,
			Description: "Description3",
			OwnerID:     "",
			RemindIn:    0,
		},
		"4": &storage.Event{
			ID:          "4",
			Title:       "Title4",
			StartDate:   10,
			EndDate:     12,
			Description: "Description4",
			OwnerID:     "",
			RemindIn:    0,
		},
		"5": &storage.Event{
			ID:          "5",
			Title:       "Title5",
			StartDate:   15,
			EndDate:     20,
			Description: "Description5",
			OwnerID:     "",
			RemindIn:    0,
		},
	}
}

func (m *MemStoreSuite) TestInsertNewEventSuccess() {
	newEvent := storage.Event{
		ID:          "6",
		Title:       "Title6",
		StartDate:   100500,
		EndDate:     200200,
		Description: "Description6",
		OwnerID:     "",
		RemindIn:    150,
	}

	err := m.store.NewEvent(context.Background(), newEvent)

	m.Require().NoError(err)

	saved := m.store.events["6"]

	m.Require().NotNil(saved)
	m.Require().Equal(newEvent.Title, saved.Title)
}

func (m *MemStoreSuite) TestInsertNewEventWithFail() {
	newEvent := storage.Event{
		ID:          "1",
		Title:       "Title6",
		StartDate:   100500,
		EndDate:     200200,
		Description: "Description6",
		OwnerID:     "",
		RemindIn:    150,
	}

	err := m.store.NewEvent(context.Background(), newEvent)

	m.Require().Error(err)
	m.Require().EqualError(ErrEventAlreadyExist, err.Error())
}

func (m *MemStoreSuite) TestUpdateEventSuccess() {
	toUpdate := storage.Event{
		ID:          "1",
		Title:       "TitleUpdated",
		StartDate:   1,
		EndDate:     2,
		Description: "DescriptionUpdate",
		OwnerID:     "",
		RemindIn:    5,
	}
	err := m.store.UpdateEvent(context.Background(), toUpdate)

	m.Require().NoError(err)

	updated := m.store.events["1"]

	m.Require().NotNil(updated)
	m.Require().Equal(toUpdate.Title, updated.Title)
	m.Require().Equal(toUpdate.Description, updated.Description)
}

func (m *MemStoreSuite) TestUpdateEventWithError() {
	toUpdate := storage.Event{
		ID:          "6",
		Title:       "Title6",
		StartDate:   100500,
		EndDate:     200200,
		Description: "Description6",
		OwnerID:     "",
		RemindIn:    150,
	}

	err := m.store.UpdateEvent(context.Background(), toUpdate)

	m.Require().Error(err)
	m.Require().EqualError(ErrEventDoesNotExist, err.Error())
}

func (m *MemStoreSuite) TestRemoveEventSuccess() {
	err := m.store.RemoveEvent(context.Background(), "1")

	m.Require().NoError(err)

	deleted := m.store.events["1"]

	m.Require().Nil(deleted)
}

func (m *MemStoreSuite) TestRemoveEventWithError() {
	err := m.store.RemoveEvent(context.Background(), "NaN")

	m.Require().Error(err)
	m.Require().EqualError(ErrEventDoesNotExist, err.Error())
}

func (m *MemStoreSuite) TestEventListSuccess() {
	list, err := m.store.EventList(context.Background(), 3, 10)

	m.Require().NoError(err)
	m.Require().Len(list, 3)
	m.Require().Contains(list, *m.store.events["2"])
	m.Require().Contains(list, *m.store.events["3"])
	m.Require().Contains(list, *m.store.events["4"])
}

func (m *MemStoreSuite) TestEventListWithError() {
	list, err := m.store.EventList(context.Background(), 500, 700)

	m.Require().Error(err)
	m.Require().EqualError(ErrNoEvents, err.Error())
	m.Require().Nil(list)
}

func (m *MemStoreSuite) TestAsyncOperations() {
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		for i := 10; i < 20; i++ {
			newEvent := storage.Event{
				ID:    fmt.Sprint(i),
				Title: fmt.Sprintf("Title%d", i),
			}

			err := m.store.NewEvent(context.Background(), newEvent)

			m.Require().NoError(err)

			time.Sleep(30 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		toUpdate := storage.Event{
			ID:          "1",
			Title:       "TitleUpdate",
			StartDate:   1,
			EndDate:     2,
			Description: "DescriptionUpdate",
			OwnerID:     "",
			RemindIn:    5,
		}

		err := m.store.UpdateEvent(context.Background(), toUpdate)

		m.Require().NoError(err)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(150 * time.Millisecond)
		list, err := m.store.EventList(context.Background(), 6, 10)

		m.Require().NoError(err)
		m.Require().Len(list, 2)
	}()

	go func() {
		defer wg.Done()
		err := m.store.RemoveEvent(context.Background(), "5")
		m.Require().NoError(err)
	}()

	wg.Wait()
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(MemStoreSuite))
}
