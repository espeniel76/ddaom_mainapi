package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/memdb"
	"ddaom/tools"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/appleboy/go-fcm"
	"gorm.io/gorm"
)

func GetMyLogDbMaster(allocated int8) *gorm.DB {

	var myLogDb *gorm.DB
	switch allocated {
	case 1:
		myLogDb = db.List[define.Mconn.DsnLog1Master]
	case 2:
		myLogDb = db.List[define.Mconn.DsnLog2Master]
	}
	return myLogDb
}

func GetMyLogDbSlave(allocated int8) *gorm.DB {
	// fmt.Println("allocated: ", allocated)
	var myLogDb *gorm.DB
	switch allocated {
	case 1:
		myLogDb = db.List[define.Mconn.DsnLog1Slave]
	case 2:
		myLogDb = db.List[define.Mconn.DsnLog2Slave]
	}
	return myLogDb
}

func corm(o *gorm.DB, res *domain.CommonResponse) bool {
	isError := false
	if o.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = o.Error.Error()
		isError = true
	}
	return isError
}

func getUserLogDbMaster(_db *gorm.DB, seqMember int64) *gorm.DB {
	allocatedDb := 1
	_db.Model(schemas.Member{}).
		Select("allocated_db").
		Where("seq_member = ?", seqMember).Scan(&allocatedDb)
	ldb := GetMyLogDbMaster(int8(allocatedDb))
	return ldb
}

func getUserLogDbSlave(_db *gorm.DB, seqMember int64) *gorm.DB {
	allocatedDb := 1
	_db.Model(schemas.Member{}).
		Select("allocated_db").
		Where("seq_member = ?", seqMember).Scan(&allocatedDb)
	ldb := GetMyLogDbSlave(int8(allocatedDb))
	return ldb
}

func getMyLike(userToken *domain.UserToken, step int8, seqNovel int64) bool {
	if userToken == nil {
		return false
	}
	ldb := GetMyLogDbSlave(userToken.Allocated)
	likeYn := false
	switch step {
	case 1:
		ldb.Model(schemas.MemberLikeStep1{}).Select("like_yn").Where("seq_member = ? AND seq_novel_step1 = ?", userToken.SeqMember, seqNovel).Scan(&likeYn)
		return likeYn
	case 2:
		ldb.Model(schemas.MemberLikeStep2{}).Select("like_yn").Where("seq_member = ? AND seq_novel_step2 = ?", userToken.SeqMember, seqNovel).Scan(&likeYn)
		return likeYn
	case 3:
		ldb.Model(schemas.MemberLikeStep3{}).Select("like_yn").Where("seq_member = ? AND seq_novel_step3 = ?", userToken.SeqMember, seqNovel).Scan(&likeYn)
		return likeYn
	case 4:
		ldb.Model(schemas.MemberLikeStep4{}).Select("like_yn").Where("seq_member = ? AND seq_novel_step4 = ?", userToken.SeqMember, seqNovel).Scan(&likeYn)
		return likeYn
	}
	return false
}

func getMySubscribe(userToken *domain.UserToken, seqMember int64) string {
	if userToken == nil {
		return "NONE"
	}
	ldb := GetMyLogDbSlave(userToken.Allocated)
	status := "NONE"
	ldb.Model(schemas.MemberSubscribe{}).Select("status").Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, seqMember).Scan(&status)
	return status
}

func getSeqKeyword(step int8, seqNovel int64) int64 {
	sdb := db.List[define.Mconn.DsnSlave]
	var seq int64
	switch step {
	case 1:
		sdb.Raw("SELECT seq_keyword FROM novel_step1 WHERE seq_novel_step1 = ?", seqNovel).Scan(&seq)
	case 2:
		sdb.Raw("SELECT ns1.seq_keyword FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 WHERE ns2.seq_novel_step2 = ?", seqNovel).Scan(&seq)
	case 3:
		sdb.Raw("SELECT ns1.seq_keyword FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 WHERE ns3.seq_novel_step3 = ?", seqNovel).Scan(&seq)
	case 4:
		sdb.Raw("SELECT ns1.seq_keyword FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 WHERE ns4.seq_novel_step4 = ?", seqNovel).Scan(&seq)
	}
	return seq
}

