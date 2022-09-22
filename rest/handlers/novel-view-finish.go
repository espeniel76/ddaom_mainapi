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
		created_at,
		cnt_reply_step1,
		cnt_reply_step2,
		cnt_reply_step3,
		cnt_reply_step4
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
	var cnt int64
	myLike := false
	myBookmark := false
	query = `
	SELECT
		(SELECT COUNT(*) FROM member_like_step1 mls WHERE seq_member = ? AND seq_novel_step1 = ? AND like_yn = true) AS cnt1,
		(SELECT COUNT(*) FROM member_like_step2 mls WHERE seq_member = ? AND seq_novel_step2 = ? AND like_yn = true) AS cnt2,
		(SELECT COUNT(*) FROM member_like_step3 mls WHERE seq_member = ? AND seq_novel_step3 = ? AND like_yn = true) AS cnt3,
		(SELECT COUNT(*) FROM member_like_step4 mls WHERE seq_member = ? AND seq_novel_step4 = ? AND like_yn = true) AS cnt4
	`
	var cntAll CntAll
	ldb.Raw(query,
		userToken.SeqMember, n.SeqNovelStep1,
		userToken.SeqMember, n.SeqNovelStep2,
		userToken.SeqMember, n.SeqNovelStep3,
		userToken.SeqMember, n.SeqNovelStep4).Scan(&cntAll)
	if cntAll.Cnt1 > 0 || cntAll.Cnt2 > 0 || cntAll.Cnt3 > 0 || cntAll.Cnt4 > 0 {
		myLike = true
	}

	query = `
	SELECT
		(SELECT blocked_yn FROM members WHERE seq_member = ?) AS step1,
		(SELECT blocked_yn FROM members WHERE seq_member = ?) AS step2,
		(SELECT blocked_yn FROM members WHERE seq_member = ?) AS step3,
		(SELECT blocked_yn FROM members WHERE seq_member = ?) AS step4
	`
	var blockedAll BlockedYnAll
	sdb.Raw(query,
		n.SeqMemberStep1,
		n.SeqMemberStep2,
		n.SeqMemberStep3,
		n.SeqMemberStep4).Scan(&blockedAll)

	result = ldb.Model(schemas.MemberBookmark{}).Where("seq_member = ? AND seq_novel_finish = ? AND bookmark_yn = true", userToken.SeqMember, _seqNovelFinish).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt > 0 {
		myBookmark = true
	}

	mdb := db.List[define.Mconn.DsnMaster]
	mdb.Exec("UPDATE novel_finishes SET cnt_view = cnt_view + 1 WHERE seq_novel_finish = ?", _seqNovelFinish)

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
			BlockedYn bool   "json:\"blocked_yn\""
			BlockYn   bool   "json:\"block_yn\""
			Content   string "json:\"content\""
			CntReply  int64  "json:\"cnt_reply\""
		}{
			SeqMember: n.SeqMemberStep1,
			NickName:  n.NickNameStep1,
			DeletedYn: n.DeletedYnStep1,
			BlockedYn: blockedAll.Step1,
			BlockYn:   isBlockMember(userToken.Allocated, userToken.SeqMember, n.SeqMemberStep1),
			Content:   n.ContentStep1,
			CntReply:  n.CntReplyStep1,
		},
		Step2: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			BlockedYn bool   "json:\"blocked_yn\""
			BlockYn   bool   "json:\"block_yn\""
			Content   string "json:\"content\""
			CntReply  int64  "json:\"cnt_reply\""
		}{
			SeqMember: n.SeqMemberStep2,
			NickName:  n.NickNameStep2,
			DeletedYn: n.DeletedYnStep2,
			BlockedYn: blockedAll.Step2,
			BlockYn:   isBlockMember(userToken.Allocated, userToken.SeqMember, n.SeqMemberStep2),
			Content:   n.ContentStep2,
			CntReply:  n.CntReplyStep2,
		},
		Step3: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			BlockedYn bool   "json:\"blocked_yn\""
			BlockYn   bool   "json:\"block_yn\""
			Content   string "json:\"content\""
			CntReply  int64  "json:\"cnt_reply\""
		}{
			SeqMember: n.SeqMemberStep3,
			NickName:  n.NickNameStep3,
			DeletedYn: n.DeletedYnStep3,
			BlockedYn: blockedAll.Step3,
			BlockYn:   isBlockMember(userToken.Allocated, userToken.SeqMember, n.SeqMemberStep3),
			Content:   n.ContentStep3,
			CntReply:  n.CntReplyStep3,
		},
		Step4: struct {
			SeqMember int64  "json:\"seq_member\""
			NickName  string "json:\"nick_name\""
			DeletedYn bool   "json:\"deleted_yn\""
			BlockedYn bool   "json:\"blocked_yn\""
			BlockYn   bool   "json:\"block_yn\""
			Content   string "json:\"content\""
			CntReply  int64  "json:\"cnt_reply\""
		}{
			SeqMember: n.SeqMemberStep4,
			NickName:  n.NickNameStep4,
			DeletedYn: n.DeletedYnStep4,
			BlockedYn: blockedAll.Step4,
			BlockYn:   isBlockMember(userToken.Allocated, userToken.SeqMember, n.SeqMemberStep4),
			Content:   n.ContentStep4,
			CntReply:  n.CntReplyStep4,
		},
	}

	res.Data = novelViewFinishRes

	cacheMainPopular()

	return res
}

