package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

func MypageListLive(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	masterDB := db.List[define.DSN_MASTER]
	// myLogDb := GetMyLogDb(userToken.Allocated)

	// var list []int64
	// result := myLogDb.
	// 	Raw("SELECT seq_novel_finish FROM member_bookmarks mb WHERE seq_member = ?", userToken.SeqMember).
	// 	Scan(&list)
	// if result.Error != nil {
	// 	res.ResultCode = define.DB_ERROR_ORM
	// 	res.ErrorDesc = result.Error.Error()
	// 	return res
	// }

	var totalData int64
	var cnt int64
	result := masterDB.Model(schemas.NovelStep1{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	totalData += cnt
	result = masterDB.Model(schemas.NovelStep2{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	totalData += cnt
	result = masterDB.Model(schemas.NovelStep3{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	totalData += cnt
	result = masterDB.Model(schemas.NovelStep4{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	totalData += cnt

	novelMyListLiveRes := NovelMyListLiveRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	var query bytes.Buffer
	query.WriteString(`
	(SELECT
		seq_novel_step1,
		0 AS seq_novel_step2,
		0 AS seq_novel_step3,
		0 AS seq_novel_step4,
		seq_keyword,
		seq_genre,
		seq_image,
		seq_color,
		title,
		content,
		cnt_like,
		UNIX_TIMESTAMP(created_at) *1000 AS created_at,
		1 AS step
	FROM ddaom.novel_step1 ns WHERE seq_member = ?)
	UNION 
	(SELECT
		0 AS seq_novel_step1,
		ns2.seq_novel_step2,
		0 AS seq_novel_step3,
		0 AS seq_novel_step4,
		ns1.seq_keyword,
		ns1.seq_genre,
		ns1.seq_image,
		ns1.seq_color,
		ns1.title,
		ns2.content,
		ns2.cnt_like,
		UNIX_TIMESTAMP(ns2.created_at) *1000 AS created_at,
		2 AS step
	FROM ddaom.novel_step2 ns2 INNER JOIN ddaom.novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 
	WHERE ns2.seq_member = ?)
	UNION 
	(SELECT
		0 AS seq_novel_step1,
		0 AS seq_novel_step2,
		ns3.seq_novel_step3,
		0 AS seq_novel_step4,
		ns1.seq_keyword,
		ns1.seq_genre,
		ns1.seq_image,
		ns1.seq_color,
		ns1.title,
		ns3.content,
		ns3.cnt_like,
		UNIX_TIMESTAMP(ns3.created_at) *1000 AS created_at,
		3 AS step
	FROM ddaom.novel_step3 ns3 INNER JOIN ddaom.novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 
	WHERE ns3.seq_member = ?)
	UNION 
	(SELECT
		0 AS seq_novel_step1,
		0 AS seq_novel_step2,
		0 AS seq_novel_step3,
		ns4.seq_novel_step4,
		ns1.seq_keyword,
		ns1.seq_genre,
		ns1.seq_image,
		ns1.seq_color,
		ns1.title,
		ns4.content,
		ns4.cnt_like,
		UNIX_TIMESTAMP(ns4.created_at) *1000 AS created_at,
		4 AS step
	FROM ddaom.novel_step4 ns4 INNER JOIN ddaom.novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 
	WHERE ns4.seq_member = ?)
	`)
	query.WriteString(" ORDER BY created_at DESC")
	query.WriteString(" LIMIT ?, ?")
	result = masterDB.Raw(query.String(), userToken.SeqMember, userToken.SeqMember, userToken.SeqMember, userToken.SeqMember, limitStart, _sizePerPage).Find(&novelMyListLiveRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = novelMyListLiveRes

	return res
}

type NovelMyListLiveRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		Step          int8    `json:"step"`
		SeqNovelStep1 int64   `json:"seq_novel_step1"`
		SeqNovelStep2 int64   `json:"seq_novel_step2"`
		SeqNovelStep3 int64   `json:"seq_novel_step3"`
		SeqNovelStep4 int64   `json:"seq_novel_step4"`
		SeqKeyword    int64   `json:"seq_keyword"`
		SeqGenre      int64   `json:"seq_genre"`
		SeqImage      int64   `json:"seq_image"`
		SeqColor      int64   `json:"seq_color"`
		Title         string  `json:"title"`
		Content       string  `json:"content"`
		CntLike       int64   `json:"cnt_like"`
		CreatedAt     float64 `json:"created_at"`
	} `json:"list"`
}
