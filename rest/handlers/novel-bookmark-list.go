package handlers

import (
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

	sdb := db.List[define.DSN_SLAVE]
	ldb := GetMyLogDbSlave(userToken.Allocated)

	var list []int64
	result := ldb.
		Model(schemas.MemberBookmark{}).
		Select("seq_novel_finish").
		Where("seq_member = ? AND bookmark_yn = true", userToken.SeqMember).
		Scan(&list)
	if corm(result, &res) {
		return res
	}

	var totalData int64
	result = sdb.
		Model(schemas.NovelFinish{}).
		Where("seq_novel_finish IN (?)", list).
		Count(&totalData)
	if corm(result, &res) {
		return res
	}
	novelBookmarkListRes := NovelBookmarkListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	query := `
	SELECT
		nf.seq_novel_finish,
		ns.seq_keyword,
		ns.seq_genre,
		ns.seq_image,
		ns.seq_color,
		ns.title
	FROM novel_finishes nf 
	INNER JOIN novel_step1 ns ON nf.seq_novel_step1 = ns.seq_novel_step1
	WHERE ns.active_yn = true AND nf.seq_novel_finish IN (?)
	ORDER BY nf.created_at DESC
	LIMIT ?, ?
	`
	result = sdb.
		Raw(query, list, limitStart, _sizePerPage).
		Scan(&novelBookmarkListRes.List)
	if corm(result, &res) {
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
		SeqNovelFinish int    `json:"seq_novel_finish"`
		SeqKeyword     int    `json:"seq_keyword"`
		SeqGenre       int    `json:"seq_genre"`
		SeqImage       int    `json:"seq_image"`
		SeqColor       int    `json:"seq_color"`
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
	ldb := GetMyLogDbMaster(userToken.Allocated)
	result := ldb.
		Where("seq_member = ? AND seq_novel_finish IN (?)", userToken.SeqMember, tmpList).
		Delete(&schemas.MemberBookmark{})
	if corm(result, &res) {
		return res
	}

	return res
}
