package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
)

func MypageListLive(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqMember := CpInt64(req.Parameters, "seq_member")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	itsMe := false
	if userToken != nil {
		if userToken.SeqMember == _seqMember {
			itsMe = true
		}
	}

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	sdb := db.List[define.Mconn.DsnSlave]

	var totalData int64
	seq := _seqMember
	query := ""
	if itsMe {
		query = `
		SELECT SUM(cnt1) + SUM(cnt2) + SUM(cnt3) + SUM(cnt4) AS cnt
		FROM
		(
			(
				SELECT COUNT(*) AS cnt1, 0 AS cnt2,0 AS cnt3, 0 AS cnt4 
				FROM novel_step1 ns1 INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns1.seq_member = ? AND ns1.active_yn = true AND ns1.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
			UNION ALL
			(
				SELECT 0 AS cnt1, COUNT(*) AS cnt2, 0 AS cnt3, 0 AS cnt4 
				FROM novel_step2 ns2
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
			UNION ALL
			(
				SELECT 0 AS cnt1, 0 AS cnt2, COUNT(*) AS cnt3, 0 AS cnt4 
				FROM novel_step3 ns3
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
			UNION ALL
			(
				SELECT 0 AS cnt1, 0 AS cnt2, 0 AS cnt3, COUNT(*) AS cnt4 
				FROM novel_step4 ns4
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
		) AS s
		`
	} else {
		query = `
		SELECT SUM(cnt1) + SUM(cnt2) + SUM(cnt3) + SUM(cnt4) AS cnt
		FROM
		(
			(
				SELECT COUNT(*) AS cnt1, 0 AS cnt2,0 AS cnt3, 0 AS cnt4 
				FROM novel_step1 ns1 INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns1.seq_member = ? AND ns1.active_yn = true AND ns1.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns1.deleted_yn = false
			)
			UNION ALL
			(
				SELECT 0 AS cnt1, COUNT(*) AS cnt2, 0 AS cnt3, 0 AS cnt4 
				FROM novel_step2 ns2
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns2.deleted_yn = false
			)
			UNION ALL
			(
				SELECT 0 AS cnt1, 0 AS cnt2, COUNT(*) AS cnt3, 0 AS cnt4 
				FROM novel_step3 ns3
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns3.deleted_yn = false
			)
			UNION ALL
			(
				SELECT 0 AS cnt1, 0 AS cnt2, 0 AS cnt3, COUNT(*) AS cnt4 
				FROM novel_step4 ns4
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns4.deleted_yn = false
			)
		) AS s
		`
	}

	result := sdb.Raw(query, seq, seq, seq, seq).Scan(&totalData)
	novelMyListLiveRes := NovelMyListLiveRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	if itsMe {
		query = `
			(
				SELECT
					ns1.seq_novel_step1,
					0 AS seq_novel_step2,
					0 AS seq_novel_step3,
					0 AS seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns1.content,
					ns1.cnt_like,
					UNIX_TIMESTAMP(ns1.created_at) *1000 AS created_at,
					ns1.updated_at,
					1 AS step,
					ns1.deleted_yn
				FROM novel_step1 ns1 INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns1.seq_member = ? AND ns1.active_yn = true AND ns1.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
			UNION ALL
			(
				SELECT
					0 AS seq_novel_step1,
					ns2.seq_novel_step2,
					0 AS seq_novel_step3,
					0 AS seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns2.content,
					ns2.cnt_like,
					UNIX_TIMESTAMP(ns2.created_at) *1000 AS created_at,
					ns2.updated_at,
					2 AS step,
					ns2.deleted_yn
				FROM novel_step2 ns2
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
			UNION ALL
			(
				SELECT
					0 AS seq_novel_step1,
					0 AS seq_novel_step2,
					ns3.seq_novel_step3,
					0 AS seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns3.content,
					ns3.cnt_like,
					UNIX_TIMESTAMP(ns3.created_at) *1000 AS created_at,
					ns3.updated_at,
					3 AS step,
					ns3.deleted_yn
				FROM novel_step3 ns3
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
			UNION ALL
			(
				SELECT
					0 AS seq_novel_step1,
					0 AS seq_novel_step2,
					0 AS seq_novel_step3,
					ns4.seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns4.content,
					ns4.cnt_like,
					UNIX_TIMESTAMP(ns4.created_at) *1000 AS created_at,
					ns4.updated_at,
					4 AS step,
					ns4.deleted_yn
				FROM novel_step4 ns4
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date
			)
		ORDER BY updated_at DESC
		LIMIT ?, ?
		`
	} else {
		query = `
		(
				SELECT
					ns1.seq_novel_step1,
					0 AS seq_novel_step2,
					0 AS seq_novel_step3,
					0 AS seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns1.content,
					ns1.cnt_like,
					UNIX_TIMESTAMP(ns1.created_at) *1000 AS created_at,
					ns1.updated_at,
					1 AS step,
					ns1.deleted_yn
				FROM novel_step1 ns1 INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns1.seq_member = ? AND ns1.active_yn = true AND ns1.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns1.deleted_yn = false
			)
			UNION ALL
			(
				SELECT
					0 AS seq_novel_step1,
					ns2.seq_novel_step2,
					0 AS seq_novel_step3,
					0 AS seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns2.content,
					ns2.cnt_like,
					UNIX_TIMESTAMP(ns2.created_at) *1000 AS created_at,
					ns2.updated_at,
					2 AS step,
					ns2.deleted_yn
				FROM novel_step2 ns2
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns2.deleted_yn = false
			)
			UNION ALL
			(
				SELECT
					0 AS seq_novel_step1,
					0 AS seq_novel_step2,
					ns3.seq_novel_step3,
					0 AS seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns3.content,
					ns3.cnt_like,
					UNIX_TIMESTAMP(ns3.created_at) *1000 AS created_at,
					ns3.updated_at,
					3 AS step,
					ns3.deleted_yn
				FROM novel_step3 ns3
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns3.deleted_yn = false
			)
			UNION ALL
			(
				SELECT
					0 AS seq_novel_step1,
					0 AS seq_novel_step2,
					0 AS seq_novel_step3,
					ns4.seq_novel_step4,
					ns1.seq_keyword,
					ns1.seq_genre,
					ns1.seq_image,
					ns1.seq_color,
					ns1.title,
					ns4.content,
					ns4.cnt_like,
					UNIX_TIMESTAMP(ns4.created_at) *1000 AS created_at,
					ns4.updated_at,
					4 AS step,
					ns4.deleted_yn
				FROM novel_step4 ns4
				INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1
				INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
				WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = false AND NOW() BETWEEN k.start_date AND k.end_date AND ns4.deleted_yn = false
			)
		ORDER BY updated_at DESC
		LIMIT ?, ?
		`
	}

	result = sdb.
		Raw(query, seq, seq, seq, seq, limitStart, _sizePerPage).
		Scan(&novelMyListLiveRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = novelMyListLiveRes

	return res
}

type NovelMyListLiveRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		Step          int8    `json:"step"`
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqNovelStep2 int64   `json:"seq_novel_step2"`
		SeqNovelStep3 int64   `json:"seq_novel_step3"`
		SeqNovelStep4 int64   `json:"seq_novel_step4"`
		SeqKeyword    int64   `json:"seq_keyword"`
		SeqGenre      int64   `json:"seq_genre"`
		SeqImage      int64   `json:"seq_image"`
		SeqColor      int64   `json:"seq_color"`
		Title         string  `json:"title"`
		Content       string  `json:"content"`
		CntLike       int64   `json:"cnt_like"`
		CreatedAt     float64 `json:"created_at"`
		DeletedYn     bool    `json:"deleted_yn"`
	} `json:"list"`
}
