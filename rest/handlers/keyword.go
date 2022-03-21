package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"fmt"
	"time"
)

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
		WHERE k.active_yn = true
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
		}{SeqKeyword: o.SeqKeyword, Keyword: o.Keyword, IsToday: isToday, StartDate: o.StartDate.UnixMilli(), EndDate: o.EndDate.UnixMilli(), CntTotal: int64(o.EndDate.Nanosecond())})
	}

	res.Data = keywordRes

	return res
}

type KeywordDate struct {
	SeqKeyword int64     `json:"seq_keyword"`
	Keyword    string    `json:"keyword"`
	ViewDate   string    `json:"view_date"`
	IsToday    bool      `json:"is_today"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CntTotal   int64     `json:"cnt_total"`
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
