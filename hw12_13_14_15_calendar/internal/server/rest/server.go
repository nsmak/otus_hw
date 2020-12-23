package rest

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
)

type ServerError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *ServerError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *ServerError) Unwrap() error {
	return e.Err
}

type Application interface {
	// TODO
}

type ServerAPI interface {
	Routes() []Route
}

type Route struct {
	Name   string
	Method string
	Path   string
	Func   http.HandlerFunc
}

type Server struct {
	Address string
	public  ServerAPI
	server  *http.Server
	log     app.Logger
}

func NewServer(public ServerAPI, host, port string, logger app.Logger) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		public:  public,
		log:     logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{ // nolint: exhaustivestruct
		Addr:         s.Address,
		Handler:      s.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return &ServerError{Message: "start server error", Err: err}
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("rest server is nil")
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return &ServerError{Message: "stop server error", Err: err}
	}
	return nil
}

func (s *Server) router() *mux.Router {
	router := mux.NewRouter()
	for _, route := range s.public.Routes() {
		handler := alice.New(s.loggingMiddleware).ThenFunc(route.Func)
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
