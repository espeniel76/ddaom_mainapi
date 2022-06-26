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

	mdb := db.List[define.Mconn.DsnSlave]
	if pushToken != "<nil>" && pushToken != "" {
		mdb.Model(schemas.Member{}).Where("seq_member = ?", userToken.SeqMember).Update("push_token", pushToken)
		setPushToken(userToken.SeqMember, pushToken)
	}

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

	memberDetail := schemas.MemberDetail{}
	result := mdb.Select("nick_name").Where("seq_member = ?", userToken.SeqMember).Find(&memberDetail)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	m := make(map[string]interface{})
	m["access_token"] = accessToken
	m["refresh_token"] = refreshToken
	m["http_server"] = define.Mconn.HTTPServer
	m["nick_name"] = memberDetail.NickName

	res.Data = m

	return res
}
