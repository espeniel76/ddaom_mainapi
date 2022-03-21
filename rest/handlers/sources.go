package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"
	"time"
)

func Assets(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	today := tools.TodayFormattedDate()
	masterDB := db.List[define.DSN_MASTER]
	query := `
			SELECT
			k.seq_keyword,
			k.keyword,
			kt.view_date,
			k.start_date,
			k.end_date,
			k.cnt_total
		FROM ddaom.keywords AS k
		INNER JOIN ddaom.keyword_todays AS kt ON k.seq_keyword = kt.seq_keyword
		WHERE k.active_yn = true AND NOW() BETWEEN k.start_date AND k.end_date
		ORDER BY kt.view_date ASC
	`
	keywordDate := []KeywordDate{}
	result := masterDB.Raw(query).Scan(&keywordDate)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	assetRes := AssetRes{}
	for i := 0; i < len(keywordDate); i++ {
		o := keywordDate[i]
		isToday := false
		if o.ViewDate == today {
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

	result = masterDB.Model(&schemas.Image{}).Where("active_yn = true").Find(&assetRes.ListImage)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	result = masterDB.Model(&schemas.Color{}).Where("active_yn = true").Find(&assetRes.ListColor)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	result = masterDB.Model(&schemas.Genre{}).Where("active_yn = true").Find(&assetRes.ListGenre)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
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
}

func Keyword(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	// 오늘 날짜
	today := tools.TodayFormattedDate()
	fmt.Println(today)

	// 데이터 가져온다.
	masterDB := db.List[define.DSN_MASTER]
	query := `
			SELECT
			k.seq_keyword,
			k.keyword,
			kt.view_date,
			k.start_date,
			k.end_date,
			k.cnt_total
		FROM ddaom.keywords AS k
		INNER JOIN ddaom.keyword_todays AS kt ON k.seq_keyword = kt.seq_keyword
		WHERE k.active_yn = true AND NOW() BETWEEN k.start_date AND k.end_date
		ORDER BY kt.view_date ASC
	`
	keywordDate := []KeywordDate{}
	result := masterDB.Raw(query).Scan(&keywordDate)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}
	keywordRes := KeywordRes{}
	for i := 0; i < len(keywordDate); i++ {
		o := keywordDate[i]
		isToday := false
		if o.ViewDate == today {
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
	SeqKeyword int64  `json:"seq_keyword"`
	Keyword    string `json:"keyword"`
	ViewDate   string `json:"view_date"`
	IsToday    bool   `json:"is_today"`
	StartDate  time.Time
	EndDate    time.Time
	CntTotal   int64 `json:"cnt_total"`
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
	masterDB := db.List[define.DSN_MASTER]
	skinRes := SkinRes{}
	result := masterDB.Model(&schemas.Image{}).Where("active_yn = true").Find(&skinRes.List)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
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
	masterDB := db.List[define.DSN_MASTER]
	genreRes := GenreRes{}
	result := masterDB.Model(&schemas.Genre{}).Where("active_yn = true").Find(&genreRes.List)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
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
