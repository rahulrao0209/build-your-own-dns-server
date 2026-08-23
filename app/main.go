package main

import (
	"encoding/binary"
	"fmt"
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

		packetId := binary.BigEndian.Uint16(buf[0:2])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		message := &types.DNSMessage{
			Header: types.Header{
				PacketIdentifier:      packetId,
				ResponseIndicator:     1,
				OperationCode:         0,
				AuthoritativeAnswer:   0,
				Truncation:            0,
				RecursionDesired:      0,
				RecursionAvailable:    0,
				Reserved:              0,
				ResponseCode:          0,
				QuestionCount:         0,
				AnswerRecordCount:     0,
				AuthorityRecordCount:  0,
				AdditionalRecordCount: 0,
			},
		}

		marshalledDNSMessage, err := message.MarshalBinary()
		if err != nil {
			return
		}

		response := marshalledDNSMessage

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
