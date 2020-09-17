package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GShamian/tavern-of-games/internal/app/model"
	"github.com/GShamian/tavern-of-games/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "tavern_of_games"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

// Server object
type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

// newServer func. Constructor for a server. It creates new
// server instance with mux router, logger and our imported
// session store and store.
func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

// ServeHTTP func. Wrap for ServeHTTP function
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter func. Configuring router
func (s *server) configureRouter() {
	// Appending middleware func setRequestID to the router chain
	s.router.Use(s.setRequestID)
	// Appending middleware func logRequest to the router chain
	s.router.Use(s.logRequest)
	// Applying the CORS middleware to our router
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	// Registering a new route for url /users for our router
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	// Registering a new route for url /sessions for our router
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	// Registering a new route for /private url path prefix and
	// creating a subrouter for the route.
	private := s.router.PathPrefix("/private").Subrouter()
	// Appending middleware func authenticateUser to the router chain
	private.Use(s.authenticateUser)
	// Registering a new route for url /whoami for our router
	private.HandleFunc("/whoami", s.handleWhoami())
}

// setRequestID func. Middleware func for http handler, that sets id in
// request header.
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Creating new request id
		id := uuid.New().String()
		// Setting id in request header
		w.Header().Set("X-Request-ID", id)
		// Handling the request with context with value
		// that store request id.
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// logRequest func. Middleware func for http handler, that
// logs requests.
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Creating new logger that has two fields
		logger := s.logger.WithFields(logrus.Fields{
			// Writing the address that sent network request
			"remote_addr": r.RemoteAddr,
			// Writing the request id from context with value
			"request_id": r.Context().Value(ctxKeyRequestID),
		})
		// Logging request method and unmodified request-target
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		// Getting local time
		start := time.Now()
		// Making response with statut ok
		rw := &responseWriter{w, http.StatusOK}
		// Serving response
		next.ServeHTTP(rw, r)
		// Logger variable
		var level logrus.Level
		switch {
		case rw.code >= 500:
			// Assigning logger level with error
			level = logrus.ErrorLevel
		case rw.code >= 400:
			// Assigning logger level with warning
			level = logrus.WarnLevel
		default:
			// Assigning logger level with information
			level = logrus.InfoLevel
		}
		// Main logging process
		logger.Logf(
			// Logger level
			level,
			"completed with %d %s in %v",
			// Logging http status
			rw.code,
			// Logging formatted http status
			http.StatusText(rw.code),
			// Logging time duration
			time.Now().Sub(start),
		)
	})
}

// authenticateUser func. Middleware func for http handler, that
// autentificates user.
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Creating request session entity
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		// Getting user id
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		// Checking for user autentification
		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		// Serving response with context
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

// handleUsersCreate func. Middleware func for http handler, that
// handles user creation.
func (s *server) handleUsersCreate() http.HandlerFunc {
	// Creating request object
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Creating request entity
		req := &request{}
		// Decoding json from request to our entity
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// Creating user model entity with email and
		// password from the request.
		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		// Adding user model to DB
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Clearing users password field
		u.Sanitize()
		// Creating response with status 201 (User created)
		s.respond(w, r, http.StatusCreated, u)
	}
}

// handleSessionsCreate func. Middleware func for http handler, that
// handles session creation.
func (s *server) handleSessionsCreate() http.HandlerFunc {
	// Creating request object
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Creating request entity
		req := &request{}
		// Decoding json from request to our entity
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// Finding user with Email from the request
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		// Creating request session entity
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		// Adding user id to session values
		session.Values["user_id"] = u.ID
		// Saving session in the underlying store
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		// Creating response with status 200 (OK status)
		s.respond(w, r, http.StatusOK, nil)
	}
}

// handleWhoami func. Middleware func for http handler, that
// handles operation of getting information about actual user of
// this session.
func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Creating response with context which stores information about user
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

// error func. Function that create a response with status code and
// a map which is initialised with and error.
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond func. Function that writes status code in response header
// and encodes data in json for response.
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
