package types

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

type Question struct {
	Name  string // content
	Type  uint16 // 2 bytes
	Class uint16 // 2  bytes
}

type DNSMessage struct {
	Header   Header
	Question Question
}
