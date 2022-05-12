<?php
/**
 * ### 완결 소설 배치 ###
 * 예외 case 정리 (1차 대결 : 스텝별)
 * 1) 하나의 step에 좋아요 수가 같은 경우 발생 (ex. step3 에서 좋아요 200개가 2개 이상인 경우)
 *     => 배치 돌 때 등록일시가 가장 이른 소설에 좋아요 카운팅 1개 자동으로 추가
 * 2) 좋아요가 모두 0개인 step 발생 (ex. step3에 등록된 글이 10개인데 모두 좋아요가 0개인 경우)
 *     => 탈락
 * 3) 작성되지 않은 step 발생 (ex. step4에 등록된 글이 0개인 경우)
 *     => 탈락
 */
error_reporting(E_ALL ^ E_NOTICE);

$conn = mysqli_connect("localhost", "espeniel", "anjgkrp", "ddaom");

// 0. 배치 실행 로그
$sql = "INSERT INTO novel_finish_batch_run_logs (created_at) VALUES (NOW())";
mysqli_query($conn, $sql);
$lastInsertId = mysqli_insert_id($conn);

// 1. 종료 된 키워드 조회
$sql =
	"SELECT * FROM keywords k WHERE active_yn = true AND end_date < NOW() AND finish_yn = false ORDER BY created_at DESC LIMIT 1";
$result = mysqli_query($conn, $sql);
$keyword = mysqli_fetch_assoc($result);
if (!$keyword) {
	// 없으면 종료
	mysqli_close($conn);
	exit();
}

// 추출이 되든 안되든, 1회만 실행 시켜야 하므로 플래그 업데이트
// $sql = "UPDATE keywords SET finish_yn = true, finished_at = NOW() WHERE seq_keyword = {$keyword["seq_keyword"]}";
// mysqli_query($conn, $sql);
// 일단 주석처리

$list = [];

