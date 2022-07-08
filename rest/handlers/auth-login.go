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
		// } else if tmpMember.SeqMember > 0 && tmpMember.DeletedYn {
		// 	fmt.Println("탈퇴회원")
		// 	// mdb.Where("seq_member = ?", tmpMember.SeqMember).Delete(schemas.Member{})
		// 	// mdb.Where("seq_member = ?", tmpMember.SeqMember).Delete(schemas.MemberDetail{})

		// 	// 블록 상태인가?
		// 	if member.BlockedYn {
		// 		// 블록 상태가 된지 90일이 지났는가? (이건 배치를 믿어본다)
		// 		res.ResultCode = define.UJDTBLOCKED
		// 		return res
		// 	}

		// 	mdb.Create(member)

		// 2. 탈퇴 여부 확인
		memberBackup := schemas.MemberBackup{}
		query := "SELECT * FROM member_backups WHERE email = ? ORDER BY created_at DESC LIMIT 1"
		result = mdb.Raw(query, email).Scan(&memberBackup)
		if corm(result, &res) {
			return res
		}
		// 2.1. 데이터가 있다
		if memberBackup.SeqMember > 0 {
			fmt.Println("탈퇴이력 있음")
			// members 에는 데이터가 없고, member_backups 에는 있다 면 탈퇴 후 재 가입으로 간주
			// seq_member 로, 블랙 처리 여부 판단
			blockedYn := false
			result = mdb.Model(schemas.Member{}).Select("blocked_yn").Where("seq_member = ?", memberBackup.SeqMember).Scan(&blockedYn)
			if corm(result, &res) {
				return res
			}
			if blockedYn {
				fmt.Println("블랙인데 탈퇴한 유저임, 가입 안됨")
				res.ResultCode = define.UJDTBLOCKED
				return res
			}
		}

		// 3. 기존회원인가
	} else {
		fmt.Println("기존회원")
		result = mdb.Find(&member, "email", email)
		if corm(result, &res) {
			return res
		}

		// 휴면 상태 회원인가
		if member.DormacyYn {
			res.ResultCode = define.DORMANCY
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

	go setPushToken(member.SeqMember, pushToken)

	var myLogDB *gorm.DB
	var allocatedDb int8
	var lastAllocatedDb int8
	ldb1 := db.List[define.Mconn.DsnLog1Master]
	ldb2 := db.List[define.Mconn.DsnLog2Master]
	if !isExist {
		// var count1, count2 int64
		// ldb1.Table("member_exists").Count(&count1)
		// ldb2.Table("member_exists").Count(&count2)
		// if count1 > count2 {
		// 	myLogDB = ldb2
		// 	allocatedDb = 2
		// } else {
		// 	myLogDB = ldb1
		// 	allocatedDb = 1
		// }
		// myLogDB.Create(&schemas.MemberExist{SeqMember: member.SeqMember})

		mdb.Raw("SELECT allocated_db FROM members ORDER BY seq_member DESC LIMIT 1").Scan(&lastAllocatedDb)
		if lastAllocatedDb == 1 {
			allocatedDb = 2
			myLogDB = ldb2
		} else {
			allocatedDb = 1
			myLogDB = ldb1
		}

		mdb.Model(&member).Update("allocated_db", allocatedDb).Where("seq_member = ?", member.SeqMember)
	} else {
		result = mdb.Find(&member, "email", email)
		if corm(result, &res) {
			return res
		}
		allocatedDb = member.AllocatedDb
		if allocatedDb == 1 {
			myLogDB = ldb1
		} else if allocatedDb == 2 {
			myLogDB = ldb2
		} else {
			// 1,2 다 아니면 다시 할당 한다
			// 원래 이 상황이 생기면 안 되는데, 없으면 안 되므로
			// 에러방지용
			mdb.Raw("SELECT allocated_db FROM members ORDER BY seq_member DESC LIMIT 1").Scan(&lastAllocatedDb)
			if lastAllocatedDb == 1 {
				allocatedDb = 2
				myLogDB = ldb2
			} else {
				allocatedDb = 1
				myLogDB = ldb1
			}
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
	m["blocked_yn"] = member.BlockedYn

	res.Data = m

	return res
}

func setPushToken(seqMember int64, pushToken string) {
	if len(pushToken) > 100 && seqMember > 0 {
		mdb := db.List[define.Mconn.DsnMaster]
		mdb.Model(schemas.Member{}).Where("seq_member = ?", seqMember).Update("push_token", pushToken)
		query := `
			INSERT INTO member_push_tokens (seq_member, push_token, created_at, updated_at)
			VALUES (?, ?, NOW(), NOW())
			ON DUPLICATE KEY UPDATE updated_at = NOW()
		`
		mdb.Exec(query, seqMember, pushToken)
	}
}
