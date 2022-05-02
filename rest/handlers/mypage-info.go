package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
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

	sdb := db.List[define.DSN_SLAVE]

	// 닉네임, 프로필
	result := sdb.Model(schemas.MemberDetail{}).
		Where("seq_member = ?", _seqMember).
		Select("nick_name, profile_photo, seq_member").Scan(&data)
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

	// 구독현황
	ldb := getUserLogDb(sdb, _seqMember)
	fmt.Println(ldb)
	listStatus := []string{}
	result = ldb.Model(&schemas.MemberSubscribe{}).Select("status").
		Where("seq_member = ?", _seqMember).Scan(&listStatus)
	if corm(result, &res) {
		return res
	}
	cntFollower := 0
	cntFollowing := 0
	for _, v := range listStatus {
		switch v {
		case define.FOLLOWING:
			cntFollowing++
		case define.FOLLOWER:
			cntFollower++
		case define.BOTH:
			cntFollowing++
			cntFollower++
		}
	}
	data["cnt_following"] = cntFollowing
	data["cnt_follower"] = cntFollower

	data["is_new_alarm"] = true
	data["my_subscribe"] = false

	res.Data = data

	return res
}
