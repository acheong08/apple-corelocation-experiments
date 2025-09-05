package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"wloc/lib"

	"google.golang.org/protobuf/encoding/protowire"
)

func decodeRawProtobuf(data []byte) (map[string]any, error) {
	result := make(map[string]any)

	for len(data) > 0 {
		fieldNum, wireType, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, fmt.Errorf("invalid tag")
		}
		data = data[n:]

		fieldName := fmt.Sprintf("%d", fieldNum)

		var value any

		switch wireType {
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(data)
			if n < 0 {
				return nil, fmt.Errorf("invalid varint")
			}
			value = v
			data = data[n:]

		case protowire.Fixed64Type:
			v, n := protowire.ConsumeFixed64(data)
			if n < 0 {
				return nil, fmt.Errorf("invalid fixed64")
			}
			value = v
			data = data[n:]

		case protowire.BytesType:
			v, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return nil, fmt.Errorf("invalid bytes")
			}

			// Try to decode as nested message
			if nested, err := decodeRawProtobuf(v); err == nil && len(nested) > 0 {
				value = nested
			} else {
				// Treat as string if valid UTF-8, otherwise as bytes
				if isValidUTF8(v) {
					value = string(v)
				} else {
					value = fmt.Sprintf("bytes[%d]: %x", len(v), v)
				}
			}
			data = data[n:]

		case protowire.Fixed32Type:
			v, n := protowire.ConsumeFixed32(data)
			if n < 0 {
				return nil, fmt.Errorf("invalid fixed32")
			}
			value = v
			data = data[n:]

		default:
			return nil, fmt.Errorf("unknown wire type: %d", wireType)
		}

		// Handle repeated fields by creating arrays
		if existing, exists := result[fieldName]; exists {
			// Convert to array if not already
			if arr, isArray := existing.([]any); isArray {
				result[fieldName] = append(arr, value)
			} else {
				result[fieldName] = []any{existing, value}
			}
		} else {
			result[fieldName] = value
		}
	}

	return result, nil
}

func isValidUTF8(data []byte) bool {
	for _, b := range data {
		if b < 32 && b != 9 && b != 10 && b != 13 {
			return false
		}
		if b > 126 {
			return false
		}
	}
	return true
}

func tryDecodeProtobuf(data []byte) {
	if result, err := decodeRawProtobuf(data); err == nil {
		fmt.Println("=== Raw Protobuf Structure ===")
		j, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(j))
	} else {
		fmt.Printf("Failed to decode as protobuf: %v\n", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: printbin <file>")
	}

	filePath := os.Args[1]
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	arpcData := lib.ArpcRequest{}
	if err := arpcData.Deserialize(b); err != nil {
		log.Printf("Failed to parse as ARPC: %v\n", err)
		log.Println("Trying direct protobuf decode...")
		tryDecodeProtobuf(b)
		return
	}

	fmt.Println("=== ARPC Wrapper ===")
	fmt.Printf("Version: %s\n", arpcData.Version)
	fmt.Printf("Locale: %s\n", arpcData.Locale)
	fmt.Printf("App Identifier: %s\n", arpcData.AppIdentifier)
	fmt.Printf("OS Version: %s\n", arpcData.OsVersion)
	fmt.Printf("Function ID: %d\n", arpcData.FunctionId)
	fmt.Printf("Payload Length: %d bytes\n", len(arpcData.Payload))

	if slices.Contains(os.Args, "-hex") {
		fmt.Printf("Payload hex: %x\n", arpcData.Payload)
	}

	if slices.Contains(os.Args, "-proto") {
		fmt.Println("\n=== Payload Analysis ===")
		tryDecodeProtobuf(arpcData.Payload)
	}
}
