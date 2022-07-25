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

/**
 * step1에 새로 등록할때 제목 중복 체크 기준 수정 요청
 * => 완결로 등록된 소설 제목, 연재중 소설 제목만 중복 체크되도록 하고,
 * 그 외 종료된 주제어의 소설 제목들은 중복 체크에서 제외"
 */
func NovelCheckTitle(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	_title := Cp(req.Parameters, "title")

	sdb := db.List[define.Mconn.DsnSlave]
	var cnt int64
	isExist := false
	result := sdb.Model(schemas.NovelStep1{}).Where("title = ?", _title).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	// 현재 진행중인 소설
	if cnt > 0 {
		isExist = true
	} else {
		// 완결로 등록된 소설
		result = sdb.Model(schemas.NovelFinish{}).Where("title = ?", _title).Count(&cnt)
		if corm(result, &res) {
			return res
		}
		if cnt > 0 {
			isExist = true
		}

	}
	data := make(map[string]bool)
	data["is_exist"] = isExist
	res.Data = data

	return res
}

func NovelCheckBlocked(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}
	data := make(map[string]bool)
	data["blocked_yn"] = isBlocked(userToken.SeqMember)
	res.Data = data

	return res
}

func NovelWriteStep1(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqKeyword := CpInt64(req.Parameters, "seq_keyword")
	_seqGenre := CpInt64(req.Parameters, "seq_genre")
	_seqImage := CpInt64(req.Parameters, "seq_image")
	_seqColor := CpInt64(req.Parameters, "seq_color")
	_title := Cp(req.Parameters, "title")
	_content := Cp(req.Parameters, "content")
	_tempYn := CpBool(req.Parameters, "temp_yn")

	// 블록처리된 유저 여부
	if isBlocked(userToken.SeqMember) {
		res.ResultCode = define.BLOCKED_ME
		return res
	}

	// 존재하는 닉네임 여부
	mdb := db.List[define.Mconn.DsnMaster]
	var cnt int64
	result := mdb.Model(schemas.MemberDetail{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_NICK
		return res
	}

	// 가용 키워드/이미지/컬러 검사
	if isAbleKeyword(_seqKeyword) != true {
		res.ResultCode = define.INACTIVE_KEYWORD
		return res
	}
	if isAbleImage(_seqImage) != true {
		res.ResultCode = define.INACTIVE_IMAGE
		return res
	}
	if isAbleColor(_seqColor) != true {
		res.ResultCode = define.INACTIVE_COLOR
		return res
	}
	if isAbleGenre(_seqGenre) != true {
		res.ResultCode = define.INACTIVE_GENRE
		return res
	}

	if _seqNovelStep1 == 0 { // 신규 작성
		novelWriteStep1 := schemas.NovelStep1{
			SeqKeyword: _seqKeyword,
			SeqImage:   _seqImage,
			SeqColor:   _seqColor,
			SeqGenre:   _seqGenre,
			SeqMember:  userToken.SeqMember,
			Title:      _title,
			Content:    _content,
			TempYn:     _tempYn,
			DeletedAt:  time.Now(),
		}

		result = mdb.Model(&novelWriteStep1).Where("title = ?", _title).Count(&cnt)
		if corm(result, &res) {
			return res
		}
		if cnt > 0 {
			res.ResultCode = define.ALREADY_EXISTS_TITLE
			return res
		}
		result = mdb.Model(&novelWriteStep1).Create(&novelWriteStep1)
		if corm(result, &res) {
			return res
		}
		addKeywordCnt(_seqKeyword)
		_seqNovelStep1 = novelWriteStep1.SeqNovelStep1

	} else { // 업데이트

		novelStep := schemas.NovelStep1{}
		result := mdb.Model(&novelStep).Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&novelStep)
		if corm(result, &res) {
			return res
		}
		if novelStep.SeqNovelStep1 == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		if novelStep.SeqMember != userToken.SeqMember {
			res.ResultCode = define.OTHER_USER
			return res
		}

		query := `UPDATE novel_step1 SET
			seq_genre = ?,
			seq_image = ?,
			seq_color = ?,
			seq_keyword = ?,
			title = ?,
			content = ?,
			temp_yn = ?,
			updated_at = NOW()
		WHERE seq_novel_step1 = ?
		`
		result = mdb.Exec(query, _seqGenre, _seqImage, _seqColor, _seqKeyword, _title, _content, _tempYn, _seqNovelStep1)
		if corm(result, &res) {
			return res
		}
	}

	if !_tempYn {
		go cacheMainLive(_seqKeyword)
		go pushWriteTopic(userToken, 1, _seqNovelStep1)
		go educeImage(_seqColor, _seqImage, _seqNovelStep1)
	}

	return res
}

