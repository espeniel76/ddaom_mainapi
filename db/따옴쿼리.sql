use ddaom;

SELECT title, content, cnt_like, created_at FROM novel_step1 WHERE seq_novel_step1 = 1;

SELECT * FROM novel_step1 ns;
SELECT * FROM novel_step2 ns;
SELECT * FROM novel_step3 ns;
SELECT * FROM novel_step4 ns;

desc notices;
insert into notices (seq_member_admin,title,content,created_at,updated_at)
values (1,'�޴��� ����Ȯ�� ���� ���� �߻� ���� �ȳ�','�� �̾�� �ʾƵ� �˴ϴ�.',NOW(),NOW());
desc category_faqs;
INSERT into category_faqs (category_faq,created_at,updated_at,active_yn,creator) values
('�з�1',NOW(),NOW(),true,'���������쵿'),
('�з�2',NOW(),NOW(),true,'���������쵿'),
('�з�3',NOW(),NOW(),true,'���������쵿'),
('�з�4',NOW(),NOW(),true,'���������쵿'),
('�з�5',NOW(),NOW(),true,'���������쵿');
desc faqs;
INSERT into faqs (seq_member_admin,seq_category_faq,title,content,active_yn,created_at) VALUES
(1,1,'�Ҽ��� �� �̾��� �ϳ���1?','�� �̾�� �ʾƵ� �˴ϴ� �����',true,NOW()),
(1,2,'�Ҽ��� �� �̾��� �ϳ���2?','�� �̾�� �ʾƵ� �˴ϴ� �����',true,NOW()),
(1,3,'�Ҽ��� �� �̾��� �ϳ���3?','�� �̾�� �ʾƵ� �˴ϴ� �����',true,NOW()),
(1,4,'�Ҽ��� �� �̾��� �ϳ���4?','�� �̾�� �ʾƵ� �˴ϴ� �����',true,NOW()),
(1,5,'�Ҽ��� �� �̾��� �ϳ���5?','�� �̾�� �ʾƵ� �˴ϴ� �����',true,NOW());