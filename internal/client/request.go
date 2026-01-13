package client

import "encoding/binary"

const (
	mbapHeaderSize = 7
)

type Request struct {
	txID uint16
}

func NewRequest() *Request {
	return &Request{}
}

// TxID returns the most recently used transaction id.
func (r *Request) TxID() uint16 {
	return r.txID
}

// BuildReadRequest builds a Modbus TCP Read request (FC 1–4).
// buf must be at least 12 bytes.
func (r *Request) BuildReadRequest(
	buf []byte,
	unitID uint8,
	functionCode uint8,
	address uint16,
	quantity uint16,
) []byte {

	r.txID++

	// MBAP Header
	binary.BigEndian.PutUint16(buf[0:2], r.txID)
	binary.BigEndian.PutUint16(buf[2:4], 0) // Protocol ID
	binary.BigEndian.PutUint16(buf[4:6], 6) // Length = UnitID + PDU(5)
	buf[6] = unitID

	// PDU
	buf[7] = functionCode
	binary.BigEndian.PutUint16(buf[8:10], address)
	binary.BigEndian.PutUint16(buf[10:12], quantity)

	return buf[:12]
}
