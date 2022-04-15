package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

func ServiceInquiry(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_title := Cp(req.Parameters, "title")
	_content := Cp(req.Parameters, "content")
	_emailYn := CpBool(req.Parameters, "email_yn")

	mdb := db.List[define.DSN_MASTER]
	m := schemas.ServiceInquiry{
		SeqMember: userToken.SeqMember,
		Title:     _title,
		Content:   _content,
		EmailYn:   _emailYn,
	}
	result := mdb.Model(&m).Create(&m)
	if corm(result, &res) {
		return res
	}

	return res
}

func ServiceInquiryList(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.DSN_SLAVE1]
	result := sdb.Model(schemas.ServiceInquiry{}).Count(&totalData)
	if corm(result, &res) {
		return res
	}

	m := ServiceInquiryListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	result = sdb.Model(schemas.ServiceInquiry{}).
		Select(`
		seq_service_inquiry,
		title,
		content,
		status,
		UNIX_TIMESTAMP(created_at) * 1000 AS created_at,
		UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at
		`).
		Order("seq_service_inquiry DESC").
		Offset(int(limitStart)).
		Limit(int(_sizePerPage)).
		Scan(&m.List)

	if corm(result, &res) {
		return res
	}

	res.Data = m

	return res
}

type ServiceInquiryListRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqServiceInquiry int64   `json:"seq_service_inquiry"`
		Title             string  `json:"title"`
		Content           string  `json:"content"`
		Status            int8    `json:"status"`
		CreatedAt         float64 `json:"created_at"`
		UpdatedAt         float64 `json:"updated_at"`
	} `json:"list"`
}
