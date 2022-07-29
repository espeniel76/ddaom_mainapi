package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func NovelViewStep(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_seqNovelStep4 := CpInt64(req.Parameters, "seq_novel_step4")
	query := ""
	sdb := db.List[define.Mconn.DsnSlave]
	o := NovelViewStepRes{}
	if _seqNovelStep1 > 0 {
		result := sdb.Model(schemas.NovelStep1{}).
			Select("title, content, cnt_like, UNIX_TIMESTAMP(created_at) * 1000 AS created_at, seq_member").
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
			UNIX_TIMESTAMP(ns2.created_at) * 1000 AS created_at,
			ns2.seq_member
		FROM novel_step1 ns1
		INNER JOIN novel_step2 ns2 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
		WHERE ns2.seq_novel_step2 = ?
		`
		result := sdb.Raw(query, _seqNovelStep2).Scan(&o)
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
			UNIX_TIMESTAMP(ns3.created_at) * 1000 AS created_at,
			ns3.seq_member
		FROM novel_step1 ns1
		INNER JOIN novel_step3 ns3 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
		WHERE ns3.seq_novel_step3 = ?
		`
		result := sdb.Raw(query, _seqNovelStep3).Scan(&o)
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
			UNIX_TIMESTAMP(ns4.created_at) * 1000 AS created_at,
			ns4.seq_member
		FROM novel_step1 ns1
		INNER JOIN novel_step4 ns4 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
		WHERE ns3.seq_novel_step4 = ?
		`
		result := sdb.Raw(query, _seqNovelStep3).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 4, _seqNovelStep4)
		o.Step = 4
	}
	bm := getBlockMember(userToken.Allocated, userToken.SeqMember, o.SeqMember)
	o.BlockYn = bm.BlockYn

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
	SeqMember int64   `json:"seq_member"`
	BlockYn   bool    `json:"block_yn"`
}
