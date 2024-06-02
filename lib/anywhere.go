package lib

import (
	"errors"
	"log"
	"math"
	"wloc/lib/distance"
	"wloc/lib/mac"
	"wloc/lib/morton"
	"wloc/lib/spiral"
)

const ErrInvalidInput = "invalid input"

func SearchProximity(lat, long float64, limit uint8) ([]distance.Point, error) {
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
		tile, err := GetTile(morton.Pack(mLat, mLong))
		if err != nil {
			continue
		}
		for _, r := range tile.GetRegion() {
			for _, d := range r.GetDevices() {
				if d == nil || d.GetBssid() == 0 {
					continue
				}
				closest = distance.Closer(&target, closest, &distance.Point{
					Id: mac.Decode(d.GetBssid()),
					Y:  float64(d.GetEntry().GetLat()) * math.Pow10(-7),
					X:  float64(d.GetEntry().GetLong()) * math.Pow10(-7),
				})
			}
		}
		break
	}
	var points []distance.Point
	for {
		devices, err := QueryBssid([]string{closest.Id}, true)
		if err != nil {
			log.Println(closest)
			return nil, err
		}
		if len(devices.GetWifiDevices()) == 0 {
			return nil, errors.New("could not find given BSSID")
		}
		points = make([]distance.Point, len(devices.GetWifiDevices()))
		for i, device := range devices.GetWifiDevices() {
			points[i] = distance.Point{
				Id: device.GetBssid(),
				Y:  float64(*device.GetLocation().Latitude) * math.Pow10(-8),
				X:  float64(*device.GetLocation().Longitude) * math.Pow10(-8),
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
