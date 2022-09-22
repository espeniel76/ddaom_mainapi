package handlers

import (
	"bytes"
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"fmt"
	"strings"
)

func NovelViewReReplyList(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_seqReply := CpInt64(req.Parameters, "seq_reply")
	_page := CpInt64(req.Parameters, "page")
	_sizePerPage := CpInt64(req.Parameters, "size_per_page")

	if _page < 1 || _sizePerPage < 1 {
		res.ResultCode = define.REQUIRE_OVER_1
		return res
	}
	limitStart := (_page - 1) * _sizePerPage

	var totalData int64
	sdb := db.List[define.Mconn.DsnSlave]
	var query bytes.Buffer
	query.WriteString(`SELECT COUNT(nr.seq_re_reply) FROM novel_re_replies nr WHERE nr.seq_reply = ?`)
	result := sdb.Raw(query.String(), _seqReply).Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelListLReReplyRes := NovelListReReplyRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	query.Reset()
	query.WriteString(`
		SELECT
			nr.seq_re_reply,
			nr.seq_member,
			md.nick_name,
			md.profile_photo,
			nr.cnt_like,
			nr.likes,
			nr.contents,
			nr.seq_re_reply_org,
			nr.seq_member_org,
			(SELECT nick_name FROM member_details md2 WHERE seq_member = nr.seq_member_org) AS nick_name_org,
			UNIX_TIMESTAMP(nr.created_at) * 1000 AS created_at
		FROM
			novel_re_replies nr
		INNER JOIN member_details md
		ON nr.seq_member = md.seq_member
		WHERE nr.seq_reply = ?
		ORDER BY nr.updated_at DESC
		LIMIT ?, ?
	`)
	result = sdb.Raw(query.String(), _seqReply, limitStart, _sizePerPage).Find(&novelListLReReplyRes.List)
	if result.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = result.Error.Error()
		return res
	}

	// 로그인 상태면, 내가 좋아요 한 게시물 여부
	userToken, _ := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if userToken != nil {

		// 내 seq 문자화
		mySeqMemberString := fmt.Sprint(userToken.SeqMember)

		for i := 0; i < len(novelListLReReplyRes.List); i++ {

			novelListLReReplyRes.List[i].MyLike = false
			list := strings.Split(novelListLReReplyRes.List[i].Likes, ",")

			for _, seqMember := range list {
				if len(seqMember) > 0 {
					if seqMember == mySeqMemberString {
						novelListLReReplyRes.List[i].MyLike = true
						break
					}
				}
			}

		}
	}

	res.Data = novelListLReReplyRes

	return res
}

type NovelListReReplyRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqReReply    int64   `json:"seq_re_reply"`
		SeqMember     int64   `json:"seq_member"`
		NickName      string  `json:"nick_name"`
		ProfilePhoto  string  `json:"profile_photo"`
		CntLike       int64   `json:"cnt_like"`
		Likes         string  `json:"likes"`
		MyLike        bool    `json:"my_like"`
		Contents      string  `json:"contents"`
		SeqReReplyOrg int64   `json:"seq_re_reply_org"`
		SeqMemberOrg  int64   `json:"seq_member_org"`
		NickNameOrg   string  `json:"nick_name_org"`
		CreatedAt     float64 `json:"created_at"`
	} `json:"list"`
}
