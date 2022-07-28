package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
)

func NovelReportLive(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	mdb := db.List[define.Mconn.DsnMaster]

	_step := CpInt64(req.Parameters, "step")
	_seqNovel := CpInt64(req.Parameters, "seq_novel")
	_reasonType := CpInt64(req.Parameters, "reason_type")
	_reason := Cp(req.Parameters, "reason")

	fmt.Println(_step, _seqNovel, _reasonType, _reason)

	novelReport := schemas.NovelReport{
		SeqMember:  userToken.SeqMember,
		NovelType:  "LIVE",
		Step:       int8(_step),
		SeqNovel:   _seqNovel,
		ReasonType: int8(_reasonType),
		Reason:     _reason,
	}
	result := mdb.Create(&novelReport)
	if corm(result, &res) {
		return res
	}

	return res
}

func NovelReportFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	mdb := db.List[define.Mconn.DsnMaster]

	_step := CpInt64(req.Parameters, "step")
	_seqNovelFinish := CpInt64(req.Parameters, "seq_novel_finish")
	_reasonType := CpInt64(req.Parameters, "reason_type")
	_reason := Cp(req.Parameters, "reason")

	novelReport := schemas.NovelReport{
		SeqMember:  userToken.SeqMember,
		NovelType:  "FINISH",
		Step:       int8(_step),
		SeqNovel:   _seqNovelFinish,
		ReasonType: int8(_reasonType),
		Reason:     _reason,
	}
	result := mdb.Create(&novelReport)
	if corm(result, &res) {
		return res
	}

	return res
}

func MypageUserReport(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	mdb := db.List[define.Mconn.DsnMaster]

	_seqMember := CpInt64(req.Parameters, "seq_member")
	_reasonType := CpInt64(req.Parameters, "reason_type")
	_reason := Cp(req.Parameters, "reason")

	novelReport := schemas.MemberReport{
		SeqMember:   userToken.SeqMember,
		SeqMemberTo: _seqMember,
		ReasonType:  int8(_reasonType),
		Reason:      _reason,
	}
	result := mdb.Create(&novelReport)
	if corm(result, &res) {
		return res
	}

	return res
}
