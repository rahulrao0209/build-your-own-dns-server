// Hardcoded data for DNS server
package data

type answer struct {
	Name   string  // (variable) Domain name encoded as a sequence of labels
	Type   uint16  // (2 bytes) Record type (ex: 1 for A record, 5 for CNAME)
	Class  uint16  // (2 bytes) Record class (ex: 1 for IN)
	TTL    uint32  // (4 bytes) Time to live is the duration in seconds a record can be cached before requerying (ex: 60)
	Length uint16  // (2 bytes) Length of the RDATA field in bytes; RDLENGTH
	Data   [4]byte // (4 bytes) Data specific to the record type; RDATA (ex: \x08\x08\x08\x08 - 4 byte encoding of 8.8.8.8)
}

var Mock = map[string]*answer{
	"banana.com": {
		Name:   "banana.com",
		Type:   1,
		Class:  1,
		TTL:    60,
		Length: 4,
		Data:   [4]byte{8, 8, 8, 8},
	},
	"apple.com": {
		Name:   "apple.com",
		Type:   1,
		Class:  1,
		TTL:    60,
		Length: 4,
		Data:   [4]byte{1, 1, 1, 1},
	},
	"mango.com": {
		Name:   "mango.com",
		Type:   1,
		Class:  1,
		TTL:    60,
		Length: 4,
		Data:   [4]byte{2, 2, 2, 2},
	},
}
