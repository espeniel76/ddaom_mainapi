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

	_seqMember := CpInt64(req.Parameters, "seq_member")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.JWT_ACCESS_SECRET)
	isLogin := false
	itsMine := false
	if userToken != nil {
		isLogin = true
		if _seqMember == 0 {
			itsMine = true
			_seqMember = userToken.SeqMember
		}
	}

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	sdb := db.List[define.DSN_SLAVE]
	ldb := getUserLogDb(sdb, _seqMember)

	// 구독현황
	var totalData int64
	info := []SubscribeInfo{}
	result := ldb.Model(&schemas.MemberSubscribe{}).Select("seq_member, seq_member_opponent, status, created_at").
		Where("seq_member = ?", _seqMember).Scan(&info)
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
		Where("seq_member = ?", _seqMember).
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

	// 닉네임 딕셔너리
	nicks := []TmpMembinfo{}
	result = sdb.Raw("SELECT seq_member, nick_name FROM member_details WHERE seq_member IN (?)", seqs).Scan(&nicks)
	if corm(result, &res) {
		return res
	}

	// status  딕셔너리 (로그인 했고, 남의 구독 상태 볼 때)
	statuses := []TmpStatinfo{}
	if isLogin && !itsMine {
		// 내 구독 목록 가져옴
		myLdb := getUserLogDb(sdb, userToken.SeqMember)
		result = myLdb.Raw("SELECT seq_member_opponent, status FROM member_subscribes WHERE seq_member = ?", userToken.SeqMember).Scan(&statuses)
		if corm(result, &res) {
			return res
		}
		fmt.Println(statuses)
	}

	nickName := ""
	isYou := false
	mySubscribe := ""
	for _, v := range info {
		for _, k := range nicks {
			if v.SeqMemberOpponent == k.SeqMember {
				nickName = k.NickName
				break
			}
		}
		isYou = false
		if userToken != nil {
			if v.SeqMemberOpponent == userToken.SeqMember {
				isYou = true
			}
		}
		if isLogin && !itsMine {
			isExist := false
			for _, o := range statuses {
				if o.SeqMemberOpponent == v.SeqMemberOpponent {
					isExist = true
					mySubscribe = o.Status
					break
				}
			}
			if !isExist {
				mySubscribe = define.NONE
			}
		} else {
			mySubscribe = v.Status
		}

		o.List = append(o.List, struct {
			SeqMember   int64  "json:\"seq_member\""
			NickName    string "json:\"nick_name\""
			IsYou       bool   "json:\"is_you\""
			MySubscribe string "json:\"my_subscribe\""
		}{
			SeqMember:   int64(v.SeqMemberOpponent),
			NickName:    nickName,
			IsYou:       isYou,
			MySubscribe: mySubscribe,
		})
	}

	res.Data = o

	return res
}

type TmpStatinfo struct {
	SeqMemberOpponent int64
	Status            string
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
		IsYou       bool   `json:"is_you"`
		MySubscribe string `json:"my_subscribe"`
	} `json:"list"`
}
