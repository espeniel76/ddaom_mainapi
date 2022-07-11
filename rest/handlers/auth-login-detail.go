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
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
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

	mdb := db.List[define.Mconn.DsnMaster]
	memberDetail := &schemas.MemberDetail{}
	memberDetailBackup := &schemas.MemberDetailBackup{}

	// 기존 사용중인 닉네임 검색
	result := mdb.Where("nick_name = ?", _nickName).Find(&memberDetail)
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

	// 탈퇴 사용자 닉네임 검색
	result = mdb.Where("nick_name = ?", _nickName).Find(&memberDetailBackup)
	if corm(result, &res) {
		return res
	}
	if memberDetailBackup.SeqMember > 0 {
		res.ResultCode = define.ALREADY_EXISTS_NICKNAME
		res.ErrorDesc = "Nickname that already exists"
		return res
	}

	member := &schemas.Member{}
	result = mdb.Find(&member, "email", userToken.Email)
	if corm(result, &res) {
		return res
	}

	result = mdb.Find(&memberDetail, "seq_member", userToken.SeqMember)
	if corm(result, &res) {
		return res
	}

	if memberDetail.SeqMember > 0 {
		result = mdb.Model(memberDetail).
			Where("seq_member = ?", userToken.SeqMember).
			Update("nick_name", _nickName)
		if corm(result, &res) {
			return res
		}

	} else {
		memberDetail.SeqMember = userToken.SeqMember
		memberDetail.Email = userToken.Email
		memberDetail.ProfilePhoto = define.Mconn.DefaultProfile
		memberDetail.AuthenticationAt = time.Now()
		memberDetail.NickName = _nickName
		memberDetail.DeletedAt = time.Now()
		result = mdb.Create(&memberDetail)
		if corm(result, &res) {
			return res
		}
	}

	go cacheMainPopularWriter()

	return res
}