func NovelWriteStep2(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep1 := CpInt64(req.Parameters, "seq_novel_step1")
	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_content := Cp(req.Parameters, "content")
	_tempYn := CpBool(req.Parameters, "temp_yn")

	// 블록처리된 유저 여부
	if isBlocked(userToken.SeqMember) {
		res.ResultCode = define.BLOCKED_ME
		return res
	}

	mdb := db.List[define.Mconn.DsnMaster]
	var cnt int64
	result := mdb.Model(schemas.MemberDetail{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_NICK
		return res
	}

	// step 2 단계 글 신규
	if _seqNovelStep2 == 0 {
		// 상위글 삭제 여부 검사
		isDeleted := false
		result = mdb.Model(schemas.NovelStep1{}).Select("deleted_yn").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&isDeleted)
		if corm(result, &res) {
			return res
		}
		if isDeleted {
			res.ResultCode = define.DELETED_PARENT
			return res
		}

		// 가용 키워드 검사
		seqKeyword := getSeqKeyword(1, int64(_seqNovelStep1))
		if isAbleKeyword(seqKeyword) != true {
			res.ResultCode = define.INACTIVE_KEYWORD
			return res
		}

		result = mdb.Model(schemas.NovelStep1{}).Where("seq_novel_step1 = ?", _seqNovelStep1).Count(&cnt)
		if corm(result, &res) {
			return res
		}
		if cnt == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		if !_tempYn {
			mdb.Exec("UPDATE novel_step1 SET cnt_step2 = cnt_step2 + 1 WHERE seq_novel_step1 = ?", _seqNovelStep1)
		}

		novelStep2 := schemas.NovelStep2{
			SeqNovelStep1: _seqNovelStep1,
			SeqMember:     userToken.SeqMember,
			Content:       _content,
			TempYn:        _tempYn,
			DeletedAt:     time.Now(),
		}
		result = mdb.Save(&novelStep2)
		if corm(result, &res) {
			return res
		}
		addKeywordCnt(seqKeyword)
		// step 2 단계 글 기존
	} else {
		seqKeyword := getSeqKeyword(2, int64(_seqNovelStep2))
		if isAbleKeyword(seqKeyword) != true {
			res.ResultCode = define.INACTIVE_KEYWORD
			return res
		}

		novelStep := schemas.NovelStep2{}
		result := mdb.Model(&novelStep).Where("seq_novel_step2 = ?", _seqNovelStep2).Scan(&novelStep)
		if corm(result, &res) {
			return res
		}
		if novelStep.SeqNovelStep2 == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		// 상위글 삭제 여부 검사
		// for push
		_seqNovelStep1 = novelStep.SeqNovelStep1
		isDeleted := false
		result = mdb.Model(schemas.NovelStep1{}).Select("deleted_yn").Where("seq_novel_step1 = ?", _seqNovelStep1).Scan(&isDeleted)
		if corm(result, &res) {
			return res
		}
		if isDeleted {
			res.ResultCode = define.DELETED_PARENT
			return res
		}
		if novelStep.SeqMember != userToken.SeqMember {
			res.ResultCode = define.OTHER_USER
			return res
		}

		result = mdb.Model(&novelStep).
			Where("seq_novel_step2 = ?", _seqNovelStep2).
			Updates(map[string]interface{}{"content": _content, "temp_yn": _tempYn, "updated_at": time.Now()})
		if corm(result, &res) {
			return res
		}
	}

	if !_tempYn {
		go pushWriteTopic(userToken, 2, _seqNovelStep1)
	}

	return res
}

