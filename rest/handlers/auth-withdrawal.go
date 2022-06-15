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

	// 각종 탈퇴 프로세스 처리
	/*
			- 개인정보 파기 : 수집한 이메일, 회원정보수정으로 최종 저장된 이메일 (즉시)
		    - 재가입 즉시/반복 가능
		    - 닉네임 재사용 불가
		    - 진행중인 주제어에 작성한 글이 있는 경우, 삭제 처리 (admin에서 확인 가능 / 삭제된 소설로 처리)
		    - 완결 소설은 삭제하지 않으며, admin에서 삭제된 소설로 처리하지 않음
		    - 완결 소설 하단에 작가 정보에서 클릭 불가하도록 비활성화 처리
		    - Front에서는 완결 소설, 작가명만 노출됨 / admin에서 탈퇴 회원의 정보 개인정보 제외한 전체 확인 가능
		    - 메인 화면의 ‘인기 작가 리스트’에 있는 경우, 삭제 처리
		    - 탈퇴 회원이 다른 일반 회원을 구독한 건이 있는 경우, 삭제 처리 (받은 구독 수 제외, 마이페이지 구독 상세 리스트에서 삭제)
		    - 다른 일반 회원이 탈퇴 회원을 구독한 건이 있는 경우, 삭제 처리 (보낸 구독 수 제외, 마이페이지 구독 상세 리스트에서 삭제)
	*/

	return res
}
