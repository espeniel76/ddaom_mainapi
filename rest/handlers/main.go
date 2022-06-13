package handlers

import (
	"ddaom/domain"
	"ddaom/memdb"
	"encoding/json"
)

func Main(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	// _seqKeyword, _ := strconv.Atoi(req.Vars["seq_keyword"])
	// sdb := db.List[define.DSN_SLAVE]
	// 고도화 DB -> Redis
	mainRes := MainRes{}

	// 연재중인 소설 (오늘 주제어 키워드)
	list, _ := memdb.Get("CACHES:MAIN:LIST_LIVE:" + req.Vars["seq_keyword"])
	json.Unmarshal([]byte(list), &mainRes.ListLive)

	// 인기작
	list, _ = memdb.Get("CACHES:MAIN:LIST_POPULAR")
	json.Unmarshal([]byte(list), &mainRes.ListPopular)

	// 완결작
	list, _ = memdb.Get("CACHES:MAIN:LIST_FINISH")
	json.Unmarshal([]byte(list), &mainRes.ListFinish)

	// 인기작가
	list, _ = memdb.Get("CACHES:MAIN:LIST_POPULAR_WRITER")
	json.Unmarshal([]byte(list), &mainRes.ListPopularWriter)

	res.Data = mainRes

	return res
}

// func Main(req *domain.CommonRequest) domain.CommonResponse {

// 	var res = domain.CommonResponse{}
// 	_seqKeyword, _ := strconv.Atoi(req.Vars["seq_keyword"])

// 	sdb := db.List[define.DSN_SLAVE]
// 	mainRes := MainRes{}

// 	// 연재중인 소설 (오늘 주제어 키워드)
// 	// today := tools.TodayFormattedDate()
// 	query := `
// 	SELECT
// 		seq_novel_step1,
// 		seq_image,
// 		seq_color,
// 		title
// 	FROM novel_step1
// 	WHERE seq_keyword = ? AND active_yn = true AND temp_yn = false AND deleted_yn = false
// 	ORDER BY created_at DESC
// 	LIMIT 10
// 	`
// 	result := sdb.Raw(query, _seqKeyword).Scan(&mainRes.ListLive)
// 	if corm(result, &res) {
// 		return res
// 	}

// 	// 인기작
// 	query = `
// 	SELECT
// 		seq_novel_finish,
// 		seq_image,
// 		seq_color,
// 		title,
// 		cnt_like + cnt_view AS cnt_sum
// 	FROM novel_finishes
// 	WHERE active_yn = true
// 	ORDER BY cnt_sum DESC
// 	LIMIT 10
// 	`
// 	result = sdb.Raw(query).Scan(&mainRes.ListPopular)
// 	if corm(result, &res) {
// 		return res
// 	}

// 	// 완결작
// 	query = `
// 	SELECT
// 		seq_novel_finish,
// 		seq_image,
// 		seq_color,
// 		title
// 	FROM novel_finishes
// 	WHERE active_yn = true
// 	ORDER BY created_at DESC
// 	LIMIT 10
// 	`
// 	result = sdb.Raw(query).Scan(&mainRes.ListFinish)
// 	if corm(result, &res) {
// 		return res
// 	}

// 	// 인기작가
// 	query = "SELECT seq_member, nick_name, profile_photo FROM member_details ORDER BY cnt_like DESC LIMIT 10"
// 	result = sdb.Raw(query).Scan(&mainRes.ListPopularWriter)
// 	if corm(result, &res) {
// 		return res
// 	}

// 	res.Data = mainRes

// 	return res
// }

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
	// _seqKeyword, _ := strconv.Atoi(req.Vars["seq_keyword"])

	// sdb := db.List[define.DSN_SLAVE]
	mainRes := ListLiveRes{}

	// 연재중인 소설 (오늘 주제어 키워드)
	list, _ := memdb.Get("CACHES:MAIN:LIST_LIVE:" + req.Vars["seq_keyword"])
	json.Unmarshal([]byte(list), &mainRes.ListLive)

	// query := `
	// SELECT
	// 	seq_novel_step1,
	// 	seq_image,
	// 	seq_color,
	// 	title
	// FROM novel_step1
	// WHERE active_yn = true AND seq_keyword = ? AND temp_yn = false AND deleted_yn = false
	// ORDER BY created_at DESC
	// LIMIT 10
	// `
	// result := sdb.Raw(query, _seqKeyword).Scan(&mainRes.ListLive)
	// if corm(result, &res) {
	// 	return res
	// }
	res.Data = mainRes

	return res
}

type ListLiveRes struct {
	ListLive []ListLive `json:"list_live"`
}

type ListLive struct {
	SeqNovelStep1 int64  `json:"seq_novel_step1"`
	SeqImage      int64  `json:"seq_image"`
	SeqColor      int64  `json:"seq_color"`
	Title         string `json:"title"`
}
type ListPopular struct {
	SeqNovelFinish int64  `json:"seq_novel_finish"`
	SeqImage       int64  `json:"seq_image"`
	SeqColor       int64  `json:"seq_color"`
	Title          string `json:"title"`
}

type ListPopularWriter struct {
	SeqMember    int64  `json:"seq_member"`
	NickName     string `json:"nick_name"`
	ProfilePhoto string `json:"profile_photo"`
}
