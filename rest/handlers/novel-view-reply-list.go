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

func NovelViewReplyList(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}

	_step := CpInt64(req.Parameters, "step")
	_seqNovel := CpInt64(req.Parameters, "seq_novel")
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
	query.WriteString(`SELECT COUNT(nr.seq_reply) FROM novel_replies nr WHERE nr.step = ? AND seq_novel = ?`)
	result := sdb.Raw(query.String(), _step, _seqNovel).Count(&totalData)
	if result.Error != nil {
		res.ResultCode = define.OK
		res.ErrorDesc = result.Error.Error()
		return res
	}

	novelListLReplyRes := NovelListReplyRes{
		NowPage:   int(_page),
		TotalPage: tools.GetTotalPage(totalData, _sizePerPage),
		TotalData: int(totalData),
	}

	query.Reset()
	query.WriteString(`
		SELECT
			nr.seq_reply,
			nr.seq_member,
			md.nick_name,
			md.profile_photo,
			nr.cnt_like,
			nr.likes,
			nr.contents,
			nr.cnt_re_reply,
			UNIX_TIMESTAMP(nr.created_at) * 1000 AS created_at,
			UNIX_TIMESTAMP(nr.updated_at) * 1000 AS updated_at
		FROM
			novel_replies nr
		INNER JOIN member_details md
		ON nr.seq_member = md.seq_member
		WHERE nr.step = ? AND seq_novel = ?
		ORDER BY nr.updated_at DESC
		LIMIT ?, ?
	`)
	result = sdb.Raw(query.String(), _step, _seqNovel, limitStart, _sizePerPage).Find(&novelListLReplyRes.List)
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

		for i := 0; i < len(novelListLReplyRes.List); i++ {

			novelListLReplyRes.List[i].MyLike = false
			list := strings.Split(novelListLReplyRes.List[i].Likes, ",")

			for _, seqMember := range list {
				if len(seqMember) > 0 {
					if seqMember == mySeqMemberString {
						novelListLReplyRes.List[i].MyLike = true
						break
					}
				}
			}
		}
	}

	res.Data = novelListLReplyRes

	return res
}

type NovelListReplyRes struct {
	NowPage   int `json:"now_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
	List      []struct {
		SeqReply     int64   `json:"seq_reply"`
		SeqMember    int64   `json:"seq_member"`
		NickName     string  `json:"nick_name"`
		ProfilePhoto string  `json:"profile_photo"`
		CntLike      int64   `json:"cnt_like"`
		Likes        string  `json:"likes"`
		MyLike       bool    `json:"my_like"`
		Contents     string  `json:"contents"`
		CntReReply   int64   `json:"cnt_re_reply"`
		CreatedAt    float64 `json:"created_at"`
	} `json:"list"`
}
