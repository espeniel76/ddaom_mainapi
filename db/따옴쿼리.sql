use ddaom;

SELECT title, content, cnt_like, created_at FROM novel_step1 WHERE seq_novel_step1 = 1;

SELECT ns1.title, ns2.content, ns2.cnt_like, UNIX_TIMESTAMP(ns2.created_at) * 1000 AS created_at
FROM novel_step1 as ns1 INNER JOIN novel_step2 ns2 ON ns1.seq_novel_step1 = ns2.seq_novel_step2
WHERE ns2.seq_novel_step2 = 1;

SELECT ns1.title, ns2.content, ns2.cnt_like, UNIX_TIMESTAMP(ns2.created_at) *1000 AS created_at,
FROM novel_step1 as ns1 INNER JOIN novel_step2 ns2 ON ns1.seq_novel_step1 = ns2.seq_novel_step2
WHERE ns2.seq_novel_step2 = 1;