// Websocket server package.
// This package starts separate websocket server and transfers all
// messages from redis channel "ws" to client browser.
//
// Running server
//   web_socket.StartServer()
// Example code to send data to redis:
//	c, _ := redis.Dial("tcp", ":6379")
//  msg := web_socket.RedisMessage{"apiKey", map[string]interface{}{
//	  "log_per_second": 24,
//    "kbs_per_second": 128,
//  }}
//  msg.Send(c)
package web_socket

import (
	"encoding/json"
	"log"
	"net/http"

	c "github.com/firstrow/logvoyage/configuration"
	"github.com/firstrow/logvoyage/common"

	"code.google.com/p/go.net/websocket"
	"github.com/garyburd/redigo/redis"
)

const (
	redisChannel = "ws"
)

// Represents data to be sent to user by its apiKey
type RedisMessage struct {
	ApiKey string
	Data   interface{}
}

func (m *RedisMessage) Send(r redis.Conn) (interface{}, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return r.Do("PUBLISH", redisChannel, string(j))
}

// Store connected clients: [apikey]Connection
var clients = make(map[string]*websocket.Conn)

func StartServer() {
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
	c, err := redis.Dial("tcp", c.ReadConf().Redis.GetURI())
	checkError(err)
	c.Send("SUBSCRIBE", redisChannel)
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
	// websocket.Message.Send(ws, "Hello dear user!")

	for {
		// Message received from client
		var message string

		// Read messages from client
		// Code blocks here, after any message received
		// will resume execution.
		if err := websocket.Message.Receive(ws, &message); err != nil {
			log.Println("Error receiving message. Closing connection.")
			return
		}

		// Register user
		// TODO: Cache user
		user, _ := common.FindUserByApiKey(message)
		if user != nil {
			log.Println("Registering apiKey", user.ApiKey)
			clients[user.ApiKey] = ws
		} else {
			log.Println("Error registering user", message)
		}
	}
}
