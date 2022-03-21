package schemas

import (
	"time"
)

type MemberLoginLog struct {
	SeqMemberLoginLog int64  `gorm:"primaryKey;autoIncrement:true" json:"seq_member_login_log"`
	SeqMember         int64  `gorm:"index"`
	Token             string `gorm:"type:varchar(1024)"`
	LoginAt           time.Time
}

type MemberExist struct {
	SeqMember int64 `gorm:"unique" json:"seq_member"`
}

type MemberSubscribe struct {
	SeqMemberSubscribe int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_subscribe"`
	SeqMember          int64     `gorm:"index" json:"seq_member"`
	SeqMemberFollowing int64     `gorm:"index" json:"seq_member_following"`
	CreatedAt          time.Time `json:"created_at"`
}

type MemberBookmark struct {
	SeqMemberBookmark int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_bookmark"`
	SeqMember         int64     `gorm:"index" json:"seq_member"`
	SeqNovelStep1     int64     `gorm:"index" json:"seq_novel_step1"`
	CreatedAt         time.Time `json:"created_at"`
}

type MemberLike struct {
	SeqMemberLike int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_like"`
	SeqMember     int64     `gorm:"index" json:"seq_member"`
	SeqNovelStep1 int64     `gorm:"index" json:"seq_novel_step1"`
	SeqNovelStep2 int64     `gorm:"index" json:"seq_novel_step2"`
	SeqNovelStep3 int64     `gorm:"index" json:"seq_novel_step3"`
	SeqNovelStep4 int64     `gorm:"index" json:"seq_novel_step4"`
	CreatedAt     time.Time `json:"created_at"`
}
