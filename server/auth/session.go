package auth

import (
	"encoding/json"
	"sync"

	"go-app/server/config"
	"go-app/server/storage"

	// "go-app/server/storage/redisstorage"

	// "go-app/server/storage/redisstorage"
	"net/http"

	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

// Session defines session storage methods
type SessionAuth interface {
	NewSessionID() string
	Create(string, http.ResponseWriter) error
	Get(*http.Request) (*UserSession, error)
	Update(http.ResponseWriter, *http.Request, *UserSession) error
	Delete(*http.Request) error
	CreateAndReturn(st string, w http.ResponseWriter) (string, error)
}

// UserSession := user session representation
type UserSession struct {
	Token string                 `json:"token"`
	Other map[string]interface{} `json:"-"`
}

// ToJSON := Convering struct to json
func (us *UserSession) ToJSON() string {
	json, _ := json.Marshal(us)
	return string(json)
}

type SessionAuthImpl struct {
	Client storage.Redis
	Config *config.SessionConfig
	mux    sync.Mutex
}

type SessionAuthOpts struct {
	Client storage.Redis
	Config *config.SessionConfig
}

func NewSessionAuth(opts *SessionAuthOpts) SessionAuth {
	sai := SessionAuthImpl{
		Client: opts.Client,
		Config: opts.Config,
	}
	return &sai
}

// Get := implementing Get session method
func (s *SessionAuthImpl) Get(r *http.Request) (*UserSession, error) {
	cookie, err := r.Cookie(s.Config.CookieConfig.Name)
	if err != nil {
		return nil, err
	}
	s.mux.Lock()
	result, err := s.get(cookie.Value)
	s.mux.Unlock()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, nil
	}
	return &UserSession{Token: result}, nil
}

func (s *SessionAuthImpl) get(key string) (string, error) {
	return redis.String(s.Client.Do("GET", key))
}

// Update := implementing updates user existing session
func (s *SessionAuthImpl) Update(w http.ResponseWriter, r *http.Request, sess *UserSession) error {
	cookie, err := r.Cookie(s.Config.CookieConfig.Name)
	if err != nil {
		return err
	}
	err = s.set(cookie.Value, sess.Token)
	return err
}

func (s *SessionAuthImpl) set(key, val string) error {
	_, err := s.Client.Do("SET", key, val)
	return err
}

// Create := Creating Session
func (s *SessionAuthImpl) Create(st string, w http.ResponseWriter) error {
	sessionID := s.NewSessionID()
	err := s.set(sessionID, st)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     s.Config.CookieConfig.Name,
		Value:    sessionID,
		Path:     s.Config.CookieConfig.Path,
		HttpOnly: s.Config.CookieConfig.HttpOnly,
		Domain:   s.Config.CookieConfig.Domain,
		Secure:   s.Config.CookieConfig.Secure,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
	return err
}

// SetSessionID
func (s *SessionAuthImpl) CreateAndReturn(st string, w http.ResponseWriter) (string, error) {
	sessionID := s.NewSessionID()
	err := s.set(sessionID, st)
	if err != nil {
		return "", err
	}
	cookie := &http.Cookie{
		Name:     s.Config.CookieConfig.Name,
		Value:    sessionID,
		Path:     s.Config.CookieConfig.Path,
		HttpOnly: s.Config.CookieConfig.HttpOnly,
		Domain:   s.Config.CookieConfig.Domain,
		Secure:   s.Config.CookieConfig.Secure,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
	return sessionID, err
}

// Delete := implementing Delete session method
func (s *SessionAuthImpl) Delete(r *http.Request) error {
	cookie, err := r.Cookie(s.Config.CookieConfig.Name)
	if err != nil {
		return err
	}
	_, err = s.Client.Do("DEL", cookie.Value)
	return nil
}

// NewSessionID := return unique session ID
func (s *SessionAuthImpl) NewSessionID() string {
	return uuid.NewV4().String()
}
