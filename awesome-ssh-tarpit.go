package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type ClientManager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	metricChannel chan <- float64
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			manager.metricChannel <- 1
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				manager.metricChannel <- -1
				fmt.Println("A connection has terminated!")
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		b := make([]byte, 4096)
		length, err := client.socket.Read(b)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length == 0 {
			manager.unregister <- client
			client.socket.Close()
			break

		}
		message := string(b)
		if !strings.HasPrefix(message, "SSH") {
			manager.unregister <- client
			client.socket.Close()
			break

		}
		randStr := RandStringBytesMaskImprSrcUnsafe(32)
		client.socket.Write([]byte(randStr))
		time.Sleep(10 * time.Second)
	}
}

func listenAndServe(addr string, metricChannel chan float64) error {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		metricChannel: metricChannel,
	}
	go manager.start()
	for {
		connection, _ := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
	}
}

func main() {
	sshdAddr := flag.String("sshd", ":22", "Address to bind sshd.")
	httpAddr := flag.String("http", ":8090", "Address to bind httpd port.")
	flag.Parse()
	metricChan := make(chan float64)
	m := NewMetricServer(metricChan)
	go m.ListenAndServe(*httpAddr)
	log.Fatal(listenAndServe(*sshdAddr, metricChan))
}
