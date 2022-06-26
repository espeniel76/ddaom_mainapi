package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"time"
)

func NovelViewFinish(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqNovelFinish, _ := req.Vars["seq_novel_finish"]

	query := `
	SELECT
		seq_novel_finish,
		seq_novel_step1,
		seq_novel_step2,
		seq_novel_step3,
		seq_novel_step4,
		title,
		seq_genre,
		seq_keyword ,
		seq_image,
		seq_color ,
		cnt_like,
		cnt_view,
		seq_member_step1,
		nick_name_step1,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step1) AS deleted_yn_step1,
		content1 AS content_step1,
		seq_member_step2,
		nick_name_step2,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step2) AS deleted_yn_step2,
		content2 AS content_step2,
		seq_member_step3,
		nick_name_step3,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step3) AS deleted_yn_step3,
		content3 AS content_step3,
		seq_member_step4,
		nick_name_step4,
		(SELECT deleted_yn FROM members WHERE seq_member = seq_member_step4) AS deleted_yn_step4,
		content4 AS content_step4,
		created_at
	FROM novel_finishes nf
	WHERE seq_novel_finish = ?
	`
	sdb := db.List[define.Mconn.DsnSlave]
	n := NovelViewFinishData{}
	result := sdb.Raw(query, _seqNovelFinish).Scan(&n)
	if corm(result, &res) {
		return res
	}

	ldb := GetMyLogDbSlave(userToken.Allocated)
	var cnt1 int64
	var cnt2 int64
	var cnt3 int64
	var cnt4 int64
	var cnt int64
	myLike := false
	myBookmark := false
	ldb.Model(schemas.MemberLikeStep1{}).Where("seq_member = ? AND seq_novel_step1 = ? AND like_yn = true", userToken.SeqMember, n.SeqNovelStep1).Count(&cnt1)
	ldb.Model(schemas.MemberLikeStep2{}).Where("seq_member = ? AND seq_novel_step2 = ? AND like_yn = true", userToken.SeqMember, n.SeqNovelStep2).Count(&cnt2)
	ldb.Model(schemas.MemberLikeStep3{}).Where("seq_member = ? AND seq_novel_step3 = ? AND like_yn = true", userToken.SeqMember, n.SeqNovelStep3).Count(&cnt3)
	ldb.Model(schemas.MemberLikeStep4{}).Where("seq_member = ? AND seq_novel_step4 = ? AND like_yn = true", userToken.SeqMember, n.SeqNovelStep4).Count(&cnt4)
	if cnt1 > 0 || cnt2 > 0 || cnt3 > 0 || cnt4 > 0 {
		myLike = true
	}
	result = ldb.Model(schemas.MemberBookmark{}).Where("seq_member = ? AND seq_novel_finish = ? AND bookmark_yn = true", userToken.SeqMember, _seqNovelFinish).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt > 0 {
		myBookmark = true
	}

	mdb := db.List[define.Mconn.DsnMaster]
	query = "UPDATE novel_finishes SET cnt_view = cnt_view + 1 WHERE seq_novel_finish = ?"
	mdb.Exec(query, _seqNovelFinish)

	novelViewFinishRes := NovelViewFinishRes{
		SeqNovelFinish: n.SeqNovelFinish,
		Title:          n.Title,
		SeqGenre:       n.SeqGenre,
		SeqKeyword:     n.SeqKeyword,
		SeqImage:       n.SeqImage,
		SeqColor:       n.SeqColor,
		CntLike:        n.CntLike,
		CntView:        n.CntView + 1,
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

	cacheMainPopular()

	return res
}

type NovelViewFinishData struct {
	SeqNovelFinish int64     `json:"seq_novel_finish"`
	Title          string    `json:"title"`
	SeqNovelStep1  int64     `json:"seq_novel_step1"`
	SeqNovelStep2  int64     `json:"seq_novel_step2"`
	SeqNovelStep3  int64     `json:"seq_novel_step3"`
	SeqNovelStep4  int64     `json:"seq_novel_step4"`
	SeqGenre       int64     `json:"seq_genre"`
	SeqKeyword     int64     `json:"seq_keyword"`
	SeqImage       int64     `json:"seq_image"`
	SeqColor       int64     `json:"seq_color"`
	CntLike        int64     `json:"cnt_like"`
	CntView        int64     `json:"cnt_view"`
	MyLike         bool      `json:"my_like"`
	MyBookmark     bool      `json:"my_bookmark"`
	CreatedAt      time.Time `json:"created_at"`
	SeqMemberStep1 int64     `json:"seq_member_step1"`
	NickNameStep1  string    `json:"nick_name_step1"`
	DeletedYnStep1 bool      `json:"deleted_yn_step1"`
	ContentStep1   string    `json:"content_step1"`
	SeqMemberStep2 int64     `json:"seq_member_step2"`
	NickNameStep2  string    `json:"nick_name_step2"`
	DeletedYnStep2 bool      `json:"deleted_yn_step2"`
	ContentStep2   string    `json:"content_step2"`
	SeqMemberStep3 int64     `json:"seq_member_step3"`
	NickNameStep3  string    `json:"nick_name_step3"`
	DeletedYnStep3 bool      `json:"deleted_yn_step3"`
	ContentStep3   string    `json:"content_step3"`
	SeqMemberStep4 int64     `json:"seq_member_step4"`
	NickNameStep4  string    `json:"nick_name_step4"`
	DeletedYnStep4 bool      `json:"deleted_yn_step4"`
	ContentStep4   string    `json:"content_step4"`
}

type NovelViewFinishRes struct {
	SeqNovelFinish int64  `json:"seq_novel_finish"`
	Title          string `json:"title"`
	SeqGenre       int64  `json:"seq_genre"`
	SeqKeyword     int64  `json:"seq_keyword"`
	SeqImage       int64  `json:"seq_image"`
	SeqColor       int64  `json:"seq_color"`
	CntLike        int64  `json:"cnt_like"`
	CntView        int64  `json:"cnt_view"`
	MyLike         bool   `json:"my_like"`
	MyBookmark     bool   `json:"my_bookmark"`
	CreatedAt      int64  `json:"created_at"`
	Step1          struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		Content   string `json:"content"`
	} `json:"step1"`
	Step2 struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		Content   string `json:"content"`
	} `json:"step2"`
	Step3 struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		Content   string `json:"content"`
	} `json:"step3"`
	Step4 struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		Content   string `json:"content"`
	} `json:"step4"`
}
