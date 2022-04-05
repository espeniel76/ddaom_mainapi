package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

func MypageListFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	masterDB := db.List[define.DSN_MASTER]
	result := masterDB.Model(schemas.NovelFinish{}).Where("active_yn = true").Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelListFinishRes := NovelListFinishRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	var query bytes.Buffer
	query.WriteString(`
	SELECT
		nf.seq_novel_finish,
		ns.seq_genre,
		ns.seq_image,
		ns.seq_color,
		ns.title,
		true AS my_bookmark
	FROM novel_finishes nf
	INNER JOIN novel_step1 ns ON nf.seq_novel_step1 = ns.seq_novel_step1
	WHERE nf.active_yn = true
	AND nf.seq_member_step1 = ?
		OR nf.seq_member_step2 = ?
		OR nf.seq_member_step3 = ?
		OR nf.seq_member_step4 = ?
	`)
	query.WriteString(" ORDER BY nf.seq_novel_finish DESC")
	query.WriteString(" LIMIT ?, ?")
	result = masterDB.Raw(query.String(), userToken.SeqMember, userToken.SeqMember, userToken.SeqMember, userToken.SeqMember, limitStart, _sizePerPage).Find(&novelListFinishRes.List)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	res.Data = novelListFinishRes

	return res
}

// type NovelListFinishRes struct {
// 	NowPage   int `json:"now_page"`
// 	TotalPage int `json:"total_page"`
// 	TotalData int `json:"total_data"`
// 	List      []struct {
// 		SeqNovelFinish int64  `json:"seq_novel_finish"`
// 		SeqGenre       int64  `json:"seq_genre"`
// 		SeqImage       int64  `json:"seq_image"`
// 		SeqColor       int64  `json:"seq_color"`
// 		Title          string `json:"title"`
// 		MyBookmark     bool   `json:"my_bookmark"`
// 	} `json:"list"`
// }
