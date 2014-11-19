package main

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/garyburd/redigo/redis"
)

func startWebSocket() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal("Cannot connect to redis", err.Error())
	}
	c.Send("SUBSCRIBE", "logstats")
	c.Flush()
	psc := redis.PubSubConn{c}
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Printf("%s: message: %s\n", v.Channel, v.Data)
				for _, val := range clients {
					if err = websocket.Message.Send(val, string(v.Data)); err != nil {
						// It could not send message to a peer
						log.Println("Could not send message to ", val, err.Error())
					}
				}
			case redis.Subscription:
				log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				log.Println("Error")
			}
		}
	}()

	http.Handle("/ws", websocket.Handler(Echo))
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Store connected clients: [apikey]Connection
var clients = make(map[string]*websocket.Conn)

func Echo(ws *websocket.Conn) {
	log.Println("New client")
	defer ws.Close()

	websocket.Message.Send(ws, "Hello dear user!")
	for {
		var message string
		if err := websocket.Message.Receive(ws, &message); err != nil {
			log.Println("Error receiving message. Closing connection.")
			return
		} else {
			// Assumes that user send to us api key)))
			clients["message"] = ws
		}
	}
}
