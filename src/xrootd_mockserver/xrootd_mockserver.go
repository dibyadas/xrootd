package xrootd_mockserver

import (
		"fmt"
		"os"
		"net"
		"io"
		"encoding/binary"
		"log"
)

const ( 
		CONN_HOST = "0.0.0.0"
	 	CONN_PORT = "0" 
)


const (
	kXR_login uint16 = 3007
)


func StartServer(started chan string) {
	l, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
 
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    started <- fmt.Sprintf("%d",(l.Addr().(*net.TCPAddr).Port))
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        go handleRequest(conn)
    }
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	functions :=  map[uint16]func(net.Conn,[]byte){ 
				kXR_login: SendLoginServe,
	}

	requestHeaders := make([]byte,20)
	io.ReadFull(conn, requestHeaders)


	streamID := binary.BigEndian.Uint16(requestHeaders[0:])
	requestID := binary.BigEndian.Uint16(requestHeaders[2:])

	if streamID == 0 && requestID == 0 {
		SendHandshakeServe(conn, requestHeaders)
	} else {
		dlen_byte := make([]byte, 4)
		io.ReadFull(conn, dlen_byte)
		dlen := binary.BigEndian.Uint32(dlen_byte)
		if dlen != 0 {
			sec := make([]byte, 24+dlen)
			binary.BigEndian.PutUint32(sec[20:], dlen)
			io.ReadFull(conn, sec[24:])
			copy(sec[0:20],requestHeaders[0:])
			requestHeaders = sec
		}
		functions[requestID](conn, requestHeaders)
	}
}


func SendLoginServe(conn net.Conn, request []byte) {
	fmt.Println("Client Request in login:- ",request)

	response := make([]byte, 8)
	copy(response[0:2], request[0:2])
	if _, err := conn.Write(response); err != nil {
		log.Fatal(err)
	}
}


func SendHandshakeServe(conn net.Conn, request []byte) {
	

	fmt.Println("Client Request :- ",request)

	a := binary.BigEndian.Uint32(request[0:])
	b := binary.BigEndian.Uint32(request[4:])
	c := binary.BigEndian.Uint32(request[8:])
	d := binary.BigEndian.Uint32(request[12:])
	e := binary.BigEndian.Uint32(request[16:])


	if  a == 0 && b == 0 && c == 0  && d == 4 && e == 2012 {
		response := make([]byte, 16)
		binary.BigEndian.PutUint32(response[4:], 8)
		binary.BigEndian.PutUint32(response[8:], 784)
		binary.BigEndian.PutUint32(response[12:], 1)
		if _, err := conn.Write(response); err != nil {
			log.Fatal(err)
		}
	}
}

// func SendLoginServe(conn net.Conn) {
	
// }