package tools

import (
	"ddaom/define"
	"ddaom/domain/schemas"
	"fmt"
	"log"

	"github.com/appleboy/go-fcm"
)

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
