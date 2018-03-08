package main

import (
	"testing"
	"xrootd"
	"xrootd_mockserver"
	"fmt"
	"net"
	"os"
	// "io"
)

const (
	MOCKSEVER_HOST = xrootd_mockserver.CONN_HOST //"0.0.0.0"
)


type TestHandshakeResponse struct {
	name string
	request []byte
	serverType int
}

type TestLoginResponse struct {
	name string
	streamID [2]byte
	username string
	result int
}

var MOCKSEVER_PORT string

func TestMain(m *testing.M) {
	started := make(chan string,1)

	go xrootd_mockserver.StartServer(started)
	MOCKSEVER_PORT = <- started
	code := m.Run()
	os.Exit(code)
}


func TestSendLogin(t *testing.T) {

	testValues := []TestLoginResponse{
		{"Test1", [2]byte{190,239}, "dibya",1},
	}

	var tester func(TestLoginResponse)
	tester = func(v TestLoginResponse) {
			conn,err := net.Dial("tcp", MOCKSEVER_HOST+":"+MOCKSEVER_PORT)
			defer conn.Close()	
			if err != nil {
				fmt.Println("no")
			}
			err2 := xrootd.SendLogin(conn, [2]byte(v.streamID), v.username)
			if err2 != nil {
				t.Errorf("Fail %s",v.name)	
			}
	}

	for _, v := range testValues {
		tester(v)
	}


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
			0,0,0,0,
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
				t.Errorf("%s Fail, serverType != %d, got serverType = %d",v.name , v.serverType, serverType)	
			}
	}

	for _, v := range testValues {
		tester(v)
	}
}


