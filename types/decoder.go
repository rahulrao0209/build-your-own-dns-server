package types

import (
	"github.com/rahulrao0209/build-your-own-dns-server/constants"
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
	packetId := DecodePacketId(queryheader)
	questionCount := DecodeQuestionCount(queryheader)
	flagDecoder := DecodeFlags(queryheader)
	answerRecordCount := DecodeAnswerRecordCount(queryheader)
	authorityRecordCount := DecodeAuthorityRecordCount(queryheader)
	additionalRecordCount := DecodeAdditionalRecordCount(queryheader)

	header := Header{
		PacketIdentifier: packetId,
		Flags: Flags{
			ResponseIndicator:   flagDecoder(constants.QR),
			OperationCode:       flagDecoder(constants.OPCODE),
			AuthoritativeAnswer: flagDecoder(constants.AA),
			Truncation:          flagDecoder(constants.TC),
			RecursionDesired:    flagDecoder(constants.RD),
			RecursionAvailable:  flagDecoder(constants.RA),
			Reserved:            flagDecoder(constants.ZZZ),
			ResponseCode:        flagDecoder(constants.RCODE),
		},
		QuestionCount:         questionCount,
		AnswerRecordCount:     answerRecordCount,
		AuthorityRecordCount:  authorityRecordCount,
		AdditionalRecordCount: additionalRecordCount,
	}

	return header, nil
}

func (q *Question) UnmarshalBinary(queryQuestion []byte, questionCount int) (Question, int, error) {
	domainName, offset := DecodeDomainName(queryQuestion)
	qType, offset := DecodeType(queryQuestion, offset)
	class, offset := DecodeClass(queryQuestion, offset)

	question := Question{
		Name:  domainName,
		Type:  qType,
		Class: class,
	}

	return question, offset, nil
}
