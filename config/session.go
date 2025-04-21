package config

import "github.com/gorilla/sessions"

const SESSION_ID = "_go_session_"

var Store = sessions.NewCookieStore([]byte("hjklfasdhjklfd"))
