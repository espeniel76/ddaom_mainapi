package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
	"strconv"
	"time"
)

func Assets(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	today, _ := strconv.Atoi(tools.TodayFormattedDate())
	sdb := db.List[define.DSN_SLAVE]
	query := `
			SELECT
			k.seq_keyword,
			k.keyword,
			kt.view_start_date,
			kt.view_end_date,
			k.start_date,
			k.end_date,
			k.cnt_total
		FROM ddaom.keywords AS k
		INNER JOIN ddaom.keyword_todays AS kt ON k.seq_keyword = kt.seq_keyword
		WHERE k.active_yn = true AND NOW() BETWEEN k.start_date AND k.end_date
		ORDER BY kt.view_start_date ASC
	`
	keywordDate := []KeywordDate{}
	result := sdb.Raw(query).Scan(&keywordDate)
	if corm(result, &res) {
		return res
	}
	assetRes := AssetRes{}
	for i := 0; i < len(keywordDate); i++ {
		o := keywordDate[i]
		isToday := false
		viewStartDate, _ := strconv.Atoi(o.ViewStartDate)
		viewEndDate, _ := strconv.Atoi(o.ViewEndDate)
		if viewStartDate <= today && viewEndDate >= today {
			isToday = true
		} else {
			isToday = false
		}
		assetRes.ListKeyword = append(assetRes.ListKeyword, struct {
			SeqKeyword int64  "json:\"seq_keyword\""
			Keyword    string "json:\"keyword\""
			IsToday    bool   "json:\"is_today\""
			StartDate  int64  "json:\"start_date\""
			EndDate    int64  "json:\"end_date\""
			CntTotal   int64  "json:\"cnt_total\""
		}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, IsToday: isToday, StartDate: o.StartDate.UnixMilli(), EndDate: o.EndDate.UnixMilli(), CntTotal: o.CntTotal})
	}

	result = sdb.Model(&schemas.Image{}).
		Where("active_yn = true").
		Select("seq_image, image").
		Find(&assetRes.ListImage)
	if corm(result, &res) {
		return res
	}
	result = sdb.Model(&schemas.Color{}).
		Where("active_yn = true").
		Select("seq_color, color").
		Find(&assetRes.ListColor)
	if corm(result, &res) {
		return res
	}
	result = sdb.Model(&schemas.Genre{}).
		Where("active_yn = true").
		Select("seq_genre, genre").
		Find(&assetRes.ListGenre)
	if corm(result, &res) {
		return res
	}
	result = sdb.Model(&schemas.Slang{}).
		Where("active_yn = true").
		Select("slang").
		Find(&assetRes.ListSlang)
	if corm(result, &res) {
		return res
	}
	result = sdb.Model(&schemas.CategoryFaq{}).
		Where("active_yn = true").
		Select("seq_category_faq, category_faq").
		Find(&assetRes.ListCategoryFaq)
	if corm(result, &res) {
		return res
	}

	res.Data = assetRes

	return res
}

type AssetRes struct {
	ListKeyword []struct {
		SeqKeyword int64  `json:"seq_keyword"`
		Keyword    string `json:"keyword"`
		IsToday    bool   `json:"is_today"`
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
	UrlPrivicyPolicy  string `json:"url_privicy_policy"`
	UrlTermsOfService string `json:"url_terms_of_service"`
}

func Keyword(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	// 오늘 날짜
	today, _ := strconv.Atoi(tools.TodayFormattedDate())
	fmt.Println(today)

	// 데이터 가져온다.
	sdb := db.List[define.DSN_SLAVE]
	query := `
			SELECT
			k.seq_keyword,
			k.keyword,
			kt.view_start_date,
			kt.view_end_date,
			k.start_date,
			k.end_date,
			k.cnt_total
		FROM ddaom.keywords AS k
		INNER JOIN ddaom.keyword_todays AS kt ON k.seq_keyword = kt.seq_keyword
		WHERE k.active_yn = true AND NOW() BETWEEN k.start_date AND k.end_date
		ORDER BY kt.view_start_date ASC
	`
	keywordDate := []KeywordDate{}
	result := sdb.Raw(query).Scan(&keywordDate)
	if corm(result, &res) {
		return res
	}
	keywordRes := KeywordRes{}
	for i := 0; i < len(keywordDate); i++ {
		o := keywordDate[i]
		isToday := false
		viewStartDate, _ := strconv.Atoi(o.ViewStartDate)
		viewEndDate, _ := strconv.Atoi(o.ViewEndDate)
		if viewStartDate <= today && viewEndDate >= today {
			isToday = true
		} else {
			isToday = false
		}
		keywordRes.List = append(keywordRes.List, struct {
			SeqKeyword int64  "json:\"seq_keyword\""
			Keyword    string "json:\"keyword\""
			IsToday    bool   "json:\"is_today\""
			StartDate  int64  "json:\"start_date\""
			EndDate    int64  "json:\"end_date\""
			CntTotal   int64  "json:\"cnt_total\""
		}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, IsToday: isToday, StartDate: o.StartDate.UnixMilli(), EndDate: o.EndDate.UnixMilli(), CntTotal: o.CntTotal})
	}

	res.Data = keywordRes

	return res
}

type KeywordDate struct {
	SeqKeyword    int64  `json:"seq_keyword"`
	Keyword       string `json:"keyword"`
	ViewStartDate string `json:"view_start_date"`
	ViewEndDate   string `json:"view_end_date"`
	IsToday       bool   `json:"is_today"`
	StartDate     time.Time
	EndDate       time.Time
	CntTotal      int64 `json:"cnt_total"`
}

type KeywordRes struct {
	List []struct {
		SeqKeyword int64  `json:"seq_keyword"`
		Keyword    string `json:"keyword"`
		IsToday    bool   `json:"is_today"`
		StartDate  int64  `json:"start_date"`
		EndDate    int64  `json:"end_date"`
		CntTotal   int64  `json:"cnt_total"`
	} `json:"list"`
}

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
