package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type APIError struct {
	app.BaseError
}

func NewError(msg string, err error) *APIError {
	return &APIError{BaseError: app.BaseError{Message: msg, Err: err}}
}

type EventRemoveForm struct {
	EventID string `json:"id"`
}

type EventsQueryForm struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type API struct {
	application *app.App
}

func NewAPI(application *app.App) *API {
	return &API{application: application}
}

func (a *API) createEvent(w http.ResponseWriter, r *http.Request) {
	var event app.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		sendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.application.CreateEvent(r.Context(), event); err != nil {
		sendErrorJSON(w, r, http.StatusBadRequest, err, "can't create event")
		return
	}

	sendDataJSON(w, r, http.StatusOK, nil)
}

func (a *API) updateEvent(w http.ResponseWriter, r *http.Request) {
	var event app.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		sendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.application.UpdateEvent(r.Context(), event); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrEventDoesNotExist) {
			statusCode = http.StatusNotFound
		}
		sendErrorJSON(w, r, statusCode, err, "can't update event")
		return
	}

	sendDataJSON(w, r, http.StatusOK, nil)
}

func (a *API) removeEvent(w http.ResponseWriter, r *http.Request) {
	var form EventRemoveForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		sendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.application.RemoveEvent(r.Context(), form.EventID); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrEventDoesNotExist) {
			statusCode = http.StatusNotFound
		}
		sendErrorJSON(w, r, statusCode, err, "can't remove event")
		return
	}

	sendDataJSON(w, r, http.StatusOK, nil)
}

func (a *API) events(w http.ResponseWriter, r *http.Request) {
	var query EventsQueryForm
	if err := schema.NewDecoder().Decode(&query, r.URL.Query()); err != nil {
		sendErrorJSON(w, r, http.StatusBadRequest, err, "can't query params")
		return
	}

	events, err := a.application.Events(r.Context(), query.From, query.To)
	if err != nil {
		sendErrorJSON(w, r, http.StatusBadRequest, err, "can't get events")
		return
	}

	if len(events) == 0 {
		sendErrorJSON(w, r, http.StatusNotFound, err, "can't get events")
		return
	}

	sendDataJSON(w, r, http.StatusOK, events)
}

func (a *API) Routes() []Route {
	return []Route{
		{
			Name:   "CreateEvent",
			Method: http.MethodPost,
			Path:   "/event/create",
			Func:   a.createEvent,
		},
		{
			Name:   "UpdateEvent",
			Method: http.MethodPost,
			Path:   "/event/update",
			Func:   a.updateEvent,
		},
		{
			Name:   "RemoveEvent",
			Method: http.MethodPost,
			Path:   "/event/remove",
			Func:   a.removeEvent,
		},
		{
			Name:   "Events",
			Method: http.MethodGet,
			Path:   "/events",
			Func:   a.events,
		},
	}
}
