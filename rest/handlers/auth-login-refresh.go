package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func AuthLoginRefresh(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_refreshToken := Cp(req.Parameters, "refresh_token")
	verifyResult, err := define.VerifyToken(_refreshToken, define.Mconn.JwtRefreshSecret)
	if err != nil {
		res.ResultCode = verifyResult
		res.ErrorDesc = err.Error()
		return res
	}

	// 1. refresh 토큰 정보를 확인
	userToken, err := define.ExtractTokenMetadata(_refreshToken, define.Mconn.JwtRefreshSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	// 토큰 까서 유효한 사용자 여부 검사
	member := schemas.Member{}
	memberDetail := schemas.MemberDetail{}
	sdb := db.List[define.Mconn.DsnSlave]
	result := sdb.Model(&member).Where("email = ?", userToken.Email).Scan(&member)
	if corm(result, &res) {
		return res
	}

	// 존재하지 않는 사용자
	if member.SeqMember == 0 {
		res.ResultCode = define.NO_EXIST_USER
		return res
	}

	// 존재하나, 시퀀스가 맞지않다 (탈퇴)
	if member.SeqMember != userToken.SeqMember {
		res.ResultCode = define.WITHDRAWAL_USER
		return res
	}

	result = sdb.Model(&memberDetail).Where("seq_member = ?", userToken.SeqMember).Find(&memberDetail)
	if corm(result, &res) {
		return res
	}

	// 2. 신규 access_token 과 refresh_token 을 발급한다.
	accessToken, err := define.CreateToken(userToken, define.Mconn.JwtAccessSecret, "ACCESS")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}
	// 3. JWT 토큰 만들기 (refresh)
	_refreshToken, err = define.CreateToken(userToken, define.Mconn.JwtRefreshSecret, "REFRESH")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	if member.DormacyYn {
		res.ResultCode = define.DORMANCY
		return res
	}

	m := make(map[string]interface{})
	m["access_token"] = accessToken
	m["refresh_token"] = _refreshToken
	m["http_server"] = define.Mconn.HTTPServer
	m["nick_name"] = memberDetail.NickName
	m["blocked_yn"] = member.BlockedYn
	m["seq_member"] = member.SeqMember

	res.Data = m

	// 로그인 완료 로그
	go setUserActionLog(userToken.SeqMember, 2, "")

	return res
}
