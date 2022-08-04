package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/tools"
)

func MypageUserBlockList(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
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
	ldb := GetMyLogDbSlave(userToken.Allocated)
	memberBlocking := schemas.MemberBlocking{}
	result := ldb.Model(&memberBlocking).
		Where("seq_member = ? AND block_yn = true", userToken.SeqMember).
		Count(&totalData)
	if corm(result, &res) {
		return res
	}
	memberBlockingListRes := MemberBlockingListRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}
	query := `
	SELECT
		seq_member_to AS seq_member,
		UNIX_TIMESTAMP(updated_at) * 1000 AS updated_at,
		cnt_block
	FROM
		member_blockings
	WHERE
		seq_member = ?
		AND block_yn = true
	ORDER BY updated_at DESC
	LIMIT ?, ?
	`
	result = ldb.Raw(query, userToken.SeqMember, limitStart, _sizePerPage).Scan(&memberBlockingListRes.List)
	if corm(result, &res) {
		return res
	}

	seqs := []int64{}
	for _, v := range memberBlockingListRes.List {
		seqs = append(seqs, v.SeqMember)
	}
	// 닉네임 딕셔너리
	nicks := []TmpMembinfo{}
	sdb := db.List[define.Mconn.DsnSlave]
	sdb.Raw("SELECT seq_member, nick_name FROM member_details WHERE seq_member IN (?)", seqs).Scan(&nicks)

	// 닉네임 재 할당
	for i := 0; i < len(memberBlockingListRes.List); i++ {
		for _, v := range nicks {
			if v.SeqMember == memberBlockingListRes.List[i].SeqMember {
				memberBlockingListRes.List[i].NickName = v.NickName
				break
			}
		}
	}

	res.Data = memberBlockingListRes

	return res
}

type MemberBlockingListRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqMember int64   `json:"seq_member"`
		NickName  string  `json:"nick_name"`
		UpdatedAt float64 `json:"updated_at"`
		CntBlock  int64   `json:"cnt_block"`
	} `json:"list"`
}
