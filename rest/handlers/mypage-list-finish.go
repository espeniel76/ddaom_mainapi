package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

func MypageListFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	slaveDb := db.List[define.DSN_SLAVE1]
	seq := userToken.SeqMember
	result := slaveDb.
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
		true AS my_bookmark
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
	result = slaveDb.Raw(query, seq, seq, seq, seq, limitStart, _sizePerPage).Find(&novelListFinishRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = novelListFinishRes

	return res
}