func getNovel(step int8, seqNovel int64) GetNovelRes {
	sdb := db.List[define.Mconn.DsnSlave]
	m := GetNovelRes{}
	sql := ""
	switch step {
	case 1:
		sql = `SELECT ns1.title, ns1.seq_member, m.push_token, md.is_night_push
		FROM novel_step1 ns1 INNER JOIN members m ON ns1.seq_member = m.seq_member INNER JOIN member_details md ON ns1.seq_member = md.seq_member
		WHERE seq_novel_step1 = ? AND md.is_liked = true`
	case 2:
		sql = `SELECT ns1.title, ns2.seq_member, m.push_token, md.is_night_push
		FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 INNER JOIN members m ON ns2.seq_member = m.seq_member INNER JOIN member_details md ON ns2.seq_member = md.seq_member
		WHERE ns2.seq_novel_step2 = ? AND md.is_liked = true`
	case 3:
		sql = `SELECT ns1.title, ns3.seq_member, m.push_token, md.is_night_push
		FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 INNER JOIN members m ON ns3.seq_member = m.seq_member INNER JOIN member_details md ON ns3.seq_member = md.seq_member
		WHERE ns3.seq_novel_step3 = ? AND md.is_liked = true`
	case 4:
		sql = `SELECT ns1.title, ns4.seq_member, m.push_token, md.is_night_push
		FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 INNER JOIN members m ON ns4.seq_member = m.seq_member INNER JOIN member_details md ON ns4.seq_member = md.seq_member
		WHERE ns4.seq_novel_step4 = ? AND md.is_liked = true`
	}
	sdb.Raw(sql, seqNovel).Scan(&m)
	return m
}

type GetNovelRes struct {
	Title       string `json:"title"`
	SeqMember   int64  `json:"seq_member"`
	PushToken   string `json:"push_token"`
	IsNightPush bool   `json:"is_night_push"`
}

func isBlocked(seqMember int64) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	isBlocked := false
	sdb.Model(schemas.Member{}).Select("blocked_yn").Where("seq_member = ?", seqMember).Scan(&isBlocked)
	return isBlocked
}

func isAbleKeyword(seqKeyword int64) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	keyword := schemas.Keyword{}
	sdb.Model(&keyword).
		Where("seq_keyword = ? AND active_yn = true AND NOW() BETWEEN start_date AND end_date", seqKeyword).
		Scan(&keyword)
	if keyword.SeqKeyword > 0 {
		return true
	} else {
		return false
	}
}
func isAbleImage(seqImage int64) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	image := schemas.Image{}
	sdb.Model(&image).Where("seq_image = ? AND active_yn = true", seqImage).Scan(&image)
	if image.SeqImage > 0 {
		return true
	} else {
		return false
	}
}
func isAbleColor(seqColor int64) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	color := schemas.Color{}
	sdb.Model(&color).Where("seq_color = ? AND active_yn = true", seqColor).Scan(&color)
	if color.SeqColor > 0 {
		return true
	} else {
		return false
	}
}
func isAbleGenre(seqGenre int64) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	genre := schemas.Genre{}
	sdb.Model(&genre).Where("seq_genre = ? AND active_yn = true", seqGenre).Scan(&genre)
	if genre.SeqGenre > 0 {
		return true
	} else {
		return false
	}
}

func isExistTitle(title string) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	var cnt int64
	sdb.Model(&schemas.NovelStep1{}).Where("title = ? AND temp_yn = false AND deleted_yn = false", title).Count(&cnt)
	if cnt > 0 {
		return true
	} else {
		return false
	}
}

