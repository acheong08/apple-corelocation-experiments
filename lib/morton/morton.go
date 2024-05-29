package morton

import (
	"math"

	morton "github.com/gojuno/go.morton"
)

var m morton.Morton64 = *morton.Make64(2, 32)

func Decode(tileKey int64) (lat float64, long float64) {
	mLongLat := m.Unpack(tileKey)
	long, _ = mortonToGpsLong.Predict([]float64{float64(mLongLat[0])})
	lat, _ = mortonToGpsLat.Predict([]float64{float64(mLongLat[1])})
	return lat, long
}

func Encode(lat float64, long float64) (tileKey int64) {
	mLat, _ := gpsToMortonLat.Predict([]float64{lat})
	mLong, _ := gpsToMortonLong.Predict([]float64{long})
	tileKey = m.Pack(uint64(math.Round(mLong)), uint64(math.Round(mLat)))
	return tileKey
}

func PredictAppleCoord(lat float64, long float64) (mLat, mLong int) {
	fMLat, _ := gpsToMortonLat.Predict([]float64{lat})
	fMLong, _ := gpsToMortonLong.Predict([]float64{long})
	return int(math.Round(fMLat)), int(math.Round(fMLong))
}

func Pack(mLat, mLong int) (tileKey int64) {
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

func Unpack(tileKey int64) (mLat, mLong int) {
	level := 0
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
	return row, column
}
