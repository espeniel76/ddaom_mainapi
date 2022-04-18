package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func NovelSubscribe(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqMember, _ := strconv.Atoi(req.Vars["seq_member"])

	if _seqMember == int(userToken.SeqMember) {
		res.ResultCode = define.SELF_SUBSCRIBE
		return res
	}

	mySubscribe := define.NONE
	var cnt int64

	mdb := db.List[define.DSN_MASTER]
	ldbMe := GetMyLogDb(userToken.Allocated)
	ldbYour := getUserLogDb(mdb, int64(_seqMember))

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

	// if memberSubscribe.SeqMemberSubscribe == 0 { // 존재하지 않음
	// 	// 1. 로그넣기
	// 	result = ldbMe.Create(&schemas.MemberSubscribe{
	// 		SeqMember:         userToken.SeqMember,
	// 		SeqMemberOpponent: int64(_seqMember),
	// 		Status:            define.FOLLOWING,
	// 	})
	// 	if corm(result, &res) {
	// 		return res
	// 	}

	// 	// 2. 구독 카운트 업데이트
	// 	result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
	// 	if corm(result, &res) {
	// 		return res
	// 	}
	// 	mySubscribe = true
	// } else { // 존재함
	// 	// 1. 그냥 삭제
	// 	result = ldbMe.Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMember).Delete(&memberSubscribe)
	// 	if corm(result, &res) {
	// 		return res
	// 	}
	// 	result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMember)
	// 		if corm(result, &res) {
	// 			return res
	// 		}
	// 		mySubscribe = false

	// 	if memberSubscribe.SubscribeYn {
	// 		result = ldbMe.Model(&schemas.MemberSubscribe{}).
	// 			Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).
	// 			Update("subscribe_yn", false)
	// 		if corm(result, &res) {
	// 			return res
	// 		}

	// 	} else {
	// 		result = ldbMe.Model(&schemas.MemberSubscribe{}).
	// 			Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).
	// 			Update("subscribe_yn", true)
	// 		if corm(result, &res) {
	// 			return res
	// 		}
	// 		result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
	// 		if corm(result, &res) {
	// 			return res
	// 		}
	// 		mySubscribe = true
	// 	}
	// }

	data := make(map[string]string)
	data["my_subscribe"] = mySubscribe
	res.Data = data

	return res
}
