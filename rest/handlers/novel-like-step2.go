package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func NovelLikeStep2(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelStep2, _ := strconv.Atoi(req.Vars["seq_novel_step2"])
	myLike := false
	var cnt int64
	var scanCount int64

	myLogDb := GetMyLogDb(userToken.Allocated)
	masterDB := db.List[define.DSN_MASTER]

	result := masterDB.Model(schemas.NovelStep2{}).Select("cnt_like").Where("seq_novel_step2 = ?", _seqNovelStep2).Scan(&cnt).Count(&scanCount)
	if corm(result, &res) {
		return res
	}
	if scanCount == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	MemberLikeStep2 := schemas.MemberLikeStep2{}
	result = myLogDb.Model(&MemberLikeStep2).
		Where("seq_novel_step2 = ? AND seq_member = ?", _seqNovelStep2, userToken.SeqMember).Scan(&MemberLikeStep2)
	if corm(result, &res) {
		return res
	}
	if MemberLikeStep2.SeqMemberLike == 0 { // 존재하지 않음
		result = myLogDb.Create(&schemas.MemberLikeStep2{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep2: int64(_seqNovelStep2),
			LikeYn:        true,
		})
		if corm(result, &res) {
			return res
		}

		result = masterDB.Exec("UPDATE novel_step2 SET cnt_like = cnt_like + 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
		if corm(result, &res) {
			return res
		}
		myLike = true
		cnt++
	} else { // 존재함
		if MemberLikeStep2.LikeYn {
			result = myLogDb.Model(&schemas.MemberLikeStep2{}).
				Where("seq_member = ? AND seq_novel_step2 = ?", userToken.SeqMember, _seqNovelStep2).
				Update("like_yn", false)
			if corm(result, &res) {
				return res
			}
			result = masterDB.Exec("UPDATE novel_step2 SET cnt_like = cnt_like - 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
			if corm(result, &res) {
				return res
			}
			myLike = false
			cnt--
		} else {
			result = myLogDb.Model(&schemas.MemberLikeStep2{}).
				Where("seq_member = ? AND seq_novel_step2 = ?", userToken.SeqMember, _seqNovelStep2).
				Update("like_yn", true)
			if corm(result, &res) {
				return res
			}
			result = masterDB.Exec("UPDATE novel_step2 SET cnt_like = cnt_like + 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
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
