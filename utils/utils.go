// Utilities for encoding and decoding the contents of a DNS message
package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
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

func EncodeDomainName(buf *bytes.Buffer, name string) error {
	labels := strings.Split(name, ".")
	for _, l := range labels {
		if len(l) > 63 {
			return errors.New("label too long")
		}

		// encode each label in the format: length content
		buf.WriteByte(byte(len(l)))
		buf.WriteString(l)
	}
	// terminate the label with a zero byte
	buf.WriteByte(0)
	return nil
}

func AppendNBytes(buf *bytes.Buffer, n int, data any) error {
	var tmp = make([]byte, n)

	switch v := data.(type) {
	case uint16:
		binary.BigEndian.PutUint16(tmp[:], v)
	case uint32:
		binary.BigEndian.PutUint32(tmp[:], v)
	case [4]byte:
		// Appends to the end of the a byte slice and therefore requires
		// a new zero length slice to avoid appending the data after the length.
		tmp, _ = binary.Append([]byte{}, binary.BigEndian, v)
	}

	buf.Write(tmp[:])
	return nil
}

// Decodes record type a.k.a QType from a DNS request's question section.
func DecodeType(buf []byte, offset int) (uint16, int) {
	return binary.BigEndian.Uint16(buf[offset : offset+2]), offset + 2
}

// Decodes record class from a DNS's requests question section.
func DecodeClass(buf []byte, offset int) (uint16, int) {
	return binary.BigEndian.Uint16(buf[offset : offset+2]), offset + 2
}
