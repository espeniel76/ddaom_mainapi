package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func AuthDormancy(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_email := Cp(req.Parameters, "email")

	mdb := db.List[define.Mconn.DsnMaster]
	var seqMember int64
	result := mdb.Model(schemas.Member{}).Select("seq_member").Where("email = ?", _email).Scan(&seqMember)
	if corm(result, &res) {
		return res
	}
	memberDormacy := schemas.MemberDormacy{
		SeqMember: seqMember,
		DormacyYn: false,
	}
	result = mdb.Create(&memberDormacy)
	if corm(result, &res) {
		return res
	}
	result = mdb.Model(schemas.Member{}).Where("seq_member = ?", seqMember).Update("dormacy_yn", false)
	if corm(result, &res) {
		return res
	}

	// 회원 상태 (휴면 -> 정상)
	go setUserActionLog(seqMember, 3, "")

	return res
}
