<?php
// 완결소설 배치
error_reporting(E_ALL ^ E_NOTICE);

$conn = mysqli_connect('localhost', 'espeniel', 'anjgkrp', 'ddaom');

// 1. 종료 된 키워드 조회
















// // 1. 현재 시간대에 경매 끝난 애장품 조회 (종료 된 건중, 낙찰 이 되지 않은 데이터 찾아라)
// $sql = 'SELECT seq_good, auction_start_price, auction_now_price FROM goods WHERE NOW() > end_date AND bid_price = 0';
// $result = mysqli_query($conn, $sql);
// while($row = mysqli_fetch_assoc($result)) {
// 	// 최고가 경매 찾는다.
// 	// $sql = 'SELECT * FROM auctions WHERE seq_good = '.$row['seq_good'].' ORDER BY join_price DESC LIMIT 1';
// 	$sql = 'SELECT * FROM auctions WHERE seq_good = ? ORDER BY join_price DESC LIMIT 1';
// 	$stmt = mysqli_prepare($conn, $sql);
// 	$bind = mysqli_stmt_bind_param($stmt, "i", $row['seq_good']);
// 	$exec = mysqli_stmt_execute($stmt);
// 	$retVal = mysqli_stmt_get_result($stmt);
// 	$rowSub = mysqli_fetch_assoc($retVal);

// 	// 해당 참여로그/상품/실행로그 업데이트 해준다.
// 	if ($rowSub) {

// 		// 1. 참여로그
// 		$sql = 'UPDATE auctions SET success_yn = true, bid_price = ? WHERE seq_auction = ?';
// 		$stmt = mysqli_prepare($conn, $sql);
// 		$bind = mysqli_stmt_bind_param($stmt, "ii", $rowSub['bid_price'], $rowSub['seq_auction']);
// 		mysqli_stmt_execute($stmt);

// 		// 2. 상품
// 		$sql = 'UPDATE goods SET bid_price = ?, bid_at = NOW() WHERE seq_good = ?';
// 		$stmt = mysqli_prepare($conn, $sql);
// 		$bind = mysqli_stmt_bind_param($stmt, "ii", $rowSub['bid_price'], $row['seq_good']);
// 		mysqli_stmt_execute($stmt);

// 		// 3. 실행로그
// 		$sql = 'INSERT INTO auction_join_logs (seq_good, seq_member, bid_price, bid_date, created_at) VALUES (?, ?, ?, ?, NOW())';
// 		$stmt = mysqli_prepare($conn, $sql);
// 		$bind = mysqli_stmt_bind_param($stmt, "iiiss", $row['seq_good'], $rowSub['seq_member'], $rowSub['join_price'], $rowSub['created_at']);
// 		mysqli_stmt_execute($stmt);
// 	}
// 	mysqli_free_result($retVal);
// }
// mysqli_free_result($result);