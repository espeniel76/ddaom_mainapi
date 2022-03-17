package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func AuthLoginDetail(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	var oFile domain.FileStructure
	var fullPath string
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	isExistImage := true
	nickName := req.HttpRquest.FormValue("nick_name")

	if len(nickName) < 1 {
		res.ResultCode = define.BLANK_VALUE
		res.ErrorDesc = "The nickname value is empty."
		return res
	}

	file, handler, err := req.HttpRquest.FormFile("profile_photo")
	if err != nil {
		isExistImage = false
	} else {
		oFile = domain.FileStructure{
			File:        file,
			FileName:    handler.Filename,
			ContentType: handler.Header.Get("Content-Type"),
			Size:        handler.Size,
		}
	}

	masterDB := db.List[define.DSN_MASTER]
	memberDetail := &schemas.MemberDetail{}
	result := masterDB.Where("nick_name = ?", nickName).Find(&memberDetail)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	if memberDetail.SeqMember > 0 {
		if memberDetail.SeqMember != userToken.SeqMember {
			res.ResultCode = define.ALREADY_EXISTS_NICKNAME
			res.ErrorDesc = "Nickname that already exists"
			return res
		}
	}

	if isExistImage {
		fullPath, err = SaveFile("profile", &oFile)
		if err != nil {
			res.ResultCode = define.SYSTEM_ERROR
			res.ErrorDesc = err.Error()
			return res
		}

	} else {
		member := &schemas.Member{}
		result := masterDB.Find(&member, "email", userToken.Email)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		fullPath = member.ProfileImageUrl
	}

	result = masterDB.Find(&memberDetail, "seq_member", userToken.SeqMember)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	if memberDetail.SeqMember > 0 {
		if !isExistImage && memberDetail.ProfilePhoto != "" {
			result = masterDB.Model(memberDetail).
				Where("seq_member = ?", userToken.SeqMember).
				Update("nick_name", nickName)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
		} else {
			result = masterDB.Model(memberDetail).
				Where("seq_member = ?", userToken.SeqMember).
				Updates(schemas.MemberDetail{NickName: nickName, ProfilePhoto: fullPath})
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
		}

	} else {
		result = masterDB.Model(&memberDetail).Create(&schemas.MemberDetail{SeqMember: userToken.SeqMember, NickName: nickName, ProfilePhoto: fullPath})
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
	}

	return res
}
