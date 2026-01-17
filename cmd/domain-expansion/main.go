package main

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"github.com/acheong08/apple-corelocation-experiments/lib"
	"github.com/acheong08/apple-corelocation-experiments/lib/mac"
	"github.com/acheong08/apple-corelocation-experiments/lib/morton"

	"github.com/DataDog/zstd"
	"github.com/tidwall/btree"
	_ "modernc.org/sqlite"
)

const NUM_THREADS = 8

type BeSet struct {
	btree.Set[int64]
	sync.RWMutex
}

var explored BeSet

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := sql.Open("sqlite", "seeds.db")
	if err != nil {
		panic(err)
	}
	// Start thread for database fetching (with buffered channel streaming bssids to crawl)
	log.Println("Starting database stream")
	seeds := databaseStream(ctx, db)
	log.Println("Starting writer")
	// Use separate cancellation as this must be done after all threads are done
	writerCtx, writerCancel := context.WithCancel(ctx)
	defer writerCancel()
	// go zstApWriter(ctx)
	go sqliteApWriter(writerCtx)
	threadCtx, threadCancel := context.WithCancel(ctx)
	// Start threads to process and explore the bssids
	wait := sync.WaitGroup{}
	for i := range NUM_THREADS {
		wait.Add(1)
		go func() {
			defer wait.Done()
			log.Println("Starting thread: #", i)
			exploreArea(threadCtx, seeds)
		}()
	}
	log.Println("Waiting for exit")
	die := make(chan os.Signal, 1)
	signal.Notify(die, os.Interrupt, syscall.SIGTERM)
	<-die
	threadCancel()
	wait.Wait()
}

type Seed struct {
	bssid   int64
	tilekey int64
}

const DB_BATCHSIZE = 128

func databaseStream(ctx context.Context, db *sql.DB) <-chan Seed {
	c := make(chan Seed, 512)
	// Select bssid group by tilekey
	go func() {
		offset := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				rows, err := db.Query("SELECT bssid, tilekey FROM seeds GROUP BY tilekey LIMIT ? OFFSET ?", DB_BATCHSIZE, offset*DB_BATCHSIZE)
				if err != nil {
					log.Println(err)
					return
				}

				offset++
				for rows.Next() {

					var bssid int64
					var tilekey int64
					err = rows.Scan(&bssid, &tilekey)
					if err != nil {
						panic(err)
					}
					c <- Seed{bssid, tilekey}
				}
				rows.Close()
				log.Println("Streamed 1 batch")
			}
		}
	}()
	return c
}

const SAVE_FILE = "beacons.bin.zst"

// 6 bytes for bssid, 16 for lat/lon = 22
const RECORD_SIZE = 6 + 8 + 8

type Record struct {
	bssid    int64
	lat, lon float64
}

var writeCh = make(chan Record)

func sqliteApWriter(ctx context.Context) {
	db, err := sql.Open("sqlite", "beacons.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS beacons (bssid INTEGER PRIMARY KEY, lat REAL NOT NULL, lon REAL NOT NULL)"); err != nil {
		panic(err)
	}
	bulkInsert := func(unsavedRows []Record) error {
		valueStrings := make([]string, len(unsavedRows))
		valueArgs := make([]interface{}, len(unsavedRows)*3)
		for i, post := range unsavedRows {
			valueStrings[i] = "(?,?,?)"
			valueArgs[i*3] = post.bssid
			valueArgs[i*3+1] = post.lat
			valueArgs[i*3+2] = post.lon
		}
		stmt := fmt.Sprintf("INSERT OR IGNORE INTO beacons (bssid, lat, lon) VALUES %s",
			strings.Join(valueStrings, ","))
		_, err := db.Exec(stmt, valueArgs...)
		return err
	}
	records := make([]Record, 400)
	n := 0
	for record := range writeCh {
		select {
		case <-ctx.Done():
			close(writeCh)
			if err := bulkInsert(records); err != nil {
				panic(err)
			}
			return
		default:
			if n == len(records) {
				if err := bulkInsert(records); err != nil {
					panic(err)
				}
				records = make([]Record, 400*NUM_THREADS)
				n = 0
				log.Println("Written to DB")
			}
			records[n] = record
			n++
		}
	}
}

const BUFFER_SIZE = RECORD_SIZE * 400 * NUM_THREADS

func zstApWriter(ctx context.Context) {
	var f *os.File
	var err error
	if _, err = os.Stat(SAVE_FILE); os.IsNotExist(err) {
		f, err = os.Create(SAVE_FILE)
	} else {
		f, err = os.Open(SAVE_FILE)
	}
	if err != nil {
		panic(err)
	}
	zf := zstd.NewWriter(f)
	defer zf.Close()
	var buffer []byte = make([]byte, BUFFER_SIZE)
	n := 0
	for record := range writeCh {
		select {
		case <-ctx.Done():
		default:
			if (n * RECORD_SIZE) == len(buffer) {
				_, err := zf.Write(buffer)
				if err != nil {
					panic(err)
				}
				buffer = make([]byte, BUFFER_SIZE)
				n = 0

				log.Println("Buffer written")
			}
			bRecord := make([]byte, RECORD_SIZE)
			bbssid := make([]byte, 8)
			binary.BigEndian.PutUint64(bbssid, uint64(record.bssid))
			copy(bRecord[0:6], bbssid[2:])
			binary.BigEndian.PutUint64(bRecord[6:14], math.Float64bits(record.lat))
			binary.BigEndian.PutUint64(bRecord[14:22], math.Float64bits(record.lon))
			copy(buffer[RECORD_SIZE*n:RECORD_SIZE*(n+1)], bRecord)
			n++
		}
	}
}

const MORTON_LEVEL = 20

func apsToMap(a []lib.AP, b map[int64]int64) {
	for _, ap := range a {
		bssid, _ := mac.Encode(ap.BSSID)
		code := morton.Encode(ap.Location.Lat, ap.Location.Long, MORTON_LEVEL)
		writeCh <- Record{
			lat:   ap.Location.Lat,
			lon:   ap.Location.Long,
			bssid: bssid,
		}
		explored.RLock()
		if explored.Contains(code) {
			explored.RUnlock()
			continue
		}
		explored.RUnlock()
		b[code] = bssid
		explored.Lock()
		explored.Insert(code)
		explored.Unlock()
	}
}

func getRandomFromMap[T comparable, V any](m map[T]V) (T, V, bool) {
	for t, v := range m {
		delete(m, t)
		return t, v, true
	}
	var dt T
	var dv V
	return dt, dv, false
}

func exploreArea(ctx context.Context, c <-chan Seed) {
	// Loop through channel
	for seed := range c {
		log.Println("Crawling seed: ", seed.tilekey)
		// Tilekey to bssid
		toExplore := make(map[int64]int64)
		// Populate via seed
		aps, err := lib.QueryBssid([]string{mac.Decode(seed.bssid)}, 0)
		if err != nil {
			log.Println("Failed to query bssid: ", err)
			continue
		}
		apsToMap(aps, toExplore)
	explorationLoop:
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Explore seed and keep adding to explored until there's nothing left
				if len(toExplore) == 0 {
					log.Println("broke")
					break explorationLoop
				}
				tilekey, bssid, ok := getRandomFromMap(toExplore)
				if !ok {
					break explorationLoop
				}
				delete(toExplore, tilekey)

				aps, err := lib.QueryBssid([]string{mac.Decode(bssid)}, 0)
				if err != nil {
					log.Println("Failed to query bssid: ", err)
					continue
				}
				apsToMap(aps, toExplore)
			}
		}
	}
}
