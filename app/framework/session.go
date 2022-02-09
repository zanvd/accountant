package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var ctx = context.Background()

type ConnectionConfig struct {
	Db       int
	Host     string
	Password string
	Port     string
	Username string
}

type Session struct {
	Id    string
	Clear bool
	Data  SessionData
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
	Client *redis.Client
}

func (sm *SessionManager) ClearSession(t *Tools, w http.ResponseWriter) (err error) {
	c := http.Cookie{
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Name:     "accountant-session",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // TODO Set based on environment.
		Value:    "",
	}
	http.SetCookie(w, &c)
	_, err = sm.Client.Del(ctx, t.Session.Id).Result()
	t.Session.Data = SessionData{}
	return
}

func (sm *SessionManager) Connect(cc ConnectionConfig) {
	sm.Client = redis.NewClient(&redis.Options{
		Addr:     cc.Host + ":" + cc.Port,
		DB:       cc.Db,
		Password: cc.Password,
		Username: cc.Username,
	})
}

func (sm *SessionManager) GetSession(t *Tools, r *http.Request) error {
	s := sm.createSession()
	// TODO Log errors?
	if c, err := r.Cookie("accountant-session"); err == nil {
		if sd, err := sm.Client.Get(ctx, c.Value).Result(); err == nil {
			var umd SessionData
			if err := json.Unmarshal([]byte(sd), &umd); err != nil {
				return err
			}
			s = Session{
				Id:    c.Value,
				Clear: false,
				Data:  umd,
			}
		}
	}
	t.Session = s
	return nil
}

func (sm *SessionManager) WriteSession(t *Tools, w http.ResponseWriter) (err error) {
	var md []byte
	if md, err = json.Marshal(t.Session.Data); err != nil {
		return
	}
	if err = sm.Client.Set(ctx, t.Session.Id, string(md), time.Hour).Err(); err != nil {
		return
	}
	c := http.Cookie{
		HttpOnly: true,
		MaxAge:   3600,
		Name:     "accountant-session",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // TODO Set based on environment.
		Value:    t.Session.Id,
	}
	http.SetCookie(w, &c)
	return
}

func (sm *SessionManager) createSession() Session {
	return Session{
		Id:    uuid.New().String(),
		Clear: false,
	}
}
