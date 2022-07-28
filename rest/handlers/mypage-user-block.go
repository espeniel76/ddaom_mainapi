package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func MypageUserBlock(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqMemberTo, _ := strconv.Atoi(req.Vars["seq_member"])
	myBlocking := false
	cntBlock := 0

	mdb := db.List[define.Mconn.DsnMaster]
	ldbMe := GetMyLogDbMaster(userToken.Allocated)
	ldbYour := getUserLogDbMaster(mdb, int64(_seqMemberTo))

	memberBlocking := schemas.MemberBlocking{}
	result := ldbMe.Model(&memberBlocking).
		Where("seq_member = ? AND seq_member_to = ?", userToken.SeqMember, _seqMemberTo).Scan(&memberBlocking)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	cntBlock = int(memberBlocking.CntBlock)

	// === 해당 사용자에게 내가 보낸 구독을 취소한다. ===
	memberSubscribe := schemas.MemberSubscribe{}
	result = ldbMe.Model(&memberSubscribe).
		Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMemberTo).Scan(&memberSubscribe)
	if corm(result, &res) {
		return res
	}

	if memberBlocking.SeqMemberBlocking == 0 { // 존재하지 않음
		result = ldbMe.Create(&schemas.MemberBlocking{
			SeqMember:   userToken.SeqMember,
			SeqMemberTo: int64(_seqMemberTo),
			BlockYn:     true,
			CntBlock:    1,
		})
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		// 상대방과 나와의 구독 정보가 존재 할 시
		if memberSubscribe.SeqMemberSubscribe > 0 {
			switch memberSubscribe.Status {
			// 1. 상대에게 내가 구독을 한 상태
			case define.FOLLOWING:
				// 서로 데이터 삭제
				result = ldbMe.Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMemberTo).
					Delete(&schemas.MemberSubscribe{})
				if corm(result, &res) {
					return res
				}
				result = ldbYour.Where("seq_member = ? AND seq_member_opponent = ?", _seqMemberTo, userToken.SeqMember).
					Delete(&schemas.MemberSubscribe{})
				if corm(result, &res) {
					return res
				}
				result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMemberTo)
				if corm(result, &res) {
					return res
				}

			// 2. 서로 맞 구독 상태
			case define.BOTH:
				// 상대에게 내가 구독을 받은 상태로 변경
				result = ldbMe.Model(&schemas.MemberSubscribe{}).
					Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMemberTo).
					Update("status", define.FOLLOWER)
				if corm(result, &res) {
					return res
				}
				result = ldbYour.Model(&schemas.MemberSubscribe{}).
					Where("seq_member = ? AND seq_member_opponent = ?", _seqMemberTo, userToken.SeqMember).
					Update("status", define.FOLLOWING)
				if corm(result, &res) {
					return res
				}
				result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMemberTo)
				if corm(result, &res) {
					return res
				}
			}
		}

		cntBlock++
		myBlocking = true

	} else { // 존재함

		// 1. 북마크 상태
		if memberBlocking.BlockYn {
			result = ldbMe.Model(&schemas.MemberBlocking{}).
				Where("seq_member = ? AND seq_member_to = ?", userToken.SeqMember, _seqMemberTo).
				Update("block_yn", false)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			myBlocking = false
		} else {
			sql := "UPDATE member_blockings SET block_yn = true, cnt_block = cnt_block + 1 WHERE seq_member = ? AND seq_member_to = ?"
			ldbMe.Exec(sql, userToken.SeqMember, _seqMemberTo)

			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}

			// 상대방과 나와의 구독 정보가 존재 할 시
			if memberSubscribe.SeqMemberSubscribe > 0 {

				switch memberSubscribe.Status {
				// 1. 상대에게 내가 구독을 한 상태
				case define.FOLLOWING:

					// 서로 데이터 삭제
					result = ldbMe.Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMemberTo).
						Delete(&schemas.MemberSubscribe{})
					if corm(result, &res) {
						return res
					}
					result = ldbYour.Where("seq_member = ? AND seq_member_opponent = ?", _seqMemberTo, userToken.SeqMember).
						Delete(&schemas.MemberSubscribe{})
					if corm(result, &res) {
						return res
					}
					result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMemberTo)
					if corm(result, &res) {
						return res
					}

				// 2. 서로 맞 구독 상태
				case define.BOTH:
					// 상대에게 내가 구독을 받은 상태로 변경
					result = ldbMe.Model(&schemas.MemberSubscribe{}).
						Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, _seqMemberTo).
						Update("status", define.FOLLOWER)
					if corm(result, &res) {
						return res
					}
					result = ldbYour.Model(&schemas.MemberSubscribe{}).
						Where("seq_member = ? AND seq_member_opponent = ?", _seqMemberTo, userToken.SeqMember).
						Update("status", define.FOLLOWING)
					if corm(result, &res) {
						return res
					}
					result = mdb.Exec("UPDATE member_details SET cnt_subscribe = cnt_subscribe - 1 WHERE seq_member = ?", _seqMemberTo)
					if corm(result, &res) {
						return res
					}
				}
			}

			cntBlock++
			myBlocking = true
		}
	}

	data := make(map[string]interface{})
	data["my_block"] = myBlocking
	data["cnt_block"] = cntBlock
	res.Data = data

	return res
}
