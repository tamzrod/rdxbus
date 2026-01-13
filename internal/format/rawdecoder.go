// internal/format/rawdecoder.go
package format

import (
	"encoding/binary"
	"fmt"
)

// DecodeReadValues decodes Modbus read responses (FC 1–4)
// from a raw PDU buffer.
//
// Supported layouts:
//   - [FC][ByteCount][Data...]
//   - [UnitID][FC][ByteCount][Data...]
//
// The caller must ensure:
//   - response is non-exception
//   - function code matches the request
func DecodeReadValues(
	pdu []byte,
	functionCode uint8,
	quantity uint16,
) ([]uint16, error) {

	// Locate function code (index 0 or 1)
	fcIndex := -1
	if len(pdu) >= 1 && pdu[0] == functionCode {
		fcIndex = 0
	} else if len(pdu) >= 2 && pdu[1] == functionCode {
		fcIndex = 1
	} else {
		return nil, fmt.Errorf("function code not found in PDU")
	}

	byteCountIndex := fcIndex + 1
	dataIndex := fcIndex + 2

	if byteCountIndex >= len(pdu) {
		return nil, fmt.Errorf("missing byte count")
	}

	byteCount := int(pdu[byteCountIndex])
	if dataIndex+byteCount > len(pdu) {
		return nil, fmt.Errorf("data exceeds PDU length")
	}

	data := pdu[dataIndex : dataIndex+byteCount]

	switch functionCode {

	case 3, 4:
		expected := int(quantity) * 2
		if byteCount != expected {
			return nil, fmt.Errorf(
				"register bytecount mismatch: got %d expected %d",
				byteCount, expected,
			)
		}

		out := make([]uint16, quantity)
		for i := 0; i < int(quantity); i++ {
			out[i] = binary.BigEndian.Uint16(data[i*2 : i*2+2])
		}
		return out, nil

	case 1, 2:
		expected := (int(quantity) + 7) / 8
		if byteCount != expected {
			return nil, fmt.Errorf(
				"bit bytecount mismatch: got %d expected %d",
				byteCount, expected,
			)
		}

		out := make([]uint16, quantity)
		for i := 0; i < int(quantity); i++ {
			b := data[i/8]
			out[i] = uint16((b >> uint(i%8)) & 0x01)
		}
		return out, nil

	default:
		return nil, fmt.Errorf("unsupported function code: %d", functionCode)
	}
}
