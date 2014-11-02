package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:27077")
	if err != nil {
		log.Fatal("Error connecting to logvoyage")
	}

	file, err := os.Open("/Users/andrew/Code/requests.log.2")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		conn.Write([]byte("0b137205-3291-5f5b-5832-ab2458b9936a" + scanner.Text() + "\n"))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
