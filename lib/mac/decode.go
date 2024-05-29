package mac

import "fmt"

func Decode(i int64) string {
	macHex := fmt.Sprintf("%x", i)
	if len(macHex) != 12 {
		// Fill it up to 12 with 0s in front
		for i := 0; i < 13-len(macHex); i++ {
			macHex = "0" + macHex
		}
	}
	// Insert : between every 2 hex values
	mac := ""
	for i := 0; i < len(macHex); i += 2 {
		if i+2 < len(macHex) {
			mac += macHex[i:i+2] + ":"
		} else {
			mac += macHex[i:]
		}
	}
	return mac
}
