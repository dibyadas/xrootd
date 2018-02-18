package xrootd

// import (
// 	"net"
// )

import (
	// "github.com/lunixbochs/struc"
	// "bytes"
	"fmt"
	"net"
	"encoding/binary"
)

const kXR_protocol uint16 = 3006
const kXR_ping uint16 = 3011
const kXR_login uint16 = 3007


type ClientRequest struct{      // I think using this would be a much cleaner and organized approach.
	streamID string				// This struck me later :p 
	requestid uint16
	params string
	dlen uint32
	data string
}




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
	binary.Read(conn, binary.BigEndian, &response)
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
	binary.Read(conn, binary.BigEndian, &response)
	dlen := binary.BigEndian.Uint32(response[4:])
	if dlen != 0 {
		data := make([]byte, 8+dlen)
		binary.Read(conn, binary.BigEndian, data[8:])
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
	// binary.BigEndian.PutUint32(bytesToSend[20:], 0)
	copy(bytesToSend[8:], []byte("gopher"))
	_, err := conn.Write(bytesToSend)
	if err != nil{
		return err
	}

	response := make([]byte,8)
	binary.Read(conn, binary.BigEndian, &response)
	slen := binary.BigEndian.Uint32(response[4:])
	if slen != 0 {
		sec := make([]byte, 8+slen)
		binary.Read(conn, binary.BigEndian, sec[8:])
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
	binary.Read(conn, binary.BigEndian, &response)
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


func SendHandshake(conn net.Conn) error {
	bytesToSend := PrepHandshake()
	_, err := conn.Write(bytesToSend)
	if err != nil{
		return err
	}

	response := make([]byte, 16)
	binary.Read(conn, binary.BigEndian, &response)
	serverType := binary.BigEndian.Uint32(response[12:])
	if serverType == 1 {
		fmt.Println("DataServer")
	} else if serverType == 0 {
		fmt.Println("LoadBalancer")
	}

	return err
}

// func (cr *ClientRequest) PrepHandshake(w [*bytes.Buffer]) error {
// 	return struc.Pack(w,cr)
// }
