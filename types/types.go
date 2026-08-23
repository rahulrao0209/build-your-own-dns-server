package types

import "encoding/binary"

type Header struct { // 12 bytes
	PacketIdentifier      uint16 // 16 bits
	ResponseIndicator     uint8  // 1 bit
	OperationCode         uint8  // 4 bits
	AuthoritativeAnswer   uint8  // 1 bit
	Truncation            uint8  // 1 bit
	RecursionDesired      uint8  // 1 bit
	RecursionAvailable    uint8  // 1 bit
	Reserved              uint8  // 3 bits
	ResponseCode          uint8  // 4 bits
	QuestionCount         uint16 // 16 bits
	AnswerRecordCount     uint16 // 16 bits
	AuthorityRecordCount  uint16 // 16 bits
	AdditionalRecordCount uint16 // 16 bits
}

type DNSMessage struct {
	Header Header
}

func (m *DNSMessage) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 12)
	header := m.Header

	flags := uint16(header.ResponseIndicator)<<15 |
		uint16(header.OperationCode)<<11 |
		uint16(header.AuthoritativeAnswer)<<10 |
		uint16(header.Truncation)<<9 |
		uint16(header.RecursionDesired)<<8 |
		uint16(header.ResponseCode)

	binary.BigEndian.PutUint16(buf[0:2], header.PacketIdentifier)
	binary.BigEndian.PutUint16(buf[2:4], flags)
	binary.BigEndian.PutUint16(buf[4:6], header.QuestionCount)
	binary.BigEndian.PutUint16(buf[6:8], header.AnswerRecordCount)
	binary.BigEndian.PutUint16(buf[8:10], header.AuthorityRecordCount)
	binary.BigEndian.PutUint16(buf[10:12], header.AdditionalRecordCount)

	return buf, nil
}
