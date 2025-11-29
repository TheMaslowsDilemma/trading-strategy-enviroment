package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

const cookie_name = "session_token"

type Store struct {
	mu		sync.RWMutex
	data	map[string]session
}

type session struct {
	UserId	int64
	Token	string
	Expires	time.Time
}

var store = Store{data: make(map[string]session)}

func generateId() (string, error) {
	var (
		bs	[]byte
		err	error
	)

	bs = make([]byte, 32)
	_, err = rand.Read(bs)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bs), nil
}

func Set(w http.ResponseWriter, userID int64) error {
	var (
		session_id	string
		expires		time.Time
		err			error
	)

	session_id, err = generateId()
	if err != nil {
		return err
	}

	expires = time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:		cookie_name,
		Value:		session_id,
		Expires:	expires,
		HttpOnly:	true,
		Secure:		false,
		Path:		"/",
		SameSite:	http.SameSiteLaxMode,
	})

	store.mu.Lock()
	store.data[session_id] = session{UserId: userID, Expires: expires}
	store.mu.Unlock()

	return nil
}

func Get(r *http.Request) (int64, bool) {
	var (
		cookie	*http.Cookie
		s		session
		ok		bool
		err		error
	)

	cookie, err = r.Cookie(cookie_name)
	if err != nil {
		return 0, false
	}

	store.mu.RLock()
	s, ok = store.data[cookie.Value]
	store.mu.RUnlock()

	if !ok || time.Now().After(s.Expires) {
		return 0, false
	}
	return s.UserId, true
}

func Clear(w http.ResponseWriter, r *http.Request) {
	var (
		cookie	*http.Cookie
		err		error
	)

	cookie, err = r.Cookie(cookie_name)
	if err != nil {
		return
	}
	store.mu.Lock()
	delete(store.data, cookie.Value)
	store.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   cookie_name,
		MaxAge: -1,
		Path:   "/",
	})
}