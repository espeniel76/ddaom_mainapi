package schemas

import (
	"time"
)

type MemberLoginLog struct {
	SeqMemberLoginLog int64  `gorm:"primaryKey;autoIncrement:true"`
	SeqMember         int64  `gorm:"index"`
	Token             string `gorm:"type:varchar(1024)"`
	LoginAt           time.Time
}

type MemberExist struct {
	SeqMember int64 `gorm:"unique"`
}
