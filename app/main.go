package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rahulrao0209/build-your-own-dns-server/types"
)

var _ = net.ListenUDP

func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Received DNS Query
		dnsQuery := buf[:size]

		/* Decode DNS client query and assemble a reply */
		var dnsReply *types.DNSMessage = &types.DNSMessage{}
		dnsReply, err = dnsReply.UnmarshalBinary(dnsQuery)

		// Encode DNS reply
		marshalledDNSReply, err := dnsReply.MarshalBinary()
		if err != nil {
			log.Fatal("Error encoding DNS reply")
			return
		}

		_, err = udpConn.WriteToUDP(marshalledDNSReply, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
