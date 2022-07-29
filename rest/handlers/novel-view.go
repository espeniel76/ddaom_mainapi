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
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	_seqNovelStep1, _ := req.Vars["seq_novel_step1"]
	var cntTotal int64

	sdb := db.List[define.Mconn.DsnSlave]
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
		ns.seq_color,
		m.deleted_yn
	FROM novel_step1 ns
	INNER JOIN member_details md ON ns.seq_member = md.seq_member
	INNER JOIN members m ON md.seq_member = m.seq_member
	WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false`
	step1Res := Step1Res{}
	result := sdb.Raw(query, _seqNovelStep1).Scan(&step1Res)
	if corm(result, &res) {
		return res
	}

	// step1 글이 있어야 step2 가 있다
	data := make(map[string]interface{})
	if step1Res.SeqNovelStep1 > 0 {

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
		step1["my_like"] = getMyLike(userToken, 1, step1Res.SeqNovelStep1)
		step1["deleted_yn"] = bool(step1Res.DeletedYn)
		if userToken != nil {
			bm := getBlockMember(userToken.Allocated, userToken.SeqMember, step1Res.SeqMember)
			step1["block_yn"] = bm.BlockYn
		} else {
			step1["block_yn"] = false
		}
		data["step1"] = step1

		query = `
		SELECT
			ns.seq_novel_step2,
			ns.seq_member,
			md.nick_name,
			ns.created_at,
			ns.content,
			ns.cnt_like,
			m.deleted_yn
		FROM novel_step2 ns
		INNER JOIN member_details md ON ns.seq_member = md.seq_member
		INNER JOIN members m ON md.seq_member = m.seq_member
		WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
		ORDER BY ns.updated_at DESC`
		step2Res := Step2Res{}
		result = sdb.Raw(query, _seqNovelStep1).Scan(&step2Res)
		if corm(result, &res) {
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
			step["my_like"] = getMyLike(userToken, 2, step2Res.SeqNovelStep2)
			step["deleted_yn"] = bool(step2Res.DeletedYn)
			result = sdb.Model(schemas.NovelStep2{}).Where("seq_novel_step1 = ? AND active_yn = true AND temp_yn = false AND deleted_yn = false", _seqNovelStep1).Count(&cntTotal)
			if corm(result, &res) {
				return res
			}
			step["cnt_total"] = cntTotal
			if userToken != nil {
				bm := getBlockMember(userToken.Allocated, userToken.SeqMember, step2Res.SeqMember)
				step["block_yn"] = bm.BlockYn
			} else {
				step["block_yn"] = false
			}
			data["step2"] = step

			query = `
			SELECT
				ns.seq_novel_step3,
				ns.seq_member,
				md.nick_name,
				ns.created_at,
				ns.content,
				ns.cnt_like,
				m.deleted_yn
			FROM novel_step3 ns
			INNER JOIN member_details md ON ns.seq_member = md.seq_member
			INNER JOIN members m ON md.seq_member = m.seq_member
			WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
			ORDER BY ns.updated_at DESC`
			step3Res := Step3Res{}
			result = sdb.Raw(query, _seqNovelStep1).Scan(&step3Res)
			if corm(result, &res) {
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
				step["my_like"] = getMyLike(userToken, 3, step3Res.SeqNovelStep3)
				step["deleted_yn"] = bool(step3Res.DeletedYn)
				result = sdb.Model(schemas.NovelStep3{}).Where("seq_novel_step1 = ? AND active_yn = true AND temp_yn = false AND deleted_yn = false", _seqNovelStep1).Count(&cntTotal)
				if corm(result, &res) {
					return res
				}
				step["cnt_total"] = cntTotal
				if userToken != nil {
					bm := getBlockMember(userToken.Allocated, userToken.SeqMember, step3Res.SeqMember)
					step["block_yn"] = bm.BlockYn
				} else {
					step["block_yn"] = false
				}
				data["step3"] = step

				query = `
				SELECT
					ns.seq_novel_step4,
					ns.seq_member,
					md.nick_name,
					ns.created_at,
					ns.content,
					ns.cnt_like,
					m.deleted_yn
				FROM novel_step4 ns
				INNER JOIN member_details md ON ns.seq_member = md.seq_member
				INNER JOIN members m ON md.seq_member = m.seq_member
				WHERE ns.active_yn = true AND ns.seq_novel_step1 = ? AND ns.temp_yn = false AND ns.deleted_yn = false
				ORDER BY ns.updated_at DESC`
				step4Res := Step4Res{}
				result = sdb.Raw(query, _seqNovelStep1).Scan(&step4Res)
				if corm(result, &res) {
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
					step["my_like"] = getMyLike(userToken, 4, step4Res.SeqNovelStep4)
					step["deleted_yn"] = bool(step4Res.DeletedYn)
					result = sdb.Model(schemas.NovelStep4{}).Where("seq_novel_step1 = ? AND active_yn = true AND temp_yn = false AND deleted_yn = false", _seqNovelStep1).Count(&cntTotal)
					if corm(result, &res) {
						return res
					}
					step["cnt_total"] = cntTotal
					if userToken != nil {
						bm := getBlockMember(userToken.Allocated, userToken.SeqMember, step4Res.SeqMember)
						step["block_yn"] = bm.BlockYn
					} else {
						step["block_yn"] = false
					}
					data["step4"] = step
				} else {
					data["step4"] = nil
				}
			} else {
				data["step3"] = nil
				data["step4"] = nil
			}
		} else {
			data["step2"] = nil
			data["step3"] = nil
			data["step4"] = nil
		}

	} else {
		data["step1"] = nil
		data["step2"] = nil
		data["step3"] = nil
		data["step4"] = nil
	}

	// fmt.Println(data)

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
	DeletedYn     bool      `json:"deleted_yn"`
	BlockYn       bool      `json:"block_yn"`
}
type Step2Res struct {
	SeqNovelStep2 int64     `json:"seq_novel_step2"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
	DeletedYn     bool      `json:"deleted_yn"`
	BlockYn       bool      `json:"block_yn"`
}
type Step3Res struct {
	SeqNovelStep3 int64     `json:"seq_novel_step3"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
	DeletedYn     bool      `json:"deleted_yn"`
	BlockYn       bool      `json:"block_yn"`
}
type Step4Res struct {
	SeqNovelStep4 int64     `json:"seq_novel_step4"`
	SeqMember     int64     `json:"seq_member"`
	NickName      string    `json:"nick_name"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"content"`
	CntLike       int64     `json:"cnt_like"`
	DeletedYn     bool      `json:"deleted_yn"`
	BlockYn       bool      `json:"block_yn"`
}