func getUserInfo(seqMember int64) schemas.MemberDetail {
	sdb := db.List[define.Mconn.DsnSlave]
	md := schemas.MemberDetail{}
	sdb.Model(&md).Where("seq_member = ?", seqMember).Scan(&md)
	return md
}

func isMineByStepNovelSeq(step int, seqNovel int64, seqMember int64) bool {
	sdb := db.List[define.Mconn.DsnSlave]
	var _seqMember int64
	switch step {
	case 1:
		sdb.Model(schemas.NovelStep1{}).Select("seq_member").Where("seq_novel_step1 = ?", seqNovel).Scan(&_seqMember)
	case 2:
		sdb.Model(schemas.NovelStep2{}).Select("seq_member").Where("seq_novel_step2 = ?", seqNovel).Scan(&_seqMember)
	case 3:
		sdb.Model(schemas.NovelStep3{}).Select("seq_member").Where("seq_novel_step3 = ?", seqNovel).Scan(&_seqMember)
	case 4:
		sdb.Model(schemas.NovelStep4{}).Select("seq_member").Where("seq_novel_step4 = ?", seqNovel).Scan(&_seqMember)
	}

	if _seqMember == seqMember {
		return true
	} else {
		return false
	}
}

func getUserInfoPush(seqMember int64) GetUserInfoPushRes {
	sdb := db.List[define.Mconn.DsnSlave]
	res := GetUserInfoPushRes{}
	sql := `SELECT
		m.seq_member,
		m.push_token,
		md.is_liked,
		md.is_finished,
		md.is_new_follower,
		md.is_new_following,
		md.is_night_push
	FROM members m INNER JOIN member_details md ON m.seq_member = md.seq_member WHERE m.seq_member = ?`
	sdb.Raw(sql, seqMember).Scan(&res)
	return res
}

type GetUserInfoPushRes struct {
	SeqMember      int64  `json:"seq_member"`
	PushToken      string `json:"push_token"`
	IsLiked        bool   `json:"is_liked"`
	IsFinished     bool   `json:"is_finished"`
	IsNewFollower  bool   `json:"is_new_follower"`
	IsNewFollowing bool   `json:"is_new_following"`
	IsNightPush    bool   `json:"is_night_push"`
}

func sendPush(pushToken string, alarm *schemas.Alarm) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	mdb := db.List[define.Mconn.DsnMaster]
	mdb.Create(&alarm)
	msg := &fcm.Message{
		To: pushToken,
		Data: map[string]interface{}{
			"seq_alarm":   alarm.SeqAlarm,
			"type_alarm":  alarm.TypeAlarm,
			"value_alarm": alarm.ValueAlarm,
			"step":        alarm.Step,
		},
		Notification: &fcm.Notification{
			Title: alarm.Title,
			Body:  alarm.Content,
		},
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient(define.Mconn.PushServerKey)
	if err != nil {
		// log.Fatalln(err)
		fmt.Println(err)
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		// log.Fatalln(err)
		fmt.Println(err)
	}

	log.Printf("%#v\n", response)
}

func addKeywordCnt(seqKeyword int64) {
	mdb := db.List[define.Mconn.DsnMaster]

	var totalCnt int
	mdb.Model(schemas.NovelStep1{}).
		Select("COUNT(*) + SUM(cnt_step2) + SUM(cnt_step3) + SUM(cnt_step4)").
		Where("seq_keyword = ?", seqKeyword).
		Where("deleted_yn = false AND temp_yn = false AND active_yn = true").
		Scan(&totalCnt)
	mdb.Model(schemas.Keyword{}).
		Where("seq_keyword = ?", seqKeyword).
		Update("cnt_total", totalCnt)

	// redis update (step1 novel count)
	memdb.Zadd("CACHES:ASSET:COUNT", totalCnt, seqKeyword)

}

