// Utilities for encoding and decoding the contents of a DNS message
package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/rahulrao0209/build-your-own-dns-server/constants"
)

/* ******** HEADER SECTION UTILS ******* */
func DecodePacketId(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[0:2])
}

func DecodeFlags(buf []byte) func(constants.Flag) uint8 {
	flags := binary.BigEndian.Uint16(buf[2:4])

	return func(flag constants.Flag) uint8 {
		switch flag {
		case constants.QR:
			return uint8((flags >> 15) & 0x1)
		case constants.OPCODE:
			return uint8((flags >> 11) & 0xF)
		case constants.AA:
			return uint8((flags >> 10) & 0x1)
		case constants.TC:
			return uint8((flags >> 9) & 0x1)
		case constants.RD:
			return uint8((flags >> 8) & 0x1)
		case constants.RA:
			return uint8((flags >> 7) & 0x1)
		case constants.ZZZ:
			return uint8((flags >> 4) & 0x7)
		case constants.RCODE:
			return uint8(flags & 0xF)
		default: // should ideally never hit this case
			return 0
		}
	}
}

func EncodeFlags(flags Flags) uint16 {
	return uint16(flags.ResponseIndicator)<<15 |
		uint16(flags.OperationCode)<<11 |
		uint16(flags.AuthoritativeAnswer)<<10 |
		uint16(flags.Truncation)<<9 |
		uint16(flags.RecursionDesired)<<8 |
		uint16(flags.RecursionAvailable)<<7 |
		uint16(flags.Reserved)<<4 |
		uint16(flags.ResponseCode)
}

func DecodeQuestionCount(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[4:6])
}

func DecodeAnswerRecordCount(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[6:8])
}

func DecodeAuthorityRecordCount(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[8:10])
}

func DecodeAdditionalRecordCount(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf[10:12])
}

/* ******** QUESTION SECTION  UTILS ******* */
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
