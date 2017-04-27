package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var logger *log.Logger

func main() {
	sleepTime := flag.Int64("sleepTime", 50, "sleep time between messages in ms")
	numThreads := flag.Int("numThreads", 100, "number of concurrent TCP hogs")
	messagesPerWorker := flag.Int("messagesPerWorker", 100, "number of messages a worker sends before quitting")
	flag.Parse()

	logger = log.New(os.Stdout, "", log.LstdFlags)

	address := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 13333}
	listener, err := net.ListenTCP("tcp", address)
	must(err)
	go func() {
		for {
			incoming, err := listener.Accept()
			must(err)
			go func(c net.Conn) {
				defer c.Close()
				scanner := bufio.NewScanner(c)
				for scanner.Scan() {
					logger.Println(scanner.Text())
				}
			}(incoming)
		}
	}()

	workerCh := make(chan int, 10)
	for i := 0; i < *numThreads; i++ {
		go worker(i, *messagesPerWorker, workerCh, address, time.Duration(*sleepTime))
	}

	for {
		i := <-workerCh
		go worker(i, *messagesPerWorker, workerCh, address, time.Duration(*sleepTime))
	}
}

func worker(workerID, messagesPerWorker int, done chan<- int, address *net.TCPAddr, sleepTime time.Duration) {
	logger.Printf("starting worker %d\n", workerID)
	conn, err := net.DialTCP("tcp", nil, address)
	must(err)
	defer conn.Close()
	for j := 0; j < messagesPerWorker; j++ {
		fmt.Fprintf(conn, "hi from hog %d\n", workerID)
		time.Sleep(time.Millisecond * sleepTime)
	}
	done <- workerID
	logger.Printf("finished worker %d\n", workerID)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
