package main

import (
	"flag"
	"fmt"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	PING_COMPONENT    = "ping"
	PING_MIN_INTERVAL = 0.1
)

type ndnStatic struct {
	Sent     int
	Received int
	Start    time.Time
	Min      float64
	Max      float64
	Tsum     float64
	Tsum2    float64
}

var interval = flag.Float64("i", 1.0, "set ping interval in seconds")
var count = flag.Int("c", -1, "set total number of ping")
var num = flag.Int("n", -1, "set the starting number, the number is increamented by 1")
var pt = flag.Bool("t", false, "print timestamp")
var path = flag.String("p", "/ndn/ping/", "Listen Path")
var nfdServer = flag.String("s", ":6363", "set connect nfdServer")
var keyPath = flag.String("k", "key/default.pri", "Access key path")
var sta ndnStatic

func printStatistics() {
	fmt.Println("--- ndnping statistics --", *path)
	if sta.Sent > 0 {
		lost := (float64)(sta.Sent-sta.Received) * 100 / float64(sta.Sent)
		now := time.Now()
		diffTime := now.Sub(sta.Start)
		fmt.Printf("%d Interests transmmitted, %d Data received, %.1f%% packet loss, time %s\n",
			sta.Sent, sta.Received, lost, diffTime.String())
	}

	if sta.Received > 0 {
		avg := sta.Tsum / float64(sta.Received)
		mdev := math.Sqrt(sta.Tsum2/float64(sta.Received) - avg*avg)
		fmt.Printf("rtt min/avg/max/mdev = %.3f/%.3f/%.3f/%.3f ms\n", sta.Min, avg, sta.Max, mdev)
		fmt.Printf("")
	}

	return
}

func doPing(face ndn.Face, index int, key ndn.Key) {
	interestName := *path + strconv.Itoa(index)
	f := mux.NewFetcher()
	f.Use(mux.ChecksumVerifier)
	interest := &ndn.Interest{
		Name: ndn.NewName(interestName),
	}
	interest.Selectors.MustBeFresh = true
	sta.Sent++
	start := time.Now()
	data := f.Fetch(face, interest)
	if data != nil {
		sta.Received++
		fmt.Printf("Content from %s:\tnumber=%d\trtt=%0.3fms\n", *path, index, time.Now().Sub(start).Seconds()*1000)
	} else {
		fmt.Printf("Content from %s:\tnumber=%d\t Timeout", *path, index)
	}

}

func main() {
	flag.Parse()
	fmt.Println(*count)
	if *path == "" {
		fmt.Println("Please set Listen path by using -p params")
		return
	}
	if *interval < PING_MIN_INTERVAL {
		fmt.Println("Interval is less than min interval")
		return
	}

	if *num == -1 {
		rand.Seed(time.Now().Unix())
		*num = rand.Int() % 10000
	}

	//SIGINT deal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			printStatistics()
			os.Exit(0)
		}
	}()

	conn, err := packet.Dial("tcp", *nfdServer)
	if err != nil {
		log.Fatalln(err)
	}

	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	pem, err := os.Open(*keyPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer pem.Close()

	key, _ := ndn.DecodePrivateKey(pem)

	sta.Start = time.Now()
	durString := strconv.FormatFloat(*interval, 'f', 2, 64) + "s"
	d, err := time.ParseDuration(durString)
	if err != nil {
		fmt.Println("Can't parse interval time ")
		return
	}

	for *count != 0 {
		*count--
		go doPing(face, *num, key)
		*num++
		time.Sleep(d)
	}
	printStatistics()

}
