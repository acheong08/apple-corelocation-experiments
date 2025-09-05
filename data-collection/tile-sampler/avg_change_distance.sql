SELECT
  AVG(distance) AS avg_distance,
  (
    SELECT AVG(distance)
    FROM (
      SELECT distance
      FROM (
        SELECT
          6371000 * 2 * 
          ASIN(
            SQRT(
              POWER(SIN(RADIANS((new_lat - old_lat) / 2)), 2) +
              COS(RADIANS(old_lat)) * COS(RADIANS(new_lat)) *
              POWER(SIN(RADIANS((new_long - old_long) / 2)), 2)
            )
          ) AS distance
        FROM location_changes
      )
      ORDER BY distance
      LIMIT 2 - (SELECT COUNT(*) FROM location_changes) % 2
      OFFSET (SELECT (COUNT(*) - 1) / 2 FROM location_changes)
    )
  ) AS median_distance,
  MAX(distance) AS max_distance,
  MIN(distance) AS min_distance
FROM (
  SELECT
    6371000 * 2 * 
    ASIN(
      SQRT(
        POWER(SIN(RADIANS((new_lat - old_lat) / 2)), 2) +
        COS(RADIANS(old_lat)) * COS(RADIANS(new_lat)) *
        POWER(SIN(RADIANS((new_long - old_long) / 2)), 2)
      )
    ) AS distance
  FROM location_changes
)
