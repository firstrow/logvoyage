package main

import (
	"log"
	"os"
)

func toBacklog(message string) {
	path, err := os.Getwd()
	if err != nil {
		log.Println("Can't define current directory")
	}
	fullPath := path + string(os.PathSeparator) + "back.log"
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Println("Can not open backlog file")
		return
	}
	defer file.Close()
	file.Write([]byte(message))
}

func resend() {

}
