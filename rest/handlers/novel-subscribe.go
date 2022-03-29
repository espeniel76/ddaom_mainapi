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

	myLogDb := GetMyLogDb(userToken.Allocated)
	masterDB := db.List[define.DSN_MASTER]

	// 사용자 존재 여부
	result := masterDB.Model(schemas.MemberDetail{}).Where("seq_member = ?", _seqMember).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 구독 상태 확인
	memberSubscribe := schemas.MemberSubscribe{}
	result = myLogDb.Model(&memberSubscribe).
		Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).Scan(&memberSubscribe)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if memberSubscribe.SeqMemberSubscribe == 0 { // 존재하지 않음
		// 1. 로그넣기
		result = myLogDb.Create(&schemas.MemberSubscribe{
			SeqMember:          userToken.SeqMember,
			SeqMemberFollowing: int64(_seqMember),
			SubscribeYn:        true,
		})
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		// 2. 구독 카운트 업데이트
		result = masterDB.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		mySubscribe = true
	} else { // 존재함
		// 1. 구독 상태
		if memberSubscribe.SubscribeYn {
			result = myLogDb.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).
				Update("subscribe_yn", false)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			result = masterDB.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMember)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			mySubscribe = false
		} else {
			result = myLogDb.Model(&schemas.MemberSubscribe{}).
				Where("seq_member = ? AND seq_member_following = ?", userToken.SeqMember, _seqMember).
				Update("subscribe_yn", true)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			result = masterDB.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe + 1 WHERE seq_member = ?", _seqMember)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
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
