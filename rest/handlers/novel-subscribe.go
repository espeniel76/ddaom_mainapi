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

func NovelSubscribe(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqMember, _ := strconv.Atoi(req.Vars["seq_member"])
	// 블록처리된 유저 여부 (보내는 사람, 받는사람 둘다)
	if isBlocked(userToken.SeqMember) {
		res.ResultCode = define.BLOCKED_ME
		return res
	}
	if isBlocked(int64(_seqMember)) {
		res.ResultCode = define.BLOCKED_USER
		return res
	}

	if _seqMember == int(userToken.SeqMember) {
		res.ResultCode = define.SELF_SUBSCRIBE
		return res
	}

	mySubscribe := define.NONE
	var cnt int64

	mdb := db.List[define.Mconn.DsnMaster]
	ldbMe := GetMyLogDbMaster(userToken.Allocated)
	ldbYour := getUserLogDbMaster(mdb, int64(_seqMember))

	// 사용자 존재 여부
	result := mdb.Model(schemas.MemberDetail{}).Where("seq_member = ?", _seqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 내가 너를 구독 한적 있는가
	memberSubscribe := schemas.MemberSubscribe{}
	result = ldbMe.Model(&memberSubscribe).
		Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMember).Scan(&memberSubscribe)
	if corm(result, &res) {
		return res
	}

	// 데이터가 없다
	if memberSubscribe.SeqMemberSubscribe == 0 {

		// 1. 내가 너를 구독 한다는 데이터 입력
		result = ldbMe.Create(&schemas.MemberSubscribe{
			SeqMember:         userToken.SeqMember,
			SeqMemberOpponent: int64(_seqMember),
			Status:            define.FOLLOWING,
		})
		if corm(result, &res) {
			return res
		}

		// 2. 너는 나의 구독을 받았다 데이터 입력
		result = ldbYour.Create(&schemas.MemberSubscribe{
			SeqMember:         int64(_seqMember),
			SeqMemberOpponent: userToken.SeqMember,
			Status:            define.FOLLOWER,
		})
		if corm(result, &res) {
			return res
		}

		result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
		if corm(result, &res) {
			return res
		}

		mySubscribe = define.FOLLOWING

		// 데이터가 있다
	} else {
		switch memberSubscribe.Status {
		case define.FOLLOWING: // 1. 상대에게 내가 구독을 한 상태
			// 서로 데이터 삭제
			result = ldbMe.Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMember).
				Delete(&schemas.MemberSubscribe{})
			if corm(result, &res) {
				return res
			}
			result = ldbYour.Where("seq_member = ? AND seq_member_opponent = ?", _seqMember, userToken.SeqMember).
				Delete(&schemas.MemberSubscribe{})
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMember)
			if corm(result, &res) {
				return res
			}

			mySubscribe = define.NONE

		case define.FOLLOWER: // 2. 상대에게 내가 구독을 받은 상태
			// 서로 BOTH 로 변경
			result = ldbMe.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMember).
				Update("status", define.BOTH)
			if corm(result, &res) {
				return res
			}
			result = ldbYour.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_opponent = ?", _seqMember, userToken.SeqMember).
				Update("status", define.BOTH)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
			if corm(result, &res) {
				return res
			}
			mySubscribe = define.BOTH

		case define.BOTH: // 3. 서로가 구독중인 상태
			// 상대에게 내가 구독을 받은 상태로 변경
			result = ldbMe.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMember).
				Update("status", define.FOLLOWER)
			if corm(result, &res) {
				return res
			}
			result = ldbYour.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_opponent = ?", _seqMember, userToken.SeqMember).
				Update("status", define.FOLLOWING)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMember)
			if corm(result, &res) {
				return res
			}
			mySubscribe = define.FOLLOWER
		}
	}

	data := make(map[string]string)
	data["my_subscribe"] = mySubscribe
	res.Data = data

	switch mySubscribe {
	case define.FOLLOWING:
		fallthrough
	case define.BOTH:
		go pushSubscribeTopic(userToken.SeqMember, int64(_seqMember))
	}

	go cacheMainPopularWriter()

	return res
}

func pushSubscribeTopic(seqMemberFrom int64, seqMemberTo int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	userInfoFrom := getUserInfo(seqMemberFrom)
	userInfoTo := getUserInfoPush(seqMemberTo)
	isNight := tools.IsNight()

	if userInfoTo.SeqMember > 0 && userInfoTo.IsNewFollower {
		mdb := db.List[define.Mconn.DsnMaster]
		alarm := schemas.Alarm{
			SeqMember:  userInfoTo.SeqMember,
			Title:      "따옴",
			TypeAlarm:  4,
			ValueAlarm: int(seqMemberFrom),
			Step:       0,
			Content:    userInfoFrom.NickName + "님이 나를 구독하였습니다",
		}
		mdb.Create(&alarm)
		push := InfoPushTopic{}
		isPush := false
		query := "SELECT seq_member, is_night_push FROM member_details WHERE seq_member = ? AND is_new_follower = true"
		mdb.Raw(query, userInfoTo.SeqMember).Scan(&push)
		fmt.Println(push)
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

func pushSubscribe(seqMemberFrom int64, seqMemberTo int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	userInfoFrom := getUserInfo(seqMemberFrom)
	userInfoTo := getUserInfoPush(seqMemberTo)
	isNight := tools.IsNight()

	if userInfoTo.SeqMember > 0 && userInfoTo.IsNewFollower {
		mdb := db.List[define.Mconn.DsnMaster]
		alarm := schemas.Alarm{
			SeqMember:  userInfoTo.SeqMember,
			Title:      "따옴",
			TypeAlarm:  4,
			ValueAlarm: int(seqMemberFrom),
			Step:       0,
			Content:    userInfoFrom.NickName + "님이 나를 구독하였습니다",
		}
		mdb.Create(&alarm)

		// 발송 대상 추출
		listPush := []InfoPush{}
		listFinalPush := []InfoPush{}
		query := "SELECT mpt.push_token, mpt.seq_member, md.is_night_push FROM member_push_tokens mpt LEFT JOIN member_details md ON md.seq_member = mpt.seq_member WHERE md.seq_member = ? AND md.is_new_follower = true"
		mdb.Raw(query, userInfoTo.SeqMember).Scan(&listPush)
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
