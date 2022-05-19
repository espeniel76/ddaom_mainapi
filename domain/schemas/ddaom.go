package schemas

import (
	"time"
)

type Member struct {
	SeqMember       int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member"`
	Email           string    `gorm:"type:varchar(255);unique" json:"email"`
	Token           string    `gorm:"type:varchar(5120)" json:"token"`
	ProfileImageUrl string    `gorm:"type:varchar(512)" json:"profile_image_url"`
	SnsType         string    `gorm:"type:ENUM('KAKAO','NAVER','FACEBOOK','GOOGLE','APPLE'); DEFAULT:'GOOGLE'" json:"sns_type"`
	ActiveYn        bool      `gorm:"default:false" json:"active_yn"`
	UserLevel       int8      `gorm:"default:5" json:"user_level"`
	AllocatedDb     int8      `json:"allocted_db"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushToken       string    `gorm:"type:varchar(5120)" json:"push_token"`
}

type MemberDetail struct {
	SeqMemberDetail  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_detail"`
	SeqMember        int64     `gorm:"unique" json:"seq_member"`
	Email            string    `gorm:"type:varchar(1024)" json:"email"`
	Name             string    `gorm:"type:varchar(50)" json:"name"`
	NickName         string    `gorm:"type:varchar(50);unique;column:nick_name" json:"nick_name"`
	ProfilePhoto     string    `gorm:"type:varchar(1024)" json:"profile_photo"`
	Tel              string    `gorm:"type:varchar(50)" json:"tel"`
	MobileCompany    int8      `gorm:"default:0" json:"mobile_company"`
	Mobile           string    `gorm:"type:varchar(50)" json:"mobile"`
	Address          string    `gorm:"type:varchar(1024)" json:"address"`
	AddressDetail    string    `gorm:"type:varchar(1024)" json:"address_detail"`
	Zipcode          string    `gorm:"type:varchar(50)" json:"zipcode"`
	AuthenticationCi string    `gorm:"type:varchar(255)" json:"authentication_ci"`
	AuthenticationAt time.Time `json:"authentication_at"`
	IsNewKeyword     bool      `gorm:"default:false" json:"is_new_keyword"`
	IsLiked          bool      `gorm:"default:false" json:"is_liked"`
	IsFinished       bool      `gorm:"default:false" json:"is_finished"`
	IsNewFollower    bool      `gorm:"default:false" json:"is_new_follower"`
	IsNewFollowing   bool      `gorm:"default:false" json:"is_new_following"`
	IsNightPush      bool      `gorm:"default:false" json:"is_night_push"`
	CntSubscribe     int64     `gorm:"default:0" json:"cnt_subscribe"`
	CntLike          int64     `gorm:"default:0" json:"cnt_like"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type MemberAdmin struct {
	SeqMemberAdmin int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_admin"`
	UserId         string    `gorm:"type:varchar(50);unique" json:"user_id"`
	Password       string    `gorm:"type:varchar(1024)" json:"password"`
	ActiveYn       bool      `gorm:"default:false" json:"active_yn"`
	Name           string    `gorm:"type:varchar(50)" json:"name"`
	NickName       string    `gorm:"type:varchar(50);unique;column:nick_name" json:"nick_name"`
	UserLevel      int8      `gorm:"default:5" json:"user_level"`
	ProfileImage   string    `gorm:"type:varchar(1024)" json:"profile_image"`
	CreatedAt      time.Time `json:"created_at"`
	Creator        string    `gorm:"type:varchar(50)" json:"creator"`
	UpdatedAt      time.Time `json:"updated_at"`
	Updator        string    `gorm:"type:varchar(50)" json:"updator"`
}

type MemberAdminLoginLog struct {
	SeqMemberAdminLoginLog int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_admin_login_log"`
	SeqMemberAdmin         int64     `gorm:"index" json:"seq_member_admin"`
	Token                  string    `gorm:"type:varchar(1024)" json:"token"`
	LoginAt                time.Time `json:"login_at"`
}

// 주저에
type KeywordAlarmLog struct {
	SeqKeywordAlarmLog int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword_alarm_log"`
	SeqKeyword         string    `gorm:"index" json:"seq_keyword"`
	CreatedAt          time.Time `json:"created_at"`
	CntPush            int64     `gorm:"default:0" json:"cnt_push"`
}

