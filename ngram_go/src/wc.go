package main


import (
	"fmt"
	"os"
	"bufio"
	"maps"
	"runtime"
	"sync"
	"strings"
	"flag"
	"time"
)

var (
	c1       chan string
	c2       chan string
	c3       chan string
	c4       chan string
	lineChan chan string

	m1 map[string]uint32
	m2 maps.SMap
	m3 maps.SMap
	m4 maps.SMap

	p int // Go MAXPROCS
	b int // Buffer size

	partCounts chan int
	sumOut     chan int
	waitGrp    *sync.WaitGroup

	debug bool = true

	wordCounter = 0
)


//--------
const cutSet = " \r\n\t.;,:?!+=|()*[]\\|/<>&^%$#@`~\"'"


var cutSet_map = make(map[int]bool)

func init() {
	for _, rune := range cutSet {
		cutSet_map[rune] = true
	}
}

func cutSet_func(rune int) bool {
	return rune < 256 && cutSet_map[rune]
}

func splitLine() {
	for {
		line := <-lineChan
		for _, word := range strings.Split(line, " ", -1) {
			if strings.TrimFunc(word, cutSet_func) != "" {
				//if debug {fmt.Printf("word = %s\n",word)}
				wordCounter++
				c1 <- word
				c2 <- word
				c3 <- word
				c4 <- word
			}
		}
		waitGrp.Done()
	}
}
//--------


//---

func ngram1() {
	m1 = make(map[string]uint32)
	for p_word := range c1 {
		//p_word := <-c1
		n, Ok := m1[p_word]
		if Ok {
			n++
			m1[p_word] = n
		} else {
			m1[p_word] = 1
		}
		//fmt.Printf("ng1: %s\n", p_word)
	}
}

//---


type Ngram2 [2]string

func (b Ngram2) String() string {
	return b[0] + " " + b[1]
}

func ngram2() {
	m2 = maps.NewSMap()
	var ng Ngram2
	ng[0], ng[1] = "", ""
	for p_word := range c2 {
		//p_word := <-c2
		ng[1] = p_word
		ng = Ngram2{ng[0], ng[1]}

		if ng[0] != "" {
			r, Ok := m2.Get(ng)
			if Ok {
				i := r.(uint32)
				i++
				m2.Insert(ng, i)
				//fmt.Printf("k,v = %v, %v\n", ng,  r)
			} else {
				m2.Insert(ng, uint32(1))
			}
		}
		ng[0] = ng[1]
		//fmt.Printf("ng2: %s\n", *p_word)
	}
}


//---


type Ngram3 [3]string

func (b Ngram3) String() string {
	return b[0] + " " + b[1] + " " + b[2]
}

func ngram3() {
	m3 = maps.NewSMap()
	var ng Ngram3
	ng[0], ng[1], ng[2] = "", "", ""
	for p_word := range c3 {
		//p_word := <-c3
		ng[2] = p_word
		ng = Ngram3{ng[0], ng[1], ng[2]}

		if ng[0] != "" && ng[1] != "" {
			r, Ok := m3.Get(ng)
			if Ok {
				i := r.(uint32)
				i++
				m3.Insert(ng, i)
				//fmt.Printf("k,v = %v, %v\n", ng,  r)
			} else {

				m3.Insert(ng, uint32(1))
			}
		}
		ng[0], ng[1] = ng[1], ng[2]
		//fmt.Printf("ng2: %s\n", *p_word)
	}
}


//---

type Ngram4 [4]string

func (b Ngram4) String() string {
	return b[0] + " " + b[1] + " " + b[2] + " " + b[3]
}

func ngram4() {
	m4 = maps.NewSMap()
	var ng Ngram4
	ng[0], ng[1], ng[2], ng[3] = "", "", "", ""
	for p_word := range c4 {
		//p_word := <-c4
		ng[3] = p_word
		ng = Ngram4{ng[0], ng[1], ng[2], ng[3]}

		if ng[0] != "" && ng[1] != "" && ng[2] != "" {
			r, Ok := m4.Get(ng)
			if Ok {
				i := r.(uint32)
				i++
				m4.Insert(ng, i)
				//fmt.Printf("k,v = (%v), %v\n", ng, r)
			} else {
				m4.Insert(ng, uint32(1))
			}
		}
		ng[0], ng[1], ng[2] = ng[1], ng[2], ng[3]
		//fmt.Printf("ng4: %s\n", p_word)
	}
}

//---

func main() {

	if debug {
		println("Starting")
	}

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

	//creating waitGrp
	waitGrp = new(sync.WaitGroup)

	c1 = make(chan string, b)
	c2 = make(chan string, b)
	c3 = make(chan string, b)
	c4 = make(chan string, b)
	lineChan = make(chan string, b)

	go ngram1()
	go ngram2()
	go ngram3()
	go ngram4()
	go splitLine()

	if debug {
		println("Opening file")
	}

	// Open file in buffered mode
	file, err := os.Open(fname, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Printf("Input error: %s\n", err)
		return
	}

	// Read data line-by-line and send each line to a separate coroutine
	fileReader := bufio.NewReader(file)
	var line string
	for err == nil {
		line, err = fileReader.ReadString('\n')
		if len(line) > 0 && line != "\n" {
			waitGrp.Add(1)
			lineChan <- line
		}
	}
	waitGrp.Wait()

	// closing channels to filters when iteration over lines is finished
	close(c1)
	close(c2)
	close(c3)
	close(c4)

	// Stop timer and print results
	stopTime := time.Nanoseconds()
	runElapsed := float64(stopTime-runTime) / 1000000000
	fmt.Printf("Done: %d words in %f seconds\n", wordCounter, runElapsed)

	return
}