// novel-save (완)
func cacheMainLive(seqKeyword int64) {
	sdb := db.List[define.Mconn.DsnSlave]
	listLive := []ListLive{}
	query := `
	SELECT
		seq_novel_step1,
		seq_image,
		seq_color,
		title
	FROM novel_step1
	WHERE seq_keyword = ? AND active_yn = true AND temp_yn = false AND deleted_yn = false
	ORDER BY created_at DESC
	LIMIT 10
	`
	sdb.Raw(query, seqKeyword).Scan(&listLive)
	j, _ := json.Marshal(listLive)
	memdb.Set("CACHES:MAIN:LIST_LIVE:"+strconv.FormatInt(seqKeyword, 10), string(j))
}

// novel-view-finish (view cnt, 완)
// complete batch (like cnt, 완)
func cacheMainPopular() {
	sdb := db.List[define.Mconn.DsnSlave]
	listPopular := []ListPopular{}
	query := `
	SELECT
		seq_novel_finish,
		seq_image,
		seq_color,
		title
	FROM novel_finishes
	WHERE active_yn = true
	ORDER BY cnt_like DESC, cnt_view DESC
	LIMIT 10
	`
	sdb.Raw(query).Scan(&listPopular)
	j, _ := json.Marshal(listPopular)
	memdb.Set("CACHES:MAIN:LIST_POPULAR", string(j))
}

// novel-subscribe (subscribe cnt)
// like step1~4 (like cnt)
/*
*
인기작가 순위 기준 수정
- 받은 구독 수 + 북마크 수 + 연재중 좋아요 수
(연재중일때는 반영하고, 주제어 종료되면 다시 삭제되는 방식?으로 하여 계속 변동이 있게? or 주제어 종료되어도 유지)
*/
func cacheMainPopularWriter() {
	sdb := db.List[define.Mconn.DsnSlave]
	listPopularWriter := []ListPopularWriter{}
	query := `
		SELECT
			md.seq_member, md.nick_name, md.profile_photo, md.cnt_subscribe + md.cnt_bookmark AS cnt_subscribe_bookmark
		FROM
			member_details md INNER JOIN members m ON md.seq_member = m.seq_member
		WHERE
			m.deleted_yn = false AND md.cnt_subscribe > 0 OR md.cnt_bookmark > 0
		ORDER BY cnt_subscribe_bookmark DESC
		LIMIT 10`
	sdb.Raw(query).Scan(&listPopularWriter)
	j, _ := json.Marshal(listPopularWriter)
	memdb.Set("CACHES:MAIN:LIST_POPULAR_WRITER", string(j))
}

func cacheMainPopularWriterLike() {
	sdb := db.List[define.Mconn.DsnSlave]
	listPopularWriterLike := []ListPopularWriterLIke{}
	query := `
	SELECT
		A.seq_keyword, A.seq_member, A.nick_name, A.profile_photo, SUM(A.cnt) AS cnt
	FROM
	(	
	SELECT k.seq_keyword, ns1.seq_member, md.nick_name, md.profile_photo, SUM(ns1.cnt_like) AS cnt
		FROM novel_step1 ns1 INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
		INNER JOIN member_details md ON ns1.seq_member = md.seq_member
		WHERE NOW() BETWEEN k.start_date AND k.end_date AND ns1.cnt_like > 0
		GROUP BY k.seq_keyword, ns1.seq_member, md.nick_name, md.profile_photo
		UNION ALL
		SELECT k.seq_keyword, ns2.seq_member, md.nick_name, md.profile_photo, SUM(ns2.cnt_like) AS cnt
		FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns2.seq_novel_step1 = ns1.seq_novel_step1
		INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
		INNER JOIN member_details md ON ns2.seq_member = md.seq_member
		WHERE NOW() BETWEEN k.start_date AND k.end_date AND ns2.cnt_like > 0
		GROUP BY k.seq_keyword, ns2.seq_member, md.nick_name, md.profile_photo
		UNION ALL
		SELECT k.seq_keyword, ns3.seq_member, md.nick_name, md.profile_photo, SUM(ns3.cnt_like) AS cnt
		FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns3.seq_novel_step1 = ns1.seq_novel_step1
		INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
		INNER JOIN member_details md ON ns3.seq_member = md.seq_member
		WHERE NOW() BETWEEN k.start_date AND k.end_date AND ns3.cnt_like > 0
		GROUP BY k.seq_keyword, ns3.seq_member, md.nick_name, md.profile_photo
		UNION ALL
		SELECT k.seq_keyword, ns4.seq_member, md.nick_name, md.profile_photo, SUM(ns4.cnt_like) AS cnt
		FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns4.seq_novel_step1 = ns1.seq_novel_step1
		INNER JOIN keywords k ON ns1.seq_keyword = k.seq_keyword
		INNER JOIN member_details md ON ns4.seq_member = md.seq_member
		WHERE NOW() BETWEEN k.start_date AND k.end_date AND ns4.cnt_like > 0
		GROUP BY k.seq_keyword, ns4.seq_member, md.nick_name, md.profile_photo
	) AS A
	GROUP BY A.seq_keyword, A.seq_member, A.nick_name, A.profile_photo
	`
	sdb.Raw(query).Scan(&listPopularWriterLike)
	j, _ := json.Marshal(listPopularWriterLike)
	memdb.Set("CACHES:MAIN:LIST_POPULAR_WRITER_LIKE", string(j))
}

