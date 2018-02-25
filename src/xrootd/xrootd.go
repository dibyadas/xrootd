package xrootd

import (
	"fmt"
	"net"
	"encoding/binary"
	"io"
)

const kXR_protocol uint16 = 3006
const kXR_ping uint16 = 3011
const kXR_login uint16 = 3007


func SendProtocol(conn net.Conn, streamID [2]byte) error {
	bytesToSend := make([]byte,24)

	copy(bytesToSend[0:],streamID[0:])

	binary.BigEndian.PutUint16(bytesToSend[2:], kXR_protocol)
	binary.BigEndian.PutUint32(bytesToSend[4:], 784)  // 784 is the protocol number
	
	_, err := conn.Write(bytesToSend)
	if err != nil{
		return err
	}

	response := make([]byte,16)
	if _, readErr := io.ReadFull(conn,response); readErr != nil {
		return readErr
	}

	fmt.Println(response)
	return err
}

func SendInvalid(conn net.Conn, streamID [2]byte) error {
	bytesToSend := make([]byte,24)
	copy(bytesToSend[0:],streamID[0:])
	binary.BigEndian.PutUint16(bytesToSend[2:], 0)  // for invalid request

	_, err := conn.Write(bytesToSend)
	if err != nil{
		return err
	}

	response := make([]byte,8)
	if _, readErr := io.ReadFull(conn,response); readErr != nil {
		return readErr
	}

	dlen := binary.BigEndian.Uint32(response[4:])
	if dlen != 0 {
		data := make([]byte, 8+dlen)
		if _, readErr2 := io.ReadFull(conn,data[8:]); readErr2 != nil {
			return readErr2
		}
		copy(data[0:8],response[0:])
		response = data
		fmt.Println(response)
		fmt.Println("Invalid Status Code:- ",binary.BigEndian.Uint16(response[2:]))
		fmt.Println("Server message:- ",string(response[12:])) // why from 12 and not 8?? 
		return err	
	}

	fmt.Println(response)
	return err
}


func SendLogin(conn net.Conn, streamID [2]byte) error {
	bytesToSend := make([]byte, 24)
	copy(bytesToSend[0:], streamID[0:])

	binary.BigEndian.PutUint16(bytesToSend[2:], kXR_login)

	copy(bytesToSend[8:], []byte("gopher"))
	_, err := conn.Write(bytesToSend)
	if err != nil{
		return err
	}

	response := make([]byte,8)
	if _, readErr := io.ReadFull(conn,response); readErr != nil {
		return readErr
	}

	slen := binary.BigEndian.Uint32(response[4:])
	if slen != 0 {
		sec := make([]byte, 8+slen)
		if _, readErr2 := io.ReadFull(conn,sec[8:]); readErr2 != nil {
			return readErr2
		}
		copy(sec[0:8],response[0:])
		response = sec
		fmt.Println(response)
		return err	
	}
	fmt.Println(response)
	return err
}


func SendPing(conn net.Conn, streamID [2]byte) error {
	bytesToSend := make([]byte, 24)

	copy(bytesToSend[0:], streamID[0:])

	binary.BigEndian.PutUint16(bytesToSend[2:], kXR_ping)
	binary.BigEndian.PutUint32(bytesToSend[20:], 0)
	
	_, err := conn.Write(bytesToSend)
	if err != nil{
		return err
	}

	response := make([]byte,8)
	if _, readErr := io.ReadFull(conn,response); readErr != nil {
		return readErr
	}
	fmt.Println(response)
	return err	

}


func PrepHandshake() []byte {
	bytesToSend := make([]byte,20)  // 20 bytes byte array
	handshakeValues := []uint32{0,0,0,4,2012}    // handshake protocol 
	for i, v  := range handshakeValues {
		binary.BigEndian.PutUint32(bytesToSend[4*i:], v)
	}
	return bytesToSend
}


func SendHandshake(conn net.Conn, bytesToSend []byte) (int, error) {
	_, err := conn.Write(bytesToSend)
	if err != nil{
		return -1,err
	}

	response := make([]byte, 16)
	if _, readErr := io.ReadFull(conn,response); readErr != nil {
		return -1,readErr
	}
	serverType := binary.BigEndian.Uint32(response[12:])
	if serverType == 1 {
		fmt.Println("DataServer")
	} else if serverType == 0 {
		fmt.Println("LoadBalancer")
	}
	return int(serverType),err
}
