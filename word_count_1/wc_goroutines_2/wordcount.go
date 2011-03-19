package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"
	"os"
	"bufio"
	"strings"
)

const cutSet = " \r\n\t.;,:?!-+=|()*[]\\|/<>&^%$#@`~\"'"

var (
	p int // Go MAXPROCS
	b int // Buffer size

	partCounts chan int
	sumOut     chan int
	waitGrp    *sync.WaitGroup
)

var cutSet_map = make(map[int]bool)

func init() {
	for _, rune := range cutSet {
		cutSet_map[rune] = true
	}
}

func cutSet_func(rune int) bool {
	return rune < 256 && cutSet_map[rune]
}

func countWords(line *string) {
	cnt := 0
	for _, word := range strings.Split(*line, " ", -1) {
		if strings.TrimFunc(word, cutSet_func) != "" {
			cnt++
		}
	}
	partCounts <- cnt
	waitGrp.Done()
}

func reduce() {
	sum := 0
	for cnt := 0; cnt >= 0; cnt = <-partCounts {
		sum += cnt
	}
	sumOut <- sum
}

func main() {
	// Parse args
	flag.IntVar(&p, "p", 1, "Max. number of Go processes")
	flag.IntVar(&b, "b", 1, "Channel buffer size")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("No input file name specified\n")
		return
	}
	fname := flag.Arg(0)

	// Go MAXPROCS tweak
	runtime.GOMAXPROCS(p)

	// Set up timer
	runTime := time.Nanoseconds()

	// Initialize channels and run reducer

	//partCounts = make(chan int, b)
	//sumOut = make(chan int, b)
	partCounts = make(chan int)
	sumOut = make(chan int)
	waitGrp = new(sync.WaitGroup)

	go reduce()

	// Open file in buffered mode
	file, err := os.Open(fname, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Input error: %s\n", err)
		return
	}

	fileReader := bufio.NewReader(file)

	// Read data line-by-line and send each line to a separate coroutine
	var line string
	for err == nil {
		line, err = fileReader.ReadString('\n')
		if len(line) > 0 && line != "\n" {
			waitGrp.Add(1)
			go countWords(&line)
		}
	}

	file.Close()

	// Wait for all goroutines to finish
	waitGrp.Wait()

	// Terminate reducer and count sum
	partCounts <- -1
	count := <-sumOut

	// Stop timer and print results
	stopTime := time.Nanoseconds()
	runElapsed := float64(stopTime-runTime) / 1000000000
	fmt.Printf("Done: %d words in %f seconds\n", count, runElapsed)
}

