package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"strconv"
)

func Main(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	slaveDb := db.List[define.DSN_SLAVE1]
	mainRes := MainRes{}

	// 연재중인 소설 (오늘 주제어 키워드)
	today := tools.TodayFormattedDate()
	query := `
	SELECT
		ns.seq_novel_step1,
		ns.seq_image,
		ns.seq_color,
		ns.title
	FROM novel_step1 ns
	INNER JOIN keywords k ON ns.seq_keyword = k.seq_keyword 
	INNER JOIN keyword_todays kt ON kt.seq_keyword = k.seq_keyword
	WHERE view_date = ? AND ns.active_yn = true
	ORDER BY ns.created_at DESC
	LIMIT 10
	`
	result := slaveDb.Raw(query, today).Scan(&mainRes.ListLive)
	if corm(result, &res) {
		return res
	}

	// 인기작
	query = `
	SELECT
		nf.seq_novel_finish,
		ns.seq_image,
		ns.seq_color,
		ns.title
	FROM novel_finishes nf
	INNER JOIN novel_step1 ns ON nf.seq_novel_step1 = ns.seq_novel_step1
	WHERE nf.active_yn
	ORDER BY nf.cnt_like DESC
	LIMIT 10
	`
	result = slaveDb.Raw(query).Scan(&mainRes.ListPopular)
	if corm(result, &res) {
		return res
	}

	// 완결작
	query = `
	SELECT
		nf.seq_novel_finish,
		ns.seq_image,
		ns.seq_color,
		ns.title
	FROM novel_finishes nf
	INNER JOIN novel_step1 ns ON nf.seq_novel_step1 = ns.seq_novel_step1
	WHERE nf.active_yn = true
	ORDER BY nf.created_at DESC
	LIMIT 10
	`
	result = slaveDb.Raw(query).Scan(&mainRes.ListFinish)
	if corm(result, &res) {
		return res
	}

	// 인기작가
	query = "SELECT seq_member, nick_name, profile_photo FROM member_details ORDER BY cnt_like DESC LIMIT 10"
	result = slaveDb.Raw(query).Scan(&mainRes.ListPopularWriter)
	if corm(result, &res) {
		return res
	}

	res.Data = mainRes

	return res
}

type MainRes struct {
	ListLive    []ListLive `json:"list_live"`
	ListPopular []struct {
		SeqNovelFinish int64  `json:"seq_novel_finish"`
		SeqImage       int64  `json:"seq_image"`
		SeqColor       int64  `json:"seq_color"`
		Title          string `json:"title"`
	} `json:"list_popular"`
	ListFinish []struct {
		SeqNovelFinish int64  `json:"seq_novel_finish"`
		SeqImage       int64  `json:"seq_image"`
		SeqColor       int64  `json:"seq_color"`
		Title          string `json:"title"`
	} `json:"list_finish"`
	ListPopularWriter []struct {
		SeqMember    int64  `json:"seq_member"`
		NickName     string `json:"nick_name"`
		ProfilePhoto string `json:"profile_photo"`
	} `json:"list_popular_writer"`
}

func MainKeyword(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	_seqKeyword, _ := strconv.Atoi(req.Vars["seq_keyword"])

	slaveDb := db.List[define.DSN_SLAVE1]
	mainRes := []ListLive{}

	query := `
	SELECT
		seq_novel_step1,
		seq_image,
		seq_color,
		title
	FROM novel_step1
	WHERE active_yn = true AND seq_keyword = ?
	ORDER BY created_at DESC
	LIMIT 10
	`
	result := slaveDb.Raw(query, _seqKeyword).Scan(&mainRes)
	if corm(result, &res) {
		return res
	}
	res.Data = mainRes

	return res
}

type ListLive struct {
	SeqNovelStep1 int64  `json:"seq_novel_step1"`
	SeqImage      int64  `json:"seq_image"`
	SeqColor      int64  `json:"seq_color"`
	Title         string `json:"title"`
}
