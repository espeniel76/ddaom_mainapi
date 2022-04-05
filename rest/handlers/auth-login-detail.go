package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"time"
)

func AuthLoginDetail(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_nickName := req.HttpRquest.FormValue("nick_name")
	if len(_nickName) < 1 {
		res.ResultCode = define.BLANK_VALUE
		res.ErrorDesc = "The nickname value is empty."
		return res
	}

	masterDB := db.List[define.DSN_MASTER]
	memberDetail := &schemas.MemberDetail{}
	result := masterDB.Where("nick_name = ?", _nickName).Find(&memberDetail)
	if corm(result, &res) {
		return res
	}

	if memberDetail.SeqMember > 0 {
		if memberDetail.SeqMember != userToken.SeqMember {
			res.ResultCode = define.ALREADY_EXISTS_NICKNAME
			res.ErrorDesc = "Nickname that already exists"
			return res
		}
	}

	member := &schemas.Member{}
	result = masterDB.Find(&member, "email", userToken.Email)
	if corm(result, &res) {
		return res
	}

	result = masterDB.Find(&memberDetail, "seq_member", userToken.SeqMember)
	if corm(result, &res) {
		return res
	}

	if memberDetail.SeqMember > 0 {
		result = masterDB.Model(memberDetail).
			Where("seq_member = ?", userToken.SeqMember).
			Update("nick_name", _nickName)
		if corm(result, &res) {
			return res
		}

	} else {
		memberDetail.SeqMember = userToken.SeqMember
		memberDetail.Email = userToken.Email
		memberDetail.ProfilePhoto = define.DEFAULT_PROFILE
		memberDetail.AuthenticationAt = time.Now()
		memberDetail.NickName = _nickName
		result = masterDB.Create(&memberDetail)
		if corm(result, &res) {
			return res
		}
	}

	return res
}
