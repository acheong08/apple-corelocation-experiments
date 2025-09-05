package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"wloc/lib"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func decodeGenericProtobuf(data []byte) (map[string]any, error) {
	msg := &dynamicpb.Message{}

	decoder := proto.UnmarshalOptions{
		AllowPartial: true,
	}

	if data == nil {
		panic("data is nil")
	}
	err := decoder.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)

	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fieldNum := int(fd.Number())
		fieldName := fmt.Sprintf("field_%d", fieldNum)

		switch fd.Kind() {
		case protoreflect.StringKind:
			result[fieldName] = v.String()
		case protoreflect.BytesKind:
			result[fieldName] = fmt.Sprintf("bytes[%d]: %x", len(v.Bytes()), v.Bytes())
		case protoreflect.Int32Kind, protoreflect.Int64Kind:
			result[fieldName] = v.Int()
		case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			result[fieldName] = v.Uint()
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			result[fieldName] = v.Float()
		case protoreflect.BoolKind:
			result[fieldName] = v.Bool()
		case protoreflect.MessageKind:
			if submsg, ok := v.Message().(*dynamicpb.Message); ok {
				subdata, err := proto.Marshal(submsg)
				if err == nil {
					if subresult, err := decodeGenericProtobuf(subdata); err == nil {
						result[fieldName] = subresult
					} else {
						result[fieldName] = fmt.Sprintf("nested_message[%d_bytes]", proto.Size(submsg))
					}
				}
			}
		default:
			result[fieldName] = fmt.Sprintf("unknown_type_%s", fd.Kind())
		}
		return true
	})

	return result, nil
}

func tryDecodeProtobuf(data []byte) {
	if result, err := decodeGenericProtobuf(data); err == nil {
		fmt.Println("=== Generic Protobuf Structure ===")
		j, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(j))
	} else {
		fmt.Printf("Failed to decode as protobuf: %v\n", err)
		fmt.Printf("Raw hex: %x\n", data)
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

	arpcData, err := lib.ParseArpc(b)
	if err != nil {
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

	fmt.Println("\n=== Payload Analysis ===")
	tryDecodeProtobuf(arpcData.Payload)
}
