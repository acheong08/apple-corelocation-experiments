package lib

import (
	"errors"
	"log"
	"wloc/lib/distance"
	"wloc/lib/morton"
	"wloc/lib/spiral"
)

const ErrInvalidInput = "invalid input"

func SearchProximity(lat, long float64, limit uint8, options ...Modifier) ([]distance.Point, error) {
	if options == nil {
		options = make([]Modifier, 0)
	}
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return nil, errors.New(ErrInvalidInput)
	}
	if limit == 0 {
		// We don't want infinite search
		return nil, errors.New(ErrInvalidInput)
	}
	mLat, mLong := morton.Unpack(morton.Encode(lat, long))
	sp := spiral.NewSpiral(mLat, mLong)
	target := distance.Point{
		Y: lat,
		X: long,
	}
	var closest *distance.Point
	for i := 0; i < int(limit); i++ {
		mLat, mLong = sp.Next()
		tile, err := GetTile(morton.Pack(mLat, mLong), options...)
		if err != nil {
			continue
		}
		for _, d := range tile {
			closest = distance.Closer(&target, closest, &distance.Point{
				Id: d.BSSID,
				Y:  d.Location.Lat,
				X:  d.Location.Long,
			})
		}
		break
	}
	if closest == nil {
		return nil, errors.New("no devices found")
	}
	var points []distance.Point
	for {
		devices, err := QueryBssid([]string{closest.Id}, 0, options...)
		if err != nil {
			log.Println(closest)
			return nil, err
		}
		if len(devices) == 0 {
			return nil, errors.New("could not find given BSSID")
		}
		points = make([]distance.Point, len(devices))
		for i, device := range devices {
			points[i] = distance.Point{
				Id: device.BSSID,
				Y:  device.Location.Lat,
				X:  device.Location.Long,
			}
		}
		newClosest := distance.Closest(target, points)
		if newClosest.Id == closest.Id {
			break
		}
		closest = &newClosest
	}
	// ensure closest is #0 in points
	points[0] = *closest // We can kill a point since we have so many
	return points, nil
}
