package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	log.Printf("Listen at port: 6060")
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	for {
		_ = fmt.Sprint("test sprint")
		time.Sleep(time.Millisecond)
	}
}
