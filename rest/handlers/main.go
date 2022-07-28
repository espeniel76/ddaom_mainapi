package handlers

import (
	"ddaom/domain"
	"ddaom/memdb"
	"encoding/json"
)

func Main(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	mainRes := MainRes{}

	// 연재중인 소설 (오늘 주제어 키워드)
	list, err := memdb.Get("CACHES:MAIN:LIST_LIVE:" + req.Vars["seq_keyword"])
	if err == nil {
		json.Unmarshal([]byte(list), &mainRes.ListLive)
	}

	// 인기작
	list, err = memdb.Get("CACHES:MAIN:LIST_POPULAR")
	if err == nil {
		json.Unmarshal([]byte(list), &mainRes.ListPopular)
	}

	// 완결작
	list, err = memdb.Get("CACHES:MAIN:LIST_FINISH")
	if err == nil {
		json.Unmarshal([]byte(list), &mainRes.ListFinish)
	}

	// 인기작가
	list, err = memdb.Get("CACHES:MAIN:LIST_POPULAR_WRITER")
	if err == nil {
		json.Unmarshal([]byte(list), &mainRes.ListPopularWriter)
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
	mainRes := ListLiveRes{}

	list, _ := memdb.Get("CACHES:MAIN:LIST_LIVE:" + req.Vars["seq_keyword"])
	json.Unmarshal([]byte(list), &mainRes.ListLive)
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
