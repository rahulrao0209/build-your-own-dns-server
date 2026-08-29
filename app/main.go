package main

import (
	"fmt"
	"net"

	"github.com/rahulrao0209/build-your-own-dns-server/types"
	"github.com/rahulrao0209/build-your-own-dns-server/utils"
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

		/* Decode DNS client messasge */
		// from header section
		packetId := utils.DecodePacketId(buf[:size])
		questionCount := utils.DecodeQuestionCount(buf[:size])

		// question section
		domain, qTypeIdx := utils.DecodeDomainName(buf[12:])
		qtype := utils.DecodeType(buf[12:], qTypeIdx)  // QTYPE occupies the next 2 bytes after domain name
		class := utils.DecodeClass(buf[12:], qTypeIdx) // CLASS occupies the next 2 bytes after QTYPE

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
				QuestionCount:         questionCount,
				AnswerRecordCount:     0,
				AuthorityRecordCount:  0,
				AdditionalRecordCount: 0,
			},
			Question: types.Question{
				Name:  domain,
				Type:  qtype,
				Class: class,
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
