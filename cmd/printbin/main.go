package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"wloc/pb"

	"google.golang.org/protobuf/proto"
)

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
	var wloc pb.AppleWLoc
	// Loop through removing starting bytes until it works
	i := 0
	for i = 0; i < len(b); i += 1 {
		err = proto.Unmarshal(b[i:], &wloc)
		if err == nil {
			break
		}
	}
	log.Println("Num removed before valid: ", i)
	log.Printf("Removed: %x\n", b[:i])
	log.Println("Content length: ", len(b[i:]))
	if slices.Contains(os.Args, "-json") {
		j, _ := json.MarshalIndent(&wloc, "", " ")
		fmt.Println(string(j))
	}
	if slices.Contains(os.Args, "-hex") {
		fmt.Printf("%x\n", b[i:])
	}
}
