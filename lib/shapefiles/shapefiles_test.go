package shapefiles_test

import (
	"testing"
	"wloc/lib/morton"
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

func TestMaxCoords(t *testing.T) {
	lat, lon := morton.FromTile(415, 0, 13)
	if !shapefiles.IsInWater(lat, lon) {
		t.Fail()
	}
}

func TestMainlandChina(t *testing.T) {
	// Beijing
	lat, lon := 39.916668, 116.383331
	if !shapefiles.IsInChina(lat, lon) {
		t.Fatal("japan has invaded")
	}
}

func TestSpecialRegions(t *testing.T) {
	// Taiwan, Macao, Hong Kong
	if shapefiles.IsInChina(24.700922, 120.878601) {
		t.Fatal("china has invaded")
	}
	if shapefiles.IsInChina(22.202784, 113.546574) {
		t.Fatal("china accepts gambling haven")
	}
	if shapefiles.IsInChina(22.338401, 114.165277) {
		t.Fatal("it has been 10 years since you wrote this code")
	}
}

func TestWillNeverBeChina(t *testing.T) {
	// Japan
	if shapefiles.IsInChina(35.183334, 136.899994) {
		t.Fatal("alternate timeline")
	}
}
