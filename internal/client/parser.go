package client

import (
	"encoding/binary"
	"fmt"
)

type ModbusExceptionError struct {
	Function uint8
	Code     uint8
}

func (e *ModbusExceptionError) Error() string {
	return fmt.Sprintf("modbus exception fc=%d code=%d", e.Function, e.Code)
}

func IsModbusException(err error) (*ModbusExceptionError, bool) {
	if err == nil {
		return nil, false
	}
	me, ok := err.(*ModbusExceptionError)
	return me, ok
}

type ResponseParser struct {
	strict bool
}

func NewResponseParser(strict bool) *ResponseParser {
	return &ResponseParser{strict: strict}
}

func (p *ResponseParser) Parse(
	conn *Connection,
	expectedTxID uint16,
	expectedFC uint8,
	hdr []byte,
	pduBuf []byte,
) error {

	// MBAP
	txID := binary.BigEndian.Uint16(hdr[0:2])
	if txID != expectedTxID {
		return fmt.Errorf("txid mismatch: got %d expected %d", txID, expectedTxID)
	}
	if binary.BigEndian.Uint16(hdr[2:4]) != 0 {
		return fmt.Errorf("invalid protocol id")
	}

	length := binary.BigEndian.Uint16(hdr[4:6])
	if length < 1 {
		return fmt.Errorf("invalid mbap length")
	}

	// STRICT: read exactly length bytes and validate unitID+FC framing
	if p.strict {
		if int(length) > len(pduBuf) {
			return fmt.Errorf("pdu buffer too small")
		}
		if err := conn.ReadFull(pduBuf[:length]); err != nil {
			return err
		}
		return validateStrictPDU(pduBuf[:length], expectedFC)
	}

	// LENIENT: auto-detect whether payload begins with FC or UnitID.
	// We will read enough bytes to extract: FC (+ optional bytecount) or exception code.

	// Read first 2 bytes of "payload"
	if err := conn.ReadFull(pduBuf[:2]); err != nil {
		return err
	}

	b0 := pduBuf[0]
	b1 := pduBuf[1]

	// Case A: FC-first (b0 is fc or exception fc)
	if b0 == expectedFC || b0 == (expectedFC|0x80) {
		fc := b0
		if fc&0x80 != 0 {
			// b1 is exception code
			return &ModbusExceptionError{Function: fc & 0x7F, Code: b1}
		}

		// normal: b1 is bytecount
		byteCount := int(b1)
		if byteCount > 0 {
			if err := conn.ReadFull(pduBuf[2 : 2+byteCount]); err != nil {
				return err
			}
		}
		return nil
	}

	// Case B: UnitID + FC (b1 is fc or exception fc)
	fc := b1
	if fc&0x80 != 0 {
		// need 1 more byte for exception code
		if err := conn.ReadFull(pduBuf[2:3]); err != nil {
			return err
		}
		return &ModbusExceptionError{Function: fc & 0x7F, Code: pduBuf[2]}
	}

	if fc != expectedFC {
		return fmt.Errorf("function code mismatch: got %d expected %d", fc, expectedFC)
	}

	// need 1 more byte for bytecount
	if err := conn.ReadFull(pduBuf[2:3]); err != nil {
		return err
	}
	byteCount := int(pduBuf[2])

	if byteCount > 0 {
		if err := conn.ReadFull(pduBuf[3 : 3+byteCount]); err != nil {
			return err
		}
	}

	return nil
}

func validateStrictPDU(pdu []byte, expectedFC uint8) error {
	if len(pdu) < 2 {
		return fmt.Errorf("pdu too short")
	}
	// pdu[0] = UnitID (ignored)
	fc := pdu[1]

	if fc&0x80 != 0 {
		if len(pdu) < 3 {
			return fmt.Errorf("malformed exception response")
		}
		return &ModbusExceptionError{Function: fc & 0x7F, Code: pdu[2]}
	}

	if fc != expectedFC {
		return fmt.Errorf("function code mismatch: got %d expected %d", fc, expectedFC)
	}
	return nil
}
