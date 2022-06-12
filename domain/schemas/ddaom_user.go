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

// type MemberSubscribe struct {
// 	SeqMemberSubscribe int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_subscribe"`
// 	SeqMember          int64     `gorm:"index" json:"seq_member"`
// 	SeqMemberFollowing int64     `gorm:"index" json:"seq_member_following"`
// 	SubscribeYn        bool      `gorm:"default:false" json:"subscribe_yn"`
// 	CreatedAt          time.Time `json:"created_at"`
// }

type MemberSubscribe struct {
	SeqMemberSubscribe int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_subscribe"`
	SeqMember          int64     `gorm:"index" json:"seq_member"`
	SeqMemberOpponent  int64     `gorm:"index" json:"seq_member_opponent"`
	Status             string    `gorm:"type:ENUM('FOLLOWER','FOLLOWING','BOTH'); DEFAULT:'FOLLOWER'" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type MemberBookmark struct {
	SeqMemberBookmark int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_bookmark"`
	SeqMember         int64     `gorm:"index" json:"seq_member"`
	SeqNovelFinish    int64     `gorm:"index" json:"seq_novel_finish"`
	BookmarkYn        bool      `gorm:"default:false" json:"bookmark_yn"`
	CreatedAt         time.Time `json:"created_at"`
}

type MemberLikeStep1 struct {
	SeqMemberLike int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_like"`
	SeqMember     int64     `gorm:"index:idx_like" json:"seq_member"`
	SeqNovelStep1 int64     `gorm:"index:idx_like" json:"seq_novel_step1"`
	LikeYn        bool      `gorm:"default:false" json:"like_yn"`
	CreatedAt     time.Time `json:"created_at"`
}
type MemberLikeStep2 struct {
	SeqMemberLike int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_like"`
	SeqMember     int64     `gorm:"index:idx_like" json:"seq_member"`
	SeqNovelStep2 int64     `gorm:"index:idx_like" json:"seq_novel_step2"`
	LikeYn        bool      `gorm:"default:false" json:"like_yn"`
	CreatedAt     time.Time `json:"created_at"`
}
type MemberLikeStep3 struct {
	SeqMemberLike int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_like"`
	SeqMember     int64     `gorm:"index:idx_like" json:"seq_member"`
	SeqNovelStep3 int64     `gorm:"index:idx_like" json:"seq_novel_step3"`
	LikeYn        bool      `gorm:"default:false" json:"like_yn"`
	CreatedAt     time.Time `json:"created_at"`
}
type MemberLikeStep4 struct {
	SeqMemberLike int64     `gorm:"primaryKey;autoIncrement:true" json:"seq_member_like"`
	SeqMember     int64     `gorm:"index:idx_like" json:"seq_member"`
	SeqNovelStep4 int64     `gorm:"index:idx_like" json:"seq_novel_step4"`
	LikeYn        bool      `gorm:"default:false" json:"like_yn"`
	CreatedAt     time.Time `json:"created_at"`
}
