package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
)

type ParseResult struct {
	HeaderLength int
	BodyLength   uint16
	Header       []byte
	LengthBytes  []byte
	Body         []byte
	Valid        bool
}

func reverseParseHeaderLength(data []byte) *ParseResult {
	var results []ParseResult

	// Try finding length field at different positions
	// Look for 4-byte length fields (common in ARPC)
	for lengthPos := range len(data) - 4 {
		// Try 4-byte big-endian length
		if lengthPos+4 < len(data) {
			lengthBytes := data[lengthPos : lengthPos+4]
			bodyLength := binary.BigEndian.Uint32(lengthBytes)

			// Body starts right after length field
			bodyStart := lengthPos + 4
			expectedEnd := bodyStart + int(bodyLength)

			// Check if this matches the actual data length
			valid := expectedEnd == len(data)

			if valid || bodyLength < 1000 { // Include reasonable lengths even if not perfect match
				result := ParseResult{
					HeaderLength: lengthPos,
					BodyLength:   uint16(bodyLength), // Note: truncating for display
					Valid:        valid,
					LengthBytes:  lengthBytes,
				}

				if lengthPos > 0 {
					result.Header = data[:lengthPos]
				}

				if bodyStart < len(data) {
					result.Body = data[bodyStart:]
				}

				results = append(results, result)
			}
		}

		// Also try 2-byte length fields
		if lengthPos+2 < len(data) {
			lengthBytes := data[lengthPos : lengthPos+2]
			bodyLength := binary.BigEndian.Uint16(lengthBytes)

			bodyStart := lengthPos + 2
			expectedEnd := bodyStart + int(bodyLength)

			valid := expectedEnd == len(data)

			if valid || bodyLength < 1000 {
				result := ParseResult{
					HeaderLength: lengthPos,
					BodyLength:   bodyLength,
					Valid:        valid,
					LengthBytes:  lengthBytes,
				}

				if lengthPos > 0 {
					result.Header = data[:lengthPos]
				}

				if bodyStart < len(data) {
					result.Body = data[bodyStart:]
				}

				results = append(results, result)
			}
		}
	}

	return dedupeResults(results)
}

func dedupeResults(results []ParseResult) *ParseResult {
	if len(results) == 0 {
		return nil
	}
	bestResult := ParseResult{}
	for i := range results {
		if results[i].Valid && results[i].BodyLength > bestResult.BodyLength || (results[i].BodyLength == bestResult.BodyLength && results[i].HeaderLength > bestResult.HeaderLength) {
			bestResult = results[i]
		}
	}

	return &bestResult
}

func printResults(result *ParseResult) {
	fmt.Printf("  Header Length: %d bytes\n", result.HeaderLength)
	fmt.Printf("  Body Length: %d bytes (0x%04x)\n", result.BodyLength, result.BodyLength)

	if len(result.Header) > 0 {
		fmt.Printf("  Header: %s\n", hex.EncodeToString(result.Header))
	}
	fmt.Printf("  Length Bytes: %s\n", hex.EncodeToString(result.LengthBytes))
	if len(result.Body) > 0 && len(result.Body) <= 100 {
		fmt.Printf("  Body: %s\n", hex.EncodeToString(result.Body))
	} else if len(result.Body) > 100 {
		fmt.Printf("  Body: %s (%d bytes total)\n", hex.EncodeToString(result.Body), len(result.Body))
	}
	fmt.Println()
}

func main() {
	var (
		hexData = flag.String("hex", "", "Hex string to parse")
		file    = flag.String("file", "", "File containing binary data")
	)
	flag.Parse()

	var data []byte
	var err error

	if *hexData != "" {
		data, err = hex.DecodeString(*hexData)
		if err != nil {
			log.Fatalf("Failed to decode hex string: %v", err)
		}
	} else if *file != "" {
		data, err = os.ReadFile(*file)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
	} else {
		log.Fatal("Must provide either -hex or -file argument")
	}

	if len(data) < 3 {
		log.Fatal("Data must be at least 3 bytes long")
	}

	fmt.Printf("Analyzing %d bytes of data...\n", len(data))

	results := reverseParseHeaderLength(data)
	if results != nil {
		fmt.Println("✅ Found valid parsing configuration(s)")
		printResults(results)
	} else {
		fmt.Printf("❌ No valid parsing configurations found\n")
	}
}
