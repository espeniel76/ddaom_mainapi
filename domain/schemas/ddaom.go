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
	SeqMemberAdmin int64  `gorm:"primaryKey;autoIncrement:true" json:"seq_member_admin"`
	UserId         string `gorm:"type:varchar(50);unique" json:"user_id"`
	Password       string `gorm:"type:varchar(1024)" json:"password"`
	ActiveYn       bool   `gorm:"default:false" json:"active_yn"`
	Name           string `gorm:"type:varchar(50)" json:"name"`
	NickName       string `gorm:"type:varchar(50);unique;column:nick_name" json:"nick_name"`
	UserLevel      int8   `gorm:"default:5" json:"user_level"`
	ProfileImage   string `gorm:"type:varchar(1024)" json:"profile_image"`
	CreatedAt      time.Time
	Creator        string `gorm:"type:varchar(50)" json:"creator"`
	UpdatedAt      time.Time
	Updator        string `gorm:"type:varchar(50)" json:"updator"`
}

type MemberAdminLoginLog struct {
	SeqMemberAdminLoginLog int64  `gorm:"primaryKey;autoIncrement:true" json:"seq_member_admin_login_log"`
	SeqMemberAdmin         int64  `gorm:"index" json:"seq_member_admin"`
	Token                  string `gorm:"type:varchar(1024)" json:"token"`
	LoginAt                time.Time
}
