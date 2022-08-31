package handlers

import (
	"ddaom/define"
	"ddaom/domain"
	"ddaom/memdb"
	"encoding/json"
	"strconv"
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
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	list, _ = memdb.Get("CACHES:MAIN:LIST_POPULAR_WRITER")

	listPopularWriter := []ListPopularWriter{}
	json.Unmarshal([]byte(list), &listPopularWriter)

	// 차단 작가 제외 로직
	if userToken != nil {
		var seqs []int64
		_list, err := memdb.Get("CACHES:USERS:BLOCK:" + strconv.FormatInt(userToken.SeqMember, 10))
		if err == nil {
			json.Unmarshal([]byte(_list), &seqs)
		}

		for i := 0; i < len(listPopularWriter); i++ {
			o := listPopularWriter[i]
			isExist := false
			for _, v := range seqs {
				if v == o.SeqMember {
					isExist = true
					break
				}
			}
			if !isExist {
				mainRes.ListPopularWriter = append(mainRes.ListPopularWriter, o)
			}
		}
	} else {
		if err == nil {
			json.Unmarshal([]byte(list), &mainRes.ListPopularWriter)
		}
	}

	mainRes.IsNewAlarm = true
	res.Data = mainRes

	return res
}

type MainRes struct {
	IsNewAlarm  bool       `json:"is_new_alarm"`
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
		SeqMember            int64  `json:"seq_member"`
		NickName             string `json:"nick_name"`
		ProfilePhoto         string `json:"profile_photo"`
		CntSubscribeBookmark int64  `json:"cnt_subscribe_bookmark"`
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
	SeqMember            int64  `json:"seq_member"`
	NickName             string `json:"nick_name"`
	ProfilePhoto         string `json:"profile_photo"`
	CntSubscribeBookmark int64  `json:"cnt_subscribe_bookmark"`
}

type ListPopularWriterLIke struct {
	SeqKeyword   int64  `json:"seq_keyword"`
	SeqMember    int64  `json:"seq_member"`
	NickName     string `json:"nick_name"`
	ProfilePhoto string `json:"profile_photo"`
	Cnt          int64  `json:"cnt"`
}
