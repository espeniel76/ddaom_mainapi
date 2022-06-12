package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"time"
)

func AuthAuthentication(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	memberDetail := schemas.MemberDetail{}
	masterDB := db.List[define.DSN_MASTER]
	masterDB.
		Model(schemas.MemberDetail{}).
		Where("seq_member", userToken.SeqMember).
		Find(&memberDetail)
	isAuthentication := false
	var authenticationAt int64
	if len(memberDetail.AuthenticationCi) > 0 {
		isAuthentication = true
		authenticationAt = memberDetail.AuthenticationAt.UnixMilli()
	}
	data := make(map[string]interface{})
	data["is_authentication"] = isAuthentication
	data["authentication_at"] = authenticationAt
	data["authentication_ci"] = memberDetail.AuthenticationCi
	res.Data = data

	return res
}

func AuthAuthenticationSet(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_authenticationCi := Cp(req.Parameters, "authentication_ci")

	memberDetail := schemas.MemberDetail{
		AuthenticationCi: _authenticationCi,
		AuthenticationAt: time.Now(),
	}
	masterDB := db.List[define.DSN_MASTER]
	result := masterDB.
		Model(schemas.MemberDetail{}).
		Where("seq_member", userToken.SeqMember).
		Updates(&memberDetail)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	data := make(map[string]interface{})
	data["is_authentication"] = true
	data["authentication_at"] = time.Now().UnixMilli()
	data["authentication_ci"] = _authenticationCi
	res.Data = data

	return res
}
