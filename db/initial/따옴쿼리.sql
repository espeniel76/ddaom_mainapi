TRUNCATE table ddaom_user1.member_subscribes; 
TRUNCATE table ddaom_user2.member_subscribes;

SELECT A.seq_member, A.seq_member_following
FROM
(
        (SELECT seq_member, seq_member_following, created_at FROM ddaom_user1.member_subscribes WHERE subscribe_yn = true AND (seq_member = 8 OR seq_member_following = 8))
        UNION ALL
        (SELECT seq_member, seq_member_following, created_at FROM ddaom_user2.member_subscribes WHERE subscribe_yn = true AND (seq_member = 8 OR seq_member_following = 8))
) AS A;

SELECT A.seq_member, A.seq_member_following, A.created_at
FROM
(
        (SELECT seq_member, seq_member_following, created_at FROM ddaom_user1.member_subscribes WHERE subscribe_yn = true AND (seq_member = 8 OR seq_member_following = 8))
        UNION ALL
        (SELECT seq_member, seq_member_following, created_at FROM ddaom_user2.member_subscribes WHERE subscribe_yn = true AND (seq_member = 8 OR seq_member_following = 8))
) AS A
ORDER BY A.created_at DESC
LIMIT 0, 10;

SELECT * FROM ddaom_user1.member_subscribes;
SELECT * FROM ddaom_user2.member_subscribes;

(SELECT seq_member, seq_member_following, created_at FROM ddaom_user1.member_subscribes WHERE subscribe_yn = true AND (seq_member = 8 OR seq_member_following = 8))
UNION ALL
(SELECT seq_member, seq_member_following, created_at FROM ddaom_user2.member_subscribes WHERE subscribe_yn = true AND (seq_member = 8 OR seq_member_following = 8));