func NovelWriteStep3(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep2 := CpInt64(req.Parameters, "seq_novel_step2")
	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_content := Cp(req.Parameters, "content")
	_tempYn := CpBool(req.Parameters, "temp_yn")

	// 블록처리된 유저 여부
	if isBlocked(userToken.SeqMember) {
		res.ResultCode = define.BLOCKED_ME
		return res
	}

	mdb := db.List[define.Mconn.DsnMaster]
	var cnt int64
	result := mdb.Model(schemas.MemberDetail{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_NICK
		return res
	}

	if _seqNovelStep3 == 0 {
		// 상위글 삭제 여부 검사
		isDeleted := false
		result = mdb.Model(schemas.NovelStep2{}).Select("deleted_yn").Where("seq_novel_step2 = ?", _seqNovelStep2).Scan(&isDeleted)
		if corm(result, &res) {
			return res
		}
		if isDeleted {
			res.ResultCode = define.DELETED_PARENT
			return res
		}

		seqKeyword := getSeqKeyword(2, int64(_seqNovelStep2))
		if isAbleKeyword(seqKeyword) != true {
			res.ResultCode = define.INACTIVE_KEYWORD
			return res
		}

		result = mdb.Model(schemas.NovelStep2{}).Where("seq_novel_step2 = ?", _seqNovelStep2).Count(&cnt)
		if corm(result, &res) {
			return res
		}
		if cnt == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		var seqNovelStep1 int64
		result = mdb.Model(schemas.NovelStep2{}).Where("seq_novel_step2 = ?", _seqNovelStep2).
			Pluck("seq_novel_step1", &seqNovelStep1)
		if corm(result, &res) {
			return res
		}

		if !_tempYn {
			result = mdb.Exec("UPDATE novel_step1 SET cnt_step3 = cnt_step3 + 1 WHERE seq_novel_step1 = ?", seqNovelStep1)
			result = mdb.Exec("UPDATE novel_step2 SET cnt_step3 = cnt_step3 + 1 WHERE seq_novel_step2 = ?", _seqNovelStep2)
		}

		novelStep3 := schemas.NovelStep3{
			SeqNovelStep1: seqNovelStep1,
			SeqNovelStep2: _seqNovelStep2,
			SeqMember:     userToken.SeqMember,
			Content:       _content,
			TempYn:        _tempYn,
			DeletedAt:     time.Now(),
		}
		result = mdb.Save(&novelStep3)
		if corm(result, &res) {
			return res
		}
		addKeywordCnt(seqKeyword)
		_seqNovelStep3 = novelStep3.SeqNovelStep3
	} else {
		seqKeyword := getSeqKeyword(3, int64(_seqNovelStep3))
		if isAbleKeyword(seqKeyword) != true {
			res.ResultCode = define.INACTIVE_KEYWORD
			return res
		}

		novelStep := schemas.NovelStep3{}
		result := mdb.Model(&novelStep).Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&novelStep)
		if corm(result, &res) {
			return res
		}
		if novelStep.SeqNovelStep3 == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		// 상위글 삭제 여부 검사
		_seqNovelStep2 = novelStep.SeqNovelStep2
		isDeleted := false
		result = mdb.Model(schemas.NovelStep2{}).Select("deleted_yn").Where("seq_novel_step2 = ?", _seqNovelStep2).Scan(&isDeleted)
		if corm(result, &res) {
			return res
		}
		if isDeleted {
			res.ResultCode = define.DELETED_PARENT
			return res
		}
		if novelStep.SeqMember != userToken.SeqMember {
			res.ResultCode = define.OTHER_USER
			return res
		}

		result = mdb.Model(&novelStep).
			Where("seq_novel_step3 = ?", _seqNovelStep3).
			Updates(map[string]interface{}{"content": _content, "temp_yn": _tempYn, "updated_at": time.Now()})
		if corm(result, &res) {
			return res
		}
	}

	var _seqNovelStep1 int64
	mdb.Model(schemas.NovelStep3{}).Select("seq_novel_step1").Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&_seqNovelStep1)
	if !_tempYn {
		go pushWriteTopic(userToken, 3, _seqNovelStep1)
	}

	return res
}

