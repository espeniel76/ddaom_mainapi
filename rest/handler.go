package rest

import (
	"github.com/gorilla/mux"

	"ddaom/rest/handlers"
)

func Handlers(mux *mux.Router) {

	// Initialize DB
	mux.HandleFunc("/db/initialize", common(handlers.InitialDb)) // 1
	// auth
	mux.HandleFunc("/auth/login", common(handlers.AuthLogin))                // 2
	mux.HandleFunc("/auth/login/detail", common(handlers.AuthLoginDetail))   // 3
	mux.HandleFunc("/auth/login/refresh", common(handlers.AuthLoginRefresh)) // 4
	mux.HandleFunc("/auth/info/update", common(handlers.AuthInfoUpdate))     // 5
	mux.HandleFunc("/auth/info", common(handlers.AuthInfo))                  // 6
	mux.HandleFunc("/auth/withdrawal", common(handlers.AuthWithdrawal))      // 9

	// 리소스
	mux.HandleFunc("/assets", common(handlers.Assets))            // 10
	mux.HandleFunc("/keyword", common(handlers.Keyword))          // 11
	mux.HandleFunc("/asset/images", common(handlers.ImageColors)) // 12
	mux.HandleFunc("/asset/genres", common(handlers.Genres))      // 13
	mux.HandleFunc("/asset/slangs", common(handlers.Slangs))      // 14

	// 소설 등록
	mux.HandleFunc("/novel/check/title", common(handlers.NovelCheckTitle)) // 15
	mux.HandleFunc("/novel/write/step1", common(handlers.NovelWriteStep1)) // 16
	mux.HandleFunc("/novel/write/step2", common(handlers.NovelWriteStep2)) // 17
	mux.HandleFunc("/novel/write/step3", common(handlers.NovelWriteStep3)) // 18
	mux.HandleFunc("/novel/write/step4", common(handlers.NovelWriteStep4)) // 19

	// 진행중 소설
	mux.HandleFunc("/novel/list/live", common(handlers.NovelListLive))                 // 20
	mux.HandleFunc("/novel/list/step2", common(handlers.NovelListStep2))               // 21
	mux.HandleFunc("/novel/list/step3", common(handlers.NovelListStep3))               // 22
	mux.HandleFunc("/novel/list/step4", common(handlers.NovelListStep4))               // 23
	mux.HandleFunc("/novel/view/{seq_novel_step1:[0-9]+}", common(handlers.NovelView)) // 24
	mux.HandleFunc("/novel/view/step", common(handlers.NovelViewStep))                 // 25

	// 완결 소설
	mux.HandleFunc("/novel/list/finish", common(handlers.NovelListFinish))                           // 26
	mux.HandleFunc("/novel/view/finish/{seq_novel_finish:[0-9]+}", common(handlers.NovelViewFinish)) // 27

	// 소설 좋아요
	mux.HandleFunc("/novel/like/step1/{seq_novel_step1:[0-9]+}", common(handlers.NovelLikeStep1)) // 28
	mux.HandleFunc("/novel/like/step2/{seq_novel_step2:[0-9]+}", common(handlers.NovelLikeStep2)) // 29
	mux.HandleFunc("/novel/like/step3/{seq_novel_step3:[0-9]+}", common(handlers.NovelLikeStep3)) // 30
	mux.HandleFunc("/novel/like/step4/{seq_novel_step4:[0-9]+}", common(handlers.NovelLikeStep4)) // 31

	// 소설 북마크
	mux.HandleFunc("/novel/bookmark/{seq_novel_finish:[0-9]+}", common(handlers.NovelBookmark)) // 32

	// 작가 구독
	mux.HandleFunc("/novel/subscribe/{seq_member:[0-9]+}", common(handlers.NovelSubscribe)) // 33

	// 즐겨 찾기
	mux.HandleFunc("/novel/bookmark/list", common(handlers.NovelBookmarkList))     // 34
	mux.HandleFunc("/novel/bookmark/delete", common(handlers.NovelBookmarkDelete)) // 35

	// 마이페이지
	mux.HandleFunc("/mypage/info/{seq_member:[0-9]+}", common(handlers.MypageInfo))                               // 36
	mux.HandleFunc("/mypage/list/live", common(handlers.MypageListLive))                                          // 37
	mux.HandleFunc("/mypage/list/finish", common(handlers.MypageListFinish))                                      // 38
	mux.HandleFunc("/mypage/list/temp", common(handlers.MypageListTemp))                                          // 39
	mux.HandleFunc("/mypage/list/complete", common(handlers.MypageListComplete))                                  // 40
	mux.HandleFunc("/mypage/view/complete/{step:[0-9]+}/{seq_novel:[0-9]+}", common(handlers.MypageViewComplete)) // 41
	mux.HandleFunc("/mypage/view/live/{step:[0-9]+}/{seq_novel:[0-9]+}", common(handlers.MypageViewLive))         // 42
	mux.HandleFunc("/mypage/view/finish/{step:[0-9]+}/{seq_novel:[0-9]+}", common(handlers.MypageViewFinish))     // 43
	mux.HandleFunc("/mypage/list/live", common(handlers.MypageListLive))                                          // 44
	mux.HandleFunc("/mypage/list/subscribe", common(handlers.MypageListSubscribe))                                // 45
	mux.HandleFunc("/mypage/list/alarm", common(handlers.MypageListAlarm))                                        // 46
	mux.HandleFunc("/mypage/alarm/{seq_alarm:[0-9]+}", common(handlers.MypageAlarmReceiveSet))                    // 47

	// 메인
	mux.HandleFunc("/main/{seq_keyword:[0-9]+}", common(handlers.Main))                               // 48
	mux.HandleFunc("/main/keyword/{seq_keyword:[0-9]+}", common(handlers.MainKeyword))                // 49
	mux.HandleFunc("/main/temp/info/{step:[0-9]+}/{seq_novel:[0-9]+}", common(handlers.MainTempInfo)) // 50

	// 설정
	mux.HandleFunc("/config/alarm", common(handlers.ConfigAlarm))        // 51
	mux.HandleFunc("/config/alarm/get", common(handlers.ConfigAlarmGet)) // 52

	// server
	mux.HandleFunc("/service/inquiry", common(handlers.ServiceInquiry))              // 53
	mux.HandleFunc("/service/inquiry/edit", common(handlers.ServiceInquiryEdit))     // 54
	mux.HandleFunc("/service/inquiry/delete", common(handlers.ServiceInquiryDelete)) // 55
	mux.HandleFunc("/service/inquiry/list", common(handlers.ServiceInquiryList))     // 56
	mux.HandleFunc("/service/notice/list", common(handlers.ServiceNoticeList))       // 57
	mux.HandleFunc("/service/faq/list", common(handlers.ServiceFaqList))             // 58
}
