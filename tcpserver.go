package main

import (
	"log"
	"net"
)

type publishMessage struct {
	message   []byte
	sessionID int
}

func main() {
	newConnections := make(chan net.Conn, 128)
	deadConnections := make(chan int, 128)
	publishes := make(chan publishMessage, 128)
	connections := make(map[int]net.Conn)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// New incoming connection
	go checkForNewIncomingConnections(listener, newConnections)
	connectionCounter := 0 // Used to generate session IDs
	for {
		select {
		case connection := <-newConnections:
			connections[connectionCounter] = connection
			go newConnectionSession(connection, publishes, deadConnections, connectionCounter)
			connectionCounter++
			log.Print("Number of connections: ", len(connections))
		case deadConnection := <-deadConnections:
			_ = connections[deadConnection].Close()
			delete(connections, deadConnection)
			log.Print("Number of connections: ", len(connections))
		case publish := <-publishes:
			for session, connection := range connections {
				if publish.sessionID != session {
					go newPublish(publish, connection, deadConnections)
				}
			}
		}
	}
	listener.Close()
}

func newPublish(publish publishMessage, connection net.Conn, deadConnections chan int) {
	totalWritten := 0
	for totalWritten < len(publish.message) {
		writtenThisCall, err := connection.Write(publish.message[totalWritten:])
		if err != nil {
			deadConnections <- publish.sessionID
			break
		}
		totalWritten += writtenThisCall
	}
}

func newConnectionSession(connection net.Conn, publishes chan publishMessage, deadConnections chan int, id int) {
	buf := make([]byte, 1024)
	// Wait for incoming events
	for {
		numberOfBytes, err := connection.Read(buf)
		if err != nil {
			deadConnections <- id
			break
		} else {
			var messageData publishMessage
			messageData.message = make([]byte, numberOfBytes)
			copy(messageData.message, buf[:numberOfBytes])
			messageData.sessionID = id
			publishes <- messageData
		}
	}
}

func checkForNewIncomingConnections(listener net.Listener, newConnections chan net.Conn) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		newConnections <- connection
	}
}
