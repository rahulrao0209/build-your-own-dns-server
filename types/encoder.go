package types

import (
	"bytes"
	"encoding/binary"

	"github.com/rahulrao0209/build-your-own-dns-server/data"
)

func (m *DNSMessage) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Encode DNS reply question section
	question, err := m.Question.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Encode DNS reply answer section
	answer, err := m.Answer.MarshalBinary(&m.Header, m.Question)
	if err != nil {
		return nil, err
	}

	// Encode DNS reply header section
	// Note: header is encoded at the end because some sections of the header
	// depend upon the question and answer sections
	// ex: The answerRecordCount in the header depends on whether the server
	// has an answer for the query which is determined in the encoding method for the answer.
	header, err := m.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Write all sections
	buf.Write(header)
	buf.Write(question)
	buf.Write(answer)

	return buf.Bytes(), nil
}

func (h *Header) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 12)

	// marshal header

	// set flags
	h.ResponseIndicator = 1 // Reply packet
	h.AuthoritativeAnswer = 0
	h.Truncation = 0
	h.RecursionAvailable = 0
	h.Reserved = 0

	if h.OperationCode == 0 {
		h.ResponseCode = 0 // no error if opcode is 0
	} else {
		h.ResponseCode = 4 // not implemented
	}

	flags := EncodeFlags(h.Flags)

	binary.BigEndian.PutUint16(buf[0:2], h.PacketIdentifier)
	binary.BigEndian.PutUint16(buf[2:4], flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QuestionCount)
	binary.BigEndian.PutUint16(buf[6:8], h.AnswerRecordCount)
	binary.BigEndian.PutUint16(buf[8:10], h.AuthorityRecordCount)
	binary.BigEndian.PutUint16(buf[10:12], h.AdditionalRecordCount)

	return buf, nil
}

func (q *Questions) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	for _, qu := range *q {
		// encode domain name
		err := EncodeDomainName(&buf, qu.Name)
		if err != nil {
			return nil, err
		}

		// encode type
		err = AppendNBytes(&buf, 2, qu.Type)

		// encode class
		err = AppendNBytes(&buf, 2, qu.Class)

		if err != nil {
			return nil, err
		}

	}

	return buf.Bytes(), nil
}

func (a *Answer) MarshalBinary(h *Header, q Questions) ([]byte, error) {
	var buf bytes.Buffer

	for _, qu := range q {
		domainName := qu.Name
		qtype := qu.Type
		class := qu.Class

		var ttl uint32
		var rDataLength uint16
		var rData [4]byte
		dname, ok := data.Mock[domainName]
		if ok {
			h.AnswerRecordCount++
			ttl = dname.TTL
			rDataLength = dname.Length
			rData = dname.Data
		}

		// encode values
		EncodeDomainName(&buf, domainName)
		AppendNBytes(&buf, 2, qtype)
		AppendNBytes(&buf, 2, class)
		AppendNBytes(&buf, 4, ttl)
		AppendNBytes(&buf, 2, rDataLength)
		err := AppendNBytes(&buf, 4, rData)

		if err != nil {
			return nil, err
		}

	}

	return buf.Bytes(), nil
}
