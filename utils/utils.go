package utils

import (
	"encoding/binary"
	"strings"
)

// Decodes packet ID from a DNS request's header section.
func DecodePacketId(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[0:2])
}

func DecodeQuestionCount(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[4:6])
}

// Decodes domain name from a DNS request's question section
// ex: banana.com
// 06 b  a  n  a  n  a  03 c  o  m
// 06 62 61 6e 61 6e 61 03 63 6f 6d 00
// 0. 1. 2. 3. 4. 5. 6. 7. 8. 9. 10. 11
func DecodeDomainName(buf []byte) (string, int) {
	var labels []string
	var offset int

	for i := 0; i < len(buf); {
		length := int(buf[i])
		i++

		if length == 0 { // end of domain name is marked by a zero byte
			offset = i
			break
		}

		labels = append(labels, string(buf[i:length+i]))
		i += length
	}
	return strings.Join(labels, "."), offset
}

// Decodes record type a.k.a QType from a DNS request's question section.
func DecodeType(buf []byte, offset int) (uint16, int) {
	return binary.BigEndian.Uint16(buf[offset : offset+2]), offset + 2
}

// Decodes record class from a DNS's requests question section.
func DecodeClass(buf []byte, offset int) (uint16, int) {
	return binary.BigEndian.Uint16(buf[offset : offset+2]), offset + 2
}
