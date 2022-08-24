package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
	"strconv"
)

func NovelLikeStep1(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	// 블록처리된 유저 여부 (보내는 사람, 받는사람 둘다)
	if isBlocked(userToken.SeqMember) {
		res.ResultCode = define.BLOCKED_ME
		return res
	}

	_seqNovelStep1, _ := strconv.Atoi(req.Vars["seq_novel_step1"])
	myLike := false
	var cnt int64
	var scanCount int64

	ldb := GetMyLogDbMaster(userToken.Allocated)
	mdb := db.List[define.Mconn.DsnMaster]

	// 소설 존재/삭제 여부
	novelStep := schemas.NovelStep1{}
	result := mdb.Model(&novelStep).Select("cnt_like, deleted_yn, seq_member").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&novelStep).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if isBlocked(novelStep.SeqMember) {
		res.ResultCode = define.BLOCKED_USER
		return res
	}

	cnt = novelStep.CntLike
	if novelStep.DeletedYn {
		res.ResultCode = define.DELETED_NOVEL
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 좋아요 상태 확인
	memberLikeStep1 := schemas.MemberLikeStep1{}
	result = ldb.Model(&memberLikeStep1).
		Where("seq_novel_step1 = ? AND seq_member = ?", _seqNovelStep1, userToken.SeqMember).Scan(&memberLikeStep1)
	seqKeyword := getSeqKeyword(1, int64(_seqNovelStep1))
	if corm(result, &res) {
		return res
	}

	// 작가 seq
	var seqMemberWriter int64
	mdb.Model(schemas.NovelStep1{}).Where("seq_novel_step1 = ?", _seqNovelStep1).Select("seq_member").Scan(&seqMemberWriter)

	if memberLikeStep1.SeqMemberLike == 0 { // 존재하지 않음
		// 1. 로그넣기
		result = ldb.Create(&schemas.MemberLikeStep1{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep1: int64(_seqNovelStep1),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		// 2. 좋아요 카운트 업데이트
		result = mdb.Exec("UPDATE novel_step1 SET cnt_like = cnt_like + 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
		updateKeywordMemberLike(seqMemberWriter, seqKeyword, "PLUS")
	} else { // 존재함
		// 1. 좋아요 상태
		if memberLikeStep1.LikeYn {
			result = ldb.Model(&schemas.MemberLikeStep1{}).
				Where("seq_member = ? AND seq_novel_step1 = ?", userToken.SeqMember, _seqNovelStep1).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step1 SET cnt_like = cnt_like - 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
			updateKeywordMemberLike(seqMemberWriter, seqKeyword, "MINUS")
		} else {
			result = ldb.Model(&schemas.MemberLikeStep1{}).
				Where("seq_member = ? AND seq_novel_step1 = ?", userToken.SeqMember, _seqNovelStep1).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step1 SET cnt_like = cnt_like + 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
			if corm(result, &res) {
				return res
			}
			myLike = true
			cnt++
			mdb.Exec("UPDATE keywords SET cnt_like = cnt_like + 1 WHERE seq_keyword = ?", seqKeyword)
			updateKeywordMemberLike(seqMemberWriter, seqKeyword, "PLUS")
		}

	}

	data := make(map[string]interface{})
	data["my_like"] = myLike
	data["cnt_like"] = cnt
	res.Data = data

	// push 날리기
	if myLike {
		go pushLikeTopic(1, int64(_seqNovelStep1), userToken.SeqMember)
	}

	go cacheMainPopularWriterLike()

	return res
}

func updateKeywordMemberLike(seqMember int64, seqKeyword int64, direction string) {
	mdb := db.List[define.Mconn.DsnMaster]
	if direction == "PLUS" {
		mdb.Exec("UPDATE keywords SET cnt_like = cnt_like + 1 WHERE seq_keyword = ?", seqKeyword)
		mdb.Exec("UPDATE member_details SET cnt_like = cnt_like + 1 WHERE seq_member = ?", seqMember)
	} else {
		mdb.Exec("UPDATE keywords SET cnt_like = cnt_like - 1 WHERE seq_keyword = ?", seqKeyword)
		mdb.Exec("UPDATE member_details SET cnt_like = cnt_like - 1 WHERE seq_member = ?", seqMember)
	}

}

// push 방식 변경
func pushLikeTopic(step int8, seqNovel int64, seqMember int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	info := getNovel(step, seqNovel)
	isNight := tools.IsNight()
	if info.SeqMember > 0 {
		mdb := db.List[define.Mconn.DsnMaster]
		userInfoFrom := getUserInfo(seqMember)
		alarm := schemas.Alarm{
			SeqMember:  info.SeqMember,
			Title:      "따옴",
			TypeAlarm:  2,
			ValueAlarm: int(seqNovel),
			Step:       step,
			Content:    "\"" + info.Title + " - step" + fmt.Sprintf("%d", step) + "\"를 " + userInfoFrom.NickName + " 님이 좋아합니다.",
		}
		mdb.Create(&alarm)
		push := InfoPushTopic{}
		isPush := false
		query := "SELECT seq_member, is_night_push FROM member_details WHERE seq_member = ? AND is_liked = true"
		mdb.Raw(query, info.SeqMember).Scan(&push)
		if isNight {
			if push.IsNightPush {
				isPush = true
			}
		} else {
			isPush = true
		}
		if isPush {
			go tools.SendPushMessageTopic(&alarm)
		}
	}
}

// 사용자 개별 토큰으로 push (사용하지 않음)
func pushLike(step int8, seqNovel int64, seqMember int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	info := getNovel(step, seqNovel)
	isNight := tools.IsNight()
	if info.SeqMember > 0 {
		mdb := db.List[define.Mconn.DsnMaster]

		// 로그 쌓기
		userInfoFrom := getUserInfo(seqMember)
		alarm := schemas.Alarm{
			SeqMember:  info.SeqMember,
			Title:      "따옴",
			TypeAlarm:  2,
			ValueAlarm: int(seqNovel),
			Step:       step,
			Content:    "\"" + info.Title + " - step" + fmt.Sprintf("%d", step) + "\"를 " + userInfoFrom.NickName + " 님이 좋아합니다.",
		}
		mdb.Create(&alarm)
		listPush := []InfoPush{}
		listFinalPush := []InfoPush{}
		query := "SELECT mpt.push_token, mpt.seq_member, md.is_night_push FROM member_push_tokens mpt LEFT JOIN member_details md ON md.seq_member = mpt.seq_member WHERE md.seq_member = ? AND md.is_liked = true"
		mdb.Raw(query, info.SeqMember).Scan(&listPush)
		for _, o := range listPush {
			if isNight {
				if o.IsNightPush {
					listFinalPush = append(listFinalPush, o)
				}
			} else {
				listFinalPush = append(listFinalPush, o)
			}
		}
		for _, o := range listFinalPush {
			go tools.SendPushMessage(o.PushToken, &alarm)
		}
	}
}

type InfoPushTopic struct {
	SeqMember   int64
	IsNightPush bool
}

// 발송 대상 추출
type InfoPush struct {
	PushToken   string
	SeqMember   int64
	IsNightPush bool
}
