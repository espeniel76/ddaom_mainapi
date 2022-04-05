package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"time"
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
	_email := req.HttpRquest.FormValue("email")
	_isDefaultImage := req.HttpRquest.FormValue("is_default_image")
	file, handler, err := req.HttpRquest.FormFile("profile_photo")
	fullPath := ""

	if _isDefaultImage == "N" {
		var profilePhoto domain.FileStructure
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

		if isExistImage {
			fullPath, err = SaveFile("profile", &profilePhoto)
			if err != nil {
				res.ResultCode = define.SYSTEM_ERROR
				res.ErrorDesc = err.Error()
				return res
			}
		}
	} else {
		fullPath = define.DEFAULT_PROFILE
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

	result = masterDB.Model(&memberDetail).Where("seq_member = ?", userToken.SeqMember).Scan(&memberDetail)
	if corm(result, &res) {
		return res
	}
	isExistMember := false
	if memberDetail.SeqMember > 0 {
		isExistMember = true
	}
	if len(_nickName) > 0 {
		memberDetail.NickName = _nickName
	}
	if _isDefaultImage == "Y" {
		memberDetail.ProfilePhoto = define.DEFAULT_PROFILE
	} else {
		if isExistImage {
			memberDetail.ProfilePhoto = fullPath
		}
	}
	if len(_email) > 0 {
		memberDetail.Email = _email
	}
	if isExistMember {
		result = masterDB.Model(&memberDetail).
			Where("seq_member = ?", userToken.SeqMember).
			Updates(&memberDetail)
		if corm(result, &res) {
			return res
		}
	} else {
		memberDetail.SeqMember = userToken.SeqMember
		memberDetail.Email = userToken.Email
		if isExistImage {
			memberDetail.ProfilePhoto = fullPath
		} else {
			memberDetail.ProfilePhoto = define.DEFAULT_PROFILE
		}
		memberDetail.AuthenticationAt = time.Now()
		result = masterDB.Create(&memberDetail)
		if corm(result, &res) {
			return res
		}
	}

	return res
}
