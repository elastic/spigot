package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

var ACTIONS = [...]string{"ACCEPT", "REJECT"}
var STATUSES = [...]string{"OK", "SKIPDATA", "NODATA"}
var COUNTS = [...]int{1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072}

type flow struct {
	//	version     string
	//	accountId   string
	//	interfaceId string
	srcAddr   net.IP
	dstAddr   net.IP
	srcPort   int
	dstPort   int
	protocol  int
	packets   int
	bytes     int
	start     int64
	end       int64
	action    string
	logStatus string
}

func (f flow) String() string {
	return fmt.Sprintf("2 123456789010 eni-1235b8ca123456789 %s %s %s %s %s %s %s %s %s %s %s\n",
		f.srcAddr.String(),
		f.dstAddr.String(),
		strconv.Itoa(f.srcPort),
		strconv.Itoa(f.dstPort),
		strconv.Itoa(f.protocol),
		strconv.Itoa(f.packets),
		strconv.Itoa(f.bytes),
		strconv.FormatInt(f.start, 10),
		strconv.FormatInt(f.end, 10),
		f.action,
		f.logStatus)
}

func generateAddr() net.IP {
	u32 := rand.Uint32()
	return net.IPv4(byte(u32&0xff), byte((u32>>8)&0xff), byte((u32>>16)&0xff), byte((u32>>24)&0xff))
}

func newFlow() flow {
	f := flow{}
	f.srcAddr = generateAddr()
	f.dstAddr = generateAddr()
	f.srcPort = rand.Intn(65536)
	f.dstPort = rand.Intn(65536)
	f.protocol = rand.Intn(256)
	f.packets = rand.Intn(1048576)
	f.bytes = f.packets * 1500
	f.start = time.Now().Unix() - int64(rand.Intn(60))
	f.end = f.start + int64(f.bytes/800)
	f.action = ACTIONS[rand.Intn(2)]
	if f.packets == 0 {
		f.logStatus = STATUSES[2]
	} else {
		f.logStatus = STATUSES[rand.Intn(2)]
	}
	return f
}

func generateFile(workerNum int, fChan <-chan int, rChan chan<- int) {
	for fileNum := range fChan {
		path, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current directory: %v", err)
			rChan <- 1
			continue
		}
		fp, err := os.CreateTemp(path, "vpcflowlogs_*.gz")
		if err != nil {
			log.Printf("Error creating temp file: %v", err)
			rChan <- 1
			continue
		}

		gw := gzip.NewWriter(fp)
		bw := bufio.NewWriter(gw)

		lines := COUNTS[rand.Intn(len(COUNTS))]
		for i := 0; i < lines; i++ {
			f := newFlow()
			_, err := bw.WriteString(f.String())
			if err != nil {
				log.Printf("Error writing string: %v", err)
				continue
			}
		}
		bw.Flush()
		gw.Close()
		fp.Close()
		rChan <- fileNum
	}

}

func main() {
	var numFiles int
	var numWorkers int
	flag.IntVar(&numFiles, "n", 1, "number of files to generate")
	flag.IntVar(&numWorkers, "w", 1, "number of workers")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	fChan := make(chan int, numFiles)
	rChan := make(chan int, numFiles)

	for i := 0; i < numWorkers; i++ {
		go generateFile(i, fChan, rChan)
	}

	for i := 0; i < numFiles; i++ {
		fChan <- i
	}
	for i := 0; i < numFiles; i++ {
		<-rChan
	}
}
