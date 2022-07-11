package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func AuthLoginRefresh(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	refreshToken := Cp(req.Parameters, "refresh_token")
	pushToken := Cp(req.Parameters, "push_token")
	verifyResult, err := define.VerifyToken(refreshToken, define.Mconn.JwtRefreshSecret)
	if err != nil {
		res.ResultCode = verifyResult
		res.ErrorDesc = err.Error()
		return res
	}

	// 1. refresh 토큰 정보를 확인
	userToken, err := define.ExtractTokenMetadata(refreshToken, define.Mconn.JwtRefreshSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	go setPushToken(userToken.SeqMember, pushToken)

	// 2. 신규 access_token 과 refresh_token 을 발급한다.
	accessToken, err := define.CreateToken(userToken, define.Mconn.JwtAccessSecret, "ACCESS")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}
	// 3. JWT 토큰 만들기 (refresh)
	refreshToken, err = define.CreateToken(userToken, define.Mconn.JwtRefreshSecret, "REFRESH")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	member := schemas.Member{}
	memberDetail := schemas.MemberDetail{}
	sdb := db.List[define.Mconn.DsnSlave]
	result := sdb.Select("blocked_yn, dormacy_yn").Where("seq_member = ?", userToken.SeqMember).Find(&member)
	if corm(result, &res) {
		return res
	}
	result = sdb.Select("nick_name").Where("seq_member = ?", userToken.SeqMember).Find(&memberDetail)
	if corm(result, &res) {
		return res
	}

	if member.DormacyYn {
		res.ResultCode = define.DORMANCY
		return res
	}

	m := make(map[string]interface{})
	m["access_token"] = accessToken
	m["refresh_token"] = refreshToken
	m["http_server"] = define.Mconn.HTTPServer
	m["nick_name"] = memberDetail.NickName
	m["blocked_yn"] = member.BlockedYn

	res.Data = m

	return res
}
