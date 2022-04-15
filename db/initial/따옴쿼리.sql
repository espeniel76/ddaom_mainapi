use ddaom;
SELECT * FROM novel_step1 nf;
SELECT * FROM novel_finishes nf;
SELECT * from novel_step4 ns;

desc novel_finishes;
INSERT INTO novel_finishes (
seq_novel_step1,
seq_novel_step2,
seq_novel_step3,
seq_novel_step4,
active_yn,
created_at,
cnt_bookmark,
seq_member_step1,
seq_member_step2,
seq_member_step3,
seq_member_step4,
cnt_like
) VALUES (1,5,2,1,true,now(),50,8,8,9,9,3120);

(
	SELECT A.seq_member_following, false AS is_follower 
	FROM
	(
		(SELECT seq_member, seq_member_following FROM ddaom_user1.member_subscribes ms)
		UNION ALL
		(SELECT seq_member, seq_member_following FROM ddaom_user2.member_subscribes ms)
	) AS A
	WHERE A.seq_member = 8
)
UNION ALL
(
	SELECT A.seq_member, true AS is_follower
	FROM
	(
		(SELECT seq_member, seq_member_following FROM ddaom_user1.member_subscribes ms)
		UNION ALL
		(SELECT seq_member, seq_member_following FROM ddaom_user2.member_subscribes ms)
	) AS A
	WHERE A.seq_member_following = 8
);

SELECT A.seq_member, A.seq_member_following, created_at
FROM
(
	(SELECT seq_member, seq_member_following, created_at FROM ddaom_user1.member_subscribes WHERE subscribe_yn = true)
	UNION ALL
	(SELECT seq_member, seq_member_following, created_at FROM ddaom_user2.member_subscribes WHERE subscribe_yn = true)
) AS A 
WHERE A.seq_member = 8 OR A.seq_member_following = 8
;

SELECT * FROM ddaom_user1.member_subscribes ms;