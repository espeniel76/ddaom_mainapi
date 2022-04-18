package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"fmt"
)

func MainTempInfo(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	fmt.Println(userToken)

	sdb := db.List[define.DSN_SLAVE]

	// 1. 스텝별로 데이터 존재 체크 후 최상의 존재 가져옴
	query := `
	SELECT step, seq_novel_step1, seq_novel_step2, seq_novel_step3, seq_novel_step4, created_at
	FROM
	(
		(SELECT 1 AS step, seq_novel_step1, 0 AS seq_novel_step2, 0 AS seq_novel_step3, 0 AS seq_novel_step4, created_at
		FROM novel_step1 ns WHERE seq_member = ? AND temp_yn = true ORDER BY created_at DESC LIMIT 1)
		UNION ALL
		(SELECT 2 AS step, 0 AS seq_novel_step1, seq_novel_step2, 0 AS seq_novel_step3, 0 AS seq_novel_step4, created_at
		FROM novel_step2 ns WHERE seq_member = ? AND temp_yn = true ORDER BY created_at DESC LIMIT 1)
		UNION ALL
		(SELECT 3 AS step, 0 AS seq_novel_step1, 0 AS seq_novel_step2, seq_novel_step3, 0 AS seq_novel_step4, created_at
		FROM novel_step3 ns WHERE seq_member = ? AND temp_yn = true ORDER BY created_at DESC LIMIT 1)
		UNION ALL
		(SELECT 4 AS step, 0 AS seq_novel_step1, 0 AS seq_novel_step2, 0 AS seq_novel_step3, seq_novel_step4, created_at
		FROM novel_step4 ns WHERE seq_member = ? AND temp_yn = true ORDER BY created_at DESC LIMIT 1)
	) AS s
	ORDER BY s.created_at DESC
	LIMIT 1
	`
	seq := userToken.SeqMember
	tmp := make(map[string]interface{})
	result := sdb.Raw(query, seq, seq, seq, seq).Scan(&tmp)
	if corm(result, &res) {
		return res
	}
	mainTempInfoRes := MainTempInfoRes{}
	if tmp["step"] == nil {
		mainTempInfoRes.Step = 0
		res.Data = mainTempInfoRes
		return res
	} else {
		mainTempInfoRes.Step, _ = tmp["step"].(int64)
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
			result = sdb.Raw(qstep1, tmp["seq_novel_step1"]).Scan(&mainTempInfoRes.Step1)
			if corm(result, &res) {
				return res
			}
		case 2:
			result = sdb.Raw(qstep2, tmp["seq_novel_step2"]).Scan(&mainTempInfoRes.Step2)
			if corm(result, &res) {
				return res
			}
			result = sdb.Raw(qstep1, mainTempInfoRes.Step2.SeqNovelStep1).Scan(&mainTempInfoRes.Step1)
			if corm(result, &res) {
				return res
			}
		case 3:
			result = sdb.Raw(qstep3, tmp["seq_novel_step3"]).Scan(&mainTempInfoRes.Step3)
			if corm(result, &res) {
				return res
			}
			result = sdb.Raw(qstep2, mainTempInfoRes.Step3.SeqNovelStep2).Scan(&mainTempInfoRes.Step2)
			if corm(result, &res) {
				return res
			}
		case 4:
			result = sdb.Raw(qstep4, tmp["seq_novel_step4"]).Scan(&mainTempInfoRes.Step4)
			if corm(result, &res) {
				return res
			}
			result = sdb.Raw(qstep3, mainTempInfoRes.Step4.SeqNovelStep3).Scan(&mainTempInfoRes.Step3)
			if corm(result, &res) {
				return res
			}
		}
	}

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
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqNovelStep2 int64   `json:"seq_novel_step2"`
		SeqMember     int64   `json:"seq_member"`
		NickName      string  `json:"nick_name"`
		CreatedAt     float64 `json:"created_at"`
		Content       string  `json:"content"`
	} `json:"step2"`
	Step3 struct {
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqNovelStep2 int64   `json:"seq_novel_step2"`
		SeqNovelStep3 int64   `json:"seq_novel_step3"`
		SeqMember     int64   `json:"seq_member"`
		NickName      string  `json:"nick_name"`
		CreatedAt     float64 `json:"created_at"`
		Content       string  `json:"content"`
	} `json:"step3"`
	Step4 struct {
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqNovelStep2 int64   `json:"seq_novel_step2"`
		SeqNovelStep3 int64   `json:"seq_novel_step3"`
		SeqNovelStep4 int64   `json:"seq_novel_step4"`
		SeqMember     int64   `json:"seq_member"`
		NickName      string  `json:"nick_name"`
		CreatedAt     float64 `json:"created_at"`
		Content       string  `json:"content"`
	} `json:"step4"`
}
