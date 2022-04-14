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

	mySubscribe := false
	var cnt int64

	ldb := GetMyLogDb(userToken.Allocated)
	mdb := db.List[define.DSN_MASTER]

	// 사용자 존재 여부
	result := mdb.Model(schemas.MemberDetail{}).Where("seq_member = ?", _seqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 구독 상태 확인
	memberSubscribe := schemas.MemberSubscribe{}
	result = ldb.Model(&memberSubscribe).
		Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).Scan(&memberSubscribe)
	if corm(result, &res) {
		return res
	}
	if memberSubscribe.SeqMemberSubscribe == 0 { // 존재하지 않음
		// 1. 로그넣기
		result = ldb.Create(&schemas.MemberSubscribe{
			SeqMember:          userToken.SeqMember,
			SeqMemberFollowing: int64(_seqMember),
			SubscribeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		// 2. 구독 카운트 업데이트
		result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
		if corm(result, &res) {
			return res
		}
		mySubscribe = true
	} else { // 존재함
		// 1. 구독 상태
		if memberSubscribe.SubscribeYn {
			result = ldb.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).
				Update("subscribe_yn", false)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMember)
			if corm(result, &res) {
				return res
			}
			mySubscribe = false
		} else {
			result = ldb.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).
				Update("subscribe_yn", true)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
			if corm(result, &res) {
				return res
			}
			mySubscribe = true
		}
	}

	data := make(map[string]bool)
	data["my_subscribe"] = mySubscribe
	res.Data = data

	return res
}