func cacheMyBlockUser(userToken *domain.UserToken) {
	ldb := GetMyLogDbSlave(userToken.Allocated)
	var seqs []int64
	ldb.
		Model(schemas.MemberBlocking{}).
		Select("seq_member_to").
		Where("seq_member = ? AND block_yn = true", userToken.SeqMember).
		Scan(&seqs)
	j, _ := json.Marshal(seqs)
	memdb.Set("CACHES:USERS:BLOCK:"+strconv.FormatInt(userToken.SeqMember, 10), string(j))
}

func educeImage(seqColor int64, seqImage int64, seqNovelStep1 int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	s3 := tools.S3Info{
		AwsProfileName: "ddaom",
		AwsS3Region:    define.Mconn.AwsS3Region,
		AwsSecretKey:   define.Mconn.AwsSecretKey,
		AwsAccessKey:   define.Mconn.AwsAccessKey,
		BucketName:     define.Mconn.AwsBucketName,
	}

	imageName := strconv.Itoa(int(seqColor)) + "_" + strconv.Itoa(int(seqImage)) + ".jpg"
	savePath := define.Mconn.ReplacePath + "/thumb/"

	mdb := db.List[define.Mconn.DsnMaster]

	fmt.Println("공유이미지 만들기", imageName)

	// 1. 해당 조합의 DB 데이터가 있는지 확인
	var cnt int64
	mdb.Model(schemas.NovelStep1{}).Select("COUNT(*)").Where("endure_image = ?", imageName).Scan(&cnt)
	if cnt > 0 {
		return
	}

	// 2. 없으면, DB 에서 경로 가져옴
	var imgPath string
	var hexValue string
	mdb.Model(schemas.Image{}).Select("image").Where("seq_image = ?", seqImage).Scan(&imgPath)
	imgSrc := define.Mconn.ReplacePath + imgPath
	mdb.Model(schemas.Color{}).Select("color").Where("seq_color = ?", seqColor).Scan(&hexValue)

	// 3. 가져온 경로로 이미지 다운 (AWS 일 시)
	if define.Mconn.HTTPServer == "https://s3.ap-northeast-2.amazonaws.com/image.ttaom.com" {
		err := s3.SetS3ConfigByKey()
		if err != nil {
			return
		}

		s3.DownloadFile("/tmp/thumb", strings.Replace(imgPath, "/", "", 1))
		fileNames := strings.Split(imgPath, "/")
		fileName := fileNames[len(fileNames)-1]
		imgSrc = "/tmp/thumb/" + fileName
	}

	// 4. MERGE 작업
	imgSource, _ := os.Open(imgSrc)
	imgLayer, _ := png.Decode(imgSource)
	defer imgSource.Close()

	b := imgLayer.Bounds()
	imgResult := image.NewRGBA(b)

	m := image.NewRGBA(b)
	colorRgba, _ := parseHexColor(hexValue)
	draw.Draw(m, m.Bounds(), &image.Uniform{colorRgba}, image.Point{}, draw.Src)

	draw.Draw(imgResult, b, m, image.Point{}, draw.Src)
	draw.Draw(imgResult, b, imgLayer, image.Point{}, draw.Over)

	resultPath := savePath + imageName
	// fmt.Println(resultPath)
	third, err := os.Create(resultPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	jpeg.Encode(third, imgResult, &jpeg.Options{90})
	defer third.Close()

	// 5. DB save
	mdb.Model(schemas.NovelStep1{}).Where("seq_novel_step1 = ?", seqNovelStep1).Update("endure_image", imageName)

	// 6. MERGE 한 파일 업로드 (AWS 일 시)
	if define.Mconn.HTTPServer == "https://s3.ap-northeast-2.amazonaws.com/image.ttaom.com" {
		// fmt.Println("/tmp/thumb/" + imageName)
		s3.UploadFileByFileName("/tmp/thumb/"+imageName, "thumb/"+imageName, "image/jpeg")
	}
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

func setUserActionLog(_seqMember int64, _type int8, _contents string) {
	mdb := db.List[define.Mconn.DsnMaster]
	mdb.Create(&schemas.MemberLog{
		SeqMember: _seqMember,
		Type:      _type,
		Contents:  _contents,
	})
}

// 사용자 시퀀스로 나와 사용자의 차단을 리턴 받는다.
func getBlockMember(allocatedDb int8, seqMember int64, memberTo int64) MemberBlock {
	sdb := GetMyLogDbSlave(allocatedDb)
	memberBlock := MemberBlock{}

	// // 1. 나의 로그에서 차단 유저 가져오기
	sdb.Model(schemas.MemberBlocking{}).
		Select("seq_member_to AS seq_member, block_yn").
		Where("seq_member = ?", seqMember).
		Where("seq_member_to = ?", memberTo).
		Where("block_yn = true").
		Scan(&memberBlock)
	return memberBlock
}

func isBlockMember(allocatedDb int8, seqMember int64, memberTo int64) bool {
	sdb := GetMyLogDbSlave(allocatedDb)
	isBlock := false

	// // 1. 나의 로그에서 차단 유저 가져오기
	sdb.Model(schemas.MemberBlocking{}).
		Select("block_yn").
		Where("seq_member = ?", seqMember).
		Where("seq_member_to = ?", memberTo).
		Where("block_yn = true").
		Scan(&isBlock)
	return isBlock
}

// 사용자 시퀀스로 나와 사용자의 차단 목록을 리턴 받는다.
func getBlockMemberList(allocatedDb int8, seqMember int64, listMemberTo []int64) []MemberBlock {
	keys := make(map[int64]bool)
	var list []int64
	for _, value := range listMemberTo {
		if _, saveValue := keys[value]; !saveValue {
			keys[value] = true
			list = append(list, value)
		}
	}
	sdb := GetMyLogDbSlave(allocatedDb)
	listMemberBlock := []MemberBlock{}

	// // 1. 나의 로그에서 차단 유저 가져오기
	sdb.Model(schemas.MemberBlocking{}).
		Select("seq_member_to AS seq_member, block_yn").
		Where("seq_member = ?", seqMember).
		Where("seq_member_to IN (?)", list).
		Where("block_yn = true").
		Scan(&listMemberBlock)
	return listMemberBlock
}

type MemberBlock struct {
	SeqMember int64 `json:"seq_member"`
	BlockYn   bool  `json:"block_yn"`
}
