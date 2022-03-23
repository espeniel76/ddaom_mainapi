package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"time"

	"gorm.io/gorm"
)

func AuthLogin(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	masterDB := db.List[define.DSN_MASTER]

	email := Cp(req.Parameters, "email")
	token := Cp(req.Parameters, "token")
	snsType := Cp(req.Parameters, "sns_type")
	var result *gorm.DB
	nickName := ""

	if email == "<nil>" || len(email) == 0 {
		res.ResultCode = define.NO_PARAMETER
		res.ErrorDesc = "no date email"
		return res
	}
	if token == "<nil>" || len(token) == 0 {
		res.ResultCode = define.NO_PARAMETER
		res.ErrorDesc = "no date token"
		return res
	}
	if snsType == "<nil>" || len(snsType) == 0 {
		res.ResultCode = define.NO_PARAMETER
		res.ErrorDesc = "no date snsType"
		return res
	}

	isExist := db.ExistRow(masterDB, "members", "email", email)
	member := &schemas.Member{
		Email:   email,
		Token:   token,
		SnsType: snsType,
	}
	if !isExist {
		result = masterDB.Create(member)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
	} else {
		result = masterDB.Find(&member, "email", email)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		memberDetail := schemas.MemberDetail{}
		masterDB.Select("nick_name").Where("seq_member = ?", member.SeqMember).Find(&memberDetail)
		nickName = memberDetail.NickName
	}

	var myLogDB *gorm.DB
	var allocatedDb int8
	logDB1 := db.List[define.DSN_LOG1]
	logDB2 := db.List[define.DSN_LOG2]
	if !isExist {
		var count1, count2 int64
		result = logDB1.Table("member_exists").Count(&count1)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		result = logDB2.Table("member_exists").Count(&count2)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		if count1 > count2 {
			myLogDB = logDB2
			allocatedDb = 2
		} else {
			myLogDB = logDB1
			allocatedDb = 1
		}
		result = myLogDB.Create(&schemas.MemberExist{SeqMember: member.SeqMember})
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		result = masterDB.Model(&member).Update("allocated_db", allocatedDb)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
	} else {
		result = masterDB.Find(&member, "email", email)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		allocatedDb = member.AllocatedDb
		if allocatedDb == 1 {
			myLogDB = logDB1
		} else {
			myLogDB = logDB2
		}
	}

	userToken := domain.UserToken{
		Authorized: true,
		SeqMember:  member.SeqMember,
		Email:      email,
		UserLevel:  5,
		Allocated:  allocatedDb,
	}
	accessToken, err := define.CreateToken(&userToken, define.JWT_ACCESS_SECRET, "ACCESS")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}
	refreshToken, err := define.CreateToken(&userToken, define.JWT_REFRESH_SECRET, "REFRESH")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	result = myLogDB.Create(&schemas.MemberLoginLog{
		SeqMember: member.SeqMember,
		Token:     refreshToken,
		LoginAt:   time.Now(),
	})
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	m := make(map[string]interface{})
	m["access_token"] = accessToken
	m["refresh_token"] = refreshToken
	m["nick_name"] = nickName
	m["http_server"] = define.HTTP_SERVER

	res.Data = m

	return res
}
