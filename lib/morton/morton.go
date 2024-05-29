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

func PackAppleCoord(mLat, mLong int) (tileKey int64) {
	return m.Pack(uint64(mLong), uint64(mLat))
}

func UnpackAppleCoord(tileKey int64) (mLat, mLong int) {
	mLongLat := m.Unpack(tileKey)
	return int(mLongLat[1]), int(mLongLat[0])
}
