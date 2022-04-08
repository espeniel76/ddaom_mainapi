package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func MypageInfo(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqMember, _ := strconv.ParseInt(req.Vars["seq_member"], 10, 64)
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	var seqMemberToken int64
	if userToken != nil {
		seqMemberToken = userToken.SeqMember
	}
	data := make(map[string]interface{})
	if seqMemberToken == _seqMember {
		data["is_you"] = true
	} else {
		if _seqMember == 0 && userToken != nil {
			_seqMember = userToken.SeqMember
			data["is_you"] = true
		} else {
			data["is_you"] = false
		}
	}

	sdb := db.List[define.DSN_SLAVE1]

	// 닉네임, 프로필
	result := sdb.Model(schemas.MemberDetail{}).
		Where("seq_member = ?", _seqMember).
		Select("nick_name, profile_photo").Scan(&data)
	if corm(result, &res) {
		return res
	}
	if result.RowsAffected == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}

	// 임시저장, 작성완료
	var isTemps []bool
	query := `
		(SELECT temp_yn FROM novel_step1 WHERE seq_member = ? AND active_yn = true)
		UNION ALL
		(SELECT temp_yn FROM novel_step2 WHERE seq_member = ? AND active_yn = true)
		UNION ALL
		(SELECT temp_yn FROM novel_step3 WHERE seq_member = ? AND active_yn = true)
		UNION ALL
		(SELECT temp_yn FROM novel_step4 WHERE seq_member = ? AND active_yn = true)
	`
	result = sdb.Raw(query, _seqMember, _seqMember, _seqMember, _seqMember).Scan(&isTemps)
	if corm(result, &res) {
		return res
	}
	cntTemp := 0
	cntWrited := 0
	for _, v := range isTemps {
		if v == true {
			cntTemp++
		} else {
			cntWrited++
		}
	}
	data["cnt_temp"] = cntTemp
	data["cnt_writed"] = cntWrited

	// 팔로잉, 팔로워
	ldb := getUserLogDb(sdb, _seqMember)

	// 팔로잉
	var cntFollowing int64
	result = ldb.Model(schemas.MemberSubscribe{}).
		Where("seq_member = ?", _seqMember).
		Count(&cntFollowing)
	if corm(result, &res) {
		return res
	}
	data["cnt_following"] = cntFollowing

	// 팔로워
	var cntFollower int64
	result = ldb.Model(schemas.MemberSubscribe{}).
		Where("seq_member_following = ?", _seqMember).
		Count(&cntFollower)
	if corm(result, &res) {
		return res
	}
	data["cnt_follower"] = cntFollower

	res.Data = data

	return res
}
