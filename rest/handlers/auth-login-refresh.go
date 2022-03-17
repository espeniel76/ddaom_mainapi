package handlers

import (
	"ddaom/define"
	"ddaom/domain"
)

func AuthLoginRefresh(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	refreshToken := Cp(req.Parameters, "refresh_token")
	verifyResult, err := define.VerifyToken(refreshToken, define.JWT_REFRESH_SECRET)
	if err != nil {
		res.ResultCode = verifyResult
		res.ErrorDesc = err.Error()
		return res
	}

	// 1. refresh 토큰 정보를 확인
	userToken, err := define.ExtractTokenMetadata(refreshToken, define.JWT_REFRESH_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	// 2. 신규 access_token 과 refresh_token 을 발급한다.
	accessToken, err := define.CreateToken(userToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}
	// 3. JWT 토큰 만들기 (refresh)
	refreshToken, err = define.CreateToken(userToken, define.JWT_REFRESH_SECRET)
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	m := make(map[string]string)
	m["access_token"] = accessToken
	m["refresh_token"] = refreshToken

	res.Data = m

	return res
}
