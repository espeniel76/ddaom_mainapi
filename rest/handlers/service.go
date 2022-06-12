package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
	"fmt"

	"gorm.io/gorm"
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
	fmt.Println(_emailYn)

	mdb := db.List[define.DSN_MASTER]
	m := schemas.ServiceInquiry{
		SeqMember: userToken.SeqMember,
		Title:     _title,
		Content:   _content,
		EmailYn:   _emailYn,
	}
	fmt.Println(m)
	result := mdb.Model(&m).Create(&m)
	if corm(result, &res) {
		return res
	}

	return res
}

func ServiceInquiryEdit(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqServiceInquiry := CpInt64(req.Parameters, "seq_service_inquiry")
	_title := Cp(req.Parameters, "title")
	_content := Cp(req.Parameters, "content")
	_emailYn := CpBool(req.Parameters, "email_yn")

	mdb := db.List[define.DSN_MASTER]
	m := schemas.ServiceInquiry{}
	result := mdb.Model(&m).Where("seq_service_inquiry = ?", _seqServiceInquiry).Scan(&m)
	if corm(result, &res) {
		return res
	}
	if m.SeqServiceInquiry == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}
	if m.SeqMember != userToken.SeqMember {
		res.ResultCode = define.OTHER_USER
		return res
	}
	result = mdb.Exec("UPDATE service_inquiries SET title = ?, content = ?, email_yn = ? WHERE seq_service_inquiry = ?",
		_title, _content, _emailYn, _seqServiceInquiry)
	if corm(result, &res) {
		return res
	}

	return res
}

func ServiceInquiryDelete(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	_seqServiceInquiry := CpInt64(req.Parameters, "seq_service_inquiry")

	mdb := db.List[define.DSN_MASTER]
	m := schemas.ServiceInquiry{}
	result := mdb.Model(&m).Where("seq_service_inquiry = ?", _seqServiceInquiry).Scan(&m)
	if corm(result, &res) {
		return res
	}
	if m.SeqServiceInquiry == 0 {
		res.ResultCode = define.NO_EXIST_DATA
		return res
	}
	if m.SeqMember != userToken.SeqMember {
		res.ResultCode = define.OTHER_USER
		return res
	}
	result = mdb.Exec("DELETE FROM service_inquiries WHERE seq_service_inquiry = ?", _seqServiceInquiry)
	if corm(result, &res) {
		return res
	}

	return res
}

func ServiceInquiryList(req *domain.CommonRequest) domain.CommonResponse {

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
	result := sdb.Model(schemas.ServiceInquiry{}).
		Where("seq_member = ?", userToken.SeqMember).
		Count(&totalData)
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
			answer,
			UNIX_TIMESTAMP(created_at) * 1000 AS created_at,
			UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at
		`).
		Where("seq_member = ?", userToken.SeqMember).
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
		Answer            string  `json:"answer"`
		CreatedAt         float64 `json:"created_at"`
		UpdatedAt         float64 `json:"updated_at"`
	} `json:"list"`
}

func ServiceNoticeList(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.DSN_SLAVE]
	result := sdb.Model(schemas.Notice{}).Count(&totalData)
	if corm(result, &res) {
		return res
	}

	m := ServiceNoticeListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	result = sdb.Model(schemas.Notice{}).
		Select(`
			seq_notice,
			title,
			content,
			UNIX_TIMESTAMP(created_at) * 1000 AS created_at,
			UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at
		`).
		Where("active_yn = true").
		Order("seq_notice DESC").
		Offset(int(limitStart)).
		Limit(int(_sizePerPage)).
		Scan(&m.List)

	if corm(result, &res) {
		return res
	}

	res.Data = m

	return res
}

type ServiceNoticeListRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqNotice int64   `json:"seq_notice"`
		Title     string  `json:"title"`
		Content   string  `json:"content"`
		CreatedAt float64 `json:"created_at"`
		UpdatedAt float64 `json:"updated_at"`
	} `json:"list"`
}

func ServiceFaqList(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqCategoryFaq := CpInt64(req.Parameters, "seq_category_faq")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.DSN_SLAVE]
	var result *gorm.DB
	if _seqCategoryFaq > 0 {
		result = sdb.Model(schemas.Faq{}).Where("seq_category_faq = ?", _seqCategoryFaq).Count(&totalData)
	} else {
		result = sdb.Model(schemas.Faq{}).Count(&totalData)
	}
	if corm(result, &res) {
		return res
	}

	m := ServiceFaqListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	if _seqCategoryFaq > 0 {
		result = sdb.Model(schemas.Faq{}).
			Select(`
			seq_faq,
			seq_category_faq,
			title,
			content,
			UNIX_TIMESTAMP(created_at) * 1000 AS created_at,
			UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at
		`).
			Where("active_yn = true").
			Where("seq_category_faq = ?", _seqCategoryFaq).
			Order("seq_faq DESC").
			Offset(int(limitStart)).
			Limit(int(_sizePerPage)).
			Scan(&m.List)
	} else {
		result = sdb.Model(schemas.Faq{}).
			Select(`
			seq_faq,
			seq_category_faq,
			title,
			content,
			UNIX_TIMESTAMP(created_at) * 1000 AS created_at,
			UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at
		`).
			Where("active_yn = true").
			Order("seq_faq DESC").
			Offset(int(limitStart)).
			Limit(int(_sizePerPage)).
			Scan(&m.List)
	}

	if corm(result, &res) {
		return res
	}

	res.Data = m

	return res
}

type ServiceFaqListRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqFaq         int64   `json:"seq_faq"`
		SeqCategoryFaq int64   `json:"seq_category_faq"`
		Title          string  `json:"title"`
		Content        string  `json:"content"`
		CreatedAt      float64 `json:"created_at"`
		UpdatedAt      float64 `json:"updated_at"`
	} `json:"list"`
}
