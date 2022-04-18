package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
)

func ConfigAlarm(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_isNewKeyword := CpBool(req.Parameters, "is_new_keyword")
	_isLiked := CpBool(req.Parameters, "is_liked")
	_isFinished := CpBool(req.Parameters, "is_finished")
	_isNewFollower := CpBool(req.Parameters, "is_new_follower")
	_isNewFollowing := CpBool(req.Parameters, "is_new_following")
	_isNightPush := CpBool(req.Parameters, "is_night_push")

	mdb := db.List[define.DSN_MASTER]
	m := schemas.MemberDetail{}
	result := mdb.Model(&m).Where("seq_member = ?", userToken.SeqMember).Scan(&m)
	if corm(result, &res) {
		return res
	}
	m.IsNewKeyword = _isNewKeyword
	m.IsLiked = _isLiked
	m.IsFinished = _isFinished
	m.IsNewFollower = _isNewFollower
	m.IsNewFollowing = _isNewFollowing
	m.IsNightPush = _isNightPush
	result = mdb.Model(&m).Updates(&m)
	if corm(result, &res) {
		return res
	}

	return res
}

func ConfigAlarmGet(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	sdb := db.List[define.DSN_SLAVE]
	m := schemas.MemberDetail{}
	result := sdb.Model(&m).Where("seq_member = ?", userToken.SeqMember).Scan(&m)
	if corm(result, &res) {
		return res
	}
	o := make(map[string]bool)
	o["is_new_keyword"] = m.IsNewKeyword
	o["is_liked"] = m.IsLiked
	o["is_finished"] = m.IsFinished
	o["is_new_follower"] = m.IsNewFollower
	o["is_new_following"] = m.IsNewFollowing
	o["is_night_push"] = m.IsNightPush
	res.Data = o
	fmt.Println(o)

	return res
}
