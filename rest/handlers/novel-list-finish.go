package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
	"strconv"
)

func NovelListFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqGenre := CpInt64(req.Parameters, "seq_genre")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	_sort := Cp(req.Parameters, "sort")
	fmt.Println(_sort)

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.DSN_SLAVE]
	result := sdb.Model(schemas.NovelFinish{}).Where("active_yn = true").Count(&totalData)
	if corm(result, &res) {
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
	WHERE nf.active_yn = true`)
	if _seqGenre > 0 {
		query.WriteString(" AND ns.seq_genre = " + strconv.Itoa(int(_seqGenre)))
	}
	query.WriteString(" ORDER BY nf.seq_novel_finish DESC")
	query.WriteString(" LIMIT ?, ?")
	result = sdb.Raw(query.String(), limitStart, _sizePerPage).Find(&novelListFinishRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = novelListFinishRes

	return res
}

type NovelListFinishRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNovelFinish int64  `json:"seq_novel_finish"`
		SeqGenre       int64  `json:"seq_genre"`
		SeqImage       int64  `json:"seq_image"`
		SeqColor       int64  `json:"seq_color"`
		Title          string `json:"title"`
		MyBookmark     bool   `json:"my_bookmark"`
	} `json:"list"`
}
