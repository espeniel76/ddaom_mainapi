package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
)

func MypageViewLiveStep(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	// userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	// if err != nil {
	// 	res.ResultCode = define.INVALID_TOKEN
	// 	res.ErrorDesc = err.Error()
	// 	return res
	// }
	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_seqNovelStep4 := CpInt64(req.Parameters, "seq_novel_step4")
	var cntTotal int64

	ldb := db.List[define.DSN_SLAVE1]
	if _seqNovelStep1 > 0 {
		// 2,3,4
		result := ldb.Model(schemas.NovelStep2{}).Select("seq_novel_step2").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep2)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep3{}).Select("seq_novel_step3").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep3)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep4{}).Select("seq_novel_step4").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep4)
		if corm(result, &res) {
			return res
		}

	} else if _seqNovelStep2 > 0 {
		// 1,3,4
		result := ldb.Model(schemas.NovelStep1{}).Select("seq_novel_step1").Where("seq_novel_step2 = ?", _seqNovelStep2).Scan(&_seqNovelStep1)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep3{}).Select("seq_novel_step3").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep3)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep4{}).Select("seq_novel_step4").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep4)
		if corm(result, &res) {
			return res
		}
	} else if _seqNovelStep3 > 0 {
		// 1,2,4
		result := ldb.Model(schemas.NovelStep1{}).Select("seq_novel_step1").Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&_seqNovelStep3)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep2{}).Select("seq_novel_step2").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep2)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep4{}).Select("seq_novel_step4").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep4)
		if corm(result, &res) {
			return res
		}
	} else if _seqNovelStep4 > 0 {
		// 1,2,3
		result := ldb.Model(schemas.NovelStep4{}).Select("seq_novel_step1").Where("seq_novel_step4 = ?", _seqNovelStep4).Scan(&_seqNovelStep1)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep2{}).Select("seq_novel_step2").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep2)
		if corm(result, &res) {
			return res
		}
		result = ldb.Model(schemas.NovelStep3{}).Select("seq_novel_step3").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&_seqNovelStep3)
		if corm(result, &res) {
			return res
		}
	}
	fmt.Println(_seqNovelStep1, _seqNovelStep2, _seqNovelStep3, _seqNovelStep4)

	query := `
	SELECT
		ns.seq_novel_step1,
		ns.title,
		ns.created_at,
		ns.seq_genre,
		ns.seq_keyword,
		ns.seq_member,
		md.nick_name,
		ns.content,
		ns.cnt_like,
		ns.seq_image,
		ns.seq_color
	FROM novel_step1 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.seq_novel_step1 = ?`
	step1Res := Step1Res{}
	result := ldb.Raw(query, _seqNovelStep1).Scan(&step1Res)
	if corm(result, &res) {
		return res
	}
	data := make(map[string]interface{})
	data["title"] = step1Res.Title
	data["created_at"] = step1Res.CreatedAt.UnixMilli()
	data["seq_genre"] = step1Res.SeqGenre
	data["seq_keyword"] = step1Res.SeqKeyword
	data["seq_image"] = step1Res.SeqImage
	data["seq_color"] = step1Res.SeqColor

	step1 := make(map[string]interface{})
	step1["seq_novel_step1"] = step1Res.SeqNovelStep1
	step1["seq_member"] = step1Res.SeqMember
	step1["nick_name"] = step1Res.NickName
	step1["content"] = step1Res.Content
	step1["cnt_like"] = step1Res.CntLike
	step1["my_like"] = true
	data["step1"] = step1

	query = `
	SELECT
		ns.seq_novel_step2,
		ns.seq_member,
		md.nick_name,
		ns.created_at,
		ns.content,
		ns.cnt_like
	FROM novel_step2 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.seq_novel_step2 = ?`
	step2Res := Step2Res{}
	result = ldb.Raw(query, _seqNovelStep2).Scan(&step2Res)
	if corm(result, &res) {
		return res
	}
	if step2Res.SeqNovelStep2 > 0 {
		step := make(map[string]interface{})
		step["seq_novel_step2"] = step2Res.SeqNovelStep2
		step["seq_member"] = step2Res.SeqMember
		step["nick_name"] = step2Res.NickName
		step["created_at"] = step2Res.CreatedAt.UnixMilli()
		step["content"] = step2Res.Content
		step["cnt_like"] = step2Res.CntLike
		step["my_like"] = true
		result = ldb.Model(schemas.NovelStep2{}).Where("seq_novel_step1 = ?", _seqNovelStep2).Count(&cntTotal)
		if corm(result, &res) {
			return res
		}
		step["cnt_total"] = cntTotal
		data["step2"] = step
	} else {
		data["step2"] = nil
	}

	query = `
	SELECT
		ns.seq_novel_step3,
		ns.seq_member,
		md.nick_name,
		ns.created_at,
		ns.content,
		ns.cnt_like
	FROM novel_step3 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.seq_novel_step3 = ?`
	step3Res := Step3Res{}
	result = ldb.Raw(query, _seqNovelStep3).Scan(&step3Res)
	if corm(result, &res) {
		return res
	}
	if step3Res.SeqNovelStep3 > 0 {
		step := make(map[string]interface{})
		step["seq_novel_step3"] = step3Res.SeqNovelStep3
		step["seq_member"] = step3Res.SeqMember
		step["nick_name"] = step3Res.NickName
		step["created_at"] = step3Res.CreatedAt.UnixMilli()
		step["content"] = step3Res.Content
		step["cnt_like"] = step3Res.CntLike
		step["my_like"] = true
		result = ldb.Model(schemas.NovelStep3{}).Where("seq_novel_step1 = ?", _seqNovelStep3).Count(&cntTotal)
		if corm(result, &res) {
			return res
		}
		step["cnt_total"] = cntTotal
		data["step3"] = step
	} else {
		data["step3"] = nil
	}

	query = `
	SELECT
		ns.seq_novel_step4,
		ns.seq_member,
		md.nick_name,
		ns.created_at,
		ns.content,
		ns.cnt_like
	FROM novel_step4 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.seq_novel_step4 = ?`
	step4Res := Step4Res{}
	result = ldb.Raw(query, _seqNovelStep4).Scan(&step4Res)
	if corm(result, &res) {
		return res
	}
	if step4Res.SeqNovelStep4 > 0 {
		step := make(map[string]interface{})
		step["seq_novel_step4"] = step4Res.SeqNovelStep4
		step["seq_member"] = step4Res.SeqMember
		step["nick_name"] = step4Res.NickName
		step["created_at"] = step4Res.CreatedAt.UnixMilli()
		step["content"] = step4Res.Content
		step["cnt_like"] = step4Res.CntLike
		step["my_like"] = true
		result = ldb.Model(schemas.NovelStep4{}).Where("seq_novel_step1 = ?", _seqNovelStep4).Count(&cntTotal)
		if corm(result, &res) {
			return res
		}
		step["cnt_total"] = cntTotal
		data["step4"] = step
	} else {
		data["step4"] = nil
	}

	res.Data = data

	return res
}
