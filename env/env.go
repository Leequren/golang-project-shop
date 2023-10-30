package env

import "github.com/gorilla/sessions"

var (
	Key   = []byte("super-secret-key")
	Store = sessions.NewCookieStore(Key)
)
