package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func NovelLikeStep4(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep4, _ := strconv.Atoi(req.Vars["seq_novel_step4"])
	myLike := false
	var cnt int64
	var scanCount int64

	ldb := GetMyLogDbMaster(userToken.Allocated)
	mdb := db.List[define.DSN_MASTER]

	result := mdb.Model(schemas.NovelStep4{}).Select("cnt_like").Where("seq_novel_step4 = ?", _seqNovelStep4).Scan(&cnt).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	MemberLikeStep4 := schemas.MemberLikeStep4{}
	result = ldb.Model(&MemberLikeStep4).
		Where("seq_novel_step4 = ? AND seq_member = ?", _seqNovelStep4, userToken.SeqMember).Scan(&MemberLikeStep4)

	seqKeyword := getSeqKeyword(4, int64(_seqNovelStep4))
	if corm(result, &res) {
		return res
	}

	// 작가 seq
	var seqMemberWriter int64
	mdb.Model(schemas.NovelStep4{}).Where("seq_novel_step4 = ?", _seqNovelStep4).Select("seq_member").Scan(&seqMemberWriter)

	if MemberLikeStep4.SeqMemberLike == 0 { // 존재하지 않음
		result = ldb.Create(&schemas.MemberLikeStep4{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep4: int64(_seqNovelStep4),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		result = mdb.Exec("UPDATE novel_step4 SET cnt_like = cnt_like + 1 WHERE seq_novel_step4 = ?", _seqNovelStep4)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
		updateKeywordMemberLike(seqMemberWriter, seqKeyword, "PLUS")
	} else { // 존재함
		if MemberLikeStep4.LikeYn {
			result = ldb.Model(&schemas.MemberLikeStep4{}).
				Where("seq_member = ? AND seq_novel_step4 = ?", userToken.SeqMember, _seqNovelStep4).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step4 SET cnt_like = cnt_like - 1 WHERE seq_novel_step4 = ?", _seqNovelStep4)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
			updateKeywordMemberLike(seqMemberWriter, seqKeyword, "MINUS")
		} else {
			result = ldb.Model(&schemas.MemberLikeStep4{}).
				Where("seq_member = ? AND seq_novel_step4 = ?", userToken.SeqMember, _seqNovelStep4).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step4 SET cnt_like = cnt_like + 1 WHERE seq_novel_step4 = ?", _seqNovelStep4)
			if corm(result, &res) {
				return res
			}
			myLike = true
			cnt++
			updateKeywordMemberLike(seqMemberWriter, seqKeyword, "PLUS")
		}
	}

	data := make(map[string]interface{})
	data["my_like"] = myLike
	data["cnt_like"] = cnt
	res.Data = data

	// push 날리기
	if myLike {
		pushLike(4, int64(_seqNovelStep4), userToken.SeqMember)
	}

	return res
}
