package main

import (
	"fmt"
	"log"
	"net/http"

	stomp "github.com/go-stomp/stomp"
)

func main() {

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{

		//If this is not used, there seems to be a connection drop frequently
		stomp.ConnOpt.HeartBeatGracePeriodMultiplier(5),
	}

	//ActiveMQ IP and TCP Port for data transmission
	var connection_string string = "IP:Port"

	conn, err := stomp.Dial("tcp", connection_string, options...)
	if err != nil {
		fmt.Println(err)
	}

	buffer := make(chan string, 10)

	// Launch Go Routine for pushing data from channel to the Queue
	go Producer(buffer, conn)

	http.HandleFunc("/insms", func(w http.ResponseWriter, r *http.Request) {

		// Reading the raw query of url and passing it to the channel
		buffer <- r.URL.RawQuery

		// Responding back to the client
		w.Write([]byte("Success"))

	})

	var listner_connection string = "IP:PORT"

	// Time to listen as a server
	log.Fatal(http.ListenAndServe(listner_connection, nil))

}

func Producer(c chan string, conn *stomp.Conn) {

	// Endless loop
	for {
		message := <-c
		err := conn.Send(
			"queue/QUEUE_NAME",
			"text/plain",
			[]byte(message),

			// Should use this else observed that ActiveMQ
			// misses the deliveries
			stomp.SendOpt.Receipt)

		if err != nil {
			fmt.Println(err)
		}
	}
}
