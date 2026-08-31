package types

import (
	"bytes"
	"encoding/binary"

	"github.com/rahulrao0209/build-your-own-dns-server/data"
	"github.com/rahulrao0209/build-your-own-dns-server/utils"
)

func (m *DNSMessage) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Encode DNS reply header section
	header, err := m.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(header)

	// Encode DNS reply question section
	question, err := m.Question.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(question)

	// Encode DNS reply answer section
	answer, err := m.Answer.MarshalBinary(m.Question)
	if err != nil {
		return nil, err
	}
	buf.Write(answer)

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

	// encode domain name
	err := utils.EncodeDomainName(&buf, q.Name)
	if err != nil {
		return nil, err
	}

	// encode type
	err = utils.AppendNBytes(&buf, 2, q.Type)

	// encode class
	err = utils.AppendNBytes(&buf, 2, q.Class)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Answer) MarshalBinary(q Question) ([]byte, error) {
	var buf bytes.Buffer

	domainName := q.Name
	qtype := q.Type
	class := q.Class
	ttl := data.Mock[domainName].TTL
	rDataLength := data.Mock[domainName].Length
	rData := data.Mock[domainName].Data

	// encode values
	utils.EncodeDomainName(&buf, domainName)
	utils.AppendNBytes(&buf, 2, qtype)
	utils.AppendNBytes(&buf, 2, class)
	utils.AppendNBytes(&buf, 4, ttl)
	utils.AppendNBytes(&buf, 2, rDataLength)
	err := utils.AppendNBytes(&buf, 4, rData)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