type CntAll struct {
	Cnt1 int64 `json:"cnt1"`
	Cnt2 int64 `json:"cnt2"`
	Cnt3 int64 `json:"cnt3"`
	Cnt4 int64 `json:"cnt4"`
}

type BlockedYnAll struct {
	Step1 bool `json:"step1"`
	Step2 bool `json:"step2"`
	Step3 bool `json:"step3"`
	Step4 bool `json:"step4"`
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
	CntReplyStep1  int64     `json:"cnt_reply_step1"`
	SeqMemberStep2 int64     `json:"seq_member_step2"`
	NickNameStep2  string    `json:"nick_name_step2"`
	DeletedYnStep2 bool      `json:"deleted_yn_step2"`
	ContentStep2   string    `json:"content_step2"`
	CntReplyStep2  int64     `json:"cnt_reply_step2"`
	SeqMemberStep3 int64     `json:"seq_member_step3"`
	NickNameStep3  string    `json:"nick_name_step3"`
	DeletedYnStep3 bool      `json:"deleted_yn_step3"`
	ContentStep3   string    `json:"content_step3"`
	CntReplyStep3  int64     `json:"cnt_reply_step3"`
	SeqMemberStep4 int64     `json:"seq_member_step4"`
	NickNameStep4  string    `json:"nick_name_step4"`
	DeletedYnStep4 bool      `json:"deleted_yn_step4"`
	ContentStep4   string    `json:"content_step4"`
	CntReplyStep4  int64     `json:"cnt_reply_step4"`
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
		BlockedYn bool   `json:"blocked_yn"`
		BlockYn   bool   `json:"block_yn"`
		Content   string `json:"content"`
		CntReply  int64  `json:"cnt_reply"`
	} `json:"step1"`
	Step2 struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		BlockedYn bool   `json:"blocked_yn"`
		BlockYn   bool   `json:"block_yn"`
		Content   string `json:"content"`
		CntReply  int64  `json:"cnt_reply"`
	} `json:"step2"`
	Step3 struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		BlockedYn bool   `json:"blocked_yn"`
		BlockYn   bool   `json:"block_yn"`
		Content   string `json:"content"`
		CntReply  int64  `json:"cnt_reply"`
	} `json:"step3"`
	Step4 struct {
		SeqMember int64  `json:"seq_member"`
		NickName  string `json:"nick_name"`
		DeletedYn bool   `json:"deleted_yn"`
		BlockedYn bool   `json:"blocked_yn"`
		BlockYn   bool   `json:"block_yn"`
		Content   string `json:"content"`
		CntReply  int64  `json:"cnt_reply"`
	} `json:"step4"`
}
