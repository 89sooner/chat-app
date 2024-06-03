package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]bool)
	broadcast = make(chan string)
	mutex     sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started on port 8080")

	go handleMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		mutex.Lock()
		clients[conn] = true
		mutex.Unlock()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	fmt.Println("Client connected:", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("Received:", text)
		broadcast <- text
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			go func(conn net.Conn) {
				_, err := fmt.Fprintln(conn, msg)
				if err != nil {
					conn.Close()
					mutex.Lock()
					delete(clients, conn)
					mutex.Unlock()
				}
			}(client)
		}
		mutex.Unlock()
	}
}