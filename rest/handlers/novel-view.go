package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"time"
)

func NovelView(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	// userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	// if err != nil {
	// 	res.ResultCode = define.INVALID_TOKEN
	// 	res.ErrorDesc = err.Error()
	// 	return res
	// }
	_seqNovelStep1, _ := req.Vars["seq_novel_step1"]
	// fmt.Println(_seqNovelStep1, userToken)
	var cntTotal int64

	masterDB := db.List[define.DSN_MASTER]
	// 1단계 소설 가져오기
	query := `
	SELECT
		ns.seq_novel_step1,
		ns.title,
		ns.created_at,
		ns.seq_genre,
		ns.seq_keyword,
		ns.seq_member,
		md.nick_name,
		ns.content,
		ns.cnt_like,
		ns.seq_image,
		ns.seq_color
	FROM novel_step1 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false`
	step1Res := Step1Res{}
	result := masterDB.Raw(query, _seqNovelStep1).Scan(&step1Res)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	data := make(map[string]interface{})
	data["title"] = step1Res.Title
	data["created_at"] = step1Res.CreatedAt.UnixMilli()
	data["seq_genre"] = step1Res.SeqGenre
	data["seq_keyword"] = step1Res.SeqKeyword
	data["seq_image"] = step1Res.SeqImage
	data["seq_color"] = step1Res.SeqColor

	step1 := make(map[string]interface{})
	step1["seq_novel_step1"] = step1Res.SeqNovelStep1
	step1["seq_member"] = step1Res.SeqMember
	step1["nick_name"] = step1Res.NickName
	step1["content"] = step1Res.Content
	step1["cnt_like"] = step1Res.CntLike
	step1["my_like"] = true
	data["step1"] = step1

	query = `
	SELECT
		ns.seq_novel_step2,
		ns.seq_member,
		md.nick_name,
		ns.created_at,
		ns.content,
		ns.cnt_like
	FROM novel_step2 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false
	ORDER BY ns.seq_novel_step2 DESC`
	step2Res := Step2Res{}
	result = masterDB.Raw(query, _seqNovelStep1).Scan(&step2Res)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if step2Res.SeqNovelStep2 > 0 {
		step := make(map[string]interface{})
		step["seq_novel_step2"] = step2Res.SeqNovelStep2
		step["seq_member"] = step2Res.SeqMember
		step["nick_name"] = step2Res.NickName
		step["created_at"] = step2Res.CreatedAt.UnixMilli()
		step["content"] = step2Res.Content
		step["cnt_like"] = step2Res.CntLike
		step["my_like"] = true
		result = masterDB.Model(schemas.NovelStep2{}).Where("seq_novel_step1 = ?", _seqNovelStep1).Count(&cntTotal)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		step["cnt_total"] = cntTotal
		data["step2"] = step
	} else {
		data["step2"] = nil
	}

	query = `
	SELECT
		ns.seq_novel_step3,
		ns.seq_member,
		md.nick_name,
		ns.created_at,
		ns.content,
		ns.cnt_like
	FROM novel_step3 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false
	ORDER BY ns.seq_novel_step3 DESC`
	step3Res := Step3Res{}
	result = masterDB.Raw(query, _seqNovelStep1).Scan(&step3Res)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if step3Res.SeqNovelStep3 > 0 {
		step := make(map[string]interface{})
		step["seq_novel_step3"] = step3Res.SeqNovelStep3
		step["seq_member"] = step3Res.SeqMember
		step["nick_name"] = step3Res.NickName
		step["created_at"] = step3Res.CreatedAt.UnixMilli()
		step["content"] = step3Res.Content
		step["cnt_like"] = step3Res.CntLike
		step["my_like"] = true
		result = masterDB.Model(schemas.NovelStep3{}).Where("seq_novel_step1 = ?", _seqNovelStep1).Count(&cntTotal)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		step["cnt_total"] = cntTotal
		data["step3"] = step
	} else {
		data["step3"] = nil
	}

	query = `
	SELECT
		ns.seq_novel_step4,
		ns.seq_member,
		md.nick_name,
		ns.created_at,
		ns.content,
		ns.cnt_like
	FROM novel_step4 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false
	ORDER BY ns.seq_novel_step4 DESC`
	step4Res := Step4Res{}
	result = masterDB.Raw(query, _seqNovelStep1).Scan(&step4Res)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	if step4Res.SeqNovelStep4 > 0 {
		step := make(map[string]interface{})
		step["seq_novel_step4"] = step4Res.SeqNovelStep4
		step["seq_member"] = step4Res.SeqMember
		step["nick_name"] = step4Res.NickName
		step["created_at"] = step4Res.CreatedAt.UnixMilli()
		step["content"] = step4Res.Content
		step["cnt_like"] = step4Res.CntLike
		step["my_like"] = true
		result = masterDB.Model(schemas.NovelStep4{}).Where("seq_novel_step1 = ?", _seqNovelStep1).Count(&cntTotal)
		if result.Error != nil {
			res.ResultCode = define.DB_ERROR_ORM
			res.ErrorDesc = result.Error.Error()
			return res
		}
		step["cnt_total"] = cntTotal
		data["step4"] = step
	} else {
		data["step4"] = nil
	}

	res.Data = data

	return res
}

type Step1Res struct {
	SeqNovelStep1 int64     `json:"seq_novel_step1"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"created_at"`
	SeqGenre      int64     `json:"seq_genre"`
	SeqKeyword    int64     `json:"seq_keyword"`
	SeqImage      int64     `json:"seq_image"`
	SeqColor      int64     `json:"seq_color"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
}
type Step2Res struct {
	SeqNovelStep2 int64     `json:"seq_novel_step2"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
}
type Step3Res struct {
	SeqNovelStep3 int64     `json:"seq_novel_step3"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
}
type Step4Res struct {
	SeqNovelStep4 int64     `json:"seq_novel_step4"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
}