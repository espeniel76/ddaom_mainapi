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

	// 아무나 실행 못 시키는 장치가 필요하다...

	mdb := db.List[define.Mconn.DsnMaster]
	ldb1 := db.List[define.Mconn.DsnLog1Master]
	ldb2 := db.List[define.Mconn.DsnLog2Master]

	mdb.AutoMigrate(schemas.Member{})
	if result != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error()
		return res
	}

	mdb.AutoMigrate(schemas.Member{})
	mdb.AutoMigrate(schemas.MemberBlock{})
	mdb.AutoMigrate(schemas.MemberDormacy{})
	mdb.AutoMigrate(schemas.MemberDetail{})
	mdb.AutoMigrate(schemas.MemberPushToken{})
	mdb.AutoMigrate(schemas.MemberLog{})

	mdb.AutoMigrate(schemas.MemberBackup{})
	mdb.AutoMigrate(schemas.MemberDetailBackup{})

	ldb1.AutoMigrate(schemas.MemberLoginLog{})
	ldb2.AutoMigrate(schemas.MemberLoginLog{})

	mdb.AutoMigrate(schemas.MemberAdmin{})
	mdb.AutoMigrate(schemas.MemberAdminLoginLog{})

	mdb.AutoMigrate(schemas.Keyword{})
	mdb.AutoMigrate(schemas.KeywordToday{})
	mdb.AutoMigrate(schemas.Genre{})
	mdb.AutoMigrate(schemas.Image{})
	mdb.AutoMigrate(schemas.Color{})
	mdb.AutoMigrate(schemas.Slang{})

	mdb.AutoMigrate(schemas.NovelStep1{})
	mdb.AutoMigrate(schemas.NovelStep2{})
	mdb.AutoMigrate(schemas.NovelStep3{})
	mdb.AutoMigrate(schemas.NovelStep4{})

	mdb.AutoMigrate(schemas.ServiceInquiry{})
	mdb.AutoMigrate(schemas.Notice{})

	ldb1.AutoMigrate(schemas.MemberSubscribe{})
	ldb2.AutoMigrate(schemas.MemberSubscribe{})
	ldb1.AutoMigrate(schemas.MemberBookmark{})
	ldb2.AutoMigrate(schemas.MemberBookmark{})

	ldb1.AutoMigrate(schemas.MemberLikeStep1{})
	ldb2.AutoMigrate(schemas.MemberLikeStep1{})
	ldb1.AutoMigrate(schemas.MemberLikeStep2{})
	ldb2.AutoMigrate(schemas.MemberLikeStep2{})
	ldb1.AutoMigrate(schemas.MemberLikeStep3{})
	ldb2.AutoMigrate(schemas.MemberLikeStep3{})
	ldb1.AutoMigrate(schemas.MemberLikeStep4{})
	ldb2.AutoMigrate(schemas.MemberLikeStep4{})

	mdb.AutoMigrate(schemas.NovelFinish{})
	mdb.AutoMigrate(schemas.CategoryFaq{})
	mdb.AutoMigrate(schemas.Faq{})

	mdb.AutoMigrate(schemas.KeywordChoiceFirst{})
	mdb.AutoMigrate(schemas.KeywordChoiceSecond{})
	mdb.AutoMigrate(schemas.NovelFinishBatchRunLog{})

	mdb.AutoMigrate(schemas.NovelDelete{})

	mdb.AutoMigrate(schemas.KeywordAlarmLog{})
	mdb.AutoMigrate(schemas.Alarm{})
	mdb.AutoMigrate(schemas.NovelReport{})

	mdb.AutoMigrate(schemas.MemberReport{})
	ldb1.AutoMigrate(schemas.MemberBlocking{})
	ldb2.AutoMigrate(schemas.MemberBlocking{})

	mdb.AutoMigrate(schemas.NovelReply{})
	mdb.AutoMigrate(schemas.NovelReReply{})

	return res
}
