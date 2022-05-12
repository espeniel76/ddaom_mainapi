package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

// 마이페이지 - 작성완료리스트
func MypageListComplete(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqMember := CpInt64(req.Parameters, "seq_member")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if _seqMember == 0 && userToken != nil {
		_seqMember = userToken.SeqMember
	}

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	sdb := db.List[define.DSN_SLAVE]

	var totalData int64
	seq := _seqMember
	query := `
	SELECT SUM(cnt1) + SUM(cnt2) + SUM(cnt3) + SUM(cnt4) AS cnt
	FROM
	(
		(SELECT COUNT(*) AS cnt1, 0 AS cnt2,0 AS cnt3, 0 AS cnt4 FROM novel_step1 WHERE seq_member = ? AND active_yn = true AND temp_yn = false)
		UNION ALL
		(SELECT 0 AS cnt1, COUNT(*) AS cnt2, 0 AS cnt3, 0 AS cnt4 FROM novel_step2 WHERE seq_member = ? AND active_yn = true AND temp_yn = false)
		UNION ALL
		(SELECT 0 AS cnt1, 0 AS cnt2, COUNT(*) AS cnt3, 0 AS cnt4 FROM novel_step3 WHERE seq_member = ? AND active_yn = true AND temp_yn = false)
		UNION ALL
		(SELECT 0 AS cnt1, 0 AS cnt2, 0 AS cnt3, COUNT(*) AS cnt4 FROM novel_step4 WHERE seq_member = ? AND active_yn = true AND temp_yn = false)
	) AS s
	`
	result := sdb.Raw(query, seq, seq, seq, seq).Scan(&totalData)
	if corm(result, &res) {
		return res
	}
	novelMyListCompleteRes := NovelMyListCompleteRes{
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
			ns.updated_at,
			1 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live,
			false AS my_like,
			ns.cnt_like
		FROM novel_step1 ns
		INNER JOIN keywords k ON k.seq_keyword = ns.seq_keyword
		WHERE seq_member = ? AND ns.active_yn = true AND ns.temp_yn = false
	)
	UNION 
	(
		SELECT
			0 AS seq_novel_step1,
			ns2.seq_novel_step2,
			0 AS seq_novel_step3,
			0 AS seq_novel_step4,
			ns1.title,
			UNIX_TIMESTAMP(ns2.created_at) * 1000 AS created_at,
			ns2.updated_at,
			2 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live,
			false AS my_like,
			ns2.cnt_like
		FROM novel_step2 ns2
		INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 
		INNER JOIN keywords k ON k.seq_keyword = ns1.seq_keyword
		WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns2.temp_yn = false
	)
	UNION 
	(
		SELECT
			0 AS seq_novel_step1,
			0 AS seq_novel_step2,
			ns3.seq_novel_step3,
			0 AS seq_novel_step4,
			ns1.title,
			UNIX_TIMESTAMP(ns3.created_at) * 1000 AS created_at,
			ns3.updated_at,
			3 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live,
			false AS my_like,
			ns3.cnt_like
		FROM novel_step3 ns3
		INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 
		INNER JOIN keywords k ON k.seq_keyword = ns1.seq_keyword
		WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns3.temp_yn = false
	)
	UNION 
	(
		SELECT
			0 AS seq_novel_step1,
			0 AS seq_novel_step2,
			0 AS seq_novel_step3,
			ns4.seq_novel_step4,
			ns1.title,
			UNIX_TIMESTAMP(ns4.created_at) * 1000 AS created_at,
			ns4.updated_at,
			4 AS step,
			IF (k.end_date > NOW(), true, false) AS is_live,
			false AS my_like,
			ns4.cnt_like
		FROM novel_step4 ns4
		INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 
		INNER JOIN keywords k ON k.seq_keyword = ns1.seq_keyword
		WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns4.temp_yn = false
	)
	ORDER BY updated_at DESC
	LIMIT ?, ?
	`
	result = sdb.
		Raw(query, seq, seq, seq, seq, limitStart, _sizePerPage).
		Scan(&novelMyListCompleteRes.List)
	if corm(result, &res) {
		return res
	}

	// 나의 좋아요 구한다. (로그인 시 에만)
	// 미안하다.. 무식하다...
	if userToken != nil {
		var seqNovelStep1s []int64
		var seqNovelStep2s []int64
		var seqNovelStep3s []int64
		var seqNovelStep4s []int64
		for _, v := range novelMyListCompleteRes.List {
			if v.SeqNovelStep1 > 0 {
				seqNovelStep1s = append(seqNovelStep1s, v.SeqNovelStep1)
			}
			if v.SeqNovelStep2 > 0 {
				seqNovelStep2s = append(seqNovelStep2s, v.SeqNovelStep2)
			}
			if v.SeqNovelStep3 > 0 {
				seqNovelStep3s = append(seqNovelStep3s, v.SeqNovelStep3)
			}
			if v.SeqNovelStep4 > 0 {
				seqNovelStep4s = append(seqNovelStep4s, v.SeqNovelStep4)
			}
		}
		// fmt.Println(seqNovelStep1s)
		// fmt.Println(seqNovelStep2s)
		// fmt.Println(seqNovelStep3s)
		// fmt.Println(seqNovelStep4s)
		ldb := GetMyLogDb(userToken.Allocated)
		mls1 := []schemas.MemberLikeStep1{}
		mls2 := []schemas.MemberLikeStep2{}
		mls3 := []schemas.MemberLikeStep3{}
		mls4 := []schemas.MemberLikeStep4{}
		ldb.Model(&mls1).Select("seq_novel_step1, like_yn").Where("seq_member = ? AND seq_novel_step1 IN (?)", userToken.SeqMember, seqNovelStep1s).Scan(&mls1)
		ldb.Model(&mls2).Select("seq_novel_step2, like_yn").Where("seq_member = ? AND seq_novel_step2 IN (?)", userToken.SeqMember, seqNovelStep2s).Scan(&mls2)
		ldb.Model(&mls3).Select("seq_novel_step3, like_yn").Where("seq_member = ? AND seq_novel_step3 IN (?)", userToken.SeqMember, seqNovelStep3s).Scan(&mls3)
		ldb.Model(&mls4).Select("seq_novel_step4, like_yn").Where("seq_member = ? AND seq_novel_step4 IN (?)", userToken.SeqMember, seqNovelStep4s).Scan(&mls4)

		// fmt.Println(mls1)
		// fmt.Println(mls2)
		// fmt.Println(mls3)
		// fmt.Println(mls4)
		for _, v := range mls1 {
			if v.LikeYn {
				for i := 0; i < len(novelMyListCompleteRes.List); i++ {
					if v.SeqNovelStep1 == novelMyListCompleteRes.List[i].SeqNovelStep1 {
						novelMyListCompleteRes.List[i].MyLike = true
						break
					}
				}
			}
		}
		for _, v := range mls2 {
			if v.LikeYn {
				for i := 0; i < len(novelMyListCompleteRes.List); i++ {
					if v.SeqNovelStep2 == novelMyListCompleteRes.List[i].SeqNovelStep2 {
						novelMyListCompleteRes.List[i].MyLike = true
						break
					}
				}
			}
		}
		for _, v := range mls3 {
			if v.LikeYn {
				for i := 0; i < len(novelMyListCompleteRes.List); i++ {
					if v.SeqNovelStep3 == novelMyListCompleteRes.List[i].SeqNovelStep3 {
						novelMyListCompleteRes.List[i].MyLike = true
						break
					}
				}
			}
		}
		for _, v := range mls4 {
			if v.LikeYn {
				for i := 0; i < len(novelMyListCompleteRes.List); i++ {
					if v.SeqNovelStep4 == novelMyListCompleteRes.List[i].SeqNovelStep4 {
						novelMyListCompleteRes.List[i].MyLike = true
						break
					}
				}
			}
		}
	}

	res.Data = novelMyListCompleteRes

	return res
}

type NovelMyListCompleteRes struct {
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
		MyLike        bool    `json:"my_like"`
		CntLike       int     `json:"cnt_like"`
	} `json:"list"`
}
