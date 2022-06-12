TRUNCATE table ddaom_user1.member_subscribes;
TRUNCATE table ddaom_user2.member_subscribes;
UPDATE ddaom.member_details SET cnt_subscribe = 0;

SELECT * FROM ddaom_user1.member_subscribes;
SELECT * FROM ddaom_user2.member_subscribes;

SELECT seq_member, cnt_subscribe FROM ddaom.member_details;

use ddaom;
SELECT * FROM member_admins ma;
desc member_admins;
INSERT INTO member_admins (user_id,password,active_yn,name,nick_name,user_level,created_at) VALUES
('espeniel','anjgkrp',true,'¹®º´ÁØ','´«¹°Á¥Àº¿ìµ¿',5,NOW());