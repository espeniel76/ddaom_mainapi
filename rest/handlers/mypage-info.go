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
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)

	var seqMemberToken int64
	itsMe := false
	if userToken != nil {
		seqMemberToken = userToken.SeqMember
	}
	data := make(map[string]interface{})
	if seqMemberToken == _seqMember {
		data["is_you"] = true
	} else {
		if _seqMember == 0 && userToken != nil {
			itsMe = true
			_seqMember = userToken.SeqMember
			data["is_you"] = true
		} else {
			itsMe = false
			data["is_you"] = false
		}
	}

	sdb := db.List[define.Mconn.DsnSlave]

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

	var seqKeywords []int64
	sdb.Model(schemas.Keyword{}).Select("seq_keyword").Scan(&seqKeywords)

	// 임시저장, 작성완료
	var isTemps []bool
	query := ""
	if itsMe {
		query = `
		(SELECT ns1.temp_yn FROM novel_step1 ns1 WHERE ns1.seq_member = ? AND ns1.active_yn = true AND ns1.seq_keyword IN (?))
		UNION ALL
		(SELECT ns2.temp_yn FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns1.seq_keyword IN (?))
		UNION ALL
		(SELECT ns3.temp_yn FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns1.seq_keyword IN (?))
		UNION ALL
		(SELECT ns4.temp_yn FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns1.seq_keyword IN (?))
		`
	} else {
		query = `
		(SELECT ns1.temp_yn FROM novel_step1 ns1 WHERE ns1.seq_member = ? AND ns1.active_yn = true AND ns1.seq_keyword IN (?) AND ns1.deleted_yn = false)
		UNION ALL
		(SELECT ns2.temp_yn FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 WHERE ns2.seq_member = ? AND ns2.active_yn = true AND ns1.seq_keyword IN (?) AND ns2.deleted_yn = false)
		UNION ALL
		(SELECT ns3.temp_yn FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 WHERE ns3.seq_member = ? AND ns3.active_yn = true AND ns1.seq_keyword IN (?) AND ns3.deleted_yn = false)
		UNION ALL
		(SELECT ns4.temp_yn FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 WHERE ns4.seq_member = ? AND ns4.active_yn = true AND ns1.seq_keyword IN (?) AND ns4.deleted_yn = false)
		`
	}
	result = sdb.Raw(query, _seqMember, seqKeywords, _seqMember, seqKeywords, _seqMember, seqKeywords, _seqMember, seqKeywords).Scan(&isTemps)
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
	fmt.Println("사용자고유번호: ", _seqMember)
	ldb := getUserLogDbSlave(sdb, _seqMember)
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

	data["is_new_alarm"] = false
	data["cnt_alarm"] = 0
	if itsMe {
		// 읽지 않은 메시지 조회
		var cntAlarm int64
		sdb.Model(schemas.Alarm{}).Select("COUNT(*)").Where("seq_member = ? AND is_read = false", userToken.SeqMember).Scan(&cntAlarm)
		data["cnt_alarm"] = cntAlarm
		if cntAlarm > 0 {
			data["is_new_alarm"] = true
		}
	}

	status := getMySubscribe(userToken, _seqMember)
	data["my_subscribe"] = false
	if status == define.FOLLOWING || status == define.BOTH {
		data["my_subscribe"] = true
	}

	res.Data = data

	return res
}
