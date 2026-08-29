package types

import (
	"github.com/rahulrao0209/build-your-own-dns-server/utils"
)

func (m *DNSMessage) UnmarshalBinary(query []byte) (*DNSMessage, error) {
	// Decode DNS query header
	queryHeader := query[:12]
	header, err := m.Header.UnmarshalBinary(queryHeader)
	if err != nil {
		return nil, nil
	}

	// Decode DNS query question
	queryQuestion := query[12:]
	question, _, err := m.Question.UnmarshalBinary(queryQuestion, int(header.QuestionCount))
	if err != nil {
		return nil, nil
	}

	// Assemble DNS reply message
	dnsReply := &DNSMessage{
		Header:   header,
		Question: question,
	}

	return dnsReply, nil
}

func (h *Header) UnmarshalBinary(queryheader []byte) (Header, error) {
	packetId := utils.DecodePacketId(queryheader)
	questionCount := utils.DecodeQuestionCount(queryheader)

	header := Header{
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
	}

	return header, nil
}

func (q *Question) UnmarshalBinary(queryQuestion []byte, questionCount int) (Question, int, error) {
	domainName, offset := utils.DecodeDomainName(queryQuestion)
	qType, offset := utils.DecodeType(queryQuestion, offset)
	class, offset := utils.DecodeClass(queryQuestion, offset)

	question := Question{
		Name:  domainName,
		Type:  qType,
		Class: class,
	}

	return question, offset, nil
}
