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

	// 연재중 좋아요 수
	// listLike, _ := memdb.Get("CACHES:MAIN:LIST_POPULAR_WRITER_LIKE")

	listPopularWriter := []ListPopularWriter{}
	json.Unmarshal([]byte(list), &listPopularWriter)
	// listPopularWriterLike := []ListPopularWriterLIke{}
	// json.Unmarshal([]byte(listLike), &listPopularWriterLike)

	// 연재중 키워드

	// 구독/북마크 + 연재중 좋아요 데이터 합침
	// for j := 0; j < len(listPopularWriterLike); j++ {

	// 	// 연재중 여부 확인

	// 	isExist := false
	// 	for i := 0; i < len(listPopularWriter); i++ {
	// 		if listPopularWriter[i].SeqMember == listPopularWriterLike[j].SeqMember {
	// 			isExist = true
	// 			listPopularWriter[i].CntSubscribeBookmark += listPopularWriterLike[j].Cnt
	// 			break
	// 		}
	// 	}
	// 	if !isExist {
	// 		listPopularWriter = append(listPopularWriter, ListPopularWriter{
	// 			SeqMember:            listPopularWriterLike[j].SeqMember,
	// 			NickName:             listPopularWriterLike[j].NickName,
	// 			ProfilePhoto:         listPopularWriterLike[j].ProfilePhoto,
	// 			CntSubscribeBookmark: listPopularWriterLike[j].Cnt,
	// 		})
	// 	}
	// }

	// // 데이터 정렬 (구독+북마크+연재중 좋아요)
	// sort.Slice(listPopularWriter, func(i, j int) bool {
	// 	return listPopularWriter[i].CntSubscribeBookmark > listPopularWriter[j].CntSubscribeBookmark
	// })

	// 차단 작가 제외 로직
	if userToken != nil {
		var seqs []int64
		_list, err := memdb.Get("CACHES:USERS:BLOCK:" + strconv.FormatInt(userToken.SeqMember, 10))
		if err == nil {
			json.Unmarshal([]byte(_list), &seqs)
			// fmt.Println(seqs)
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
