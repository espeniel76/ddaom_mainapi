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

func MypageListSubscribe(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	fmt.Println(userToken)

	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	sdb := db.List[define.DSN_SLAVE]
	ldb := GetMyLogDb(userToken.Allocated)

	// 구독현황
	var totalData int64
	info := []SubscribeInfo{}
	result := ldb.Model(&schemas.MemberSubscribe{}).Select("seq_member, seq_member_opponent, status, created_at").
		Where("seq_member = ?", userToken.SeqMember).Scan(&info)
	if corm(result, &res) {
		return res
	}
	cntFollower := 0
	cntFollowing := 0
	for _, v := range info {
		totalData++
		switch v.Status {
		case define.FOLLOWING:
			cntFollowing++
		case define.FOLLOWER:
			cntFollower++
		case define.BOTH:
			cntFollowing++
			cntFollower++
		}
	}
	o := MypageListSubscribeRes{
		CntFollower:  int64(cntFollower),
		CntFollowing: int64(cntFollowing),
		NowPage:      int(_page),
		TotalPage:    tools.GetTotalPage(totalData, _sizePerPage),
		TotalData:    int(totalData),
	}

	result = ldb.Model(&schemas.MemberSubscribe{}).
		Select("seq_member, seq_member_opponent, status").
		Where("seq_member = ?", userToken.SeqMember).
		Order("updated_at DESC").
		Limit(int(_sizePerPage)).
		Offset(int(limitStart)).
		Scan(&info)
	if corm(result, &res) {
		return res
	}
	seqs := []int64{}
	for _, v := range info {
		seqs = append(seqs, v.SeqMemberOpponent)
	}
	nicks := []TmpMembinfo{}
	sdb.Raw("SELECT seq_member, nick_name FROM member_details WHERE seq_member IN (?)", seqs).Scan(&nicks)
	if corm(result, &res) {
		return res
	}
	fmt.Println(nicks)

	nickName := ""
	for _, v := range info {
		for _, k := range nicks {
			if v.SeqMemberOpponent == k.SeqMember {
				nickName = k.NickName
				break
			}
		}
		o.List = append(o.List, struct {
			SeqMember   int64  "json:\"seq_member\""
			NickName    string "json:\"nick_name\""
			MySubscribe string "json:\"my_subscribe\""
		}{
			SeqMember:   int64(v.SeqMemberOpponent),
			NickName:    nickName,
			MySubscribe: v.Status,
		})
	}

	res.Data = o

	return res
}

type TmpMembinfo struct {
	SeqMember int64
	NickName  string
}
type SubscribeInfo struct {
	SeqMember         int64
	SeqMemberOpponent int64
	Status            string
	CreatedAt         time.Time
}
type MypageListSubscribeRes struct {
	CntFollower  int64 `json:"cnt_follower"`
	CntFollowing int64 `json:"cnt_following"`
	NowPage      int   `json:"now_page"`
	TotalPage    int   `json:"total_page"`
	TotalData    int   `json:"total_data"`
	List         []struct {
		SeqMember   int64  `json:"seq_member"`
		NickName    string `json:"nick_name"`
		MySubscribe string `json:"my_subscribe"`
	} `json:"list"`
}
