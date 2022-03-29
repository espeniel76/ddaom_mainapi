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

	myLogDb := GetMyLogDb(userToken.Allocated)
	masterDB := db.List[define.DSN_MASTER]

	result := masterDB.Model(schemas.NovelStep4{}).Where("seq_novel_step4 = ?", _seqNovelStep4).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	MemberLikeStep4 := schemas.MemberLikeStep4{}
	result = myLogDb.Model(&MemberLikeStep4).
		Where("seq_novel_step4 = ? AND seq_member = ?", _seqNovelStep4, userToken.SeqMember).Scan(&MemberLikeStep4)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if MemberLikeStep4.SeqMemberLike == 0 { // 존재하지 않음
		result = myLogDb.Create(&schemas.MemberLikeStep4{
			SeqMember:     userToken.SeqMember,
			SeqNovelStep4: int64(_seqNovelStep4),
			LikeYn:        true,
		})
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		result = masterDB.Exec("UPDATE novel_step4 SET cnt_like = cnt_like + 1 WHERE seq_novel_step4 = ?", _seqNovelStep4)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		myLike = true
	} else { // 존재함
		if MemberLikeStep4.LikeYn {
			result = myLogDb.Model(&schemas.MemberLikeStep4{}).
				Where("seq_member = ? AND seq_novel_step4 = ?", userToken.SeqMember, _seqNovelStep4).
				Update("like_yn", false)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			result = masterDB.Exec("UPDATE novel_step4 SET cnt_like = cnt_like - 1 WHERE seq_novel_step4 = ?", _seqNovelStep4)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			myLike = false
		} else {
			result = myLogDb.Model(&schemas.MemberLikeStep4{}).
				Where("seq_member = ? AND seq_novel_step4 = ?", userToken.SeqMember, _seqNovelStep4).
				Update("like_yn", true)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			result = masterDB.Exec("UPDATE novel_step4 SET cnt_like = cnt_like + 1 WHERE seq_novel_step4 = ?", _seqNovelStep4)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			myLike = true
		}
	}

	data := make(map[string]bool)
	data["my_like"] = myLike
	res.Data = data

	return res
}
