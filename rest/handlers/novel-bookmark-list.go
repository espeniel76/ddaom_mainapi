package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
	"strings"
)

func NovelBookmarkList(req *domain.CommonRequest) domain.CommonResponse {

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

	masterDB := db.List[define.DSN_MASTER]
	myLogDb := GetMyLogDb(userToken.Allocated)

	var list []int64
	result := myLogDb.
		Raw("SELECT seq_novel_finish FROM member_bookmarks mb WHERE seq_member = ?", userToken.SeqMember).
		Scan(&list)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	var totalData int64
	result = masterDB.
		Raw("SELECT ns.seq_novel_finish FROM novel_finishes ns WHERE ns.active_yn = true AND seq_novel_finish IN (?)", list).
		Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}
	novelBookmarkListRes := NovelBookmarkListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	var query bytes.Buffer
	query.WriteString(`
		SELECT
			nf.seq_novel_finish,
			ns.seq_keyword,
			ns.seq_genre,
			ns.seq_image,
			ns.seq_color,
			ns.title
		FROM novel_finishes nf 
		INNER JOIN novel_step1 ns ON nf.seq_novel_step1 = ns.seq_novel_step1
		WHERE ns.active_yn = true AND nf.seq_novel_finish IN (?)`)
	query.WriteString(" ORDER BY nf.created_at DESC")
	query.WriteString(" LIMIT ?, ?")
	result = masterDB.Raw(query.String(), list, limitStart, _sizePerPage).Find(&novelBookmarkListRes.List)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	res.Data = novelBookmarkListRes

	return res
}

type NovelBookmarkListRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNovelFinish int64  `json:"seq_novel_finish"`
		SeqKeyword     int64  `json:"seq_keyword"`
		SeqGenre       int64  `json:"seq_genre"`
		SeqImage       int64  `json:"seq_image"`
		SeqColor       int64  `json:"seq_color"`
		Title          string `json:"title"`
	} `json:"list"`
}

func NovelBookmarkDelete(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	tmp := fmt.Sprintf("%v", req.Parameters["seq_novel_finish_list"])
	tmp = strings.ReplaceAll(tmp, "[", "")
	tmp = strings.ReplaceAll(tmp, "]", "")
	tmpList := strings.Split(tmp, " ")
	myLogDb := GetMyLogDb(userToken.Allocated)
	result := myLogDb.
		Where("seq_member = ? AND seq_novel_finish IN (?)", userToken.SeqMember, tmpList).
		Delete(&schemas.MemberBookmark{})
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	return res
}
