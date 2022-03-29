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

	// 진행중 소설 목록
	mux.HandleFunc("/novel/list/live", common(handlers.NovelListLive))
	mux.HandleFunc("/novel/list/step2", common(handlers.NovelListStep2))
	mux.HandleFunc("/novel/list/step3", common(handlers.NovelListStep3))
	mux.HandleFunc("/novel/list/step4", common(handlers.NovelListStep4))
	mux.HandleFunc("/novel/view/{seq_novel_step1:[0-9]+}", common(handlers.NovelView))

	// 완결 소설 목록
	mux.HandleFunc("/novel/list/finish", common(handlers.NovelListFinish))

	// 소설 좋아요
	mux.HandleFunc("/novel/like/step1/{seq_novel_step1:[0-9]+}", common(handlers.NovelLikeStep1))
	mux.HandleFunc("/novel/like/step2/{seq_novel_step2:[0-9]+}", common(handlers.NovelLikeStep2))
	mux.HandleFunc("/novel/like/step3/{seq_novel_step3:[0-9]+}", common(handlers.NovelLikeStep3))
	mux.HandleFunc("/novel/like/step4/{seq_novel_step4:[0-9]+}", common(handlers.NovelLikeStep4))

	// 소설 북마크
	mux.HandleFunc("/novel/bookmark/{seq_novel_finish:[0-9]+}", common(handlers.NovelBookmark))

	// 작가 구독
	mux.HandleFunc("/novel/subscribe/{seq_member:[0-9]+}", common(handlers.NovelSubscribe))
}
