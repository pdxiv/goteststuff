package main

import (
	"log"
	"net"
)

func main() {
	newConnections := make(chan net.Conn, 128)
	deadConnections := make(chan net.Conn, 128)
	publishes := make(chan []byte, 128)
	connections := make(map[net.Conn]bool)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// New incoming connection
	go newIncomingConnection(listener, newConnections)

	for {
		select {
		case conn := <-newConnections:
			connections[conn] = true
			go func() {
				buf := make([]byte, 1024)
				for {
					nbyte, err := conn.Read(buf)
					if err != nil {
						deadConnections <- conn
						break
					} else {
						fragment := make([]byte, nbyte)
						copy(fragment, buf[:nbyte])
						publishes <- fragment
					}
				}
			}()
			log.Print("Active connections: ", len(connections))
		case deadConn := <-deadConnections:
			_ = deadConn.Close()
			delete(connections, deadConn)
			log.Print("Active connections: ", len(connections))
		case publish := <-publishes:
			for conn := range connections {
				go func(conn net.Conn) {
					totalWritten := 0
					for totalWritten < len(publish) {
						writtenThisCall, err := conn.Write(publish[totalWritten:])
						if err != nil {
							deadConnections <- conn
							break
						}
						totalWritten += writtenThisCall
					}
				}(conn)
			}
		}
	}
	listener.Close()
}

func newIncomingConnection(listener net.Listener, newConnections chan net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		newConnections <- conn
	}
}
