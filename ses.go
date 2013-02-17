package netkit

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"net/http"
	"math/big"
	"net/url"
	"time"
	"sync"
	"fmt"
)

// session cookie store
type SessionStore struct {
	session 	string
	sessions 	map[string]*Session
	gcrate 		int64
	mu 			sync.RWMutex
}

// return a new session store instance
func NewSessionStore(sessionid string) *SessionStore {
	return &SessionStore{
		session: sessionid,
		sessions: map[string]*Session{},
		gcrate: 10,
	}
}

// return an existing session, otherwise create a new session
func (self *SessionStore) GetSession(w http.ResponseWriter, r *http.Request) *Session {
	cookie, err := r.Cookie(self.session)
	if err != nil && err == http.ErrNoCookie {
		return self.newSession(w)
	}
	if cookie, ok := self.sessions[cookie.Value]; ok {
		return cookie
	}
	return self.newSession(w)
}

// terminate a session and cookie data associated
func (self *SessionStore) KillSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(self.session)
	if err != nil || cookie.Value == "" {
		return
	}
	delete(self.sessions, cookie.Value)
	cookie.MaxAge=-1
	http.SetCookie(w, cookie)
}

// create a new cookie, return a session
func (self *SessionStore) newSession(w http.ResponseWriter) *Session {
	id := url.QueryEscape(self.newSessionId())
	self.sessions[id] = &Session{
		id: id, 
		vals: map[string]interface{}{}, 
		created: time.Now(),
	}
	http.SetCookie(w, &http.Cookie{
		Name: self.session, 
		Value:id,
		Path: "/",
		Secure: false,
		HttpOnly: true,
	})
	self.gcollect()
	return self.sessions[id]
}

// return a new session id
func (self *SessionStore) newSessionId() string {
	k := make([]byte, 32)
	if _, err := cryptorand.Read(k); err != nil {
		return fmt.Sprintf(string("%x"), mathrand.Int63())
	}
	return string(k)
}

// session
type Session struct {
	id 		string
	vals 	map[string]interface{}
	created time.Time 
	flash 	[]string
	mu 		sync.Mutex
}

// set a session value
func (self *Session) Set(k string, v interface{}) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.vals[k] = v
}

// get a session value
func (self *Session) Get(k string) interface{} {
	self.mu.Lock()
	defer self.mu.Unlock()
	if v, ok := self.vals[k]; ok {
		return v
	}
	return nil
}

// delete a session value
func (self *Session) Del(k string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	if _, ok := self.vals[k]; ok {
		delete(self.vals, k)
	}
}

// set a flash value
func (self *Session) SetFlash(s string) {
	self.flash = append(self.flash, s)
}

// get a flash value
func (self *Session) GetFlash() string {
	if len(self.flash) == 0 {
		return ""	
	}
	f := self.flash[len(self.flash)-1]
	self.flash = self.flash[:len(self.flash)-1]
	return f
}

// basic garbage collector
func (self *SessionStore) gc() {
	for id, s := range self.sessions {
		if time.Now().Sub(s.created).Seconds() > 3600 {
			self.mu.Lock()
			defer self.mu.Unlock()
			self.sessions[id] = nil
		}
	}
}

// run garbage collector
func (self *SessionStore) gcollect() {
	r, err := cryptorand.Int(cryptorand.Reader, big.NewInt(self.gcrate))
	if err == nil && r.Cmp(big.NewInt(1)) == 0 {
		self.gc()
	}
}