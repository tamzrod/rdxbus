// internal/engine/modbus_engine_test.go
package engine

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"github.com/tamzrod/rdxbus/internal/format"
)

func TestModbusEngine_Execute_FC3_Success_DecodeValues(t *testing.T) {
	addr := startFakeModbusTCPServer(t, func(c net.Conn) {
		// Read the Modbus TCP request (expected 12 bytes for FC3 read)
		req := make([]byte, 12)
		if _, err := io.ReadFull(c, req); err != nil {
			t.Fatalf("server failed to read request: %v", err)
		}

		// Extract TxID to echo back
		txID := binary.BigEndian.Uint16(req[0:2])

		// Build a valid FC3 response for qty=2 registers:
		// MBAP: TxID(2) Proto(2=0) Len(2) Unit(1)
		// PDU:  FC(1) ByteCount(1) Data(4)
		unitID := req[6]
		fc := req[7]

		// Data: 0x002A, 0x002B
		data := []byte{0x00, 0x2A, 0x00, 0x2B}
		byteCount := byte(len(data))

		// Length = UnitID(1) + PDU(1+1+4) = 7
		resp := make([]byte, 7+1+1+len(data))
		binary.BigEndian.PutUint16(resp[0:2], txID)
		binary.BigEndian.PutUint16(resp[2:4], 0) // Protocol ID
		binary.BigEndian.PutUint16(resp[4:6], uint16(1+1+1+len(data)))
		resp[6] = unitID
		resp[7] = fc
		resp[8] = byteCount
		copy(resp[9:], data)

		if _, err := c.Write(resp); err != nil {
			t.Fatalf("server failed to write response: %v", err)
		}
	})

	eng := &ModbusEngine{
		TargetAddr: addr,
		Strict:     false,
	}

	req := Request{
		UnitID:       1,
		FunctionCode: 3,
		Address:      0,
		Quantity:     2,
		Timeout:      2 * time.Second,
	}

	res := eng.Execute(context.Background(), req)
	if res.Err != nil {
		t.Fatalf("Execute returned error: %v", res.Err)
	}
	if res.Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", res.Duration)
	}
	if len(res.Raw) == 0 {
		t.Fatalf("expected non-empty Raw buffer")
	}

	// Decode using the formal decoder (adapter layer)
	values, err := format.DecodeReadValues(res.Raw, req.FunctionCode, req.Quantity)
	if err != nil {
		t.Fatalf("DecodeReadValues error: %v", err)
	}

	if len(values) != 2 || values[0] != 0x002A || values[1] != 0x002B {
		t.Fatalf("unexpected values: %#v", values)
	}
}

func TestModbusEngine_Execute_ExceptionResponse_ReturnsError(t *testing.T) {
	addr := startFakeModbusTCPServer(t, func(c net.Conn) {
		// Read request (12 bytes)
		req := make([]byte, 12)
		if _, err := io.ReadFull(c, req); err != nil {
			t.Fatalf("server failed to read request: %v", err)
		}

		txID := binary.BigEndian.Uint16(req[0:2])
		unitID := req[6]
		fc := req[7]

		// Exception response:
		// PDU: (FC|0x80), ExceptionCode(1)
		excFC := fc | 0x80
		excCode := byte(0x02) // Illegal Data Address (example)

		// Length = Unit(1) + PDU(2) = 3
		resp := make([]byte, 7+2)
		binary.BigEndian.PutUint16(resp[0:2], txID)
		binary.BigEndian.PutUint16(resp[2:4], 0)
		binary.BigEndian.PutUint16(resp[4:6], 3)
		resp[6] = unitID
		resp[7] = excFC
		resp[8] = excCode

		if _, err := c.Write(resp); err != nil {
			t.Fatalf("server failed to write response: %v", err)
		}
	})

	eng := &ModbusEngine{
		TargetAddr: addr,
		Strict:     false,
	}

	req := Request{
		UnitID:       1,
		FunctionCode: 3,
		Address:      0,
		Quantity:     2,
		Timeout:      2 * time.Second,
	}

	res := eng.Execute(context.Background(), req)
	if res.Err == nil {
		t.Fatalf("expected error on exception response, got nil")
	}
}

// startFakeModbusTCPServer starts a single-shot TCP server.
// It accepts exactly one connection, runs handler, then closes.
// Returns the listener address (host:port).
func startFakeModbusTCPServer(t *testing.T, handler func(net.Conn)) string {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer ln.Close()

		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()

		handler(c)
	}()

	t.Cleanup(func() {
		_ = ln.Close()
		<-done
	})

	return ln.Addr().String()
}
