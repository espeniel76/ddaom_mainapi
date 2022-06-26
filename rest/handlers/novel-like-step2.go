package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
	"strconv"
)

func NovelLikeStep2(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep2, _ := strconv.Atoi(req.Vars["seq_novel_step2"])
	fmt.Println(_seqNovelStep2)

	myLike := false
	var cnt int64
	var scanCount int64

	ldb := GetMyLogDbMaster(userToken.Allocated)
	mdb := db.List[define.Mconn.DsnMaster]

	result := mdb.Model(schemas.NovelStep2{}).Select("cnt_like").Where("seq_novel_step2 = ?", _seqNovelStep2).Scan(&cnt).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	MemberLikeStep2 := schemas.MemberLikeStep2{}
	result = ldb.Model(&MemberLikeStep2).
		Where("seq_novel_step2 = ? AND seq_member = ?", _seqNovelStep2, userToken.SeqMember).Scan(&MemberLikeStep2)
	seqKeyword := getSeqKeyword(2, int64(_seqNovelStep2))
	if corm(result, &res) {
		return res
	}

	// 작가 seq
	var seqMemberWriter int64
	mdb.Model(schemas.NovelStep2{}).Where("seq_novel_step2 = ?", _seqNovelStep2).Select("seq_member").Scan(&seqMemberWriter)

	if MemberLikeStep2.SeqMemberLike == 0 { // 존재하지 않음
		result = ldb.Create(&schemas.MemberLikeStep2{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep2: int64(_seqNovelStep2),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		result = mdb.Exec("UPDATE novel_step2 SET cnt_like = cnt_like + 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
		updateKeywordMemberLike(seqMemberWriter, seqKeyword, "PLUS")
	} else { // 존재함
		if MemberLikeStep2.LikeYn {
			result = ldb.Model(&schemas.MemberLikeStep2{}).
				Where("seq_member = ? AND seq_novel_step2 = ?", userToken.SeqMember, _seqNovelStep2).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step2 SET cnt_like = cnt_like - 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
			updateKeywordMemberLike(seqMemberWriter, seqKeyword, "MINUS")
		} else {
			result = ldb.Model(&schemas.MemberLikeStep2{}).
				Where("seq_member = ? AND seq_novel_step2 = ?", userToken.SeqMember, _seqNovelStep2).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step2 SET cnt_like = cnt_like + 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
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
		pushLike(2, int64(_seqNovelStep2), userToken.SeqMember)
	}

	cacheMainPopularWriter()

	return res
}
