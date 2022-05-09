package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
)

func MypageListTemp(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_isLive := CpBool(req.Parameters, "is_live")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	fmt.Println(_isLive)

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	sdb := db.List[define.DSN_SLAVE]

	var seqKeywords []int64

	if _isLive {
		result := sdb.Model(schemas.Keyword{}).Select("seq_keyword").Where("NOW() BETWEEN start_date AND end_date").Scan(&seqKeywords)
		if corm(result, &res) {
			return res
		}
	} else {
		result := sdb.Model(schemas.Keyword{}).Select("seq_keyword").Scan(&seqKeywords)
		if corm(result, &res) {
			return res
		}
	}

	fmt.Println(seqKeywords)

	var totalData int64
	seq := userToken.SeqMember
	query := `
	SELECT SUM(cnt1) + SUM(cnt2) + SUM(cnt3) + SUM(cnt4) AS cnt
	FROM
	(
		(
			SELECT COUNT(*) AS cnt1, 0 AS cnt2,0 AS cnt3, 0 AS cnt4
			FROM novel_step1
			WHERE seq_member = ? AND active_yn = true AND temp_yn = true
			AND seq_keyword in (?)
		)
		UNION ALL
		(
			SELECT 0 AS cnt1, COUNT(*) AS cnt2, 0 AS cnt3, 0 AS cnt4
			FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
			WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = true
			AND ns1.seq_keyword in (?)
		)
		UNION ALL
		(
			SELECT 0 AS cnt1, 0 AS cnt2, COUNT(*) AS cnt3, 0 AS cnt4 
			FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
			WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = true
			AND ns1.seq_keyword in (?)
		)
		UNION ALL
		(
			SELECT 0 AS cnt1, 0 AS cnt2, 0 AS cnt3, COUNT(*) AS cnt4
			FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
			WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = true
			AND ns1.seq_keyword in (?)
		)
	) AS s
	`
	result := sdb.Raw(query, seq, seqKeywords, seq, seqKeywords, seq, seqKeywords, seq, seqKeywords).Scan(&totalData)
	if corm(result, &res) {
		return res
	}
	novelMyListTempRes := NovelMyListTempRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	query = `
	(
		SELECT
			seq_novel_step1,
			0 AS seq_novel_step2,
			0 AS seq_novel_step3,
			0 AS seq_novel_step4,
			title,
			UNIX_TIMESTAMP(ns.created_at) * 1000 AS created_at,
			1 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live
		FROM novel_step1 ns
		INNER JOIN keywords k ON k.seq_keyword = ns.seq_keyword
		WHERE seq_member = ? AND ns.active_yn = true AND ns.temp_yn = true
		AND ns.seq_keyword IN (?)
	)
	UNION ALL 
	(
		SELECT
			0 AS seq_novel_step1,
			ns2.seq_novel_step2,
			0 AS seq_novel_step3,
			0 AS seq_novel_step4,
			ns1.title,
			UNIX_TIMESTAMP(ns2.created_at) * 1000 AS created_at,
			2 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live
		FROM novel_step2 ns2
		INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 
		INNER JOIN keywords k ON k.seq_keyword = ns1.seq_keyword
		WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = true
		AND ns1.seq_keyword IN (?)
	)
	UNION ALL 
	(
		SELECT
			0 AS seq_novel_step1,
			0 AS seq_novel_step2,
			ns3.seq_novel_step3,
			0 AS seq_novel_step4,
			ns1.title,
			UNIX_TIMESTAMP(ns3.created_at) * 1000 AS created_at,
			3 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live
		FROM novel_step3 ns3
		INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 
		INNER JOIN keywords k ON k.seq_keyword = ns1.seq_keyword
		WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = true
		AND ns1.seq_keyword IN (?)
	)
	UNION ALL 
	(
		SELECT
			0 AS seq_novel_step1,
			0 AS seq_novel_step2,
			0 AS seq_novel_step3,
			ns4.seq_novel_step4,
			ns1.title,
			UNIX_TIMESTAMP(ns4.created_at) * 1000 AS created_at,
			4 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live
		FROM novel_step4 ns4
		INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 
		INNER JOIN keywords k ON k.seq_keyword = ns1.seq_keyword
		WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = true
		AND ns1.seq_keyword IN (?)
	)
	ORDER BY created_at DESC
	LIMIT ?, ?
	`
	result = sdb.
		Raw(query, seq, seqKeywords, seq, seqKeywords, seq, seqKeywords, seq, seqKeywords, limitStart, _sizePerPage).
		Scan(&novelMyListTempRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = novelMyListTempRes

	return res
}

type NovelMyListTempRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqNovelStep2 int64   `json:"seq_novel_step2"`
		SeqNovelStep3 int64   `json:"seq_novel_step3"`
		SeqNovelStep4 int64   `json:"seq_novel_step4"`
		Title         string  `json:"title"`
		CreatedAt     float64 `json:"created_at"`
		Step          int8    `json:"step"`
		IsLive        bool    `json:"is_live"`
	} `json:"list"`
}
