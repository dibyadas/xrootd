package main

import (
	"testing"
	"xrootd"
	"xrootd_mockserver"
	"fmt"
	"net"
	"os"
	"io"
)

const (
	MOCKSEVER_HOST = xrootd_mockserver.CONN_HOST //"0.0.0.0"
)


type TestHandshakeResponse struct {
	name string
	request []byte
	serverType int
}

var service chan string
var MOCKSEVER_PORT string

func TestMain(m *testing.M) {
	started := make(chan string)
	service = make(chan string, 1)

	go xrootd_mockserver.StartServer(started, service)
	MOCKSEVER_PORT = <- started
	code := m.Run()
	os.Exit(code)
}


func TestSendLogin(t *testing.T) {
	service <-	"SendLogin"
	conn,err := net.Dial("tcp", MOCKSEVER_HOST+":"+MOCKSEVER_PORT)
	defer conn.Close()	
	if err != nil {
		fmt.Println("no")
	}
	conn.Write([]byte{0,0,0,1})    // WIP
	resp := make([]byte, 4)
	io.ReadFull(conn, resp)
	fmt.Println(resp)
}

func TestSendHandshake(t *testing.T) {

	testValues := []TestHandshakeResponse{
		{"Test1", []byte{
			0,0,0,0,
			0,0,0,0,
			0,0,0,0,
			0,0,0,4,
			0,0,7,220,
		}, 1},  // correct
		{"Test2", []byte{
			0,0,0,0,
			0,0,1,0,
			0,0,0,0,
			0,0,0,4,
			0,0,34,220,
		}, -1}, // incorrect
		{"Test3", []byte{
			0,0,0,0,
			0,0,0,0,
			0,1,0,0,
			0,0,0,4,
			0,0,3,220,
		}, -1}, // incorrect
		{"Test4", []byte{
			0,1,0,0,
			0,0,0,0,
			0,0,0,0,
			0,0,0,4,
			0,0,7,120,
		}, -1}, // incorrect
	}

	var tester func(TestHandshakeResponse)
	tester = func(v TestHandshakeResponse) {
			conn,err := net.Dial("tcp", MOCKSEVER_HOST+":"+MOCKSEVER_PORT)
			defer conn.Close()	
			if err != nil {
				fmt.Println("no")
			}
			serverType, err := xrootd.SendHandshake(conn, v.request)
			if serverType != v.serverType {
				t.Errorf("%s Fail, serverType != %d",v.name , v.serverType)	
			}
	}

	for _, v := range testValues {
		service <- "SendHandshake"
		tester(v)
	}
}


