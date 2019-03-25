package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	readerAdress := flag.String("reader-addr", "", " Reader address. It could be find in /dev/tty*")
	flag.Parse()
	if *readerAdress != "" {
		tty, err := os.OpenFile(*readerAdress, os.O_RDWR|syscall.O_NOCTTY, 0)
		if err != nil {
			log.Fatalf("Cannot open tty port: %v\n", err)
		}

		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		go func() {
			var buf = make([]byte, 8192)
			for {
				select {
				case <-done:
					log.Println("Quit!")
					return
				default:
					nr, _ := tty.Read(buf)
					log.Printf("Write:%s(%d)\n", string(buf[:nr]), nr)
				}
			}
		}()

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			close(done)
		}()

		fmt.Println("Ctrl+C to quit")
		<-done
		fmt.Println("exiting")
	} else {
		log.Fatalln("Please provide valid reader address.")
	}

}
