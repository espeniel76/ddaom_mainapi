package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/memdb"
	"encoding/json"
	"strconv"
	"time"
)

func Assets(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	assets, err := memdb.Get("CACHES:ASSET:LIST")
	if err != nil {
		res.ResultCode = define.CACHE_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	data := AssetRes{}
	json.Unmarshal([]byte(assets), &data)

	var list []interface{}
	nowDate := time.Now().UnixMilli()
	keywordRes := KeywordRes{}
	list = append(list, "CACHES:ASSET:COUNT")
	for i := 0; i < len(data.ListKeyword); i++ {
		o := data.ListKeyword[i]

		// 실시간 카운트 위한 인덱스 추출
		list = append(list, o.SeqKeyword)

		// 실시간 주제어 정보
		if o.StartDate <= nowDate && o.EndDate > nowDate {
			keywordRes.List = append(keywordRes.List, struct {
				SeqKeyword int64  "json:\"seq_keyword\""
				Keyword    string "json:\"keyword\""
				StartDate  int64  "json:\"start_date\""
				EndDate    int64  "json:\"end_date\""
				CntTotal   int64  "json:\"cnt_total\""
			}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, StartDate: o.StartDate, EndDate: o.EndDate, CntTotal: o.CntTotal})
		}
	}

	// 실시간 카운트 정보
	data.ListKeyword = keywordRes.List
	result, err := memdb.ZMSCORE(list...)
	for i := 0; i < len(data.ListKeyword); i++ {
		num, _ := strconv.ParseInt(result[i], 10, 64)
		data.ListKeyword[i].CntTotal = num
	}

	res.Data = data

	return res
}

type AssetRes struct {
	ListKeyword []struct {
		SeqKeyword int64  `json:"seq_keyword"`
		Keyword    string `json:"keyword"`
		StartDate  int64  `json:"start_date"`
		EndDate    int64  `json:"end_date"`
		CntTotal   int64  `json:"cnt_total"`
	} `json:"list_keyword"`
	ListImage []struct {
		SeqImage int64  `json:"seq_image"`
		Image    string `json:"image"`
	} `json:"list_image"`
	ListColor []struct {
		SeqColor int64  `json:"seq_color"`
		Color    string `json:"color"`
	} `json:"list_color"`
	ListGenre []struct {
		SeqGenre int64  `json:"seq_genre"`
		Genre    string `json:"genre"`
	} `json:"list_genre"`
	ListSlang       []string `json:"list_slang"`
	ListCategoryFaq []struct {
		SeqCategoryFaq int64  `json:"seq_category_faq"`
		CategoryFaq    string `json:"category_faq"`
	} `json:"list_category_faq"`
	UrlPrivacyPolicy  string `json:"url_privacy_policy"`
	UrlTermsOfService string `json:"url_terms_of_service"`
}

func Keyword(req *domain.CommonRequest) domain.CommonResponse {
	var res = domain.CommonResponse{}
	keywords, err := memdb.Get("CACHES:ASSET:LIST:KEYWORD")
	if err != nil {
		res.ResultCode = define.CACHE_ERROR
		res.ErrorDesc = err.Error()
		return res
	}

	data := KeywordRes{}
	json.Unmarshal([]byte(keywords), &data.List)

	var list []interface{}
	nowDate := time.Now().UnixMilli()
	keywordRes := KeywordRes{}
	list = append(list, "CACHES:ASSET:COUNT")
	for i := 0; i < len(data.List); i++ {
		o := data.List[i]

		// 실시간 카운트 위한 인덱스 추출
		list = append(list, o.SeqKeyword)

		// 실시간 주제어 정보
		if o.StartDate <= nowDate && o.EndDate > nowDate {
			keywordRes.List = append(keywordRes.List, struct {
				SeqKeyword int64  "json:\"seq_keyword\""
				Keyword    string "json:\"keyword\""
				StartDate  int64  "json:\"start_date\""
				EndDate    int64  "json:\"end_date\""
				CntTotal   int64  "json:\"cnt_total\""
			}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, StartDate: o.StartDate, EndDate: o.EndDate, CntTotal: o.CntTotal})
		}
	}

	// 실시간 카운트 정보
	data.List = keywordRes.List
	result, err := memdb.ZMSCORE(list...)
	for i := 0; i < len(data.List); i++ {
		num, _ := strconv.ParseInt(result[i], 10, 64)
		data.List[i].CntTotal = num
	}
	res.Data = data

	return res
}

type KeywordDate struct {
	SeqKeyword    int64  `json:"seq_keyword"`
	Keyword       string `json:"keyword"`
	ViewStartDate string `json:"view_start_date"`
	ViewEndDate   string `json:"view_end_date"`
	StartDate     time.Time
	EndDate       time.Time
	CntTotal      int64 `json:"cnt_total"`
}

type KeywordRes struct {
	List []struct {
		SeqKeyword int64  `json:"seq_keyword"`
		Keyword    string `json:"keyword"`
		StartDate  int64  `json:"start_date"`
		EndDate    int64  `json:"end_date"`
		CntTotal   int64  `json:"cnt_total"`
	} `json:"list"`
}

// func Assets(req *domain.CommonRequest) domain.CommonResponse {

// 	var res = domain.CommonResponse{}