func NovelWriteStep4(req *domain.CommonRequest) domain.CommonResponse {

	var res = domain.CommonResponse{}
	userToken, err := define.ExtractTokenMetadata(req.JWToken, define.Mconn.JwtAccessSecret)
	if err != nil {
		res.ResultCode = define.INVALID_TOKEN
		res.ErrorDesc = err.Error()
		return res
	}

	_seqNovelStep3 := CpInt64(req.Parameters, "seq_novel_step3")
	_seqNovelStep4 := CpInt64(req.Parameters, "seq_novel_step4")
	_content := Cp(req.Parameters, "content")
	_tempYn := CpBool(req.Parameters, "temp_yn")

	// 블록처리된 유저 여부
	if isBlocked(userToken.SeqMember) {
		res.ResultCode = define.BLOCKED_ME
		return res
	}

	// 존재하는 닉네임 여부
	mdb := db.List[define.Mconn.DsnMaster]
	var cnt int64
	result := mdb.Model(schemas.MemberDetail{}).Where("seq_member = ?", userToken.SeqMember).Count(&cnt)
	if corm(result, &res) {
		return res
	}
	if cnt == 0 {
		res.ResultCode = define.NO_EXIST_NICK
		return res
	}

	if _seqNovelStep4 == 0 {
		// 상위글 삭제 여부 검사
		isDeleted := false
		result = mdb.Model(schemas.NovelStep3{}).Select("deleted_yn").Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&isDeleted)
		if corm(result, &res) {
			return res
		}
		if isDeleted {
			res.ResultCode = define.DELETED_PARENT
			return res
		}

		seqKeyword := getSeqKeyword(3, int64(_seqNovelStep3))
		if isAbleKeyword(seqKeyword) != true {
			res.ResultCode = define.INACTIVE_KEYWORD
			return res
		}

		result = mdb.Model(schemas.NovelStep3{}).Where("seq_novel_step3 = ?", _seqNovelStep3).Count(&cnt)
		if corm(result, &res) {
			return res
		}
		if cnt == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		novelStep3 := schemas.NovelStep3{}
		result = mdb.Model(schemas.NovelStep3{}).Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&novelStep3)
		if corm(result, &res) {
			return res
		}

		if !_tempYn {
			mdb.Exec("UPDATE novel_step1 SET cnt_step4 = cnt_step4 + 1 WHERE seq_novel_step1 = ?", novelStep3.SeqNovelStep1)
			mdb.Exec("UPDATE novel_step2 SET cnt_step4 = cnt_step4 + 1 WHERE seq_novel_step2 = ?", novelStep3.SeqNovelStep2)
			mdb.Exec("UPDATE novel_step3 SET cnt_step4 = cnt_step4 + 1 WHERE seq_novel_step3 = ?", _seqNovelStep3)
		}

		novelStep4 := schemas.NovelStep4{
			SeqNovelStep1: novelStep3.SeqNovelStep1,
			SeqNovelStep2: novelStep3.SeqNovelStep2,
			SeqNovelStep3: _seqNovelStep3,
			SeqMember:     userToken.SeqMember,
			Content:       _content,
			TempYn:        _tempYn,
			DeletedAt:     time.Now(),
		}
		result = mdb.Save(&novelStep4)
		if corm(result, &res) {
			return res
		}
		addKeywordCnt(seqKeyword)
		_seqNovelStep4 = novelStep4.SeqNovelStep4
	} else {
		seqKeyword := getSeqKeyword(4, int64(_seqNovelStep4))
		if isAbleKeyword(seqKeyword) != true {
			res.ResultCode = define.INACTIVE_KEYWORD
			return res
		}
		novelStep := schemas.NovelStep4{}
		result := mdb.Model(&novelStep).Where("seq_novel_step4 = ?", _seqNovelStep4).Scan(&novelStep)
		if corm(result, &res) {
			return res
		}
		if novelStep.SeqNovelStep4 == 0 {
			res.ResultCode = define.NO_EXIST_DATA
			return res
		}
		// 상위글 삭제 여부 검사
		_seqNovelStep3 = novelStep.SeqNovelStep3
		isDeleted := false
		result = mdb.Model(schemas.NovelStep3{}).Select("deleted_yn").Where("seq_novel_step3 = ?", _seqNovelStep3).Scan(&isDeleted)
		if corm(result, &res) {
			return res
		}
		if isDeleted {
			res.ResultCode = define.DELETED_PARENT
			return res
		}
		if novelStep.SeqMember != userToken.SeqMember {
			res.ResultCode = define.OTHER_USER
			return res
		}

		result = mdb.Model(&novelStep).
			Where("seq_novel_step4 = ?", _seqNovelStep4).
			Updates(map[string]interface{}{"content": _content, "temp_yn": _tempYn, "updated_at": time.Now()})
		if corm(result, &res) {
			return res
		}
	}

	var _seqNovelStep1 int64
	mdb.Model(schemas.NovelStep4{}).Select("seq_novel_step1").Where("seq_novel_step4 = ?", _seqNovelStep4).Scan(&_seqNovelStep1)
	if !_tempYn {
		go pushWriteTopic(userToken, 4, _seqNovelStep1)
	}

	return res
}

