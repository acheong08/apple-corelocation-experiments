package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gojuno/go.morton"
)

func main() {
	files, err := os.ReadDir("./points")
	if err != nil {
		panic(err)
	}
	earthMorton := morton.Make64(2, 32)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			fileName := file.Name()[0 : len(file.Name())-4]
			mortonCode, _ := strconv.Atoi(fileName)
			log.Println(mortonCode, earthMorton.Unpack(int64(mortonCode)))
		}
	}
}
