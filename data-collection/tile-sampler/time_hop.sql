WITH gaps AS (
  SELECT
    id,
    LAG(id) OVER (ORDER BY substr(change_time, 1, 23)) AS prev_id,
    change_time,
    LAG(change_time) OVER (ORDER BY substr(change_time, 1, 23)) AS prev_change_time
  FROM location_changes
)
SELECT lc.*
FROM location_changes lc
JOIN (
  SELECT id FROM gaps
  WHERE prev_change_time IS NOT NULL
    AND (strftime('%s', substr(change_time, 1, 23)) - strftime('%s', substr(prev_change_time, 1, 23))) > 60
  UNION
  SELECT prev_id FROM gaps
  WHERE prev_change_time IS NOT NULL
    AND (strftime('%s', substr(change_time, 1, 23)) - strftime('%s', substr(prev_change_time, 1, 23))) > 60
) gap_ids
ON lc.id = gap_ids.id
ORDER BY lc.id;
