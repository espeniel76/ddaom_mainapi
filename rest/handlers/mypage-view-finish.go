package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func MypageViewFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_step, _ := strconv.Atoi(req.Vars["step"])
	_seqNovel, _ := strconv.Atoi(req.Vars["seq_novel"])

	query := `
	SELECT
		nf.seq_novel_finish,
		ns1.seq_novel_step1,
		ns1.title,
		ns1.seq_genre,
		ns1.seq_keyword ,
		ns1.seq_image,
		ns1.seq_color ,
		nf.cnt_like,
		ns1.cnt_view,
		ns1.seq_member AS seq_member_step1,
		(SELECT nick_name FROM member_details WHERE seq_member = ns1.seq_member) AS nick_name_step1,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step1) AS deleted_yn_step1,
		ns1.content AS content_step1,
		ns2.seq_member AS seq_member_step2,
		(SELECT nick_name FROM member_details WHERE seq_member = ns1.seq_member) AS nick_name_step2,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step2) AS deleted_yn_step2,
		ns2.content AS content_step2,
		ns3.seq_member AS seq_member_step3,
		(SELECT nick_name FROM member_details WHERE seq_member = ns1.seq_member) AS nick_name_step3,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step3) AS deleted_yn_step3,
		ns3.content AS content_step3,
		ns4.seq_member AS seq_member_step4,
		(SELECT nick_name FROM member_details WHERE seq_member = ns1.seq_member) AS nick_name_step4,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step4) AS deleted_yn_step4,
		ns4.content AS content_step4,
		nf.created_at
	FROM novel_finishes nf
	INNER JOIN novel_step1 ns1 ON nf.seq_novel_step1 = ns1.seq_novel_step1
	INNER JOIN novel_step2 ns2 ON nf.seq_novel_step2 = ns2.seq_novel_step2
	INNER JOIN novel_step3 ns3 ON nf.seq_novel_step3 = ns3.seq_novel_step3
	INNER JOIN novel_step4 ns4 ON nf.seq_novel_step4 = ns4.seq_novel_step4
	WHERE `
	switch _step {
	case 1:
		query += "nf.seq_novel_step1 = ?"
	case 2:
		query += "nf.seq_novel_step2 = ?"
	case 3:
		query += "nf.seq_novel_step3 = ?"
	case 4:
		query += "nf.seq_novel_step4 = ?"
	}
	sdb := db.List[define.DSN_SLAVE]
	n := NovelViewFinishData{}
	result := sdb.Raw(query, _seqNovel).Scan(&n)
	if corm(result, &res) {
		return res
	}

	ldb := GetMyLogDbSlave(userToken.Allocated)
	var cnt int64
	myLike := false
	myBookmark := false
	result = ldb.
		Model(schemas.MemberLikeStep1{}).
		Where("seq_member = ? AND seq_novel_step1 = ? AND like_yn = true", userToken.SeqMember, n.SeqNovelStep1).
		Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt > 0 {
		myLike = true
	}
	result = ldb.Model(schemas.MemberBookmark{}).
		Where("seq_member = ? AND seq_novel_finish = ? AND bookmark_yn = true", userToken.SeqMember, n.SeqNovelFinish).
		Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt > 0 {
		myBookmark = true
	}

	novelViewFinishRes := NovelViewFinishRes{
		SeqNovelFinish: n.SeqNovelFinish,
		Title:          n.Title,
		SeqGenre:       n.SeqGenre,
		SeqKeyword:     n.SeqKeyword,
		SeqImage:       n.SeqImage,
		SeqColor:       n.SeqColor,
		CntLike:        n.CntLike,
		CntView:        n.CntView,
		MyLike:         myLike,
		MyBookmark:     myBookmark,
		CreatedAt:      n.CreatedAt.UnixMilli(),
		Step1: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			Content   string "json:\"content\""
		}{
			SeqMember: n.SeqMemberStep1,
			NickName:  n.NickNameStep1,
			DeletedYn: n.DeletedYnStep1,
			Content:   n.ContentStep1,
		},
		Step2: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			Content   string "json:\"content\""
		}{
			SeqMember: n.SeqMemberStep2,
			NickName:  n.NickNameStep2,
			DeletedYn: n.DeletedYnStep2,
			Content:   n.ContentStep2,
		},
		Step3: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			Content   string "json:\"content\""
		}{
			SeqMember: n.SeqMemberStep3,
			NickName:  n.NickNameStep3,
			DeletedYn: n.DeletedYnStep3,
			Content:   n.ContentStep3,
		},
		Step4: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			Content   string "json:\"content\""
		}{
			SeqMember: n.SeqMemberStep4,
			NickName:  n.NickNameStep4,
			DeletedYn: n.DeletedYnStep4,
			Content:   n.ContentStep4,
		},
	}

	res.Data = novelViewFinishRes

	return res
}