type Keyword struct {
	SeqKeyword int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword"`
	Keyword    string    `gorm:"unique;type:varchar(1024)" json:"keyword"`
	ActiveYn   bool      `gorm:"default:false" json:"active_yn"`
	FinishYn   bool      `gorm:"default:false" json:"finish_yn"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CntTotal   int64     `gorm:"default:0" json:"cnt_total"`
	CntLike    int64     `gorm:"default:0" json:"cnt_like"`
	CntFinish  int64     `gorm:"default:0" json:"cnt_finish"`
	CreatedAt  time.Time `json:"created_at"`
	Creator    string    `gorm:"type:varchar(50)" json:"creator"`
	UpdatedAt  time.Time `json:"updated_at"`
	Updator    string    `gorm:"type:varchar(50)" json:"updator"`
	FinishedAt bool      `json:"finished_at"`
}

type KeywordToday struct {
	SeqKeywordToday int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword_today"`
	SeqKeyword      int64     `gorm:"index" json:"seq_keyword"`
	ViewStartDate   string    `gorm:"type:char(8)" json:"view_start_date"`
	ViewEndDate     string    `gorm:"type:char(8)" json:"view_end_date"`
	ActiveYn        bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// 장르
type Genre struct {
	SeqGenre  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_genre"`
	Genre     string    `gorm:"unique;type:varchar(50)" json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ActiveYn  bool      `gorm:"default:false" json:"active_yn"`
	Creator   string    `gorm:"type:varchar(50)" json:"creator"`
	Updator   string    `gorm:"type:varchar(50)" json:"updator"`
}

type Image struct {
	SeqImage  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_image"`
	Name      string    `gorm:"unique;type:varchar(50)" json:"name"`
	Image     string    `gorm:"type:varchar(1024)" json:"image"`
	ActiveYn  bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt time.Time `json:"created_at"`
	Creator   string    `gorm:"type:varchar(50)" json:"creator"`
	UpdatedAt time.Time `json:"updated_at"`
	Updator   string    `gorm:"type:varchar(50)" json:"updator"`
}

type Slang struct {
	SeqSlang       int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_slang"`
	Slang          string    `gorm:"unique;type:varchar(50)" json:"slang"`
	ActiveYn       bool      `gorm:"default:false" json:"active_yn"`
	SeqMemberAdmin string    `gorm:"index" json:"seq_member_admin"`
	CreatedAt      time.Time `json:"created_at"`
	Creator        string    `gorm:"type:varchar(50)" json:"creator"`
	UpdatedAt      time.Time `json:"updated_at"`
	Updator        string    `gorm:"type:varchar(50)" json:"updator"`
}

type Color struct {
	SeqColor  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_color"`
	Name      string    `gorm:"unique;type:varchar(50)" json:"name"`
	Color     string    `gorm:"type:varchar(12)" json:"color"`
	ActiveYn  bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt time.Time `json:"created_at"`
	Creator   string    `gorm:"type:varchar(50)" json:"creator"`
	UpdatedAt time.Time `json:"updated_at"`
	Updator   string    `gorm:"type:varchar(50)" json:"updator"`
}

type NovelStep1 struct {
	SeqNovelStep1 int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_step1"`
	SeqKeyword    int64     `gorm:"index" json:"seq_keyword"`
	SeqImage      int64     `gorm:"index" json:"seq_image"`
	SeqColor      int64     `gorm:"index" json:"seq_color"`
	SeqGenre      int64     `gorm:"index" json:"seq_genre"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	Title         string    `gorm:"unique;type:varchar(1024)" json:"title"`
	Content       string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike       int64     `gorm:"default:0" json:"cnt_like"`
	CntView       int64     `gorm:"default:0" json:"cnt_view"`
	CntStep2      int64     `gorm:"default:0" json:"cnt_step2"`
	CntStep3      int64     `gorm:"default:0" json:"cnt_step3"`
	CntStep4      int64     `gorm:"default:0" json:"cnt_step4"`
	ActiveYn      bool      `gorm:"default:true" json:"active_yn"`
	TempYn        bool      `gorm:"default:false" json:"temp_yn"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NovelStep2 struct {
	SeqNovelStep2 int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_step2"`
	SeqNovelStep1 int64     `gorm:"index" json:"seq_novel_step1"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	Content       string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike       int64     `gorm:"default:0" json:"cnt_like"`
	CntView       int64     `gorm:"default:0" json:"cnt_view"`
	CntStep3      int64     `gorm:"default:0" json:"cnt_step3"`
	CntStep4      int64     `gorm:"default:0" json:"cnt_step4"`
	ActiveYn      bool      `gorm:"default:true" json:"active_yn"`
	TempYn        bool      `gorm:"default:false" json:"temp_yn"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NovelStep3 struct {
	SeqNovelStep3 int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_step3"`
	SeqNovelStep1 int64     `gorm:"index" json:"seq_novel_step1"`
	SeqNovelStep2 int64     `gorm:"index" json:"seq_novel_step2"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	Content       string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike       int64     `gorm:"default:0" json:"cnt_like"`
	CntView       int64     `gorm:"default:0" json:"cnt_view"`
	CntStep4      int64     `gorm:"default:0" json:"cnt_step4"`
	ActiveYn      bool      `gorm:"default:true" json:"active_yn"`
	TempYn        bool      `gorm:"default:false" json:"temp_yn"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NovelStep4 struct {
	SeqNovelStep4 int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_step4"`
	SeqNovelStep1 int64     `gorm:"index" json:"seq_novel_step1"`
	SeqNovelStep2 int64     `gorm:"index" json:"seq_novel_step2"`
	SeqNovelStep3 int64     `gorm:"index" json:"seq_novel_step3"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	Content       string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike       int       `gorm:"default:0" json:"cnt_like"`
	CntView       int64     `gorm:"default:0" json:"cnt_view"`
	ActiveYn      bool      `gorm:"default:true" json:"active_yn"`
	TempYn        bool      `gorm:"default:false" json:"temp_yn"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NovelFinish struct {
	SeqNovelFinish int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_finish"`
	SeqKeyword     int64     `gorm:"index" json:"seq_keyword"`
	SeqImage       int64     `gorm:"index" json:"seq_image"`
	SeqColor       int64     `gorm:"index" json:"seq_color"`
	SeqGenre       int64     `gorm:"index" json:"seq_genre"`
	SeqMemberStep1 int64     `gorm:"default:0" json:"seq_member_step1"`
	SeqMemberStep2 int64     `gorm:"default:0" json:"seq_member_step2"`
	SeqMemberStep3 int64     `gorm:"default:0" json:"seq_member_step3"`
	SeqMemberStep4 int64     `gorm:"default:0" json:"seq_member_step4"`
	NickNameStep1  string    `gorm:"type:varchar(50)" json:"nick_name_step1"`
	NickNameStep2  string    `gorm:"type:varchar(50)" json:"nick_name_step2"`
	NickNameStep3  string    `gorm:"type:varchar(50)" json:"nick_name_step3"`
	NickNameStep4  string    `gorm:"type:varchar(50)" json:"nick_name_step4"`
	Title          string    `gorm:"unique;type:varchar(1024)" json:"title"`
	Content1       string    `gorm:"type:varchar(5120)" json:"content1"`
	Content2       string    `gorm:"type:varchar(5120)" json:"content2"`
	Content3       string    `gorm:"type:varchar(5120)" json:"content3"`
	Content4       string    `gorm:"type:varchar(5120)" json:"content4"`
	CntLike        int64     `gorm:"default:0" json:"cnt_like"`
	CntBookmark    int64     `gorm:"default:0" json:"cnt_bookmark"`
	CntView        int64     `gorm:"default:0" json:"cnt_view"`
	SeqNovelStep1  int64     `json:"seq_novel_step1"`
	SeqNovelStep2  int64     `json:"seq_novel_step2"`
	SeqNovelStep3  int64     `json:"seq_novel_step3"`
	SeqNovelStep4  int64     `json:"seq_novel_step4"`
	ActiveYn       bool      `gorm:"default:true" json:"active_yn"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type NovelFinishBatchRunLog struct {
	SeqNovelFinishBatchRunLog int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_finish_batch_run_log"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type KeywordChoiceFirst struct {
	SeqKeywordChoiceFirst int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword_choice_first"`
	SeqKeyword            int64     `gorm:"index:uk_choice_first,unique" json:"seq_keyword"`
	SeqNovelStep1         int64     `gorm:"index:uk_choice_first,unique" json:"seq_novel_step1"`
	SeqNovelStep2         int64     `gorm:"default:0" json:"seq_novel_step2"`
	SeqNovelStep3         int64     `gorm:"default:0" json:"seq_novel_step3"`
	SeqNovelStep4         int64     `gorm:"default:0" json:"seq_novel_step4"`
	CntLikeStep1          int64     `gorm:"default:0" json:"cnt_like_step1"`
	CntLikeStep2          int64     `gorm:"default:0" json:"cnt_like_step2"`
	CntLikeStep3          int64     `gorm:"default:0" json:"cnt_like_step3"`
	CntLikeStep4          int64     `gorm:"default:0" json:"cnt_like_step4"`
	CntLikeTotal          int64     `gorm:"default:0" json:"cnt_like_total"`
	SuccessYn             bool      `gorm:"default:false" json:"success_yn"`
	FinishYn              bool      `gorm:"default:false" json:"finish_yn"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type KeywordChoiceSecond struct {
	SeqKeywordChoiceSecond int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword_choice_second"`
	SeqKeywordChoiceFirst  int64     `gorm:"index" json:"seq_keyword_choice_first"`
	SeqNovelFinish         int64     `gorm:"index" json:"seq_novel_finish"`
	CreatedAt              time.Time `json:"created_at"`
}

type ServiceInquiry struct {
	SeqServiceInquiry int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_service_inquiry"`
	SeqMember         int64     `gorm:"index" json:"seq_member"`
	Title             string    `gorm:"type:varchar(150)" json:"title"`
	Content           string    `gorm:"type:varchar(1024)" json:"content"`
	EmailYn           bool      `gorm:"default:false" json:"email_yn"`
	Status            int8      `gorm:"default:1" json:"status"`
	Answer            string    `gorm:"type:varchar(1024)" json:"answer"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Creator           string    `gorm:"type:varchar(50)" json:"creator"`
	Updator           string    `gorm:"type:varchar(50)" json:"updator"`
}

type Notice struct {
	SeqNotice      int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_notice"`
	SeqMemberAdmin int64     `gorm:"index" json:"seq_member_admin"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	ActiveYn       bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Creator        string    `gorm:"type:varchar(50)" json:"creator"`
	Updator        string    `gorm:"type:varchar(50)" json:"updator"`
}

type CategoryFaq struct {
	SeqCategoryFaq int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_category_faq"`
	CategoryFaq    string    `gorm:"type:varchar(50)" json:"category_faq"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ActiveYn       bool      `gorm:"default:false" json:"active_yn"`
	Creator        string    `gorm:"type:varchar(50)" json:"creator"`
	Updator        string    `gorm:"type:varchar(50)" json:"updator"`
}

type Faq struct {
	SeqFaq         int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_faq"`
	SeqMemberAdmin int64     `gorm:"index" json:"seq_member_admin"`
	SeqCategoryFaq int64     `gorm:"index" json:"seq_category_faq"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	ActiveYn       bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Creator        string    `gorm:"type:varchar(50)" json:"creator"`
	Updator        string    `gorm:"type:varchar(50)" json:"updator"`
}

type NovelDelete struct {
	SeqNovel     int64     `gorm:"index" json:"seq_novel"`
	Step         int8      `gorm:"index" json:"step"`
	SeqKeyword   int64     `gorm:"index" json:"seq_keyword"`
	SeqImage     int64     `gorm:"index" json:"seq_image"`
	SeqColor     int64     `gorm:"index" json:"seq_color"`
	SeqGenre     int64     `gorm:"index" json:"seq_genre"`
	SeqMember    int64     `gorm:"index" json:"seq_member"`
	Title        string    `gorm:"unique;type:varchar(1024)" json:"title"`
	Content      string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike      int64     `gorm:"default:0" json:"cnt_like"`
	CntView      int64     `gorm:"default:0" json:"cnt_view"`
	CntStep2     int64     `gorm:"default:0" json:"cnt_step2"`
	CntStep3     int64     `gorm:"default:0" json:"cnt_step3"`
	CntStep4     int64     `gorm:"default:0" json:"cnt_step4"`
	ActiveYn     bool      `gorm:"default:true" json:"active_yn"`
	TempYn       bool      `gorm:"default:false" json:"temp_yn"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	TypeDelete   int8      `json:"type_delete"`
	ReasonDelete string    `gorm:"type:varchar(1024)" json:"reason_delete"`
}

type Alarm struct {
	SeqAlarm   int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_alarm"`
	SeqMember  int64     `gorm:"index" json:"seq_member"`
	Title      string    `gorm:"type:varchar(50)" json:"title"`
	Content    string    `gorm:"type:varchar(1024)" json:"content"`
	TypeAlarm  int8      `json:"type_alarm"`
	ValueAlarm int       `json:"value_alarm"`
	Step       int8      `gorm:"default:0" json:"step"`
	CreatedAt  time.Time `json:"created_at"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	UpdatedAt  time.Time `json:"updated_at"`
}
