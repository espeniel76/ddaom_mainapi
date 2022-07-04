package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/appleboy/go-fcm"
)

func NovelLikeStep1(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
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
	result := mdb.Model(&novelStep).Select("cnt_like, deleted_yn").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&novelStep).Count(&scanCount)
	if corm(result, &res) {
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
		go pushLike(1, int64(_seqNovelStep1), userToken.SeqMember)
	}

	cacheMainPopularWriter()

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

func pushLike(step int8, seqNovel int64, seqMember int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	info := getNovel(step, seqNovel)
	isNight := false
	if info.SeqMember > 0 {

		// 야간 푸쉬는 받지 않는다
		if info.IsNightPush == false {
			// 낮인지 체크
			now := time.Now()
			if now.Hour() >= 9 && now.Hour() <= 20 {
				isNight = false
			} else {
				isNight = true
			}
		}
		fmt.Println(isNight)

		if !isNight {
			userInfoFrom := getUserInfo(seqMember)

			// 1. 푸쉬 테이블 삽입
			alarm := schemas.Alarm{
				SeqMember:  info.SeqMember,
				Title:      "따옴",
				TypeAlarm:  2,
				ValueAlarm: int(seqNovel),
				Step:       step,
				Content:    "\"" + info.Title + " - step" + fmt.Sprintf("%d", step) + "\"를 " + userInfoFrom.NickName + " 님이 좋아합니다.",
			}
			mdb := db.List[define.Mconn.DsnMaster]
			mdb.Create(&alarm)

			msg := &fcm.Message{
				To: info.PushToken,
				Data: map[string]interface{}{
					"seq_alarm":   alarm.SeqAlarm,
					"type_alarm":  2,
					"value_alarm": seqNovel,
					"step":        step,
				},
				Notification: &fcm.Notification{
					Title: alarm.Title,
					Body:  alarm.Content,
				},
			}

			// Create a FCM client to send the message.
			client, err := fcm.NewClient(define.Mconn.PushServerKey)
			if err != nil {
				// log.Fatalln(err)
				fmt.Println(err)
			}

			// Send the message and receive the response without retries.
			response, err := client.Send(msg)
			if err != nil {
				// log.Fatalln(err)
				fmt.Println(err)
			}

			log.Printf("%#v\n", response)
		}
	}
}
