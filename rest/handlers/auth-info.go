package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

type MemberDetailRes struct {
	Name          string `json:"name"`
	MobileCompany int8   `json:"mobile_company"`
	Mobile        string `json:"mobile"`
	Zipcode       string `json:"zipcode"`
	Address       string `json:"address"`
	AddressDetail string `json:"address_detail"`
	Email         string `json:"email"`
	NickName      string `json:"nick_name"`
	SnsType       string `json:"sns_type"`
}

func AuthInfo(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	memberDetailRes := MemberDetailRes{}
	masterDB := db.List[define.DSN_MASTER]
	query := `
	SELECT
		md.name,
		md.mobile_company,
		md.mobile,
		md.zipcode,
		md.address,
		md.address_detail,
		md.email,
		md.nick_name,
		m.sns_type 
	FROM
		members m
	INNER JOIN
		member_details md ON m.seq_member = md.seq_member
	WHERE
		m.seq_member = ?
	`
	masterDB.Raw(query, userToken.SeqMember).Scan(&memberDetailRes)
	res.Data = memberDetailRes

	return res
}

func AuthInfoUpdate(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_name := Cp(req.Parameters, "name")
	_mobileCompany, _ := strconv.ParseInt(Cp(req.Parameters, "mobile_company"), 10, 64)
	_mobile := Cp(req.Parameters, "mobile")
	_zipcode := Cp(req.Parameters, "zipcode")
	_address := Cp(req.Parameters, "address")
	_addressDetail := Cp(req.Parameters, "address_detail")
	_email := Cp(req.Parameters, "email")

	memberDetail := schemas.MemberDetail{
		Name:          _name,
		MobileCompany: int8(_mobileCompany),
		Mobile:        _mobile,
		Zipcode:       _zipcode,
		Address:       _address,
		AddressDetail: _addressDetail,
		Email:         _email,
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

	return res
}
