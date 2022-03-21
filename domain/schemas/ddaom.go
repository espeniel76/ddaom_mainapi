package schemas

import (
	"time"
)

type Member struct {
	SeqMember       int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member"`
	Email           string    `gorm:"type:varchar(255);unique" json:"email"`
	Token           string    `gorm:"type:varchar(1024)" json:"token"`
	ProfileImageUrl string    `gorm:"type:varchar(512)" json:"profile_image_url"`
	SnsType         string    `gorm:"type:ENUM('KAKAO','NAVER','FACEBOOK','GOOGLE','APPLE'); DEFAULT:'GOOGLE'" json:"sns_type"`
	ActiveYn        bool      `gorm:"default:false" json:"active_yn"`
	UserLevel       int8      `gorm:"default:5" json:"user_level"`
	AllocatedDb     int8      `json:"allocted_db"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
	AltNewEvent      bool      `gorm:"default:false" json:"alt_new_event"`
	AltSuccessfulBid bool      `gorm:"default:false" json:"alt_successful_bid"`
	AltNewContent    bool      `gorm:"default:false" json:"alt_new_content"`
	AltNightPush     bool      `gorm:"default:false" json:"alt_night_push"`
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
type Keyword struct {
	SeqKeyword int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword"`
	Keyword    string    `gorm:"unique;type:varchar(1024)" json:"keyword"`
	ActiveYn   bool      `gorm:"default:false" json:"active_yn"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CntTotal   int64     `gorm:"default:0" json:"cnt_total"`
}

type KeywordToday struct {
	SeqKeywordToday int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_keyword_today"`
	SeqKeyword      int64     `gorm:"index" json:"seq_keyword"`
	ViewDate        string    `gorm:"unique;type:char(8)" json:"view_date"`
	ActiveYn        bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// 장르
type Genre struct {
	SeqGenre  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_genre"`
	Genre     string    `gorm:"type:varchar(50)" json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ActiveYn  bool      `gorm:"default:false" json:"active_yn"`
	Creator   string    `gorm:"type:varchar(50)" json:"creator"`
	Updator   string    `gorm:"type:varchar(50)" json:"updator"`
}

type Image struct {
	SeqImage  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_image"`
	Image     string    `gorm:"type:varchar(1024)" json:"image"`
	ActiveYn  bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt time.Time `json:"created_at"`
	Creator   string    `gorm:"type:varchar(50)" json:"creator"`
}

type Color struct {
	SeqColor  int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_color"`
	Color     string    `gorm:"unique type:varchar(12)" json:"color"`
	ActiveYn  bool      `gorm:"default:false" json:"active_yn"`
	CreatedAt time.Time `json:"created_at"`
	Creator   string    `gorm:"type:varchar(50)" json:"creator"`
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
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NovelStep3 struct {
	SeqNovelStep3 int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_step3"`
	SeqNovelStep2 int64     `gorm:"index" json:"seq_novel_step2"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	Content       string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike       int64     `gorm:"default:0" json:"cnt_like"`
	CntView       int64     `gorm:"default:0" json:"cnt_view"`
	CntStep4      int64     `gorm:"default:0" json:"cnt_step4"`
	ActiveYn      bool      `gorm:"default:true" json:"active_yn"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NovelStep4 struct {
	SeqNovelStep4 int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_novel_step4"`
	SeqNovelStep3 int64     `gorm:"index" json:"seq_novel_step3"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	Content       string    `gorm:"type:varchar(5120)" json:"content"`
	CntLike       int       `gorm:"default:0" json:"cnt_like"`
	CntView       int64     `gorm:"default:0" json:"cnt_view"`
	ActiveYn      bool      `gorm:"default:true" json:"active_yn"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
