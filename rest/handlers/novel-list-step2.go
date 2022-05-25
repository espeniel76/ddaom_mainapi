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

func NovelListStep2(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	fmt.Println(_seqNovelStep1, _page, _sizePerPage)
	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	masterDB := db.List[define.DSN_MASTER]
	var query bytes.Buffer
	query.WriteString("SELECT COUNT(seq_novel_step2) FROM novel_step2 WHERE active_yn = true AND seq_novel_step1 = ? AND temp_yn = false AND deleted_yn = false")
	result := masterDB.Raw(query.String(), _seqNovelStep1).Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	query.Reset()
	query.WriteString(`
		SELECT
			ns.seq_novel_step2,
			md.seq_member,
			md.nick_name,
			ns.created_at,
			ns.cnt_like,
			ns.content
		FROM novel_step2 ns
		INNER JOIN member_details md ON ns.seq_member = md.seq_member
		WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false`)
	query.WriteString(" ORDER BY ns.seq_novel_step2 DESC")
	query.WriteString(" LIMIT ?, ?")
	step2ResTmp := []Step2ResTmp{}
	result = masterDB.Raw(query.String(), _seqNovelStep1, limitStart, _sizePerPage).Find(&step2ResTmp)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelListStep2Res := NovelListStep2Res{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	isBool := false
	var seqs []int64
	for i := 0; i < len(step2ResTmp); i++ {
		o := step2ResTmp[i]
		seqs = append(seqs, o.SeqNovelStep2)
		novelListStep2Res.List = append(novelListStep2Res.List, struct {
			SeqNovelStep2 int64  "json:\"seq_novel_step2\""
			SeqMember     int64  "json:\"seq_member\""
			NickName      string "json:\"nick_name\""
			CreatedAt     int64  "json:\"created_at\""
			CntLike       int64  "json:\"cnt_like\""
			MyLike        bool   "json:\"my_like\""
			Content       string "json:\"content\""
		}{
			SeqNovelStep2: o.SeqNovelStep2,
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
		ldb := GetMyLogDb(userToken.Allocated)
		ldb.Model(schemas.MemberLikeStep2{}).
			Select("seq_novel_step2").
			Where("seq_member = ? AND seq_novel_step2 IN (?) AND like_yn = true", userToken.SeqMember, seqs).
			Scan(&listSeq)

		for i := 0; i < len(novelListStep2Res.List); i++ {
			o := novelListStep2Res.List[i]
			for _, v := range listSeq {
				if v == o.SeqNovelStep2 {
					novelListStep2Res.List[i].MyLike = true
					break
				}
			}
		}
	}

	res.Data = novelListStep2Res

	return res
}

type Step2ResTmp struct {
	SeqNovelStep2 int64     `json:"seq_novel_step2"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	CntLike       int64     `json:"cnt_like"`
	Content       string    `json:"content"`
}

type NovelListStep2Res struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNovelStep2 int64  `json:"seq_novel_step2"`
		SeqMember     int64  `json:"seq_member"`
		NickName      string `json:"nick_name"`
		CreatedAt     int64  `json:"created_at"`
		CntLike       int64  `json:"cnt_like"`
		MyLike        bool   `json:"my_like"`
		Content       string `json:"content"`
	} `json:"list"`
}
