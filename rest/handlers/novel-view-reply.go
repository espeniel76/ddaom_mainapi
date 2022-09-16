package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
	"unicode/utf8"
)

/** 댓글 작성 */
func NovelViewReplyWrite(req *domain.CommonRequest) domain.CommonResponse {
	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelFinish := CpInt64(req.Parameters, "seq_novel_finish")
	_step := CpInt64(req.Parameters, "step")
	_seqNovel := CpInt64(req.Parameters, "seq_novel")
	_contents := Cp(req.Parameters, "contents")

	if utf8.RuneCountInString(_contents) > 150 {
		res.ResultCode = define.TEXT_LIMIT_EXCEEDED
		return res
	}

	mdb := db.List[define.Mconn.DsnMaster]

	novelReply := &schemas.NovelReply{
		Step:      int8(_step),
		SeqNovel:  _seqNovel,
		SeqMember: userToken.SeqMember,
		Contents:  _contents,
	}
	result := mdb.Create(novelReply)
	if corm(result, &res) {
		return res
	}

	go setReplyCnt(_seqNovelFinish, _step, _seqNovel)

	return res
}

/** 댓글의 답글 작성 */
func NovelViewReReplyWrite(req *domain.CommonRequest) domain.CommonResponse {
	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	fmt.Println(userToken)

	_seqNovelFinish := CpInt64(req.Parameters, "seq_novel_finish")
	_step := CpInt64(req.Parameters, "step")
	_seqNovel := CpInt64(req.Parameters, "seq_novel")
	_seqReply := CpInt64(req.Parameters, "seq_reply")
	_contents := Cp(req.Parameters, "contents")

	fmt.Println(_seqNovelFinish, _step, _seqNovel, _seqReply, _contents)

	if utf8.RuneCountInString(_contents) > 150 {
		res.ResultCode = define.TEXT_LIMIT_EXCEEDED
		return res
	}

	mdb := db.List[define.Mconn.DsnMaster]
	novelReReply := &schemas.NovelReReply{
		SeqReply:  _seqReply,
		SeqMember: userToken.SeqMember,
		Contents:  _contents,
	}
	result := mdb.Create(novelReReply)
	if corm(result, &res) {
		return res
	}

	go setReReplyCnt(_seqNovel)
	go setReplyCnt(_seqNovelFinish, _step, _seqNovel)

	return res
}

/** 해당 글의 답글 카운트 1 증가 */
func setReReplyCnt(_seqReply int64) {
	mdb := db.List[define.Mconn.DsnMaster]
	sql := "UPDATE novel_replies SET cnt_re_reply = cnt_re_reply + 1 WHERE seq_reply = ?"
	mdb.Exec(sql, _seqReply)
}

/** 해당 글의 댓글 카운트 1 증가 */
func setReplyCnt(_seqNovelFinish int64, _step int64, _seqNovel int64) {
	mdb := db.List[define.Mconn.DsnMaster]
	sql := ""
	if _seqNovelFinish == 0 {
		switch _step {
		case 1:
			sql = "UPDATE novel_step1 SET cnt_reply = cnt_reply + 1 WHERE seq_novel_step1 = ?"
		case 2:
			sql = "UPDATE novel_step2 SET cnt_reply = cnt_reply + 1 WHERE seq_novel_step2 = ?"
		case 3:
			sql = "UPDATE novel_step3 SET cnt_reply = cnt_reply + 1 WHERE seq_novel_step3 = ?"
		case 4:
			sql = "UPDATE novel_step4 SET cnt_reply = cnt_reply + 1 WHERE seq_novel_step4 = ?"
		}
		mdb.Exec(sql, _seqNovel)
	} else {
		switch _step {
		case 1:
			sql = "UPDATE novel_finishes SET cnt_reply_step1 = cnt_reply_step1 + 1 WHERE seq_novel_finish = ?"
		case 2:
			sql = "UPDATE novel_finishes SET cnt_reply_step2 = cnt_reply_step2 + 1 WHERE seq_novel_finish = ?"
		case 3:
			sql = "UPDATE novel_finishes SET cnt_reply_step3 = cnt_reply_step3 + 1 WHERE seq_novel_finish = ?"
		case 4:
			sql = "UPDATE novel_finishes SET cnt_reply_step4 = cnt_reply_step4 + 1 WHERE seq_novel_finish = ?"
		}
		mdb.Exec(sql, _seqNovelFinish)
	}
}
