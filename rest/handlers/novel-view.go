package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
)

func NovelView(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep, _ := req.Vars["seq_novel_step1"]
	fmt.Println(_seqNovelStep, userToken)

	// 1단계 소설 가져오기
	// query 기반 map interface 기반으로 다 변경
	masterDB := db.List[define.DSN_MASTER]
	novelStep1 := schemas.NovelStep1{}
	result := masterDB.Model(&novelStep1).Where("active_yn = true AND seq_novel_step1 = ?", _seqNovelStep).First(&novelStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	data := make(map[string]interface{})
	data["title"] = novelStep1.Title
	data["created_at"] = novelStep1.CreatedAt.UnixMilli()
	data["seq_genre"] = novelStep1.SeqGenre
	data["seq_keyword"] = novelStep1.SeqKeyword

	// 2단계 소설 가져오기
	// novelStep2 := schemas.NovelStep2{}
	// result = masterDB.Model(&novelStep2).Where("active_yn = true AND seq_novel_step1 = ?", _seqNovelStep).Last(&novelStep2)
	// if result.Error != nil {
	// 	res.ResultCode = define.DB_ERROR_ORM
	// 	res.ErrorDesc = result.Error.Error()
	// 	return res
	// }

	res.Data = data

	return res
}
