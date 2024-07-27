package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"wloc/lib/morton"
)

func main() {
	g := Generator{}
	g.LoadState("state.gob")
	fmt.Println(g.Current.X, g.Current.Y)
	g.Current.Y = 4356
	g.Current.X = 3845
	g.Start()
	fmt.Println(morton.Pack(g.Current.X, g.Current.Y, 13))
	g.SaveState("state.gob")
}

type Coordinate struct {
	X, Y int
}

type Generator struct {
	Current   Coordinate
	MaxTile   int
	ch        chan Coordinate
	done      chan struct{}
	stateLock sync.Mutex
}

func NewGenerator() *Generator {
	return &Generator{
		Current: Coordinate{X: 0, Y: 0},
		MaxTile: 9000,
		ch:      make(chan Coordinate),
		done:    make(chan struct{}),
	}
}

func (g *Generator) Start() {
	go func() {
		defer close(g.ch)
		for {
			select {
			case <-g.done:
				return
			case g.ch <- g.Current:
				g.stateLock.Lock()
				g.Current.X++
				if g.Current.X > g.MaxTile {
					g.Current.X = 0
					g.Current.Y++
					if g.Current.Y > g.MaxTile {
						g.stateLock.Unlock()
						return
					}
				}
				g.stateLock.Unlock()
			}
		}
	}()
}

func (g *Generator) Stop() {
	close(g.done)
}

func (g *Generator) Channel() <-chan Coordinate {
	return g.ch
}

func (g *Generator) SaveState(filename string) error {
	g.stateLock.Lock()
	defer g.stateLock.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(g.Current)
}

func (g *Generator) LoadState(filename string) error {
	g.stateLock.Lock()
	defer g.stateLock.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&g.Current)
}
