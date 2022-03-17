package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func InitialDb(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	var result error

	masterDB := db.List[define.DSN_MASTER]
	logDB1 := db.List[define.DSN_LOG1]
	logDB2 := db.List[define.DSN_LOG2]

	masterDB.AutoMigrate(schemas.Member{})
	if result != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error()
		return res
	}

	masterDB.AutoMigrate(schemas.MemberDetail{})

	logDB1.AutoMigrate(schemas.MemberExist{})
	logDB1.AutoMigrate(schemas.MemberLoginLog{})
	logDB2.AutoMigrate(schemas.MemberExist{})
	logDB2.AutoMigrate(schemas.MemberLoginLog{})

	masterDB.AutoMigrate(schemas.MemberAdmin{})
	masterDB.AutoMigrate(schemas.MemberAdminLoginLog{})

	return res
}
