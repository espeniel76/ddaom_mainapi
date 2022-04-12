use ddaom;

SELECT title, content, cnt_like, created_at FROM novel_step1 WHERE seq_novel_step1 = 1;

SELECT * FROM novel_step1 ns;
SELECT * FROM novel_step2 ns;
SELECT * FROM novel_step3 ns;
SELECT * FROM novel_step4 ns;

(SELECT seq_novel_step2, 0 AS seq_novel_step3 FROM novel_step2 WHERE seq_novel_step1 = 5 ORDER BY cnt_like DESC LIMIT 1)
UNION
(SELECT 0 AS seq_novel_step2, seq_novel_step3 FROM novel_step3 WHERE seq_novel_step1 = 5 ORDER BY cnt_like DESC LIMIT 1);

(SELECT seq_novel_step2, 0 AS seq_novel_step3, 0 AS seq_novel_step4 FROM novel_step2 WHERE seq_novel_step1 = 1 ORDER BY cnt_like DESC LIMIT 1)
UNION ALL
(SELECT 0 AS seq_novel_step2, seq_novel_step3, 0 AS seq_novel_step4 FROM novel_step3 WHERE seq_novel_step1 = 1 ORDER BY cnt_like DESC LIMIT 1)
UNION ALL
(SELECT 0 AS seq_novel_step2, 0 AS seq_novel_step3, seq_novel_step4 FROM novel_step4 WHERE seq_novel_step1 = 1 ORDER BY cnt_like DESC LIMIT 1);