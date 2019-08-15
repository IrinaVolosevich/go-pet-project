package models

import (
	"net/http"
	"net/url"
	"time"
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

const (
	sessionLength     = 24 * 3 * time.Hour
	sessionCookieName = "GophrSession"
	sessionIDLength   = 20
)

func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLength)

	session := &Session{
		ID:     GenerateID("sess", sessionIDLength),
		Expiry: expiry,
	}

	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: expiry,
	}

	http.SetCookie(w, &cookie)
	return session
}

func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)

	if cookie == nil {
		return nil
	}

	if err != nil {
		panic(err)
	}

	session, err := SessionStore.Find(GlobalSessionStore, cookie.Value)

	if err != nil {
		panic(err)
	}

	if session == nil {
		return nil
	}

	if session.Expired() {
		SessionStore.Delete(GlobalSessionStore, session)
		return nil
	}

	return session
}

func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}

func RequestUser(r *http.Request) *User {
	session := RequestSession(r)

	if session == nil || session.UserID == "" {
		return nil
	}

	user, err := GlobalUserStore.Find(session.ID)

	if err != nil {
		panic(err)
	}

	return user
}

func RequireLogin(w http.ResponseWriter, r *http.Request) {
	if RequestUser(r) != nil {
		return
	}

	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))

	http.Redirect(w, r, "/login?" + query.Encode(), http.StatusFound)
}
