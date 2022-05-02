package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"time"
)

func NovelListStep3(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	masterDB := db.List[define.DSN_MASTER]
	var query bytes.Buffer
	query.WriteString("SELECT seq_novel_step3 FROM novel_step3 WHERE active_yn = true AND seq_novel_step1 = ? AND temp_yn = false")
	result := masterDB.Raw(query.String(), _seqNovelStep1).Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
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
		WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false`)
	query.WriteString(" ORDER BY ns.seq_novel_step3 DESC")
	query.WriteString(" LIMIT ?, ?")
	step3ResTmp := []Step3ResTmp{}
	result = masterDB.Raw(query.String(), _seqNovelStep1, limitStart, _sizePerPage).Find(&step3ResTmp)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelListStep3Res := NovelListStep3Res{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	isBool := false
	for i := 0; i < len(step3ResTmp); i++ {
		o := step3ResTmp[i]
		isBool = !isBool
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
