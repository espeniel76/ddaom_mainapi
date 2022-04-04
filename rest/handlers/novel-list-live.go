package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"fmt"
	"strconv"
)

func NovelListLive(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqGenre := CpInt64(req.Parameters, "seq_genre")
	_seqKeyword := CpInt64(req.Parameters, "seq_keyword")
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
	masterDB := db.List[define.DSN_MASTER]
	var query bytes.Buffer
	query.WriteString(`
		SELECT ns.seq_novel_step1
		FROM novel_step1 ns
		INNER JOIN keywords k ON ns.seq_keyword = k.seq_keyword
		WHERE NOW() BETWEEN k.start_date AND k.end_date
		AND ns.active_yn = true AND k.active_yn = true AND k.seq_keyword = ? AND ns.temp_yn = false`)
	if _seqGenre > 0 {
		query.WriteString(" AND seq_genre = " + strconv.Itoa(int(_seqGenre)))
	}
	result := masterDB.Raw(query.String(), _seqKeyword).Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelListLiveRes := NovelListLiveRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	query.Reset()
	query.WriteString(`
		SELECT
			ns.seq_novel_step1,
			ns.seq_keyword,
			ns.seq_genre,
			ns.seq_image,
			ns.seq_color,
			ns.title,
			ns.content
		FROM novel_step1 ns
		INNER JOIN keywords k ON ns.seq_keyword = k.seq_keyword
		WHERE NOW() BETWEEN k.start_date AND k.end_date
		AND ns.active_yn = true AND k.active_yn = true AND k.seq_keyword = ? AND ns.temp_yn = false`)
	if _seqGenre > 0 {
		query.WriteString(" AND seq_genre = " + strconv.Itoa(int(_seqGenre)))
	}
	query.WriteString(" ORDER BY ns.seq_novel_step1 DESC")
	query.WriteString(" LIMIT ?, ?")
	result = masterDB.Raw(query.String(), _seqKeyword, limitStart, _sizePerPage).Find(&novelListLiveRes.List)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	res.Data = novelListLiveRes

	return res
}

type NovelListLiveRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNovelStep1 int64  `json:"seq_novel_step1"`
		SeqKeyword    int64  `json:"seq_keyword"`
		SeqGenre      int64  `json:"seq_genre"`
		SeqImage      int64  `json:"seq_image"`
		SeqColor      int64  `json:"seq_color"`
		Title         string `json:"title"`
		Content       string `json:"content"`
	} `json:"list"`
}
