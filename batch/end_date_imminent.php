<?php
/**
 *
 * 주제어 종료일 임박 알림
 * - 해당 주제어에 등록한 글이 0개에 해당하는 사용자에게 알림 (등록 후 삭제한 경우도 0 개로 간주)
 * - admin에 등록된 진행중인 주제어 종료일 기준
 * - 푸쉬 일시 : 종료일 7일 전부터 당일까지, 시간대 오전 11시 통일
 * - 클릭 시 [4. 연재중인 릴레이 소설_리스트] 해당 주제어로 선택된 리스트 화면 호출
 * 1) 종료일 7일 전 ~ 2일 전 매일 1회씩
 * - 메세지 : {"대상 주제어"} + 주제어의 마감일이 얼마 남지 않았습니다!! 작가님의 이야기 기다릴게요. (이모티콘)
 * 2) 종료일 1일 전, 당일 1회씩
 * - 메세지 : {"대상 주제어"} + 주제어가 곧 마감됩니다!!!! 작가님의 이야기 들려주실거죠? (이모티콘)
 */
error_reporting(E_ALL ^ E_NOTICE);
date_default_timezone_set("Asia/Seoul");
$conn = mysqli_connect("localhost", "espeniel", "anjgkrp", "ddaom");

// print_r($conn);

// 1. 조건에 부합하는 키워드 조회
$listToken = [];
$timeNow = date("H");
$timeNow = (int) $timeNow;
// echo $timeNow;
$sql =
	"SELECT seq_keyword, keyword, start_date, end_date, DATEDIFF(end_date, NOW()) AS remain_day FROM keywords k WHERE active_yn = true AND finish_yn = false AND DATEDIFF(end_date, NOW()) BETWEEN 1 AND 7";
$result = mysqli_query($conn, $sql);

// 1.1. 키워드 단위로 순회
while ($o = mysqli_fetch_assoc($result)) {
	print_r($o);

	// 2. 발송 대상 추출
	$sql = "
	SELECT
		m.seq_member,
		m.push_token,
		md.is_night_push
	FROM
		members m INNER JOIN member_details md ON m.seq_member = md.seq_member
	WHERE
		push_token IS NOT NULL
		AND m.seq_member NOT IN (
		SELECT DISTINCT seq_member FROM novel_step1 WHERE seq_keyword = {$o["seq_keyword"]}
		UNION
		SELECT DISTINCT ns2.seq_member FROM novel_step2 ns2 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns2.seq_novel_step1 WHERE ns1.seq_keyword = {$o["seq_keyword"]}
		UNION
		SELECT DISTINCT ns3.seq_member FROM novel_step3 ns3 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns3.seq_novel_step1 WHERE ns1.seq_keyword = {$o["seq_keyword"]}
		UNION
		SELECT DISTINCT ns4.seq_member FROM novel_step4 ns4 INNER JOIN novel_step1 ns1 ON ns1.seq_novel_step1 = ns4.seq_novel_step1 WHERE ns1.seq_keyword = {$o["seq_keyword"]}
	)";
	$ret = mysqli_query($conn, $sql);
	$listToken = [];
	while ($item = mysqli_fetch_assoc($ret)) {
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

	// 3. 발송
	for ($i = 0; $i < sizeof($listToken); $i++) {
		$item = $listToken[$i];
		// print_r($item);
		sendPushBefore(
			$conn,
			$item["seq_member"],
			$o["keyword"],
			$o["seq_keyword"],
			(int) $o["remain_day"],
			$item["push_token"]
		);
	}
}

function sendPushBefore($conn, $seqMember, $keyword, $seqKeyword, $remainDay, $token)
{
	$body = "";
	switch ($remainDay) {
		case 7:
		case 6:
		case 5:
		case 4:
			$body = "\"{$keyword}\" 주제어의 마감일이 얼마 남지 않았습니다!! 작가님의 이야기를 기다릴게요.";
			break;
		case 3:
		case 2:
			$body = "\"{$keyword}\" 주제어의 마감일 {$remainDay}일 전 입니다!! 작가님의 이야기를 기다릴게요.";
			break;
		case 1:
		case 0:
			$body = "\"{$keyword}\" 주제어가 곧 마감됩니다!!!! 작가님의 이야기 들려주실거죠?";
			break;
	}
	$sql = "INSERT INTO alarms (seq_member, title, content, type_alarm, value_alarm, created_at)
	VALUES ({$seqMember}, '따옴', '{$body}', 7, {$seqKeyword}, NOW())";
	mysqli_query($conn, $sql);
	$seqAlarm = mysqli_insert_id($conn);
	sendPush($token, $body, $seqAlarm, $seqKeyword);
}

function sendPush($token, $body, $seqAlarm, $valueAlarm)
{
	$notification = [
		"title" => "따옴",
		"body" => $body,
		"sound" => "default",
	];
	$data = ["seq_alarm" => $seqAlarm, "type_alram" => 7, "value_alarm" => $valueAlarm];
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

	print_r($result);
}
