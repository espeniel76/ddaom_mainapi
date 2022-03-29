package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
	"strconv"
)

func NovelBookmark(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelFinish, _ := strconv.Atoi(req.Vars["seq_novel_finish"])
	myBookmark := false
	var cnt int64

	myLogDb := GetMyLogDb(userToken.Allocated)
	masterDB := db.List[define.DSN_MASTER]

	// 소설 존재 여부
	fmt.Println(cnt)
	result := masterDB.Model(schemas.NovelFinish{}).Where("seq_novel_finish = ?", _seqNovelFinish).Count(&cnt)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 북마크 상태 확인
	memberBookmark := schemas.MemberBookmark{}
	result = myLogDb.Model(&memberBookmark).
		Where("seq_novel_finish = ? AND seq_member = ?", _seqNovelFinish, userToken.SeqMember).Scan(&memberBookmark)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if memberBookmark.SeqMemberBookmark == 0 { // 존재하지 않음
		// 1. 로그넣기
		result = myLogDb.Create(&schemas.MemberBookmark{
			SeqMember:      userToken.SeqMember,
			SeqNovelFinish: int64(_seqNovelFinish),
			BookmarkYn:     true,
		})
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}

		// 2. 북마크 카운트 업데이트
		result = masterDB.Exec("UPDATE novel_finishes SET cnt_bookmark = cnt_bookmark + 1 WHERE seq_novel_finish = ?", _seqNovelFinish)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		myBookmark = true
	} else { // 존재함
		// 1. 북마크 상태
		if memberBookmark.BookmarkYn {
			result = myLogDb.Model(&schemas.MemberBookmark{}).
				Where("seq_member = ? AND seq_novel_finish = ?", userToken.SeqMember, _seqNovelFinish).
				Update("bookmark_yn", false)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			result = masterDB.Exec("UPDATE novel_finishes SET cnt_bookmark = cnt_bookmark - 1 WHERE seq_novel_finish = ?", _seqNovelFinish)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			myBookmark = false
		} else {
			result = myLogDb.Model(&schemas.MemberBookmark{}).
				Where("seq_member = ? AND seq_novel_finish = ?", userToken.SeqMember, _seqNovelFinish).
				Update("bookmark_yn", true)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			result = masterDB.Exec("UPDATE novel_finishes SET cnt_bookmark = cnt_bookmark + 1 WHERE seq_novel_finish = ?", _seqNovelFinish)
			if result.Error != nil {
				res.ResultCode = define.DB_ERROR_ORM
				res.ErrorDesc = result.Error.Error()
				return res
			}
			myBookmark = true
		}
	}

	data := make(map[string]bool)
	data["my_bookmark"] = myBookmark
	res.Data = data

	return res
}
