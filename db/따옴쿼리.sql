use ddaom;

SELECT title, content, cnt_like, created_at FROM novel_step1 WHERE seq_novel_step1 = 1;

SELECT * FROM novel_step1 ns;
SELECT * FROM novel_step2 ns;
SELECT * FROM novel_step3 ns;
SELECT * FROM novel_step4 ns;

desc notices;
insert into notices (seq_member_admin,title,content,created_at,updated_at)
values (1,'휴대폰 본인확인 서비스 순단 발생 예정 안내','꼭 이어쓰지 않아도 됩니다.',NOW(),NOW());
desc category_faqs;
INSERT into category_faqs (category_faq,created_at,updated_at,active_yn,creator) values
('분류1',NOW(),NOW(),true,'눈물젖은우동'),
('분류2',NOW(),NOW(),true,'눈물젖은우동'),
('분류3',NOW(),NOW(),true,'눈물젖은우동'),
('분류4',NOW(),NOW(),true,'눈물젖은우동'),
('분류5',NOW(),NOW(),true,'눈물젖은우동');
desc faqs;
INSERT into faqs (seq_member_admin,seq_category_faq,title,content,active_yn,created_at) VALUES
(1,1,'소설을 꼭 이어써야 하나요1?','꼭 이어쓰지 않아도 됩니다 블라블라',true,NOW()),
(1,2,'소설을 꼭 이어써야 하나요2?','꼭 이어쓰지 않아도 됩니다 블라블라',true,NOW()),
(1,3,'소설을 꼭 이어써야 하나요3?','꼭 이어쓰지 않아도 됩니다 블라블라',true,NOW()),
(1,4,'소설을 꼭 이어써야 하나요4?','꼭 이어쓰지 않아도 됩니다 블라블라',true,NOW()),
(1,5,'소설을 꼭 이어써야 하나요5?','꼭 이어쓰지 않아도 됩니다 블라블라',true,NOW());