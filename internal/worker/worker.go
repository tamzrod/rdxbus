// internal/worker/worker.go
package worker

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/tamzrod/rdxbus/internal/client"
	"github.com/tamzrod/rdxbus/internal/config"
)

// ReadResult is the output of a single read execution.
type ReadResult struct {
	Values  []uint16
	Latency time.Duration
	Err     error
}

// ExecuteRead performs exactly ONE Modbus read.
// No loops, no retries, no scheduling.
func ExecuteRead(cfg config.EngineReadConfig) ReadResult {
	start := time.Now()

	// Dial per call for now (simple + correct).
	// Later, benchmark/poll modes can reuse connections upstream.
	conn, err := client.Dial(cfg.TargetAddr, cfg.Timeout)
	if err != nil {
		return ReadResult{Err: err}
	}
	defer conn.Close()

	// Build request (MBAP + PDU)
	req := client.NewRequest()
	reqBuf := make([]byte, 12) // request.go expects >=12
	frame := req.BuildReadRequest(reqBuf, cfg.UnitID, cfg.FunctionCode, cfg.Address, cfg.Quantity)
	expectedTxID := req.TxID()

	if err := conn.Write(frame); err != nil {
		return ReadResult{Err: err}
	}

	// Read MBAP header
	hdr := make([]byte, 7) // MBAP header size
	if err := conn.ReadFull(hdr); err != nil {
		return ReadResult{Err: err}
	}

	// Read/validate payload into pduBuf via parser
	// NOTE: strict mode is not yet wired into EngineReadConfig; default lenient for now.
	parser := client.NewResponseParser(false)

	// Big enough for max Modbus TCP PDU payloads we care about.
	// (Strict mode uses MBAP length, which can be up to ~253 for reads.)
	pduBuf := make([]byte, 512)

	if err := parser.Parse(conn, expectedTxID, cfg.FunctionCode, hdr, pduBuf); err != nil {
		return ReadResult{Err: err}
	}

	// Extract raw values from the bytes already read into pduBuf.
	values, err := decodeReadValues(pduBuf, cfg.FunctionCode, cfg.Quantity)
	if err != nil {
		return ReadResult{Err: err}
	}

	return ReadResult{
		Values:  values,
		Latency: time.Since(start),
		Err:     nil,
	}
}

// decodeReadValues decodes values from pduBuf for FC 1–4.
// It supports both response layouts produced by parser.Parse():
// - FC-first:        [FC][ByteCount][Data...]
// - UnitID + FC:     [UnitID][FC][ByteCount][Data...]
func decodeReadValues(pduBuf []byte, fc uint8, qty uint16) ([]uint16, error) {
	// Find where FC lives (index 0 or 1). Parser already ensured no exception and FC matches.
	fcIndex := -1
	if len(pduBuf) >= 1 && pduBuf[0] == fc {
		fcIndex = 0
	} else if len(pduBuf) >= 2 && pduBuf[1] == fc {
		fcIndex = 1
	} else {
		return nil, fmt.Errorf("cannot locate function code in response buffer")
	}

	byteCountIndex := fcIndex + 1
	dataIndex := fcIndex + 2

	if byteCountIndex >= len(pduBuf) {
		return nil, fmt.Errorf("response missing bytecount")
	}

	byteCount := int(pduBuf[byteCountIndex])
	if byteCount < 0 {
		return nil, fmt.Errorf("invalid bytecount")
	}

	if dataIndex+byteCount > len(pduBuf) {
		// Parser reads exactly what is needed, but pduBuf is larger; if we hit this, something is inconsistent.
		return nil, fmt.Errorf("response data exceeds buffer")
	}

	data := pduBuf[dataIndex : dataIndex+byteCount]

	switch fc {
	case 3, 4:
		// Register reads: byteCount should be qty*2
		expected := int(qty) * 2
		if byteCount != expected {
			return nil, fmt.Errorf("bytecount mismatch: got %d expected %d", byteCount, expected)
		}
		out := make([]uint16, qty)
		for i := 0; i < int(qty); i++ {
			out[i] = binary.BigEndian.Uint16(data[i*2 : i*2+2])
		}
		return out, nil

	case 1, 2:
		// Bit reads: byteCount should be ceil(qty/8)
		expected := (int(qty) + 7) / 8
		if byteCount != expected {
			return nil, fmt.Errorf("bytecount mismatch: got %d expected %d", byteCount, expected)
		}
		out := make([]uint16, qty)
		for i := 0; i < int(qty); i++ {
			b := data[i/8]
			bit := (b >> uint(i%8)) & 0x01
			out[i] = uint16(bit)
		}
		return out, nil

	default:
		return nil, fmt.Errorf("unsupported function code for read: %d", fc)
	}
}
