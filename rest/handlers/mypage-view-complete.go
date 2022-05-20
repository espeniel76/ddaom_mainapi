package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func MypageViewComplete(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	_step, _ := strconv.Atoi(req.Vars["step"])
	_seqNovel, _ := strconv.Atoi(req.Vars["seq_novel"])

	sdb := db.List[define.DSN_SLAVE]
	o := MypageViewCompleteRes{}
	switch _step {
	case 1:
		result := sdb.Model(schemas.NovelStep1{}).
			Select("title, content, true AS my_like, cnt_like, 1 AS step, UNIX_TIMESTAMP(created_at) * 1000 AS created_at").
			Where("seq_novel_step1 = ? AND active_yn = true AND temp_yn = false", _seqNovel).
			Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 1, int64(_seqNovel))
		res.Data = o
	case 2:
		query := `
			SELECT
				ns1.title,
				ns2.content,
				true AS my_like,
				ns2.cnt_like,
				2 AS step,
				UNIX_TIMESTAMP(ns2.created_at) * 1000 AS created_at
			FROM novel_step2 ns2 INNER JOIN novel_step1 ns1
			ON ns2.seq_novel_step1 = ns1.seq_novel_step1
			WHERE ns2.seq_novel_step2 = ? AND ns2.active_yn = true AND ns2.temp_yn = false
		`
		result := sdb.Raw(query, _seqNovel).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 2, int64(_seqNovel))
		res.Data = o
	case 3:
		query := `
			SELECT
				ns1.title,
				ns3.content,
				true AS my_like,
				ns3.cnt_like,
				3 AS step,
				UNIX_TIMESTAMP(ns3.created_at) * 1000 AS created_at
			FROM novel_step3 ns3 INNER JOIN novel_step1 ns1
			ON ns3.seq_novel_step1 = ns1.seq_novel_step1
			WHERE ns3.seq_novel_step3 = ? AND ns3.active_yn = true AND ns3.temp_yn = false
		`
		result := sdb.Raw(query, _seqNovel).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 3, int64(_seqNovel))
		res.Data = o
	case 4:
		query := `
			SELECT
				ns1.title,
				ns4.content,
				true AS my_like,
				ns4.cnt_like,
				4 AS step,
				UNIX_TIMESTAMP(ns4.created_at) * 1000 AS created_at
			FROM novel_step4 ns4 INNER JOIN novel_step1 ns1
			ON ns4.seq_novel_step1 = ns1.seq_novel_step1
			WHERE ns4.seq_novel_step4 = ? AND ns4.active_yn = true AND ns4.temp_yn = false
		`
		result := sdb.Raw(query, _seqNovel).Scan(&o)
		if corm(result, &res) {
			return res
		}
		o.MyLike = getMyLike(userToken, 4, int64(_seqNovel))
		res.Data = o
	}

	return res
}

type MypageViewCompleteRes struct {
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	CreatedAt float64 `json:"created_at"`
	MyLike    bool    `json:"my_like"`
	CntLike   int     `json:"cnt_like"`
	Step      int     `json:"step"`
}
