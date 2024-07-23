# Shapefiles

This is used to determine if a point is in water or China. We should not be fetching tiles in water and a specific API is needed for Chinese points.

- `.orb`: `[]orb.Polygon`
- `.morb`: `map[int64]orb.Polygon` - The key is level 9 morton encoded coordinates. This is used for efficient checking of whether a polygon exists
