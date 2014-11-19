// Websocket server package.
// This package starts separate websocket server and transfers all
// messages from redis channel "ws" to client browser.
//
// Example code to send data to redis:
//	c, _ := redis.Dial("tcp", ":6379")
//  msg := web_socket.RedisMessage{"apiKey", map[string]interface{}{
//	  "log_per_second": 24,
//    "kbs_per_second": 128,
//  }}
//  msg.Send(c)

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/firstrow/logvoyage/common"
	"github.com/garyburd/redigo/redis"
)

// Represents data to be sent to user by its apiKey
type RedisMessage struct {
	ApiKey string
	Data   interface{}
}

func (m *RedisMessage) Send(r redis.Conn) error {
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	r.Do("PUBLISH", "ws", string(j))
	return nil
}

// Store connected clients: [apikey]Connection
var clients = make(map[string]*websocket.Conn)

func main() {
	go startListetingRedis()

	http.Handle("/ws", websocket.Handler(wsHandler))
	err := http.ListenAndServe(":12345", nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Listen to Redis and send messages to clients
func startListetingRedis() {
	c, err := redis.Dial("tcp", ":6379")
	checkError(err)
	c.Send("SUBSCRIBE", "ws")
	c.Flush()

	log.Println("Started server and connected to redis")

	psc := redis.PubSubConn{c}
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Printf("%s: message: %s\n", v.Channel, v.Data)

				// We must recive new json-encoded RedisMessage
				var message RedisMessage
				err = json.Unmarshal(v.Data, &message)

				if err != nil {
					continue
				}

				// If client found by apiKey from message
				if wsClient, ok := clients[message.ApiKey]; ok {
					// Marshal messaged data back to json
					// and send to client
					j, _ := json.Marshal(message.Data)
					if err = websocket.Message.Send(wsClient, string(j)); err != nil {
						wsClient.Close()
						delete(clients, message.ApiKey)
						log.Println("Could not send message to ", wsClient, err.Error())
					}

				}
			case error:
				log.Println("Error occured with redis.", v)
			}
		}
	}()
}

// Connection handler. This function called after new client
// connected to websocket server.
// Also this method performs register user - client must send valid apiKey
// to receive messages from redis.
func wsHandler(ws *websocket.Conn) {
	log.Println("New client")
	defer ws.Close()
	websocket.Message.Send(ws, "Hello dear user!")

	for {
		var message string

		// Read messages from client
		// Code blocks here, after any message received it
		// will resume execution.
		if err := websocket.Message.Receive(ws, &message); err != nil {
			log.Println("Error receiving message. Closing connection.")
			return
		}

		// Register user
		// TODO: Cache user
		user := common.FindUserByApiKey(message)
		if user != nil {
			log.Println("Registering apiKey", user.ApiKey)
			clients[user.ApiKey] = ws
		} else {
			log.Println("Error registering user", message)
		}
	}
}
