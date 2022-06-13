<?php

error_reporting(E_ALL ^ E_NOTICE);
error_reporting(E_ALL ^ ~E_WARNING);

$conn = mysqli_connect("localhost", "espeniel", "anjgkrp", "ddaom");
$redis = new Redis();
$redis->connect("127.0.0.1", 6379);
$redis->auth("anjgkrp");
$redis->select(1);

// 1. 인기작 (추출 테스트)
$sql =
	"SELECT seq_novel_finish, seq_image, seq_color, title FROM novel_finishes WHERE active_yn = true ORDER BY cnt_like DESC, cnt_view DESC LIMIT 10";
$rets = mysqli_query($conn, $sql);
$list = [];
while ($main = mysqli_fetch_assoc($rets)) {
	$list[] = $main;
}
$redis->set("CACHES:MAIN:LIST_POPULAR", json_encode($list, JSON_UNESCAPED_UNICODE | JSON_NUMERIC_CHECK));

// 1. 완결작 (추출 테스트)
$sql =
	"SELECT seq_novel_finish, seq_image, seq_color, title FROM novel_finishes WHERE active_yn = true ORDER BY created_at DESC LIMIT 10";
$rets = mysqli_query($conn, $sql);
$list = [];
while ($main = mysqli_fetch_assoc($rets)) {
	$list[] = $main;
}
$redis->set("CACHES:MAIN:LIST_FINISH", json_encode($list, JSON_UNESCAPED_UNICODE | JSON_NUMERIC_CHECK));
