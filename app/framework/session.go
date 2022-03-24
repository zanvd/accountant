package framework

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id   string
	Data SessionData
}

type SessionData struct {
	LoggedIn bool        `json:"logged_in"`
	User     SessionUser `json:"user"`
}

type SessionUser struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type SessionManager struct {
	CacheManager *CacheManager
	CookieName   string
	KeyPrefix    string
	SecureCookie bool
}

func NewSessionManager(cacheManager *CacheManager, config *Config) *SessionManager {
	return &SessionManager{
		CacheManager: cacheManager,
		CookieName:   config.Session.CookieName,
		KeyPrefix:    config.Session.KeyPrefix + ":",
		SecureCookie: config.Session.SecureCookie,
	}
}

func (sm *SessionManager) ClearSession(s *Session, w http.ResponseWriter) (err error) {
	c := http.Cookie{
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Name:     sm.CookieName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   sm.SecureCookie,
		Value:    "",
	}
	http.SetCookie(w, &c)
	_, err = sm.CacheManager.Delete(s.Id)
	s.Data = SessionData{}
	return
}

func (sm *SessionManager) GetSession(r *http.Request) (Session, error) {
	s := sm.createSession()
	// TODO: Log errors?
	if c, err := r.Cookie(sm.CookieName); err == nil {
		if sd, err := sm.CacheManager.Get(sm.KeyPrefix + c.Value); err == nil {
			var umd SessionData
			if err := json.Unmarshal([]byte(sd), &umd); err != nil {
				return s, err
			}
			s = Session{
				Id:   c.Value,
				Data: umd,
			}
		}
	}
	return s, nil
}

func (sm *SessionManager) WriteSession(s Session, w http.ResponseWriter) (err error) {
	var md []byte
	if md, err = json.Marshal(s.Data); err != nil {
		return
	}
	if err = sm.CacheManager.Set(sm.KeyPrefix+s.Id, string(md), time.Hour); err != nil {
		return
	}
	c := http.Cookie{
		HttpOnly: true,
		MaxAge:   3600,
		Name:     sm.CookieName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   sm.SecureCookie,
		Value:    s.Id,
	}
	http.SetCookie(w, &c)
	return
}

func (sm *SessionManager) createSession() Session {
	return Session{
		Id: uuid.New().String(),
	}
}
