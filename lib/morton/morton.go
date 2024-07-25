package morton

import (
	"github.com/buckhx/tiles"
)

func Decode(tileKey int64) (lat float64, long float64, level int) {
	mLat, mLong, level := Unpack(tileKey)
	lat, long = FromTile(mLat, mLong, level)
	return
}

func Encode(lat float64, long float64, level int) (tileKey int64) {
	mLat, mLong := ToTile(lat, long, level)
	tileKey = Pack(mLat, mLong, level)
	return tileKey
}

func ToTile(lat, long float64, level int) (mLat, mLong int) {
	t := tiles.FromCoordinate(lat, long, level)
	p := t.ToPixel()
	t2, _ := p.ToTile()
	return t2.Y, t2.X
}

func FromTile(mLat, mLong, level int) (lat, long float64) {
	coords := tiles.Tile{Y: mLat, X: mLong, Z: level}.ToPixel().ToCoords()
	return coords.Lat, coords.Lon
}

func Pack(mLat, mLong, level int) (tileKey int64) {
	row := mLat
	column := mLong
	result := int64(powerOfTwo[level<<1])
	for i := 0; i < level; i++ {
		if column&0x1 != 0 {
			result += int64(powerOfTwo[2*i])
		}
		if row&0x1 != 0 {
			result += int64(powerOfTwo[2*i+1])
		}
		column = column >> 1
		row = row >> 1
	}
	return result
}

func Unpack(tileKey int64) (mLat, mLong, level int) {
	row := 0
	column := 0
	quadKey := tileKey
	for quadKey > 1 {
		mask := 1 << level
		if quadKey&0x1 != 0 {
			column |= mask
		}
		if quadKey&0x2 != 0 {
			row |= mask
		}
		level++
		quadKey = (quadKey - (quadKey & 0x3)) / 4
	}
	return row, column, level
}
