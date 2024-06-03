package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go readMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(conn, text)
	}
}

func readMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
