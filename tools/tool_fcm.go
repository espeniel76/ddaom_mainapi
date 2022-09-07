package tools

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain/schemas"
	"ddaom/memdb"
	"fmt"
	"log"
	"strconv"

	"github.com/appleboy/go-fcm"
)

func SendPushMessageTopic(alarm *schemas.Alarm) {
	msg := &fcm.Message{
		To: "/topics/DdaOm" + fmt.Sprint(alarm.SeqMember),
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

	go CacheMyPushCnt(alarm.SeqAlarm)

	log.Printf("%#v\n", response)
}

func CacheMyPushCnt(seqMember int64) {
	sdb := db.List[define.Mconn.DsnSlave]
	var cnt int64
	sdb.Model(schemas.Alarm{}).Where("seq_member = ? AND is_read = false", seqMember).Count(&cnt)
	memdb.Set("CACHES:USERS:PUSH_CNT:"+strconv.FormatInt(seqMember, 10), strconv.FormatInt(cnt, 10))
}

func SendPushMessage(pushToken string, alarm *schemas.Alarm) {
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