// 2. 당 키워드 관련 글 조회
// 2.1. STEP1
$sql = "SELECT * FROM novel_step1 ns WHERE seq_keyword = {$keyword["seq_keyword"]}";
$result = mysqli_query($conn, $sql);
$isRun = false;
while ($novelStep1 = mysqli_fetch_assoc($result)) {
	if ($novelStep1 && intval($novelStep1["cnt_like"]) > 0) {
		$isRun = true;
		$novelStep2 = null;
		$novelStep3 = null;
		$novelStep4 = null;

		// 2.2. STEP2
		$sql = "SELECT * FROM novel_step2 ns WHERE seq_novel_step1 = {$novelStep1["seq_novel_step1"]} ORDER BY cnt_like DESC, updated_at ASC LIMIT 1";
		$result = mysqli_query($conn, $sql);
		$novelStep2 = mysqli_fetch_assoc($result);
		if ($novelStep2 && intval($novelStep2["cnt_like"]) > 0) {
			// 2.3. STEP3
			$sql = "SELECT * FROM novel_step3 ns WHERE seq_novel_step1 = {$novelStep1["seq_novel_step1"]} ORDER BY cnt_like DESC, updated_at ASC LIMIT 1";
			$result = mysqli_query($conn, $sql);
			$novelStep3 = mysqli_fetch_assoc($result);
			if ($novelStep3 && intval($novelStep3["cnt_like"]) > 0) {
				// 2.4. STEP4
				$sql = "SELECT * FROM novel_step4 ns WHERE seq_novel_step1 = {$novelStep1["seq_novel_step1"]} ORDER BY cnt_like DESC, updated_at ASC LIMIT 1";
				$result = mysqli_query($conn, $sql);
				$novelStep4 = mysqli_fetch_assoc($result);
				if (!$novelStep4 || intval($novelStep4["cnt_like"]) == 0) {
					$novelStep4 = null;
				}
			} else {
				$novelStep3 = null;
			}
		} else {
			$novelStep2 = null;
		}

		// 3. 1차 데이터 만들기
		$seqKeyword = $keyword["seq_keyword"];
		$seqNovelStep1 = $novelStep1["seq_novel_step1"];
		$seqNovelStep2 = $novelStep2 ? $novelStep2["seq_novel_step2"] : 0;
		$seqNovelStep3 = $novelStep3 ? $novelStep3["seq_novel_step3"] : 0;
		$seqNovelStep4 = $novelStep4 ? $novelStep4["seq_novel_step4"] : 0;
		$cntLikeStep1 = intval($novelStep1["cnt_like"]);
		$cntLikeStep2 = $novelStep2 ? intval($novelStep2["cnt_like"]) : 0;
		$cntLikeStep3 = $novelStep3 ? intval($novelStep3["cnt_like"]) : 0;
		$cntLikeStep4 = $novelStep4 ? intval($novelStep4["cnt_like"]) : 0;
		$cntLikeTotal = $cntLikeStep1 + $cntLikeStep2 + $cntLikeStep3 + $cntLikeStep4;
		$successYn = 0;
		if ($cntLikeStep1 > 0 && $cntLikeStep2 > 0 && $cntLikeStep3 > 0 && $cntLikeStep4 > 0) {
			$successYn = 1;
			$list[$seqNovelStep1]["1"] = $novelStep1;
			$list[$seqNovelStep1]["2"] = $novelStep2;
			$list[$seqNovelStep1]["3"] = $novelStep3;
			$list[$seqNovelStep1]["4"] = $novelStep4;
		}
		$finishYn = 0;

		$sql = "
			INSERT INTO keyword_choice_firsts (
				seq_keyword,
				seq_novel_step1,
				seq_novel_step2,
				seq_novel_step3,
				seq_novel_step4,
				cnt_like_step1,
				cnt_like_step2,
				cnt_like_step3,
				cnt_like_step4,
				cnt_like_total,
				success_yn,
				finish_yn,
				created_at,
				updated_at
			) VALUES (
				{$seqKeyword},
				{$seqNovelStep1},
				{$seqNovelStep2},
				{$seqNovelStep3},
				{$seqNovelStep4},
				{$cntLikeStep1},
				{$cntLikeStep2},
				{$cntLikeStep3},
				{$cntLikeStep4},
				{$cntLikeTotal},
				{$successYn},
				{$finishYn},
				NOW(),
				NOW()
			)
		";
		mysqli_query($conn, $sql);
	}
}

