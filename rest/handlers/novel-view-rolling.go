package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

func NovelViewRolling(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_step := CpInt64(req.Parameters, "step")
	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.Mconn.DsnSlave]
	var query bytes.Buffer

	switch _step {
	case 2:
		query.WriteString("SELECT COUNT(seq_novel_step2) FROM novel_step2 WHERE active_yn = true AND seq_novel_step1 = ? AND temp_yn = false AND deleted_yn = false")
		result := sdb.Raw(query.String(), _seqNovelStep1).Count(&totalData)
		if corm(result, &res) {
			return res
		}
	case 3:
		// 2단계 종속여부
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
	case 4:
		// 3단계 종속여부
		if _seqNovelStep2 > 0 {
			query.WriteString("SELECT COUNT(seq_novel_step4) FROM novel_step4 WHERE active_yn = true AND seq_novel_step1 = ? AND seq_novel_step3 = ? AND temp_yn = false AND deleted_yn = false")
			result := sdb.Raw(query.String(), _seqNovelStep1, _seqNovelStep3).Count(&totalData)
			if corm(result, &res) {
				return res
			}
		} else {
			query.WriteString("SELECT COUNT(seq_novel_step4) FROM novel_step4 WHERE active_yn = true AND seq_novel_step1 = ? AND temp_yn = false AND deleted_yn = false")
			result := sdb.Raw(query.String(), _seqNovelStep1).Count(&totalData)
			if corm(result, &res) {
				return res
			}
		}
	}

	query.Reset()

	list := NovelListRollingRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
		Step:      int(_step),
	}
	switch _step {
	case 2:
		query.WriteString(`
			SELECT
				ns.seq_novel_step2 AS seq_novel,
				md.seq_member,
				md.nick_name,
				UNIX_TIMESTAMP(ns.created_at) *1000 AS created_at,
				ns.cnt_like,
				ns.content,
				false AS my_like
			FROM novel_step2 ns
			INNER JOIN member_details md ON ns.seq_member = md.seq_member
			WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
			ORDER BY ns.updated_at DESC
			LIMIT ?, ?
		`)
		result := sdb.Raw(query.String(), _seqNovelStep1, limitStart, _sizePerPage).Find(&list.List)
		if corm(result, &res) {
			return res
		}
	case 3:
		// 2단계 종속여부
		if _seqNovelStep2 > 0 {
			query.WriteString(`
				SELECT
					ns.seq_novel_step3 AS seq_novel,
					md.seq_member,
					md.nick_name,
					UNIX_TIMESTAMP(ns.created_at) *1000 AS created_at,
					ns.cnt_like,
					ns.content,
					false AS my_like
				FROM novel_step3 ns
				INNER JOIN member_details md ON ns.seq_member = md.seq_member
				WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.seq_novel_step2 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
				ORDER BY ns.updated_at DESC
				LIMIT ?, ?
			`)
			result := sdb.Raw(query.String(), _seqNovelStep1, _seqNovelStep2, limitStart, _sizePerPage).Find(&list.List)
			if corm(result, &res) {
				return res
			}
		} else {
			query.WriteString(`
				SELECT
					ns.seq_novel_step3 AS seq_novel,
					md.seq_member,
					md.nick_name,
					UNIX_TIMESTAMP(ns.created_at) *1000 AS created_at,
					ns.cnt_like,
					ns.content,
					false AS my_like
				FROM novel_step3 ns
				INNER JOIN member_details md ON ns.seq_member = md.seq_member
				WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
				ORDER BY ns.updated_at DESC
				LIMIT ?, ?
			`)
			result := sdb.Raw(query.String(), _seqNovelStep1, limitStart, _sizePerPage).Find(&list.List)
			if corm(result, &res) {
				return res
			}
		}
	case 4:
		// 3단계 종속여부
		if _seqNovelStep3 > 0 {
			query.WriteString(`
				SELECT
					ns.seq_novel_step4 AS seq_novel,
					md.seq_member,
					md.nick_name,
					UNIX_TIMESTAMP(ns.created_at) *1000 AS created_at,
					ns.cnt_like,
					ns.content,
					false AS my_like
				FROM novel_step4 ns
				INNER JOIN member_details md ON ns.seq_member = md.seq_member
				WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.seq_novel_step3 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
				ORDER BY ns.updated_at DESC
				LIMIT ?, ?
			`)
			result := sdb.Raw(query.String(), _seqNovelStep1, _seqNovelStep3, limitStart, _sizePerPage).Find(&list.List)
			if corm(result, &res) {
				return res
			}
		} else {
			query.WriteString(`
				SELECT
					ns.seq_novel_step4 AS seq_novel,
					md.seq_member,
					md.nick_name,
					UNIX_TIMESTAMP(ns.created_at) *1000 AS created_at,
					ns.cnt_like,
					ns.content,
					false AS my_like
				FROM novel_step4 ns
				INNER JOIN member_details md ON ns.seq_member = md.seq_member
				WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
				ORDER BY ns.updated_at DESC
				LIMIT ?, ?
			`)
			result := sdb.Raw(query.String(), _seqNovelStep1, limitStart, _sizePerPage).Find(&list.List)
			if corm(result, &res) {
				return res
			}
		}
	}

	var seqs []int64
	for i := 0; i < len(list.List); i++ {
		o := list.List[i]
		seqs = append(seqs, o.SeqMember)
	}

	if userToken != nil {
		var listSeq []int64
		ldb := GetMyLogDbSlave(userToken.Allocated)
		switch _step {
		case 2:
			ldb.Model(schemas.MemberLikeStep2{}).Select("seq_novel_step2").Where("seq_member = ? AND seq_novel_step2 IN (?) AND like_yn = true", userToken.SeqMember, seqs).Scan(&listSeq)
		case 3:
			ldb.Model(schemas.MemberLikeStep3{}).Select("seq_novel_step3").Where("seq_member = ? AND seq_novel_step3 IN (?) AND like_yn = true", userToken.SeqMember, seqs).Scan(&listSeq)
		case 4:
			ldb.Model(schemas.MemberLikeStep4{}).Select("seq_novel_step4").Where("seq_member = ? AND seq_novel_step4 IN (?) AND like_yn = true", userToken.SeqMember, seqs).Scan(&listSeq)
		}

		for i := 0; i < len(list.List); i++ {
			o := list.List[i]
			for _, v := range listSeq {
				if v == o.SeqNovel {
					list.List[i].MyLike = true
					break
				}
			}
		}
		listMemberBlock := getBlockMemberList(userToken.Allocated, userToken.SeqMember, seqs)
		for i := 0; i < len(list.List); i++ {
			for _, v := range listMemberBlock {
				if v.SeqMember == list.List[i].SeqMember {
					list.List[i].BlockYn = v.BlockYn
					break
				}
			}
		}
	}

	res.Data = list

	return res
}

type NovelListRollingRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	Step      int `json:"step"`
	List      []struct {
		SeqNovel  int64   `json:"seq_novel"`
		SeqMember int64   `json:"seq_member"`
		NickName  string  `json:"nick_name"`
		CreatedAt float64 `json:"created_at"`
		CntLike   int64   `json:"cnt_like"`
		MyLike    bool    `json:"my_like"`
		Content   string  `json:"content"`
		BlockYn   bool    `json:"block_yn"`
	} `json:"list"`
}
