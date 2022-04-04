package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func MypageInfo(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	masterDB := db.List[define.DSN_MASTER]
	data := make(map[string]interface{})

	// 닉네임, 프로필
	result := masterDB.Model(schemas.MemberDetail{}).Where("seq_member = ?", userToken.SeqMember).Select("nick_name, profile_photo").Scan(&data)
	if corm(result, &res) {
		return res
	}

	// 임시저장, 작성완료
	var isTemps []bool
	result = masterDB.Model(schemas.NovelStep1{}).Where("seq_member = ?", userToken.SeqMember).Select("temp_yn").Scan(&isTemps)
	if corm(result, &res) {
		return res
	}
	var cntTemp int
	var cntWrited int
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
	myLogDb := GetMyLogDb(userToken.Allocated)

	// 팔로잉
	var cntFollowing int64
	result = myLogDb.Model(schemas.MemberSubscribe{}).Where("seq_member = ?", userToken.SeqMember).Count(&cntFollowing)
	if corm(result, &res) {
		return res
	}
	data["cnt_following"] = cntFollowing

	// 팔로워
	var cntFollower int64
	result = myLogDb.Model(schemas.MemberSubscribe{}).Where("seq_member_following = ?", userToken.SeqMember).Count(&cntFollower)
	if corm(result, &res) {
		return res
	}
	data["cnt_follower"] = cntFollower

	res.Data = data

	return res
}
