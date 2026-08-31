// Types to define a DNS message packet a.k.a query or reply packet
// with methods for encoding & decoding the packet.
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
	Name  string // (variable) Domain Name encoded as a sequence of labels
	Type  uint16 // (2 bytes) QType (ex: 1 for A record, 5 for CNAME)
	Class uint16 // (2 bytes) Record class (ex: 1 for IN)
}

type Answer struct {
	Name   string  // (variable) Domain name encoded as a sequence of labels
	Type   uint16  // (2 bytes) Record type (ex: 1 for A record, 5 for CNAME)
	Class  uint16  // (2 bytes) Record class (ex: 1 for IN)
	TTL    uint32  // (4 bytes) Time to live is the duration in seconds a record can be cached before requerying (ex: 60)
	Length uint16  // (2 bytes) Length of the RDATA field in bytes; RDLENGTH
	Data   [4]byte // (4 bytes) Data specific to the record type; RDATA (ex: \x08\x08\x08\x08 - 4 byte encoding of 8.8.8.8)
}

type DNSMessage struct {
	Header   Header
	Question Question
	Answer   Answer
}
