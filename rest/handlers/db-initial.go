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

	masterDB.AutoMigrate(schemas.Member{})
	masterDB.AutoMigrate(schemas.MemberDetail{})

	logDB1.AutoMigrate(schemas.MemberExist{})
	logDB1.AutoMigrate(schemas.MemberLoginLog{})
	logDB2.AutoMigrate(schemas.MemberExist{})
	logDB2.AutoMigrate(schemas.MemberLoginLog{})

	masterDB.AutoMigrate(schemas.MemberAdmin{})
	masterDB.AutoMigrate(schemas.MemberAdminLoginLog{})

	masterDB.AutoMigrate(schemas.Keyword{})
	masterDB.AutoMigrate(schemas.KeywordToday{})
	masterDB.AutoMigrate(schemas.Genre{})
	masterDB.AutoMigrate(schemas.Image{})
	masterDB.AutoMigrate(schemas.Color{})
	masterDB.AutoMigrate(schemas.Slang{})

	masterDB.AutoMigrate(schemas.NovelStep1{})
	masterDB.AutoMigrate(schemas.NovelStep2{})
	masterDB.AutoMigrate(schemas.NovelStep3{})
	masterDB.AutoMigrate(schemas.NovelStep4{})

	logDB1.AutoMigrate(schemas.MemberSubscribe{})
	logDB2.AutoMigrate(schemas.MemberSubscribe{})
	logDB1.AutoMigrate(schemas.MemberBookmark{})
	logDB2.AutoMigrate(schemas.MemberBookmark{})
	logDB1.AutoMigrate(schemas.MemberLike{})
	logDB2.AutoMigrate(schemas.MemberLike{})

	return res
}
