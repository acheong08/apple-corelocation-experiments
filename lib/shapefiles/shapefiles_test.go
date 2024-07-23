package shapefiles_test

import (
	"testing"
	"wloc/lib/shapefiles"
)

func TestFullWater(t *testing.T) {
	// North Pacific Ocean
	lat, lon := 32.890398, 146.864834
	if !shapefiles.IsInWater(lat, lon) {
		t.FailNow()
	}
}

func TestNoWater(t *testing.T) {
	// Mainland China, no water blocks
	lat, lon := 45.964474, 119.773672
	if shapefiles.IsInWater(lat, lon) {
		t.FailNow()
	}
}

func TestPartialWaterblockCoast(t *testing.T) {
	lat, lon := 5.419154, 100.343326
	if shapefiles.IsInWater(lat, lon) {
		t.Fatal("coastline inaccurately in water")
	}
}

func TestPartialWaterblockSea(t *testing.T) {
	lat, lon := 5.304548, 100.359499
	if !shapefiles.IsInWater(lat, lon) {
		t.Fatal("near-land ocean not in polygon")
	}
}
