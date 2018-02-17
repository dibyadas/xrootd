package main

import (
	"fmt"
	"net"
	"xrootd"
	// "bytes"
	// "encoding/binary"
)

func main() {
	conn,err := net.Dial("tcp","0.0.0.0:9001")
	if err != nil {
		fmt.Println("no")
	}
	defer conn.Close()

	fmt.Println("--- Initiating Handshake ---")
	xrootd.SendHandshake(conn)
	fmt.Println("--- Done ---")
	fmt.Println()

	streamID := [2]byte{0xbe, 0xef}

	fmt.Println("--- Initiating SendProtocol ---")
	xrootd.SendProtocol(conn, streamID)
	fmt.Println("--- Done ---")
	fmt.Println()

	fmt.Println("--- Initiating SendLogin ---")
	xrootd.SendLogin(conn, streamID)
	fmt.Println("--- Done ---")
	fmt.Println()

	fmt.Println("--- Initiating SendPing ---")
	xrootd.SendPing(conn, streamID)
	fmt.Println("--- Done ---")
	fmt.Println()

	fmt.Println("--- Initiating SendInvalid ---")
	xrootd.SendInvalid(conn, streamID)
	fmt.Println("--- Done ---")
	fmt.Println()

}
