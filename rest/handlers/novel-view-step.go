package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
)

func NovelViewStep(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	fmt.Println("호출됨?")

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_seqNovelStep4 := CpInt64(req.Parameters, "seq_novel_step4")
	query := ""
	mdb := db.List[define.DSN_MASTER]
	o := NovelViewStepRes{}
	if _seqNovelStep1 > 0 {
		result := mdb.Model(schemas.NovelStep1{}).
			Select("title, content, cnt_like, UNIX_TIMESTAMP(created_at) * 1000 AS created_at").
			Where("seq_novel_step1 = ?", _seqNovelStep1).
			Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 1, _seqNovelStep1)
		o.Step = 1
	} else if _seqNovelStep2 > 0 {
		query = `
		SELECT
			ns1.title,
			ns2.content,
			ns2.cnt_like,
			UNIX_TIMESTAMP(ns2.created_at) * 1000 AS created_at
		FROM novel_step1 ns1
		INNER JOIN novel_step2 ns2 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
		WHERE ns2.seq_novel_step2 = ?
		`
		result := mdb.Raw(query, _seqNovelStep2).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 2, _seqNovelStep2)
		o.Step = 2
	} else if _seqNovelStep3 > 0 {
		query = `
		SELECT
			ns1.title,
			ns3.content,
			ns3.cnt_like,
			UNIX_TIMESTAMP(ns3.created_at) * 1000 AS created_at
		FROM novel_step1 ns1
		INNER JOIN novel_step3 ns3 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
		WHERE ns3.seq_novel_step3 = ?
		`
		result := mdb.Raw(query, _seqNovelStep3).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 3, _seqNovelStep3)
		o.Step = 3
	} else if _seqNovelStep4 > 0 {
		query = `
		SELECT
			ns1.title,
			ns4.content,
			ns4.cnt_like,
			UNIX_TIMESTAMP(ns4.created_at) * 1000 AS created_at
		FROM novel_step1 ns1
		INNER JOIN novel_step4 ns4 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
		WHERE ns3.seq_novel_step4 = ?
		`
		result := mdb.Raw(query, _seqNovelStep3).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 4, _seqNovelStep4)
		o.Step = 4
	}
	res.Data = o

	return res
}

type NovelViewStepRes struct {
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	CntLike   int     `json:"cnt_like"`
	MyLike    bool    `json:"my_like"`
	Step      int8    `json:"step"`
	CreatedAt float64 `json:"created_at"`
}