// 	today, _ := strconv.Atoi(tools.TodayFormattedDate())
// 	sdb := db.List[define.DSN_SLAVE]
// 	query := `
// 			SELECT
// 			seq_keyword,
// 			keyword,
// 			start_date AS view_start_date,
// 			end_date AS view_end_date,
// 			start_date,
// 			end_date,
// 			cnt_total
// 		FROM ddaom.keywords
// 		WHERE active_yn = true AND NOW() BETWEEN start_date AND end_date
// 		ORDER BY seq_keyword DESC
// 	`
// 	keywordDate := []KeywordDate{}
// 	result := sdb.Raw(query).Scan(&keywordDate)
// 	if corm(result, &res) {
// 		return res
// 	}
// 	assetRes := AssetRes{}
// 	for i := 0; i < len(keywordDate); i++ {
// 		o := keywordDate[i]
// 		isToday := false
// 		viewStartDate, _ := strconv.Atoi(o.ViewStartDate)
// 		viewEndDate, _ := strconv.Atoi(o.ViewEndDate)
// 		if viewStartDate <= today && viewEndDate >= today {
// 			isToday = true
// 		} else {
// 			isToday = false
// 		}
// 		assetRes.ListKeyword = append(assetRes.ListKeyword, struct {
// 			SeqKeyword int64  "json:\"seq_keyword\""
// 			Keyword    string "json:\"keyword\""
// 			IsToday    bool   "json:\"is_today\""
// 			StartDate  int64  "json:\"start_date\""
// 			EndDate    int64  "json:\"end_date\""
// 			CntTotal   int64  "json:\"cnt_total\""
// 		}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, IsToday: isToday, StartDate: o.StartDate.UnixMilli(), EndDate: o.EndDate.UnixMilli(), CntTotal: o.CntTotal})
// 	}

// 	result = sdb.Model(&schemas.Image{}).
// 		Where("active_yn = true").
// 		Select("seq_image, image").
// 		Find(&assetRes.ListImage)
// 	if corm(result, &res) {
// 		return res
// 	}
// 	result = sdb.Model(&schemas.Color{}).
// 		Where("active_yn = true").
// 		Select("seq_color, color").
// 		Find(&assetRes.ListColor)
// 	if corm(result, &res) {
// 		return res
// 	}
// 	result = sdb.Model(&schemas.Genre{}).
// 		Where("active_yn = true").
// 		Select("seq_genre, genre").
// 		Find(&assetRes.ListGenre)
// 	if corm(result, &res) {
// 		return res
// 	}
// 	result = sdb.Model(&schemas.Slang{}).
// 		Where("active_yn = true").
// 		Select("slang").
// 		Find(&assetRes.ListSlang)
// 	if corm(result, &res) {
// 		return res
// 	}
// 	result = sdb.Model(&schemas.CategoryFaq{}).
// 		Where("active_yn = true").
// 		Select("seq_category_faq, category_faq").
// 		Find(&assetRes.ListCategoryFaq)
// 	if corm(result, &res) {
// 		return res
// 	}

// 	res.Data = assetRes

// 	return res
// }

// func Keyword(req *domain.CommonRequest) domain.CommonResponse {

// 	var res = domain.CommonResponse{}

// 	// 오늘 날짜
// 	today, _ := strconv.Atoi(tools.TodayFormattedDate())
// 	fmt.Println(today)

// 	// 데이터 가져온다.
// 	sdb := db.List[define.DSN_SLAVE]
// 	query := `
// 		SELECT
// 			seq_keyword,
// 			keyword,
// 			start_date AS view_start_date,
// 			end_date AS view_end_date,
// 			start_date,
// 			end_date,
// 			cnt_total
// 		FROM ddaom.keywords
// 		WHERE active_yn = true AND NOW() BETWEEN start_date AND end_date
// 		ORDER BY seq_keyword DESC
// 	`
// 	keywordDate := []KeywordDate{}
// 	result := sdb.Raw(query).Scan(&keywordDate)
// 	if corm(result, &res) {
// 		return res
// 	}
// 	keywordRes := KeywordRes{}
// 	for i := 0; i < len(keywordDate); i++ {
// 		o := keywordDate[i]
// 		keywordRes.List = append(keywordRes.List, struct {
// 			SeqKeyword int64  "json:\"seq_keyword\""
// 			Keyword    string "json:\"keyword\""
// 			StartDate  int64  "json:\"start_date\""
// 			EndDate    int64  "json:\"end_date\""
// 			CntTotal   int64  "json:\"cnt_total\""
// 		}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, StartDate: o.StartDate.UnixMilli(), EndDate: o.EndDate.UnixMilli(), CntTotal: o.CntTotal})
// 	}

// 	res.Data = keywordRes

// 	return res
// }

func Skin(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	// 데이터 가져온다.
	sdb := db.List[define.DSN_SLAVE]
	skinRes := SkinRes{}
	result := sdb.Model(&schemas.Image{}).Where("active_yn = true").Find(&skinRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = skinRes

	return res
}

type SkinRes struct {
	List []struct {
		SeqImage int64  `json:"seq_image"`
		Image    string `json:"image"`
	} `json:"list"`
}

func Genre(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	// 데이터 가져온다.
	sdb := db.List[define.DSN_SLAVE]
	genreRes := GenreRes{}
	result := sdb.Model(&schemas.Genre{}).Where("active_yn = true").Find(&genreRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = genreRes

	return res
}

type GenreRes struct {
	List []struct {
		SeqGenre int64  `json:"seq_genre"`
		Genre    string `json:"genre"`
	} `json:"list"`
}
