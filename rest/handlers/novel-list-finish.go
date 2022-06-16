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
	if _seqGenre > 0 {
		sdb.Model(schemas.NovelFinish{}).Where("active_yn = true AND seq_genre = ?", _seqGenre).Count(&totalData)
	} else {
		sdb.Model(schemas.NovelFinish{}).Where("active_yn = true").Count(&totalData)
	}

	novelListFinishRes := NovelListFinishRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	var query bytes.Buffer
	query.WriteString(`
	SELECT
		seq_novel_finish,
		seq_genre,
		seq_image,
		seq_color,
		title,
		false AS my_bookmark
	FROM novel_finishes nf
	WHERE active_yn = true`)
	if _seqGenre > 0 {
		query.WriteString(" AND seq_genre = " + strconv.Itoa(int(_seqGenre)))
	}
	switch _sort {
	case "RECENT":
		query.WriteString(" ORDER BY seq_novel_finish DESC")
	case "LIKE":
		query.WriteString(" ORDER BY cnt_like DESC")
	}
	query.WriteString(" LIMIT ?, ?")
	result := sdb.Raw(query.String(), limitStart, _sizePerPage).Find(&novelListFinishRes.List)
	if corm(result, &res) {
		return res
	}

	//finish seq 구한다
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if userToken != nil {
		var seqNovelFinishes []int64
		var listMy []int64
		for _, v := range novelListFinishRes.List {
			seqNovelFinishes = append(seqNovelFinishes, v.SeqNovelFinish)
		}
		ldb := GetMyLogDb(userToken.Allocated)
		ldb.Model(schemas.MemberBookmark{}).
			Where("seq_member = ? AND seq_novel_finish IN (?) AND bookmark_yn = true", userToken.SeqMember, seqNovelFinishes).
			Select("seq_novel_finish").
			Scan(&listMy)
		for i := 0; i < len(novelListFinishRes.List); i++ {
			o := novelListFinishRes.List[i]
			for _, v := range listMy {
				if o.SeqNovelFinish == v {
					novelListFinishRes.List[i].MyBookmark = true
					break
				}
			}
		}
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
