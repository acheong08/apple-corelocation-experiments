WITH ordered AS (
  SELECT tilekey FROM seeds GROUP BY tilekey ORDER BY tilekey
),
groups AS (
  SELECT t1.tilekey AS group_start
  FROM ordered t1
  JOIN ordered t2 ON t2.tilekey = t1.tilekey + 1
  JOIN ordered t3 ON t3.tilekey = t1.tilekey + 2
  JOIN ordered t4 ON t4.tilekey = t1.tilekey + 3
  JOIN ordered t5 ON t5.tilekey = t1.tilekey + 4
),
sampled AS (
  SELECT group_start, ROW_NUMBER() OVER (ORDER BY group_start) AS rn, COUNT(*) OVER () AS total
  FROM groups
),
chosen_groups AS (
  SELECT group_start
  FROM sampled
  WHERE rn IN (
    SELECT CAST((total - 1) * (n / 9.0) + 1 AS INTEGER)
    FROM (SELECT 0 AS n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
          UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9)
  )
)
SELECT group_start + offset AS tilekey
FROM chosen_groups, (SELECT 0 AS offset UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4)
ORDER BY tilekey
