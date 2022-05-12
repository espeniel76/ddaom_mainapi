package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func NovelLikeStep1(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep1, _ := strconv.Atoi(req.Vars["seq_novel_step1"])
	myLike := false
	var cnt int64
	var scanCount int64

	ldb := GetMyLogDb(userToken.Allocated)
	mdb := db.List[define.DSN_MASTER]

	// 소설 존재 여부
	result := mdb.Model(schemas.NovelStep1{}).Select("cnt_like").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&cnt).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 좋아요 상태 확인
	memberLikeStep1 := schemas.MemberLikeStep1{}
	result = ldb.Model(&memberLikeStep1).
		Where("seq_novel_step1 = ? AND seq_member = ?", _seqNovelStep1, userToken.SeqMember).Scan(&memberLikeStep1)
	seqKeyword := getSeqKeyword(1, int64(_seqNovelStep1))
	if corm(result, &res) {
		return res
	}
	if memberLikeStep1.SeqMemberLike == 0 { // 존재하지 않음
		// 1. 로그넣기
		result = ldb.Create(&schemas.MemberLikeStep1{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep1: int64(_seqNovelStep1),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		// 2. 좋아요 카운트 업데이트
		result = mdb.Exec("UPDATE novel_step1 SET cnt_like = cnt_like + 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
		mdb.Exec("UPDATE keywords SET cnt_like = cnt_like + 1 WHERE seq_keyword = ?", seqKeyword)
	} else { // 존재함
		// 1. 좋아요 상태
		if memberLikeStep1.LikeYn {
			result = ldb.Model(&schemas.MemberLikeStep1{}).
				Where("seq_member = ? AND seq_novel_step1 = ?", userToken.SeqMember, _seqNovelStep1).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step1 SET cnt_like = cnt_like - 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
			mdb.Exec("UPDATE keywords SET cnt_like = cnt_like - 1 WHERE seq_keyword = ?", seqKeyword)
		} else {
			result = ldb.Model(&schemas.MemberLikeStep1{}).
				Where("seq_member = ? AND seq_novel_step1 = ?", userToken.SeqMember, _seqNovelStep1).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = mdb.Exec("UPDATE novel_step1 SET cnt_like = cnt_like + 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
			if corm(result, &res) {
				return res
			}
			myLike = true
			cnt++
			mdb.Exec("UPDATE keywords SET cnt_like = cnt_like + 1 WHERE seq_keyword = ?", seqKeyword)
		}

	}

	data := make(map[string]interface{})
	data["my_like"] = myLike
	data["cnt_like"] = cnt
	res.Data = data

	return res
}
