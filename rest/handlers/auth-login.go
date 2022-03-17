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
	existProfilePhoto := false

	// 파라미터 검증
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

	// 1. 최초 로그인 여부
	isExist := db.ExistRow(masterDB, "donuts.members", "email", email)
	member := &schemas.Member{
		Email:   email,
		Token:   token,
		SnsType: snsType,
	}
	if !isExist {
		// 1.1. 없으면, 사용자 추가
		result = masterDB.Create(member)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		existProfilePhoto = false
	} else {
		result = masterDB.Find(&member, "email", email)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		memberDetail := schemas.MemberDetail{}
		masterDB.Select("profile_photo").Where("seq_member = ?", member.SeqMember).Find(&memberDetail)
		if memberDetail.ProfilePhoto == "" {
			existProfilePhoto = false
		} else {
			existProfilePhoto = true
		}
	}

	// 2. 로그인 로그 남기기
	// 2.1. 각 LOG DB 사용자 카운트
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

		// 2.2. 데이터가 제일 작은 녀석에 데이터를 생성한다.
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

		// 2.3. 내 로그가 어느 DB 에 들어 있는지 기록한다.
		result = masterDB.Model(&member).Update("allocated_db", allocatedDb)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
	} else {
		// 2.4. 최초 로그인이 아니면, 내가 속한 DB 가져온다.
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

	// 3. JWT 토큰 만들기 (access)
	userToken := domain.UserToken{
		Authorized: true,
		SeqMember:  member.SeqMember,
		Email:      email,
		UserLevel:  5,
		Allocated:  allocatedDb,
	}
	accessToken, err := define.CreateToken(&userToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}
	// 3. JWT 토큰 만들기 (refresh)
	refreshToken, err := define.CreateToken(&userToken, define.JWT_REFRESH_SECRET)
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	// 4. login log 기록하기
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
	m["exist_profile_photo"] = existProfilePhoto
	m["http_server"] = define.HTTP_SERVER

	res.Data = m

	return res
}
