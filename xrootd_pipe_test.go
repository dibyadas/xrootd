package main

import (
	"testing"
	"xrootd"
	"xrootd_mockserver"
	// "fmt"
	"net"
	"os"
	// "io"
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


var Client, Server net.Conn


func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}


func TestSendLogin(t *testing.T) {

	testValues := []TestLoginResponse{
		{"Test1", [2]byte{190,239}, "dibya",1},
	}

	var tester func(TestLoginResponse)
	tester = func(v TestLoginResponse) {

			Client, Server = net.Pipe()
			go xrootd_mockserver.HandleRequest(Server)

			err2 := xrootd.SendLogin(Client, [2]byte(v.streamID), v.username)
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

			Client, Server = net.Pipe()
			go xrootd_mockserver.HandleRequest(Server)
			
			serverType, _ := xrootd.SendHandshake(Client, v.request)
			if serverType != v.serverType {
				t.Errorf("%s Fail, serverType != %d, got serverType = %d",v.name , v.serverType, serverType)	
			}
			
	}

	for _, v := range testValues {
		tester(v)
	}
}


