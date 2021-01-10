package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
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

func TestCreateEventSuccess(t *testing.T) {
	server := testServer(false)
	defer server.Close()

	newEvent := app.Event{
		ID:          "unique_event_id_1",
		Title:       "Event_Title_1",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_1",
		OwnerID:     "unique_owner_uid",
		RemindIn:    0,
	}
	data, err := json.Marshal(&newEvent)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/event/create", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateEventFailStore(t *testing.T) {
	server := testServer(false)
	defer server.Close()

	newEvent := app.Event{
		ID:          "unique_event_id_1",
		Title:       "Event_Title_1",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_1",
		OwnerID:     "unique_owner_uid",
		RemindIn:    0,
	}
	data, err := json.Marshal(&newEvent)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/event/create", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(server.URL+"/event/create", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestCreateEventInvalidData(t *testing.T) {
	server := testServer(false)
	defer server.Close()

	resp, err := http.Post(server.URL+"/event/create", "application/json", bytes.NewReader([]byte{}))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestUpdateEventSuccess(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	event := app.Event{
		ID:          "unique_event_id_1",
		Title:       "Event_Title_1",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_Updated",
		OwnerID:     "unique_owner_uid",
		RemindIn:    0,
	}

	data, err := json.Marshal(&event)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/event/update", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateEventFailStore(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	event := app.Event{
		ID:          "unique_event_id_100",
		Title:       "Event_Title_1",
		StartDate:   100500,
		EndDate:     300800,
		Description: "Event_Description_Updated",
		OwnerID:     "unique_owner_uid",
		RemindIn:    0,
	}

	data, err := json.Marshal(&event)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/event/update", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestUpdateEventInvalidData(t *testing.T) {
	server := testServer(false)
	defer server.Close()

	resp, err := http.Post(server.URL+"/event/update", "application/json", bytes.NewReader([]byte{}))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestRemoveEventSuccess(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	form := EventRemoveForm{
		EventID: "unique_event_id_1",
	}

	data, err := json.Marshal(&form)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/event/remove", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemoveEventFailStore(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	form := EventRemoveForm{
		EventID: "unique_event_id_10",
	}

	data, err := json.Marshal(&form)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/event/remove", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestRemoveEventInvalidData(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	resp, err := http.Post(server.URL+"/event/remove", "application/json", bytes.NewReader([]byte{}))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestEventsSuccess(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	resp, err := http.Get(server.URL + "/events?from=300800&to=500200")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsedResp struct {
		Events []app.Event `json:"data"`
		Error  JSON        `json:"error"`
	}
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Error)
	require.NotNil(t, parsedResp.Events)

	require.Len(t, parsedResp.Events, 1)
	require.Equal(t, parsedResp.Events[0].ID, "unique_event_id_2")
}

func TestEventsFailStore(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	resp, err := http.Get(server.URL + "/events?from=900800&to=10000200")
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func TestEventsInvalidData(t *testing.T) {
	server := testServer(true)
	defer server.Close()

	resp, err := http.Get(server.URL + "/events?dfsfsf=fsfsf")
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var parsedResp Response
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	require.NoError(t, err)
	require.Nil(t, parsedResp.Data)
	require.NotNil(t, parsedResp.Error)
}

func testServer(prepareData bool) *httptest.Server {
	store := memorystorage.New()
	if prepareData {
		for _, e := range mockEvents() {
			_ = store.NewEvent(context.Background(), e)
		}
	}
	a := app.New(&mockLogger{}, store)
	api := NewAPI(a)

	router := mux.NewRouter()
	for _, route := range api.Routes() {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.Func)
	}

	return httptest.NewServer(router)
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
