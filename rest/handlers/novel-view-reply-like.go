package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"fmt"
	"strings"
)

func NovelViewReplyLike(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqReply := CpInt64(req.Parameters, "seq_reply")

	mdb := db.List[define.Mconn.DsnMaster]

	// 이전에 좋아요 기록이 있는지 체크
	novelReply := schemas.NovelReply{}
	result := mdb.Model(&novelReply).Where("seq_reply = ?", _seqReply).Scan(&novelReply)
	if corm(result, &res) {
		return res
	}

	// 내 seq 문자화
	mySeqMemberString := fmt.Sprint(userToken.SeqMember)
	myLike := false
	list := strings.Split(novelReply.Likes, ",")
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 0 {
			if list[i] == mySeqMemberString {
				myLike = true
				break
			}
		}
	}

	// 현재 나의 좋아요 상태
	// fmt.Println(myLike)
	listFinal := []string{}
	if !myLike {
		// fmt.Println("없으면 추가")
		list = append(list, mySeqMemberString)
		listFinal = list
	} else {
		// fmt.Println("있으면 빼기")
		for i := 0; i < len(list); i++ {
			if len(list[i]) > 0 {
				if list[i] != mySeqMemberString {
					listFinal = append(listFinal, list[i])
				}
			}
		}
	}
	// fmt.Println(listFinal)

	// 문자열 신규 조합
	likes := ""
	cntLike := novelReply.CntLike
	for _, s := range listFinal {
		if s != "" {
			likes = likes + s + ","
		}
	}
	// fmt.Println(likes)
	if !myLike {
		cntLike += 1
	} else {
		cntLike -= 1
	}
	mdb.Exec("UPDATE novel_replies SET cnt_like = ?, likes = ? WHERE seq_reply = ?", cntLike, likes, _seqReply)

	data := make(map[string]interface{})
	data["my_like"] = !myLike
	data["cnt_like"] = cntLike
	res.Data = data

	return res
}

func NovelViewReReplyLike(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqReReply := CpInt64(req.Parameters, "seq_re_reply")

	mdb := db.List[define.Mconn.DsnMaster]

	// 이전에 좋아요 기록이 있는지 체크
	novelReReply := schemas.NovelReReply{}
	result := mdb.Model(&novelReReply).Where("seq_re_reply = ?", _seqReReply).Scan(&novelReReply)
	if corm(result, &res) {
		return res
	}

	// 내 seq 문자화
	mySeqMemberString := fmt.Sprint(userToken.SeqMember)
	myLike := false
	list := strings.Split(novelReReply.Likes, ",")
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 0 {
			if list[i] == mySeqMemberString {
				myLike = true
				break
			}
		}
	}

	// 현재 나의 좋아요 상태
	fmt.Println(myLike)
	listFinal := []string{}
	if !myLike {
		fmt.Println("없으면 추가")
		list = append(list, mySeqMemberString)
		listFinal = list
	} else {
		fmt.Println("있으면 빼기")
		for i := 0; i < len(list); i++ {
			if len(list[i]) > 0 {
				if list[i] != mySeqMemberString {
					listFinal = append(listFinal, list[i])
				}
			}
		}
	}
	fmt.Println(listFinal)

	// 문자열 신규 조합
	likes := ""
	cntLike := novelReReply.CntLike
	for _, s := range listFinal {
		if s != "" {
			likes = likes + s + ","
		}
	}
	fmt.Println(likes)
	if !myLike {
		cntLike += 1
	} else {
		cntLike -= 1
	}
	mdb.Exec("UPDATE novel_re_replies SET cnt_like = ?, likes = ? WHERE seq_re_reply = ?", cntLike, likes, _seqReReply)

	data := make(map[string]interface{})
	data["my_like"] = !myLike
	data["cnt_like"] = cntLike
	res.Data = data

	return res
}
