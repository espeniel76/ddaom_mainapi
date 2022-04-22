TRUNCATE table ddaom_user1.member_subscribes;
TRUNCATE table ddaom_user2.member_subscribes;
UPDATE ddaom.member_details SET cnt_subscribe = 0;

SELECT * FROM ddaom_user1.member_subscribes;
SELECT * FROM ddaom_user2.member_subscribes;

SELECT seq_member, cnt_subscribe FROM ddaom.member_details;

