package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
)

func AuthWithdrawal(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	mdb := db.List[define.DSN_MASTER]

	// 데이터 존재 여부 체크
	isExist := db.ExistRow(mdb, "members", "email", userToken.Email)
	if !isExist {
		res.ResultCode = define.NO_EXIST_USER
		return res
	}

	// 사용자 백업
	query := `
		INSERT INTO member_backups (
			seq_member,
			email,
			token,
			profile_image_url,
			sns_type,
			active_yn,
			user_level,
			allocated_db,
			created_at,
			updated_at,
			push_token,
			deleted_at
		)
		SELECT
			seq_member,
			email,
			token,
			profile_image_url,
			sns_type,
			active_yn,
			user_level,
			allocated_db,
			created_at,
			updated_at,
			push_token,
			NOW()
		FROM
			members
		WHERE
			seq_member = ?
	`
	result := mdb.Exec(query, userToken.SeqMember)
	if corm(result, &res) {
		return res
	}
	query = `
		INSERT INTO member_detail_backups (
			seq_member_detail,
			seq_member,
			email,
			name,
			nick_name,
			profile_photo,
			tel,
			mobile_company,
			mobile,
			address,
			address_detail,
			zipcode,
			authentication_ci,
			authentication_at,
			is_new_keyword,
			is_liked,
			is_finished,
			is_new_follower,
			is_new_following,
			is_night_push,
			is_deleted,
			cnt_subscribe,
			cnt_like,
			created_at,
			updated_at,
			deleted_at
		)
		SELECT
			seq_member_detail,
			seq_member,
			email,
			name,
			nick_name,
			profile_photo,
			tel,
			mobile_company,
			mobile,
			address,
			address_detail,
			zipcode,
			authentication_ci,
			authentication_at,
			is_new_keyword,
			is_liked,
			is_finished,
			is_new_follower,
			is_new_following,
			is_night_push,
			is_deleted,
			cnt_subscribe,
			cnt_like,
			created_at,
			updated_at,
			NOW()
		FROM
			member_details
		WHERE
			seq_member = ?
	`
	result = mdb.Exec(query, userToken.SeqMember)
	if corm(result, &res) {
		return res
	}

	// 원본 데이터 삭제
	result = mdb.Where("seq_member = ?", userToken.SeqMember).Delete(&schemas.Member{})
	if corm(result, &res) {
		return res
	}
	result = mdb.Where("seq_member = ?", userToken.SeqMember).Delete(&schemas.MemberDetail{})
	if corm(result, &res) {
		return res
	}

	return res
}
