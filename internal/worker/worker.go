package worker

import (
	"time"

	"modbus-stress/internal/client"
	"modbus-stress/internal/scheduler"
)

type Result struct {
	Latency time.Duration
	Err     error
}

type Worker struct {
	conn      *client.Connection
	req       *client.Request
	parser    *client.ResponseParser
	scheduler *scheduler.Scheduler

	unitID   uint8
	fc       uint8
	address  uint16
	quantity uint16
}

func NewWorker(
	conn *client.Connection,
	scheduler *scheduler.Scheduler,
	unitID uint8,
	fc uint8,
	address uint16,
	quantity uint16,
	strict bool,
) *Worker {
	return &Worker{
		conn:      conn,
		req:       client.NewRequest(),
		parser:    client.NewResponseParser(strict),
		scheduler: scheduler,
		unitID:    unitID,
		fc:        fc,
		address:   address,
		quantity:  quantity,
	}
}

func (w *Worker) Run(stop <-chan struct{}, results chan<- Result) {
	defer w.conn.Close()

	reqBuf := make([]byte, 12)
	hdrBuf := make([]byte, 7)
	pduBuf := make([]byte, 260)

	for {
		select {
		case <-stop:
			return
		default:
		}

		if !w.scheduler.Wait() {
			return
		}

		adu := w.req.BuildReadRequest(reqBuf, w.unitID, w.fc, w.address, w.quantity)
		txID := w.req.TxID()

		start := time.Now()

		if err := w.conn.Write(adu); err != nil {
			results <- Result{Err: err}
			return
		}

		if err := w.conn.ReadFull(hdrBuf); err != nil {
			results <- Result{Err: err}
			return
		}

		err := w.parser.Parse(w.conn, txID, w.fc, hdrBuf, pduBuf)

		results <- Result{
			Latency: time.Since(start),
			Err:     err,
		}
	}
}
