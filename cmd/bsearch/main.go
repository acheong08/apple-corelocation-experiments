package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"wloc/lib/mac"

	"github.com/schollz/progressbar/v3"
)

func main() {
	if _, err := os.Stat(sortedName); os.IsNotExist(err) {
		if _, err := os.Stat(inputName); os.IsNotExist(err) {
			panic("no input files found")
		}
		sortBin()
	}
	if len(os.Args) != 2 {
		log.Fatal("Usage: bsearch <BSSID>")
	}
	m, err := mac.Encode(os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Println(m)
	_, geocode, err := binarySearch(m)
	if err != nil {
		panic(err)
	}
	fmt.Println("Geocode: ", geocode)
}

const (
	recordSize   = 16
	int64Size    = 8
	chunkSize    = 1000000         // Number of records per chunk
	maxOpenFiles = 100             // Maximum number of files to merge at once
	bufferSize   = 4 * 1024 * 1024 // 4MB buffer for reading
	sortedName   = "sortedaps.bin"
	inputName    = "wtfps.bin"
)

func binarySearch(searchKey int64) (int64, int64, error) {
	file, err := os.Open(sortedName)
	if err != nil {
		return 0, 0, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, 0, fmt.Errorf("error getting file info: %v", err)
	}

	totalRecords := fileInfo.Size() / int64(recordSize)

	low, high := int64(0), totalRecords-1

	for low <= high {
		mid := (low + high) / 2
		offset := mid * int64(recordSize)

		_, err = file.Seek(offset, 0)
		if err != nil {
			return 0, 0, fmt.Errorf("error seeking in file: %v", err)
		}

		record := make([]byte, recordSize)
		_, err = io.ReadFull(file, record)
		if err != nil {
			return 0, 0, fmt.Errorf("error reading record: %v", err)
		}

		key := int64(binary.BigEndian.Uint64(record[:8]))
		value := int64(binary.BigEndian.Uint64(record[8:]))

		if key == searchKey {
			return key, value, nil
		} else if key < searchKey {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return 0, 0, fmt.Errorf("key not found")
}

type Record struct {
	Key   int64
	Value int64
}

func sortBin() {
	inputFile := inputName
	outputFile := sortedName

	// Get the total number of records
	totalRecords, err := getTotalRecords(inputFile)
	if err != nil {
		fmt.Printf("Error getting total records: %v\n", err)
		return
	}

	fmt.Printf("Total records: %d\n", totalRecords)

	// Step 1: Split the input file into sorted chunks
	fmt.Println("Splitting and sorting chunks...")
	bar := progressbar.Default(int64(totalRecords))
	chunkFiles, err := splitAndSortChunks(inputFile, bar)
	if err != nil {
		fmt.Printf("Error splitting and sorting chunks: %v\n", err)
		return
	}

	// Step 2: Merge the sorted chunks
	fmt.Println("\nMerging chunks...")
	bar = progressbar.Default(int64(totalRecords))
	err = mergeChunks(chunkFiles, outputFile, bar)
	if err != nil {
		fmt.Printf("Error merging chunks: %v\n", err)
		return
	}

	// Clean up temporary chunk files
	for _, file := range chunkFiles {
		os.Remove(file)
	}

	fmt.Println("\nSorting completed successfully.")
}

func getTotalRecords(filename string) (int, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return int(fileInfo.Size()) / recordSize, nil
}

func splitAndSortChunks(inputFile string, bar *progressbar.ProgressBar) ([]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunkFiles []string
	chunkNum := 0
	buffer := make([]byte, bufferSize)
	records := make([]Record, 0, chunkSize)

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for i := 0; i < bytesRead; i += recordSize {
			if len(records) >= chunkSize {
				// Sort and write the current chunk
				sort.Slice(records, func(i, j int) bool {
					return records[i].Key < records[j].Key
				})

				chunkFile := fmt.Sprintf("chunk_%d.tmp", chunkNum)
				err = writeRecords(chunkFile, records)
				if err != nil {
					return nil, err
				}

				chunkFiles = append(chunkFiles, chunkFile)
				chunkNum++
				records = records[:0] // Clear the slice
			}

			if i+recordSize <= bytesRead {
				record := Record{
					Key:   int64(binary.BigEndian.Uint64(buffer[i : i+8])),
					Value: int64(binary.BigEndian.Uint64(buffer[i+8 : i+16])),
				}
				records = append(records, record)
				_ = bar.Add(1)
			}
		}

		if err == io.EOF {
			break
		}
	}

	// Handle the last chunk if there are any remaining records
	if len(records) > 0 {
		sort.Slice(records, func(i, j int) bool {
			return records[i].Key < records[j].Key
		})

		chunkFile := fmt.Sprintf("chunk_%d.tmp", chunkNum)
		err = writeRecords(chunkFile, records)
		if err != nil {
			return nil, err
		}

		chunkFiles = append(chunkFiles, chunkFile)
	}

	return chunkFiles, nil
}

func mergeChunks(chunkFiles []string, outputFile string, bar *progressbar.ProgressBar) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for len(chunkFiles) > 0 {
		// Merge up to maxOpenFiles chunks at a time
		numFiles := min(len(chunkFiles), maxOpenFiles)
		err := mergeNWay(chunkFiles[:numFiles], out, bar)
		if err != nil {
			return err
		}
		chunkFiles = chunkFiles[numFiles:]
	}

	return nil
}

func mergeNWay(inputFiles []string, out *os.File, bar *progressbar.ProgressBar) error {
	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files provided")
	}

	files := make([]*os.File, len(inputFiles))
	records := make([]Record, len(inputFiles))
	valid := make([]bool, len(inputFiles))
	buffers := make([][]byte, len(inputFiles))

	for i, filename := range inputFiles {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", filename, err)
		}
		defer f.Close()
		files[i] = f

		// Initialize buffer with actual data
		buffers[i] = make([]byte, 0, bufferSize)
		n, err := f.Read(buffers[i][:cap(buffers[i])])
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading from file %s: %v", filename, err)
		}
		buffers[i] = buffers[i][:n]

		if len(buffers[i]) >= recordSize {
			records[i] = Record{
				Key:   int64(binary.BigEndian.Uint64(buffers[i][:8])),
				Value: int64(binary.BigEndian.Uint64(buffers[i][8:16])),
			}
			valid[i] = true
			buffers[i] = buffers[i][recordSize:]
		}
	}

	writeBuffer := make([]byte, bufferSize)
	writeBufferPos := 0

	for {
		minIdx := -1
		for i, v := range valid {
			if v && (minIdx == -1 || records[i].Key < records[minIdx].Key) {
				minIdx = i
			}
		}

		if minIdx == -1 {
			break // All files exhausted
		}

		// Write the smallest record to the buffer
		binary.BigEndian.PutUint64(writeBuffer[writeBufferPos:], uint64(records[minIdx].Key))
		binary.BigEndian.PutUint64(writeBuffer[writeBufferPos+8:], uint64(records[minIdx].Value))
		writeBufferPos += recordSize

		if writeBufferPos == bufferSize {
			// Write the full buffer to the output file
			_, err := out.Write(writeBuffer)
			if err != nil {
				return fmt.Errorf("error writing to output file: %v", err)
			}
			writeBufferPos = 0
		}

		_ = bar.Add(1)

		// Read the next record from the file we just used
		if len(buffers[minIdx]) >= recordSize {
			records[minIdx] = Record{
				Key:   int64(binary.BigEndian.Uint64(buffers[minIdx][:8])),
				Value: int64(binary.BigEndian.Uint64(buffers[minIdx][8:16])),
			}
			buffers[minIdx] = buffers[minIdx][recordSize:]
		} else {
			// Refill the buffer
			buffers[minIdx] = buffers[minIdx][:cap(buffers[minIdx])]
			n, err := files[minIdx].Read(buffers[minIdx])
			if err != nil && err != io.EOF {
				return fmt.Errorf("error reading from file %s: %v", inputFiles[minIdx], err)
			}
			buffers[minIdx] = buffers[minIdx][:n]

			if len(buffers[minIdx]) >= recordSize {
				records[minIdx] = Record{
					Key:   int64(binary.BigEndian.Uint64(buffers[minIdx][:8])),
					Value: int64(binary.BigEndian.Uint64(buffers[minIdx][8:16])),
				}
				buffers[minIdx] = buffers[minIdx][recordSize:]
			} else {
				valid[minIdx] = false
			}
		}
	}

	// Write any remaining records in the write buffer
	if writeBufferPos > 0 {
		_, err := out.Write(writeBuffer[:writeBufferPos])
		if err != nil {
			return fmt.Errorf("error writing final buffer to output file: %v", err)
		}
	}

	return nil
}

func writeRecords(filename string, records []Record) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)
	bufferPos := 0

	for _, record := range records {
		binary.BigEndian.PutUint64(buffer[bufferPos:], uint64(record.Key))
		binary.BigEndian.PutUint64(buffer[bufferPos+8:], uint64(record.Value))
		bufferPos += recordSize

		if bufferPos == bufferSize {
			_, err := file.Write(buffer)
			if err != nil {
				return err
			}
			bufferPos = 0
		}
	}

	// Write any remaining records in the buffer
	if bufferPos > 0 {
		_, err := file.Write(buffer[:bufferPos])
		if err != nil {
			return err
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
