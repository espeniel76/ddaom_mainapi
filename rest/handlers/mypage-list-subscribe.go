package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
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

	sdb := db.List[define.DSN_SLAVE1]

	// 받은/보낸구독
	query := `
	SELECT A.seq_member, A.seq_member_following 
	FROM
	(
		(SELECT seq_member, seq_member_following, created_at FROM ddaom_user1.member_subscribes WHERE subscribe_yn = true AND seq_member = ? OR seq_member_following = ?)
		UNION ALL
		(SELECT seq_member, seq_member_following, created_at FROM ddaom_user2.member_subscribes WHERE subscribe_yn = true AND seq_member = ? OR seq_member_following = ?)
	) AS A
	`
	var totalData int64
	seq := userToken.SeqMember
	info := []SubscribeInfo{}
	result := sdb.Raw(query, seq, seq, seq, seq).Scan(&info)
	if corm(result, &res) {
		return res
	}
	cntFollower := 0
	cntFollowing := 0

	for _, v := range info {
		totalData++
		if v.SeqMember == userToken.SeqMember {
			cntFollowing++
		} else {
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

	query = `
	SELECT A.seq_member, A.seq_member_following, A.created_at
	FROM
	(
		(SELECT seq_member, seq_member_following, created_at FROM ddaom_user1.member_subscribes WHERE subscribe_yn = true AND seq_member = ? OR seq_member_following = ?)
		UNION ALL
		(SELECT seq_member, seq_member_following, created_at FROM ddaom_user2.member_subscribes WHERE subscribe_yn = true AND seq_member = ? OR seq_member_following = ?)
	) AS A
	ORDER BY A.created_at DESC
	LIMIT ?, ?
	`
	info = []SubscribeInfo{}
	result = sdb.Raw(query, seq, seq, seq, seq, limitStart, _sizePerPage).Scan(&info)
	if corm(result, &res) {
		return res
	}
	seqs := []int64{}
	for _, v := range info {
		seqs = append(seqs, v.SeqMember)
		seqs = append(seqs, v.SeqMemberFollowing)
	}
	keys := make(map[int64]bool)
	ue := []int64{}
	for _, value := range seqs {
		if _, saveValue := keys[value]; !saveValue {
			keys[value] = true
			ue = append(ue, value)
		}
	}
	fmt.Println(ue)
	query = "SELECT seq_member, nick_name FROM ddaom.member_details WHERE seq_member IN (?)"
	nicks := []TmpMembinfo{}
	sdb.Raw(query, ue).Scan(&nicks)
	fmt.Println(nicks)

	isFollower := false
	seqMember := 0
	nickName := ""
	for _, v := range info {
		if v.SeqMember != userToken.SeqMember {
			seqMember = int(v.SeqMember)
			isFollower = true
		} else {
			seqMember = int(v.SeqMemberFollowing)
			isFollower = false
		}
		for _, k := range nicks {
			if isFollower {
				if v.SeqMember == k.SeqMember {
					nickName = k.NickName
					break
				}
			} else {
				if v.SeqMemberFollowing == k.SeqMember {
					nickName = k.NickName
					break
				}
			}
		}
		o.List = append(o.List, struct {
			SeqMember  int64  "json:\"seq_member\""
			NickName   string "json:\"nick_name\""
			IsFollower bool   "json:\"is_follower\""
		}{
			SeqMember:  int64(seqMember),
			NickName:   nickName,
			IsFollower: isFollower,
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
	SeqMember          int64
	SeqMemberFollowing int64
	CreatedAt          time.Time
}
type MypageListSubscribeRes struct {
	CntFollower  int64 `json:"cnt_follower"`
	CntFollowing int64 `json:"cnt_following"`
	NowPage      int   `json:"now_page"`
	TotalPage    int   `json:"total_page"`
	TotalData    int   `json:"total_data"`
	List         []struct {
		SeqMember  int64  `json:"seq_member"`
		NickName   string `json:"nick_name"`
		IsFollower bool   `json:"is_follower"`
	} `json:"list"`
}
