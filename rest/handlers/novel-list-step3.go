package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
	"time"
)

func NovelListStep3(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	_sort := Cp(req.Parameters, "sort")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.Mconn.DsnSlave]
	var query bytes.Buffer
	if _seqNovelStep2 > 0 {
		query.WriteString("SELECT COUNT(seq_novel_step3) FROM novel_step3 WHERE active_yn = true AND seq_novel_step1 = ? AND seq_novel_step2 = ? AND temp_yn = false AND deleted_yn = false")
		result := sdb.Raw(query.String(), _seqNovelStep1, _seqNovelStep2).Count(&totalData)
		if corm(result, &res) {
			return res
		}
	} else {
		query.WriteString("SELECT COUNT(seq_novel_step3) FROM novel_step3 WHERE active_yn = true AND seq_novel_step1 = ? AND temp_yn = false AND deleted_yn = false")
		result := sdb.Raw(query.String(), _seqNovelStep1).Count(&totalData)
		if corm(result, &res) {
			return res
		}
	}

	query.Reset()
	query.WriteString(`
		SELECT
			ns.seq_novel_step3,
			md.seq_member,
			md.nick_name,
			ns.created_at,
			ns.cnt_like,
			ns.content
		FROM novel_step3 ns
		INNER JOIN member_details md ON ns.seq_member = md.seq_member
	`)
	if _seqNovelStep2 > 0 {
		query.WriteString(`WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.seq_novel_step2 = ? AND ns.temp_yn = false AND ns.deleted_yn = false`)
	} else {
		query.WriteString(`WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false`)
	}
	switch _sort {
	case define.LIKE:
		query.WriteString(" ORDER BY ns.cnt_like DESC, ns.updated_at DESC")
	case define.RECENT:
		fallthrough
	default:
		query.WriteString(" ORDER BY ns.updated_at DESC")
	}
	query.WriteString(" LIMIT ?, ?")

	fmt.Println(query.String())

	step3ResTmp := []Step3ResTmp{}

	if _seqNovelStep2 > 0 {
		result := sdb.Raw(query.String(), _seqNovelStep1, _seqNovelStep2, limitStart, _sizePerPage).Find(&step3ResTmp)
		if corm(result, &res) {
			return res
		}
	} else {
		result := sdb.Raw(query.String(), _seqNovelStep1, limitStart, _sizePerPage).Find(&step3ResTmp)
		if corm(result, &res) {
			return res
		}
	}

	novelListStep3Res := NovelListStep3Res{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	isBool := false
	var seqs []int64
	for i := 0; i < len(step3ResTmp); i++ {
		o := step3ResTmp[i]
		seqs = append(seqs, o.SeqNovelStep3)
		novelListStep3Res.List = append(novelListStep3Res.List, struct {
			SeqNovelStep3 int64  "json:\"seq_novel_step3\""
			SeqMember     int64  "json:\"seq_member\""
			NickName      string "json:\"nick_name\""
			CreatedAt     int64  "json:\"created_at\""
			CntLike       int64  "json:\"cnt_like\""
			MyLike        bool   "json:\"my_like\""
			Content       string "json:\"content\""
		}{
			SeqNovelStep3: o.SeqNovelStep3,
			SeqMember:     o.SeqMember,
			NickName:      o.NickName,
			CreatedAt:     o.CreatedAt.UnixMilli(),
			CntLike:       o.CntLike,
			MyLike:        isBool,
			Content:       o.Content,
		})
	}

	if userToken != nil {
		var listSeq []int64
		ldb := GetMyLogDbSlave(userToken.Allocated)
		ldb.Model(schemas.MemberLikeStep3{}).
			Select("seq_novel_step3").
			Where("seq_member = ? AND seq_novel_step3 IN (?) AND like_yn = true", userToken.SeqMember, seqs).
			Scan(&listSeq)

		for i := 0; i < len(novelListStep3Res.List); i++ {
			o := novelListStep3Res.List[i]
			for _, v := range listSeq {
				if v == o.SeqNovelStep3 {
					novelListStep3Res.List[i].MyLike = true
					break
				}
			}
		}
	}

	res.Data = novelListStep3Res

	return res
}

type Step3ResTmp struct {
	SeqNovelStep3 int64     `json:"seq_novel_step3"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	CntLike       int64     `json:"cnt_like"`
	Content       string    `json:"content"`
}

type NovelListStep3Res struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNovelStep3 int64  `json:"seq_novel_step3"`
		SeqMember     int64  `json:"seq_member"`
		NickName      string `json:"nick_name"`
		CreatedAt     int64  `json:"created_at"`
		CntLike       int64  `json:"cnt_like"`
		MyLike        bool   `json:"my_like"`
		Content       string `json:"content"`
	} `json:"list"`
}
