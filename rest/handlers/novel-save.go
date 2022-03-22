package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func NovelWriteStep1(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqKeyword := CpInt64(req.Parameters, "seq_keyword")
	_seqGenre := CpInt64(req.Parameters, "seq_genre")
	_seqImage := CpInt64(req.Parameters, "seq_image")
	_seqColor := CpInt64(req.Parameters, "seq_color")
	_title := Cp(req.Parameters, "title")
	_content := Cp(req.Parameters, "content")

	masterDB := db.List[define.DSN_MASTER]
	novelWriteStep1 := schemas.NovelStep1{
		SeqKeyword: _seqKeyword,
		SeqImage:   _seqImage,
		SeqColor:   _seqColor,
		SeqGenre:   _seqGenre,
		SeqMember:  userToken.SeqMember,
		Title:      _title,
		Content:    _content,
	}

	// 동일 제목 검사
	var cnt int64
	result := masterDB.Model(&novelWriteStep1).Where("title = ?", _title).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt > 0 {
		res.ResultCode = define.ALREADY_EXISTS_TITLE
		return res
	}

	result = masterDB.Model(&novelWriteStep1).Create(&novelWriteStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	return res
}

func NovelWriteStep2(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_content := Cp(req.Parameters, "content")

	masterDB := db.List[define.DSN_MASTER]

	var cnt int64
	result := masterDB.Model(schemas.NovelStep1{}).Where("seq_novel_step1 = ?", _seqNovelStep1).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		res.ErrorDesc = result.Error.Error()
		return res
	}

	result = masterDB.Exec("UPDATE novel_step1 SET cnt_step2 = cnt_step2 + 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelStep2 := schemas.NovelStep2{
		SeqNovelStep1: _seqNovelStep1,
		SeqMember:     userToken.SeqMember,
		Content:       _content,
	}
	result = masterDB.Save(&novelStep2)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	return res
}

func NovelWriteStep3(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_content := Cp(req.Parameters, "content")

	masterDB := db.List[define.DSN_MASTER]
	var cnt int64
	result := masterDB.Model(schemas.NovelStep2{}).Where("seq_novel_step2 = ?", _seqNovelStep2).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		res.ErrorDesc = result.Error.Error()
		return res
	}

	// 1단계 seq 알아내기
	var seqNovelStep1 int64
	result = masterDB.Model(schemas.NovelStep2{}).Where("seq_novel_step2 = ?", _seqNovelStep2).Pluck("seq_novel_step1", &seqNovelStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	result = masterDB.Exec("UPDATE novel_step1 SET cnt_step3 = cnt_step3 + 1 WHERE seq_novel_step1 = ?", seqNovelStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	result = masterDB.Exec("UPDATE novel_step2 SET cnt_step3 = cnt_step3 + 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	// 3단계 저장
	novelStep3 := schemas.NovelStep3{
		SeqNovelStep2: _seqNovelStep2,
		SeqMember:     userToken.SeqMember,
		Content:       _content,
	}
	result = masterDB.Save(&novelStep3)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	return res
}

func NovelWriteStep4(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_content := Cp(req.Parameters, "content")

	masterDB := db.List[define.DSN_MASTER]
	var cnt int64
	result := masterDB.Model(schemas.NovelStep3{}).Where("seq_novel_step3 = ?", _seqNovelStep3).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		res.ErrorDesc = result.Error.Error()
		return res
	}

	// 1/2단계 seq 알아내기
	var seqNovelStep1 int64
	var seqNovelStep2 int64
	result = masterDB.Model(schemas.NovelStep3{}).Where("seq_novel_step3 = ?", _seqNovelStep3).Pluck("seq_novel_step2", &seqNovelStep2)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	result = masterDB.Model(schemas.NovelStep2{}).Where("seq_novel_step2 = ?", seqNovelStep2).Pluck("seq_novel_step1", &seqNovelStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	result = masterDB.Exec("UPDATE novel_step1 SET cnt_step4 = cnt_step4 + 1 WHERE seq_novel_step1 = ?", seqNovelStep1)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	result = masterDB.Exec("UPDATE novel_step2 SET cnt_step4 = cnt_step4 + 1 WHERE seq_novel_step2 = ?", seqNovelStep2)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	result = masterDB.Exec("UPDATE novel_step3 SET cnt_step4 = cnt_step4 + 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	// 4단계 저장
	novelStep4 := schemas.NovelStep4{
		SeqNovelStep3: _seqNovelStep3,
		SeqMember:     userToken.SeqMember,
		Content:       _content,
	}
	result = masterDB.Save(&novelStep4)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	return res
}
