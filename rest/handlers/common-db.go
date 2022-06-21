package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/domain/schemas"
	"ddaom/memdb"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/appleboy/go-fcm"
	"gorm.io/gorm"
)

func GetMyLogDb(allocated int8) *gorm.DB {

	var myLogDb *gorm.DB
	switch allocated {
	case 1:
		myLogDb = db.List[define.DSN_LOG1]
	case 2:
		myLogDb = db.List[define.DSN_LOG2]
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

func getUserLogDb(_db *gorm.DB, seqMember int64) *gorm.DB {
	allocatedDb := 1
	_db.Model(schemas.Member{}).
		Select("allocated_db").
		Where("seq_member = ?", seqMember).Scan(&allocatedDb)

	fmt.Println("-------------------------")
	fmt.Println(allocatedDb)
	ldb := GetMyLogDb(int8(allocatedDb))
	return ldb
}

func getMyLike(userToken *domain.UserToken, step int8, seqNovel int64) bool {
	if userToken == nil {
		return false
	}
	ldb := GetMyLogDb(userToken.Allocated)
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
	ldb := GetMyLogDb(userToken.Allocated)
	status := "NONE"
	ldb.Model(schemas.MemberSubscribe{}).Select("status").Where("seq_member = ? AND seq_member_opponent = ?", userToken.SeqMember, seqMember).Scan(&status)
	return status
}

func getSeqKeyword(step int8, seqNovel int64) int64 {
	sdb := db.List[define.DSN_SLAVE]
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
	sdb := db.List[define.DSN_SLAVE]
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

func isAbleKeyword(seqKeyword int64) bool {
	sdb := db.List[define.DSN_SLAVE]
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
	sdb := db.List[define.DSN_SLAVE]
	image := schemas.Image{}
	sdb.Model(&image).Where("seq_image = ? AND active_yn = true", seqImage).Scan(&image)
	if image.SeqImage > 0 {
		return true
	} else {
		return false
	}
}
func isAbleColor(seqColor int64) bool {
	sdb := db.List[define.DSN_SLAVE]
	color := schemas.Color{}
	sdb.Model(&color).Where("seq_color = ? AND active_yn = true", seqColor).Scan(&color)
	if color.SeqColor > 0 {
		return true
	} else {
		return false
	}
}

func getUserInfo(seqMember int64) schemas.MemberDetail {
	sdb := db.List[define.DSN_SLAVE]
	md := schemas.MemberDetail{}
	sdb.Model(&md).Where("seq_member = ?", seqMember).Scan(&md)
	return md
}

func isMineByStepNovelSeq(step int, seqNovel int64, seqMember int64) bool {
	sdb := db.List[define.DSN_SLAVE]
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
	sdb := db.List[define.DSN_SLAVE]
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
	mdb := db.List[define.DSN_MASTER]
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
	client, err := fcm.NewClient(define.PUSH_SERVER_KEY)
	if err != nil {
		log.Fatalln(err)
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
	mdb := db.List[define.DSN_MASTER]

	var totalCnt int
	mdb.Model(schemas.NovelStep1{}).
		Select("COUNT(*) + SUM(cnt_step2) + SUM(cnt_step3) + SUM(cnt_step4)").
		Where("seq_keyword = ?", seqKeyword).Scan(&totalCnt)
	mdb.Model(schemas.Keyword{}).
		Where("seq_keyword = ?", seqKeyword).
		Update("cnt_total", totalCnt)

	// redis update (step1 novel count)
	memdb.Zadd("CACHES:ASSET:COUNT", totalCnt, seqKeyword)

}

// novel-save (완)
func cacheMainLive(seqKeyword int64) {
	sdb := db.List[define.DSN_SLAVE]
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
	sdb := db.List[define.DSN_SLAVE]
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
func cacheMainPopularWriter() {
	sdb := db.List[define.DSN_SLAVE]
	listPopularWriter := []ListPopularWriter{}
	query := `
		SELECT
			md.seq_member, md.nick_name, md.profile_photo
		FROM
			member_details md INNER JOIN members m ON md.seq_member = m.seq_member
		WHERE
			m.deleted_yn = false
		ORDER BY
			md.cnt_like DESC, md.cnt_subscribe DESC
		LIMIT 10`
	sdb.Raw(query).Scan(&listPopularWriter)
	j, _ := json.Marshal(listPopularWriter)
	memdb.Set("CACHES:MAIN:LIST_POPULAR_WRITER", string(j))
}
