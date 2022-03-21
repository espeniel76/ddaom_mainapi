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

	// fmt.Println(_seqKeyword, _seqGenre, _seqImage, _seqColor, _title, _content)

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