func pushWriteTopic(userToken *domain.UserToken, step int8, seqNovel int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	if seqNovel < 1 {
		return
	}

	isNight := tools.IsNight()

	// 1. 작가 정보 로딩
	userInfo := getUserInfo(userToken.SeqMember)

	// 2. 구독자 정보 로딩
	ldb := GetMyLogDbSlave(userToken.Allocated)
	var listMsSeq []int64
	ldb.Model(schemas.MemberSubscribe{}).
		Where("seq_member = ? AND status IN ('BOTH', 'FOLLOWER')", userInfo.SeqMember).
		Select("seq_member_opponent").
		Scan(&listMsSeq)
	mdb := db.List[define.Mconn.DsnMaster]
	alarm := schemas.Alarm{
		SeqMember:  0,
		Title:      "따옴",
		TypeAlarm:  5,
		ValueAlarm: int(seqNovel),
		Step:       step,
		Content:    userInfo.NickName + "님의 신규 소설이 등록되었습니다",
	}
	for _, v := range listMsSeq {
		alarm = schemas.Alarm{
			SeqMember:  v,
			Title:      "따옴",
			TypeAlarm:  5,
			ValueAlarm: int(seqNovel),
			Step:       step,
			Content:    userInfo.NickName + "님의 신규 소설이 등록되었습니다",
		}
		mdb.Create(&alarm)
	}

	listPush := []InfoPushTopic{}
	listFinalPush := []InfoPushTopic{}
	query := "SELECT seq_member, is_night_push FROM member_details WHERE seq_member IN (?) AND is_new_following = true"
	mdb.Raw(query, listMsSeq).Scan(&listPush)
	for _, o := range listPush {
		if isNight {
			if o.IsNightPush {
				listFinalPush = append(listFinalPush, o)
			}
		} else {
			listFinalPush = append(listFinalPush, o)
		}
	}
	for _, o := range listFinalPush {
		alarm.SeqMember = o.SeqMember
		go tools.SendPushMessageTopic(&alarm)
	}
}

func pushWrite(userToken *domain.UserToken, step int8, seqNovel int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	if seqNovel < 1 {
		return
	}

	isNight := tools.IsNight()

	// 1. 작가 정보 로딩
	userInfo := getUserInfo(userToken.SeqMember)

	// 2. 구독자 정보 로딩
	ldb := GetMyLogDbSlave(userToken.Allocated)
	var listMsSeq []int64
	ldb.Model(schemas.MemberSubscribe{}).
		Where("seq_member = ? AND status IN ('BOTH', 'FOLLOWER')", userInfo.SeqMember).
		Select("seq_member_opponent").
		Scan(&listMsSeq)
	mdb := db.List[define.Mconn.DsnMaster]
	alarm := schemas.Alarm{
		SeqMember:  0,
		Title:      "따옴",
		TypeAlarm:  5,
		ValueAlarm: int(seqNovel),
		Step:       step,
		Content:    userInfo.NickName + "님의 신규 소설이 등록되었습니다",
	}
	for _, v := range listMsSeq {
		alarm.SeqMember = v
		mdb.Create(&alarm)
	}

	listPush := []InfoPush{}
	listFinalPush := []InfoPush{}
	query := "SELECT mpt.push_token, mpt.seq_member, md.is_night_push FROM member_push_tokens mpt LEFT JOIN member_details md ON md.seq_member = mpt.seq_member WHERE md.seq_member IN (?) AND md.is_new_following = true"
	mdb.Raw(query, listMsSeq).Scan(&listPush)
	for _, o := range listPush {
		if isNight {
			if o.IsNightPush {
				listFinalPush = append(listFinalPush, o)
			}
		} else {
			listFinalPush = append(listFinalPush, o)
		}
	}
	for _, o := range listFinalPush {
		alarm.SeqMember = o.SeqMember
		go tools.SendPushMessage(o.PushToken, &alarm)
	}
}

type FollowerInfo struct {
	SeqMember   int64  `json:"seq_member"`
	PushToken   string `json:"push_token"`
	IsNightPush bool   `json:"is_night_push"`
}
