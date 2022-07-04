package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Contact struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Email string             `bson:"email,omitempty"`
	Tags  []string           `bson:"tags,omitempty"`
}

func HealthCheck(req *domain.CommonRequest) domain.CommonResponse {
	var res = domain.CommonResponse{}
	res.ResultCode = define.OK
	return res
}

func AuthLogin(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	mdb := db.List[define.Mconn.DsnMaster]

	email := Cp(req.Parameters, "email")
	token := Cp(req.Parameters, "token")
	pushToken := Cp(req.Parameters, "push_token")
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

	// 회원 정보가 있나?
	isExist := db.ExistRow(mdb, "members", "email", email)

	tmpMember := schemas.Member{}
	mdb.Model(schemas.Member{}).Where("email = ?", email).Scan(&tmpMember)

	member := &schemas.Member{
		Email:     email,
		Token:     token,
		SnsType:   snsType,
		PushToken: pushToken,
		DeletedAt: time.Now(),
		BlockedAt: time.Now(),
	}

	// 1. 신규인가
	if tmpMember.SeqMember == 0 {
		fmt.Println("신규회원")
		mdb.Create(member)

		// 2. 탈퇴회원인가
	} else if tmpMember.SeqMember > 0 && tmpMember.DeletedYn {
		fmt.Println("탈퇴회원")
		// mdb.Where("seq_member = ?", tmpMember.SeqMember).Delete(schemas.Member{})
		// mdb.Where("seq_member = ?", tmpMember.SeqMember).Delete(schemas.MemberDetail{})
		mdb.Create(member)

		// 3. 기존회원인가
	} else {
		fmt.Println("기존회원")
		result = mdb.Find(&member, "email", email)
		if corm(result, &res) {
			return res
		}

		if pushToken != "<nil>" && pushToken != "" {
			result = mdb.Model(&member).Where("seq_member = ?", member.SeqMember).Update("push_token", pushToken)
			if corm(result, &res) {
				return res
			}
		}

		memberDetail := schemas.MemberDetail{}
		result = mdb.Select("nick_name").Where("seq_member = ?", member.SeqMember).Find(&memberDetail)
		if corm(result, &res) {
			return res
		}
		nickName = memberDetail.NickName
	}
	if pushToken != "<nil>" && pushToken != "" {
		setPushToken(member.SeqMember, pushToken)
	}

	var myLogDB *gorm.DB
	var allocatedDb int8
	ldb1 := db.List[define.Mconn.DsnLog1Master]
	ldb2 := db.List[define.Mconn.DsnLog2Master]
	if !isExist {
		var count1, count2 int64
		result = ldb1.Table("member_exists").Count(&count1)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		result = ldb2.Table("member_exists").Count(&count2)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		if count1 > count2 {
			myLogDB = ldb2
			allocatedDb = 2
		} else {
			myLogDB = ldb1
			allocatedDb = 1
		}
		result = myLogDB.Create(&schemas.MemberExist{SeqMember: member.SeqMember})
		if corm(result, &res) {
			return res
		}

		result = mdb.Model(&member).Update("allocated_db", allocatedDb).Where("seq_member = ?", member.SeqMember)
		if corm(result, &res) {
			return res
		}
	} else {
		result = mdb.Find(&member, "email", email)
		if corm(result, &res) {
			return res
		}
		allocatedDb = member.AllocatedDb
		if allocatedDb == 1 {
			myLogDB = ldb1
		} else {
			myLogDB = ldb2
		}
	}

	userToken := domain.UserToken{
		Authorized: true,
		SeqMember:  member.SeqMember,
		Email:      email,
		UserLevel:  5,
		Allocated:  allocatedDb,
	}
	accessToken, err := define.CreateToken(&userToken, define.Mconn.JwtAccessSecret, "ACCESS")
	if err != nil {
		res.ResultCode = define.CREATE_TOKEN_ERROR
		res.ErrorDesc = err.Error()
		return res
	}
	refreshToken, err := define.CreateToken(&userToken, define.Mconn.JwtRefreshSecret, "REFRESH")
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
	if corm(result, &res) {
		return res
	}

	m := make(map[string]interface{})
	m["access_token"] = accessToken
	m["refresh_token"] = refreshToken
	m["nick_name"] = nickName
	m["http_server"] = define.Mconn.HTTPServer

	res.Data = m

	return res
}

func setPushToken(seqMember int64, pushToken string) {
	mdb := db.List[define.Mconn.DsnMaster]
	query := `
	INSERT INTO member_push_tokens (seq_member, push_token, created_at, updated_at)
	VALUES (?, ?, NOW(), NOW())
	ON DUPLICATE KEY UPDATE updated_at = NOW()
	`
	mdb.Exec(query, seqMember, pushToken)
}
