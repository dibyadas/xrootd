package main

import (
	"testing"
	"xrootd"
	"xrootd_mockserver"
	"fmt"
	"net"
	"encoding/binary"
)

const (
	MOCKSEVER_HOST = xrootd_mockserver.CONN_HOST //"0.0.0.0"
	MOCKSEVER_PORT = xrootd_mockserver.CONN_PORT//"9005"
)


func TestSendHandshake(t *testing.T) {
	// go xrootd_mockserver.StartServer()

	// done := make(chan bool)

	testValues := []struct{
		testValue []uint32
		serverType int
	}{
		{[]uint32{0,0,0,4,2012}, 1},  // correct
		{[]uint32{0,0,0,4,1212}, -1}, // incorrect
		{[]uint32{0,1,0,4,2012}, -1}, // incorrect
		{[]uint32{0,0,1,4,2012}, -1}, // incorrect
	}

	var tester func([]byte, int)
	tester = func(bytesToSend []byte, expectedOutput int) {
			conn,err := net.Dial("tcp", MOCKSEVER_HOST+":"+MOCKSEVER_PORT)
			defer conn.Close()	
			if err != nil {
				fmt.Println("no")
			}

			serverType, err := xrootd.SendHandshake(conn, bytesToSend)
			if serverType != expectedOutput {
				t.Errorf("Fail, serverType != %d",expectedOutput)	
			}
			// done <- true
	}

	for _, v := range testValues {
		bytesToSend := make([]byte,20) 
		for i, v  := range v.testValue {
			binary.BigEndian.PutUint32(bytesToSend[4*i:], v)
		}

		tester(bytesToSend, v.serverType)
	}

	// for range testValues {
	// 	<- done
	// }
}

