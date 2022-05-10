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
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep3, _ := strconv.Atoi(req.Vars["seq_novel_step3"])
	myLike := false
	var cnt int64
	var scanCount int64

	myLogDb := GetMyLogDb(userToken.Allocated)
	masterDB := db.List[define.DSN_MASTER]

	result := masterDB.Model(schemas.NovelStep3{}).Select("cnt_like").Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&cnt).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	MemberLikeStep3 := schemas.MemberLikeStep3{}
	result = myLogDb.Model(&MemberLikeStep3).
		Where("seq_novel_step3 = ? AND seq_member = ?", _seqNovelStep3, userToken.SeqMember).Scan(&MemberLikeStep3)
	if corm(result, &res) {
		return res
	}
	if MemberLikeStep3.SeqMemberLike == 0 { // 존재하지 않음
		result = myLogDb.Create(&schemas.MemberLikeStep3{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep3: int64(_seqNovelStep3),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		result = masterDB.Exec("UPDATE novel_step3 SET cnt_like = cnt_like + 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
	} else { // 존재함
		if MemberLikeStep3.LikeYn {
			result = myLogDb.Model(&schemas.MemberLikeStep3{}).
				Where("seq_member = ? AND seq_novel_step3 = ?", userToken.SeqMember, _seqNovelStep3).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = masterDB.Exec("UPDATE novel_step3 SET cnt_like = cnt_like - 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
		} else {
			result = myLogDb.Model(&schemas.MemberLikeStep3{}).
				Where("seq_member = ? AND seq_novel_step3 = ?", userToken.SeqMember, _seqNovelStep3).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = masterDB.Exec("UPDATE novel_step3 SET cnt_like = cnt_like + 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
			if corm(result, &res) {
				return res
			}
			myLike = true
			cnt++
		}
	}

	data := make(map[string]interface{})
	data["my_like"] = myLike
	data["cnt_like"] = cnt
	res.Data = data

	return res
}
