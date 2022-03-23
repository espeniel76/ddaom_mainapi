package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
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

	isExistImage := false
	_nickName := req.HttpRquest.FormValue("nick_name")
	var profilePhoto domain.FileStructure
	file, handler, err := req.HttpRquest.FormFile("profile_photo")
	if err != nil {
		isExistImage = false
	} else {
		isExistImage = true
		profilePhoto = domain.FileStructure{
			File:        file,
			FileName:    handler.Filename,
			ContentType: handler.Header.Get("Content-Type"),
			Size:        handler.Size,
		}
	}
	fullPath := ""
	if isExistImage {
		fullPath, err = SaveFile("profile", &profilePhoto)
		if err != nil {
			res.ResultCode = define.SYSTEM_ERROR
			res.ErrorDesc = err.Error()
			return res
		}
	}

	masterDB := db.List[define.DSN_MASTER]
	if isExistImage {
		memberDetail := schemas.MemberDetail{
			NickName:     _nickName,
			ProfilePhoto: fullPath,
		}
		result := masterDB.Updates(&memberDetail).Where("seq_member = ?", userToken.SeqMember)
		if result.Error != nil {
			res.ResultCode = define.OK
			res.ErrorDesc = result.Error.Error()
			return res
		}
	} else {
		memberDetail := schemas.MemberDetail{
			NickName: _nickName,
		}
		result := masterDB.Updates(&memberDetail).Where("seq_member = ?", userToken.SeqMember)
		if result.Error != nil {
			res.ResultCode = define.OK
			res.ErrorDesc = result.Error.Error()
			return res
		}
	}

	return res
}
