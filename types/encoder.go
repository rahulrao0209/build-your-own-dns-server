package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

func (m *DNSMessage) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// write header section
	header, err := m.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(header)

	// write question section
	question, err := m.Question.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(question)

	return buf.Bytes(), nil
}

func (h *Header) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 12)

	// marshal header
	flags := uint16(h.ResponseIndicator)<<15 |
		uint16(h.OperationCode)<<11 |
		uint16(h.AuthoritativeAnswer)<<10 |
		uint16(h.Truncation)<<9 |
		uint16(h.RecursionDesired)<<8 |
		uint16(h.ResponseCode)

	binary.BigEndian.PutUint16(buf[0:2], h.PacketIdentifier)
	binary.BigEndian.PutUint16(buf[2:4], flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QuestionCount)
	binary.BigEndian.PutUint16(buf[6:8], h.AnswerRecordCount)
	binary.BigEndian.PutUint16(buf[8:10], h.AuthorityRecordCount)
	binary.BigEndian.PutUint16(buf[10:12], h.AdditionalRecordCount)

	return buf, nil
}

func (q *Question) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	labels := strings.Split(q.Name, ".")
	for _, l := range labels {
		if len(l) > 63 {
			return nil, errors.New("label too long")
		}

		// encode each label in the format: length content
		buf.WriteByte(byte(len(l)))
		buf.WriteString(l)
	}
	// terminate the label with a zero byte
	buf.WriteByte(0)

	// encode type & class
	var tmp [2]byte

	// type
	binary.BigEndian.PutUint16(tmp[:], q.Type)
	buf.Write(tmp[:])

	// class
	binary.BigEndian.PutUint16(tmp[:], q.Class)
	buf.Write(tmp[:])

	return buf.Bytes(), nil
}
