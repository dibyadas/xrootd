package xrootd_mockserver

import (
		"fmt"
		"os"
		"net"
		"io"
		"encoding/binary"
)

const ( 
		CONN_HOST = "0.0.0.0"
	 	CONN_PORT = "9005" 
)


func StartServer() {
	l, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}


func handleRequest(conn net.Conn) {
	request := make([]byte,20)
	defer conn.Close()

	io.ReadFull(conn,request)

	fmt.Println("Client Request :- ",request)

	a := binary.BigEndian.Uint32(request[0:])
	b := binary.BigEndian.Uint32(request[4:])
	c := binary.BigEndian.Uint32(request[8:])
	d := binary.BigEndian.Uint32(request[12:])
	e := binary.BigEndian.Uint32(request[16:])


	if a == 0 && b == 0 && c == 0 && d == 4 && e == 2012 {
		response := make([]byte, 16)
		binary.BigEndian.PutUint32(response[4:], 8)
		binary.BigEndian.PutUint32(response[8:], 784)
		binary.BigEndian.PutUint32(response[12:], 1)
		conn.Write(response)
	}

}