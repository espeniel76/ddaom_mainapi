package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
)

func MypageListFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqMember := CpInt64(req.Parameters, "seq_member")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.DSN_SLAVE]
	seq := _seqMember
	result := sdb.
		Model(schemas.NovelFinish{}).
		Where("active_yn = true").
		Where("seq_member_step1 = ? OR seq_member_step2 = ? OR seq_member_step3 = ? OR seq_member_step4 = ?", seq, seq, seq, seq).
		Count(&totalData)
	if corm(result, &res) {
		return res
	}

	novelListFinishRes := NovelListFinishRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	query := `
	SELECT
		nf.seq_novel_finish,
		ns.seq_genre,
		ns.seq_image,
		ns.seq_color,
		ns.title,
		false AS my_bookmark
	FROM novel_finishes nf
	INNER JOIN novel_step1 ns ON nf.seq_novel_step1 = ns.seq_novel_step1
	WHERE nf.active_yn = true
		AND nf.seq_member_step1 = ?
		OR nf.seq_member_step2 = ?
		OR nf.seq_member_step3 = ?
		OR nf.seq_member_step4 = ?
	ORDER BY nf.seq_novel_finish DESC
	LIMIT ?, ?
	`
	result = sdb.Raw(query, seq, seq, seq, seq, limitStart, _sizePerPage).Find(&novelListFinishRes.List)
	if corm(result, &res) {
		return res
	}

	if userToken != nil {
		var list []int64
		// get my bookmarks
		for _, v := range novelListFinishRes.List {
			list = append(list, v.SeqNovelFinish)
		}
		var list2 []int64
		ldb := GetMyLogDb(userToken.Allocated)
		result = ldb.Model(schemas.MemberBookmark{}).
			Select("seq_novel_finish").
			Where("seq_novel_finish IN (?) AND bookmark_yn = true AND seq_member = ?", list, userToken.SeqMember).
			Scan(&list2)
		if corm(result, &res) {
			return res
		}
		for i := 0; i < len(novelListFinishRes.List); i++ {
			for _, v := range list2 {
				if novelListFinishRes.List[i].SeqNovelFinish == v {
					novelListFinishRes.List[i].MyBookmark = true
					break
				}
			}
		}
	}

	res.Data = novelListFinishRes
	fmt.Println(novelListFinishRes)

	return res
}
