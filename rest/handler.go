package rest

import (
	"github.com/gorilla/mux"

	"ddaom/rest/handlers"
)

func Handlers(mux *mux.Router) {

	// Initialize DB
	mux.HandleFunc("/db/initialize", common(handlers.InitialDb))
	// auth
	mux.HandleFunc("/auth/login", common(handlers.AuthLogin))
	mux.HandleFunc("/auth/login/detail", common(handlers.AuthLoginDetail))
	mux.HandleFunc("/auth/login/refresh", common(handlers.AuthLoginRefresh))
	mux.HandleFunc("/auth/info/update", common(handlers.AuthInfoUpdate))
	mux.HandleFunc("/auth/info", common(handlers.AuthInfo))
	mux.HandleFunc("/auth/authentication", common(handlers.AuthAuthentication))
	mux.HandleFunc("/auth/authentication/set", common(handlers.AuthAuthenticationSet))

	// 주제어
	mux.HandleFunc("/keyword", common(handlers.Keyword))

}
