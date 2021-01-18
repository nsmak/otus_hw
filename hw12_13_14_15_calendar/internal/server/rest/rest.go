package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

var statusCtxKey = NewContextKey("Status")

type JSON map[string]interface{}

type ContextKey struct {
	name string
}

func NewContextKey(name string) *ContextKey {
	return &ContextKey{name: name}
}
func (k *ContextKey) String() string {
	return k.name
}

type Response struct {
	Data  interface{} `json:"data"`
	Error JSON        `json:"error"`
}

func status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), statusCtxKey, status))
}

func sendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string) {
	e := NewError(details, err)
	resp := Response{
		Data:  nil,
		Error: JSON{"message": e.Error()},
	}
	status(r, httpStatusCode)
	sendJSON(w, r, resp)
}

func sendDataJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, data interface{}) { // nolint: unparam
	resp := Response{Data: data, Error: nil}
	status(r, httpStatusCode)
	sendJSON(w, r, resp)
}

func sendJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	_, _ = w.Write(buf.Bytes())
}
