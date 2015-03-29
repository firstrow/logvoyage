// Accept http messages.
// bulk?apiKey=XXX&type=XXX - accepts bulk of messages separated by newline.
package backend

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("apiKey")
	if apiKey == "" {
		return
	}
	logType := r.URL.Query().Get("type")
	if logType == "" {
		return
	}

	key := fmt.Sprintf("%s@%s", apiKey, logType)

	reader := bufio.NewReader(r.Body)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			processMessage(fmt.Sprintf("%s %s", key, line))
		}
		if err != nil {
			return
		}
	}
}

func initHttpServer() {
	log.Print("Starting HTTP server at " + httpDsn)
	http.HandleFunc("/bulk", httpHandler)
	http.ListenAndServe(httpDsn, nil)
}
