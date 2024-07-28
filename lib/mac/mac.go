package mac

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

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

func Encode(mac string) (int64, error) {
	// Remove colons from the MAC address
	macHex := strings.ReplaceAll(mac, ":", "")

	// Ensure the MAC address is valid
	if len(macHex) != 12 {
		return 0, fmt.Errorf("invalid MAC address length")
	}

	// Convert hexadecimal string to int64
	b, err := hex.DecodeString(macHex)
	if err != nil {
		return 0, fmt.Errorf("invalid hex")
	}
	// Pad the start with 0s
	if len(b) != 8 {
		for i := 0; i <= 8-len(b); i++ {
			b = append([]byte{0}, b...)
		}
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

func BytesToInt64(mac []byte) int64 {
	var i int64
	for _, b := range mac {
		i = i<<8 + int64(b)
	}
	return i
}
