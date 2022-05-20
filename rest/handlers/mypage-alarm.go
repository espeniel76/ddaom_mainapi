package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"strconv"
)

func MypageListAlarm(req *domain.CommonRequest) domain.CommonResponse {

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

	var totalData int64
	sdb := db.List[define.DSN_SLAVE]
	result := sdb.
		Model(schemas.Alarm{}).
		Where("seq_member = ?", userToken.SeqMember).
		Count(&totalData)
	if corm(result, &res) {
		return res
	}

	alarmListRes := AlarmListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	query := `
	SELECT
		seq_alarm,
		title,
		content AS body,
		type_alarm,
		value_alarm,
		step,
		UNIX_TIMESTAMP(created_at) * 1000 AS created_at,
		is_read,
		UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at
	FROM alarms
	WHERE seq_member = ?
	ORDER BY is_read ASC, seq_alarm DESC
	LIMIT ?, ?
	`
	result = sdb.Raw(query, userToken.SeqMember, limitStart, _sizePerPage).Find(&alarmListRes.List)
	if corm(result, &res) {
		return res
	}

	res.Data = alarmListRes

	return res
}

func MypageAlarmReceiveSet(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqAlarm, _ := strconv.Atoi(req.Vars["seq_alarm"])
	query := "UPDATE alarms SET is_read = true, updated_at = NOW() WHERE seq_alarm = ? AND seq_member = ?"
	mdb := db.List[define.DSN_MASTER]
	result := mdb.Exec(query, _seqAlarm, userToken.SeqMember)
	if corm(result, &res) {
		return res
	}

	return res
}

type AlarmListRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqAlarm   int64   `json:"seq_alarm"`
		Title      string  `json:"title"`
		Body       string  `json:"body"`
		TypeAlarm  int8    `json:"type_alarm"`
		ValueAlarm int     `json:"value_alarm"`
		Step       int8    `json:"step"`
		CreatedAt  float64 `json:"created_at"`
		IsRead     bool    `json:"is_read"`
		UpdatedAt  float64 `json:"updated_at"`
	} `json:"list"`
}
