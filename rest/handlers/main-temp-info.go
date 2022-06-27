package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"strconv"
)

func MainTempInfo(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_step, _ := strconv.Atoi(req.Vars["step"])
	_seqNovel, _ := strconv.ParseInt(req.Vars["seq_novel"], 10, 64)

	sdb := db.List[define.Mconn.DsnSlave]

	mainTempInfoRes := MainTempInfoRes{}
	mainTempInfoRes.Step = int64(_step)
	qstep1 := `
			SELECT
				ns.seq_novel_step1,
				md.seq_member,
				md.nick_name,
				ns.seq_keyword,
				ns.seq_image,
				ns.seq_color,
				ns.seq_genre,
				UNIX_TIMESTAMP(ns.created_at) * 1000 AS created_at,
				ns.title,
				ns.content
			FROM novel_step1 ns INNER JOIN member_details md ON ns.seq_member = md.seq_member
			WHERE seq_novel_step1 = ?`
	qstep2 := `
			SELECT
				ns.seq_novel_step1,	
				ns.seq_novel_step2,
				md.seq_member,
				md.nick_name,
				UNIX_TIMESTAMP(ns.created_at) * 1000 AS created_at,
				ns.content
			FROM novel_step2 ns INNER JOIN member_details md ON ns.seq_member = md.seq_member
			WHERE seq_novel_step2 = ?`
	qstep3 := `
			SELECT
				ns.seq_novel_step1,	
				ns.seq_novel_step2,
				ns.seq_novel_step3,
				md.seq_member,
				md.nick_name,
				UNIX_TIMESTAMP(ns.created_at) * 1000 AS created_at,
				ns.content
			FROM novel_step3 ns INNER JOIN member_details md ON ns.seq_member = md.seq_member
			WHERE seq_novel_step3 = ?`
	qstep4 := `SELECT
				ns.seq_novel_step1,	
				ns.seq_novel_step2,
				ns.seq_novel_step3,
				ns.seq_novel_step4,
				md.seq_member,
				md.nick_name,
				UNIX_TIMESTAMP(ns.created_at) * 1000 AS created_at,
				ns.content
			FROM novel_step4 ns INNER JOIN member_details md ON ns.seq_member = md.seq_member
			WHERE seq_novel_step4 = ?`

	switch mainTempInfoRes.Step {
	case 1:
		result := sdb.Raw(qstep1, _seqNovel).Scan(&mainTempInfoRes.Step1)
		if corm(result, &res) {
			return res
		}
	case 2:
		result := sdb.Raw(qstep2, _seqNovel).Scan(&mainTempInfoRes.Step2)
		if corm(result, &res) {
			return res
		}
		result = sdb.Raw(qstep1, mainTempInfoRes.Step2.SeqNovelStep1).Scan(&mainTempInfoRes.Step1)
		if corm(result, &res) {
			return res
		}
		result = sdb.Model(schemas.NovelStep1{}).Select("deleted_yn").Where("seq_novel_step1 = ?", mainTempInfoRes.Step1.SeqNovelStep1).Scan(&mainTempInfoRes.Step2.IsParentDeleted)
		if corm(result, &res) {
			return res
		}

	case 3:
		result := sdb.Raw(qstep3, _seqNovel).Scan(&mainTempInfoRes.Step3)
		if corm(result, &res) {
			return res
		}
		result = sdb.Raw(qstep2, mainTempInfoRes.Step3.SeqNovelStep2).Scan(&mainTempInfoRes.Step2)
		if corm(result, &res) {
			return res
		}
		result = sdb.Raw(qstep1, mainTempInfoRes.Step2.SeqNovelStep1).Scan(&mainTempInfoRes.Step1)
		if corm(result, &res) {
			return res
		}
		result = sdb.Model(schemas.NovelStep2{}).Select("deleted_yn").Where("seq_novel_step2 = ?", mainTempInfoRes.Step2.SeqNovelStep2).Scan(&mainTempInfoRes.Step3.IsParentDeleted)
		if corm(result, &res) {
			return res
		}
	case 4:
		result := sdb.Raw(qstep4, _seqNovel).Scan(&mainTempInfoRes.Step4)
		if corm(result, &res) {
			return res
		}
		result = sdb.Raw(qstep3, mainTempInfoRes.Step4.SeqNovelStep3).Scan(&mainTempInfoRes.Step3)
		if corm(result, &res) {
			return res
		}
		result = sdb.Raw(qstep2, mainTempInfoRes.Step3.SeqNovelStep2).Scan(&mainTempInfoRes.Step2)
		if corm(result, &res) {
			return res
		}
		result = sdb.Raw(qstep1, mainTempInfoRes.Step2.SeqNovelStep1).Scan(&mainTempInfoRes.Step1)
		if corm(result, &res) {
			return res
		}
		result = sdb.Model(schemas.NovelStep3{}).Select("deleted_yn").Where("seq_novel_step3 = ?", mainTempInfoRes.Step3.SeqNovelStep3).Scan(&mainTempInfoRes.Step4.IsParentDeleted)
		if corm(result, &res) {
			return res
		}
	}

	// 상위글 삭제 여부 검수

	res.Data = mainTempInfoRes
	return res
}

type MainTempInfoRes struct {
	Step  int64 `json:"step"`
	Step1 struct {
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqMember     int64   `json:"seq_member"`
		NickName      string  `json:"nick_name"`
		SeqKeyword    int64   `json:"seq_keyword"`
		SeqImage      int64   `json:"seq_image"`
		SeqColor      int64   `json:"seq_color"`
		SeqGenre      int64   `json:"seq_genre"`
		CreatedAt     float64 `json:"created_at"`
		Title         string  `json:"title"`
		Content       string  `json:"content"`
	} `json:"step1"`
	Step2 struct {
		SeqNovelStep1   int64   `json:"seq_novel_step1"`
		SeqNovelStep2   int64   `json:"seq_novel_step2"`
		SeqMember       int64   `json:"seq_member"`
		NickName        string  `json:"nick_name"`
		CreatedAt       float64 `json:"created_at"`
		Content         string  `json:"content"`
		IsParentDeleted bool    `json:"is_parent_deleted"`
	} `json:"step2"`
	Step3 struct {
		SeqNovelStep1   int64   `json:"seq_novel_step1"`
		SeqNovelStep2   int64   `json:"seq_novel_step2"`
		SeqNovelStep3   int64   `json:"seq_novel_step3"`
		SeqMember       int64   `json:"seq_member"`
		NickName        string  `json:"nick_name"`
		CreatedAt       float64 `json:"created_at"`
		Content         string  `json:"content"`
		IsParentDeleted bool    `json:"is_parent_deleted"`
	} `json:"step3"`
	Step4 struct {
		SeqNovelStep1   int64   `json:"seq_novel_step1"`
		SeqNovelStep2   int64   `json:"seq_novel_step2"`
		SeqNovelStep3   int64   `json:"seq_novel_step3"`
		SeqNovelStep4   int64   `json:"seq_novel_step4"`
		SeqMember       int64   `json:"seq_member"`
		NickName        string  `json:"nick_name"`
		CreatedAt       float64 `json:"created_at"`
		Content         string  `json:"content"`
		IsParentDeleted bool    `json:"is_parent_deleted"`
	} `json:"step4"`
}
