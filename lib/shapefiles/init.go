package shapefiles

import (
	"bytes"
	_ "embed"
	"encoding/gob"

	"github.com/paulmach/orb"
)

//go:embed assets/china.orb
var _china []byte

//go:embed assets/water.morb
var _waters []byte

var (
	China  []orb.Polygon
	Waters map[int64][]orb.Polygon
)

func init() {
	dec := gob.NewDecoder(bytes.NewReader(_china))
	if err := dec.Decode(&China); err != nil {
		panic(err)
	}
	dec = gob.NewDecoder(bytes.NewReader(_waters))
	if err := dec.Decode(&Waters); err != nil {
		panic(err)
	}
}
