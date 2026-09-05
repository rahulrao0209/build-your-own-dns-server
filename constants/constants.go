// App wide constants
package constants

// Flag identifiers
type Flag string

const (
	QR     Flag = "ResponseIndicator"
	OPCODE Flag = "OperationCode"
	AA     Flag = "AuthoritativeAnswer"
	TC     Flag = "Truncation"
	RD     Flag = "RecursionDesired"
	RA     Flag = "RecursionAvailable"
	ZZZ    Flag = "Reserved"
	RCODE  Flag = "ResponseCode"
)
