package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func MypageListTempDelete(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	_step := CpInt64(req.Parameters, "step")
	_seqNovel := CpInt64(req.Parameters, "seq_novel")
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	mdb := db.List[define.Mconn.DsnMaster]
	switch _step {
	case 1:
		mdb.
			Where("seq_member = ? AND seq_novel_step1 = ?", userToken.SeqMember, _seqNovel).
			Delete(&schemas.NovelStep1{})
	case 2:
		mdb.
			Where("seq_member = ? AND seq_novel_step2 = ?", userToken.SeqMember, _seqNovel).
			Delete(&schemas.NovelStep2{})
	case 3:
		mdb.
			Where("seq_member = ? AND seq_novel_step3 = ?", userToken.SeqMember, _seqNovel).
			Delete(&schemas.NovelStep3{})
	case 4:
		mdb.
			Where("seq_member = ? AND seq_novel_step4 = ?", userToken.SeqMember, _seqNovel).
			Delete(&schemas.NovelStep4{})
	}

	return res
}
