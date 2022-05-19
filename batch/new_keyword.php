<?php
/**
 * ### 신규 주제어 등록 알림 ###
 * - admin에 등록된 해당 주제어의 상태가 '사용', '진행'으로 되어 진행중인 주제어로 당일 처리된 경우 알림
 * - 푸쉬 일시 : 해당 주제어 시작일의 당일 오후 1시
 * - 메세지 : {"진행중으로 등록된 주제어"} + 신규 주제어가 등록되었습니다.
 * - 클릭 시 [9. Step1 글쓰기_표지선택] 화면 호출
 * - 해당 주제어 진행이 종료된 경우 : '이미 마감되었습니다. 새로운 주제어를 확인해보 세요!'(확인) Alert 호출, (확인) 클릭 시 [2-1. Intro] 화면 호출
 */
error_reporting(E_ALL ^ E_NOTICE);
date_default_timezone_set("Asia/Seoul");
$conn = mysqli_connect("localhost", "espeniel", "anjgkrp", "ddaom");

// 1. 신규 주제어 추출
$sql = "SELECT * FROM keywords k WHERE
k.seq_keyword NOT IN (SELECT seq_keyword FROM keyword_alarm_logs)
AND k.active_yn = true
AND NOW() BETWEEN k.start_date AND k.end_date ORDER BY k.seq_keyword ASC LIMIT 1
";
$result = mysqli_query($conn, $sql);
$newKeyword = mysqli_fetch_assoc($result);
// print_r($newKeyword);
if (!$newKeyword) {
	mysqli_close($conn);
	exit();
}

// 2. 발송 대상 추출
$sql =
	"SELECT m.push_token, m.seq_member, md.is_night_push FROM member_details md INNER JOIN members m ON md.seq_member = m.seq_member WHERE md.is_new_keyword = true";
$result = mysqli_query($conn, $sql);
$listToken = [];
$timeNow = date("H");
$timeNow = (int) $timeNow;
while ($item = mysqli_fetch_assoc($result)) {
	// 2.1. 현재 시간 따져서, 야간 푸쉬 시간이면 발송하지 않는다.
	if ($item["is_night_push"] == false) {
		// 2.2. 낮인지 체크
		if ($timeNow >= 9 && $timeNow <= 20) {
			$listToken[] = $item;
		}
	} else {
		$listToken[] = $item;
	}
}
// print_r($listToken);
if (sizeof($listToken) < 1) {
	mysqli_close($conn);
	exit();
}

// 3. 발송
$body = "❝{$newKeyword["keyword"]}❞ 신규 주제어가 등록되었습니다. 작가님의 이야기를 들려주세요.";
for ($i = 0; $i < sizeof($listToken); $i++) {
	$o = $listToken[$i];

	// 4. 푸시 테이블 삽입
	$sql = "INSERT INTO alarms (seq_member, title, content, type_alarm, value_alarm, created_at)
	VALUES ({$o["seq_member"]}, '따옴', '{$body}', 1, {$newKeyword["seq_keyword"]}, NOW())";
	echo "{$sql}\n\r";
	mysqli_query($conn, $sql);
	$seqAlarm = mysqli_insert_id($conn);

	// 5. 발송
	sendPush($o["push_token"], $body, $seqAlarm, $newKeyword["seq_keyword"]);
}

// 5. 발송 로그 삽입
$sql =
	"INSERT INTO keyword_alarm_logs (seq_keyword, created_at, cnt_push) VALUES ({$newKeyword["seq_keyword"]}, NOW(), " .
	sizeof($listToken) .
	")";
// echo $sql;
mysqli_query($conn, $sql);
mysqli_close($conn);

function sendPush($token, $body, $seqAlarm, $valueAlarm)
{
	$notification = [
		"title" => "따옴",
		"body" => $body,
		"sound" => "default",
	];
	$data = ["seq_alarm" => $seqAlarm, "type_alram" => 1, "value_alarm" => $valueAlarm];
	$fcmNotification = [
		"to" => $token,
		"notification" => $notification,
		"data" => $data,
	];
	$headers = [
		"Authorization: key=AAAAs8DEFV4:APA91bHjJF63wpyefl-6IBMhJ0PVb0VPePwirNxes3PzRgMxg7wb1Q8ykTyzxnTrCVVMX8cE5ROxvjWJLLZ9cRw8pt5daXUsd-mxiK4jqgdVkR_XWaUW1snEXBSFFnebSR_D2L-Pn-wY",
		"Content-Type: application/json",
	];

	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, "https://fcm.googleapis.com/fcm/send");
	curl_setopt($ch, CURLOPT_POST, true);
	curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
	curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($fcmNotification));
	$result = curl_exec($ch);
	curl_close($ch);

	// echo $result;
}