// 4. 뭔가 돈 기록이 있다
if ($isRun) {
	$tmp = null;
	$listTmp = [];
	// 5. keyword_choice_firsts 성공인데, 미완을 찾아라
	$sql = "SELECT * FROM keyword_choice_firsts WHERE success_yn = true AND finish_yn = false";
	$result = mysqli_query($conn, $sql);
	while ($listFirst = mysqli_fetch_assoc($result)) {
		// 5.1. 최초이면 데이터 넣고
		if (!$tmp) {
			$tmp = $listFirst;
			$listTmp[] = $listFirst;

			// 5.2. 동률이면 데이터 넣고
		} elseif (intval($tmp["cnt_like_total"]) == intval($listFirst["cnt_like_total"])) {
			$tmp = $listFirst;
			$listTmp[] = $listFirst;

			// 5.3. 이상이면 기존 데이터 다 지우고 데이터 넣고
		} elseif (intval($tmp["cnt_like_total"]) < intval($listFirst["cnt_like_total"])) {
			$tmp = $listFirst;
			$listTmp = [];
			$listTmp[] = $listFirst;
		}
	}

	// 6. 실 데이터 넣기 작업
	for ($i = 0; $i < sizeof($listTmp); $i++) {
		$oChoice = $listTmp[$i];
		$oNovelInfo = $list[$oChoice["seq_novel_step1"]];

		// 6.1. novel_finishes 넣기
		$seq_keyword = $oNovelInfo["1"]["seq_keyword"];
		$seq_image = $oNovelInfo["1"]["seq_image"];
		$seq_color = $oNovelInfo["1"]["seq_color"];
		$seq_genre = $oNovelInfo["1"]["seq_genre"];
		$seq_member_step1 = $oNovelInfo["1"]["seq_member"];
		$seq_member_step2 = $oNovelInfo["2"]["seq_member"];
		$seq_member_step3 = $oNovelInfo["3"]["seq_member"];
		$seq_member_step4 = $oNovelInfo["4"]["seq_member"];
		$title = $oNovelInfo["1"]["title"];
		$content1 = $oNovelInfo["1"]["content"];
		$content2 = $oNovelInfo["2"]["content"];
		$content3 = $oNovelInfo["3"]["content"];
		$content4 = $oNovelInfo["4"]["content"];
		$cnt_like =
			intval($oNovelInfo["1"]["cnt_like"]) +
			intval($oNovelInfo["2"]["cnt_like"]) +
			intval($oNovelInfo["3"]["cnt_like"]) +
			intval($oNovelInfo["4"]["cnt_like"]);
		$cnt_bookmark = 0;
		$cnt_view = $oNovelInfo["1"]["cnt_view"];
		$seq_novel_step1 = $oNovelInfo["1"]["seq_novel_step1"];
		$seq_novel_step2 = $oNovelInfo["2"]["seq_novel_step2"];
		$seq_novel_step3 = $oNovelInfo["3"]["seq_novel_step3"];
		$seq_novel_step4 = $oNovelInfo["4"]["seq_novel_step4"];
		$active_yn = 1;
		$sql = "
		INSERT INTO novel_finishes (
			seq_keyword,
			seq_image,
			seq_color,
			seq_genre,
			seq_member_step1,
			seq_member_step2,
			seq_member_step3,
			seq_member_step4,
			title,
			content1,
			content2,
			content3,
			content4,
			cnt_like,
			cnt_bookmark,
			cnt_view,
			seq_novel_step1,
			seq_novel_step2,
			seq_novel_step3,
			seq_novel_step4,
			active_yn,
			created_at,
			updated_at
		) VALUES (
			{$seq_keyword},
			{$seq_image},
			{$seq_color},
			{$seq_genre},
			{$seq_member_step1},
			{$seq_member_step2},
			{$seq_member_step3},
			{$seq_member_step4},
			'{$title}',
			'{$content1}',
			'{$content2}',
			'{$content3}',
			'{$content4}',
			{$cnt_like},
			{$cnt_bookmark},
			{$cnt_view},
			{$seq_novel_step1},
			{$seq_novel_step2},
			{$seq_novel_step3},
			{$seq_novel_step4},
			{$active_yn},
			NOW(),
			NOW()
		)
		";
		mysqli_query($conn, $sql);
		$seqNovelFinish = mysqli_insert_id($conn);

		//  6.2. keyword_choice_seconds 넣기
		$sql = "INSERT INTO keyword_choice_seconds (
			seq_keyword_choice_first,
			seq_novel_finish,
			created_at
		) VALUES (
			{$oChoice["seq_keyword_choice_first"]},
			{$seqNovelFinish},
			NOW()
		)";
		mysqli_query($conn, $sql);

		// 6.3. flag update
		$sql = "UPDATE keyword_choice_firsts SET finish_yn = true, updated_at = NOW() WHERE seq_keyword_choice_first = {$oChoice["seq_keyword_choice_first"]}";
		mysqli_query($conn, $sql);
	}
}

$sql = "UPDATE novel_finish_batch_run_logs SET updated_at = NOW() WHERE seq_novel_finish_batch_run_log = {$lastInsertId}";
mysqli_query($conn, $sql);

mysqli_close($conn);
exit();

function println($s)
{
	echo "{$s}\n";
}
