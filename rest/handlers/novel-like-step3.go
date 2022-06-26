package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func NovelLikeStep3(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep3, _ := strconv.Atoi(req.Vars["seq_novel_step3"])
	myLike := false
	var cnt int64
	var scanCount int64

	ldb := GetMyLogDbMaster(userToken.Allocated)
	mdb := db.List[define.Mconn.DsnMaster]

	result := mdb.Model(schemas.NovelStep3{}).Select("cnt_like").Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&cnt).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	MemberLikeStep3 := schemas.MemberLikeStep3{}
	result = ldb.Model(&MemberLikeStep3).
		Where("seq_novel_step3 = ? AND seq_member = ?", _seqNovelStep3, userToken.SeqMember).Scan(&MemberLikeStep3)
	seqKeyword := getSeqKeyword(3, int64(_seqNovelStep3))
	if corm(result, &res) {
		return res
	}

	// 작가 seq
	var seqMemberWriter int64
	mdb.Model(schemas.NovelStep3{}).Where("seq_novel_step3 = ?", _seqNovelStep3).Select("seq_member").Scan(&seqMemberWriter)

	if MemberLikeStep3.SeqMemberLike == 0 { // 존재하지 않음
		result = ldb.Create(&schemas.MemberLikeStep3{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep3: int64(_seqNovelStep3),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		result = mdb.Exec("UPDATE novel_step3 SET cnt_like = cnt_like + 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
		updateKeywordMemberLike(seqMemberWriter, seqKeyword, "PLUS")
	} else { // 존재함
		if MemberLikeStep3.LikeYn {
			result = ldb.Model(&schemas.MemberLikeStep3{}).
				Where("seq_member = ? AND seq_novel_step3 = ?", userToken.SeqMember, _seqNovelStep3).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step3 SET cnt_like = cnt_like - 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
			updateKeywordMemberLike(seqMemberWriter, seqKeyword, "MINUS")
		} else {
			result = ldb.Model(&schemas.MemberLikeStep3{}).
				Where("seq_member = ? AND seq_novel_step3 = ?", userToken.SeqMember, _seqNovelStep3).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step3 SET cnt_like = cnt_like + 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
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
		pushLike(3, int64(_seqNovelStep3), userToken.SeqMember)
	}

	return res
}
