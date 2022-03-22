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
	mux.HandleFunc("/skin", common(handlers.Skin))
	mux.HandleFunc("/genre", common(handlers.Genre))
	mux.HandleFunc("/assets", common(handlers.Assets))

	// 소설 등록
	mux.HandleFunc("/novel/write/step1", common(handlers.NovelWriteStep1))
	mux.HandleFunc("/novel/write/step2", common(handlers.NovelWriteStep2))
	mux.HandleFunc("/novel/write/step3", common(handlers.NovelWriteStep3))
	mux.HandleFunc("/novel/write/step4", common(handlers.NovelWriteStep4))

	// 소설 목록
	mux.HandleFunc("/novel/list/live", common(handlers.NovelListLive))
	mux.HandleFunc("/novel/list/step2", common(handlers.NovelListStep2))
	mux.HandleFunc("/novel/list/step3", common(handlers.NovelListStep3))
	mux.HandleFunc("/novel/list/step4", common(handlers.NovelListStep4))

}
