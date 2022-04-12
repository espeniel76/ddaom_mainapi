TRUNCATE table pattern_english_user1.chapter_logs;
TRUNCATE table pattern_english_user1.sentence_logs;
TRUNCATE TABLE pattern_english_user1.member_like_sentences;
TRUNCATE TABLE pattern_english_user1.member_scrap_sentences;
TRUNCATE TABLE pattern_english_user1.chapter_middle_logs;
TRUNCATE TABLE pattern_english_user1.sentence_middle_logs;
TRUNCATE TABLE pattern_english_user1.member_scrap_chapters;
TRUNCATE table pattern_english_user2.chapter_logs;
TRUNCATE table pattern_english_user2.sentence_logs;
TRUNCATE TABLE pattern_english_user2.member_like_sentences;
TRUNCATE TABLE pattern_english_user2.member_scrap_sentences;
TRUNCATE TABLE pattern_english_user2.chapter_middle_logs;
TRUNCATE TABLE pattern_english_user2.sentence_middle_logs;
TRUNCATE TABLE pattern_english_user2.member_scrap_chapters;


SELECT * FROM pattern_english_user1.chapter_logs cl;
SELECT * FROM pattern_english_user1.sentence_logs sl;

SELECT * FROM pattern_english.chapters order by seq_chapter asc;
SELECT * FROM pattern_english.sentences order by seq_sentence DESC;
SELECT * FROM pattern_english_user1.member_like_sentences mls;
SELECT * FROM pattern_english_user2.member_like_sentences mls;
