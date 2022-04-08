package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"

	"gorm.io/gorm"
)

func VisitInfo(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	_seqMember, _ := strconv.Atoi(req.Vars["seq_member"])

	sdb := db.List[define.DSN_SLAVE1]
	data := make(map[string]interface{})

	result := sdb.Model(schemas.MemberDetail{}).
		Where("seq_member = ?", _seqMember).Select("nick_name, profile_photo").
		Scan(&data)
	if corm(result, &res) {
		return res
	}

	//작성완료
	var cntWrited int64
	result = sdb.Model(schemas.NovelStep1{}).
		Where("seq_member = ? AND temp_yn = true", _seqMember).
		Select("COUNT(*) AS cnt").Scan(&cntWrited)
	if corm(result, &res) {
		return res
	}
	data["cnt_writed"] = cntWrited

	// 팔로잉, 팔로워
	allocatedDb := 1
	result = sdb.Model(schemas.Member{}).
		Select("allocated_db").
		Where("seq_member = ?", _seqMember).Scan(&allocatedDb)
	if corm(result, &res) {
		return res
	}
	ldb := getUserLogDb(sdb, int64(_seqMember))

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
	data["seq_member"] = _seqMember

	res.Data = data

	return res
}

func getUserLogDb(_db *gorm.DB, seqMember int64) *gorm.DB {
	allocatedDb := 1
	_db.Model(schemas.Member{}).
		Select("allocated_db").
		Where("seq_member = ?", seqMember).Scan(&allocatedDb)
	ldb := GetMyLogDb(int8(allocatedDb))
	return ldb
}
