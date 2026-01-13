// internal/output/model.go
package output

import "time"

type Output struct {
	Meta    Meta
	Table   *Table
	Message string
	Error   string
}

type Meta struct {
	Mode      string
	Target    string
	UnitID    uint8
	Function  uint8
	Address   uint16
	Quantity  uint16
	Latency   time.Duration
	Timestamp time.Time
}

type Table struct {
	Columns []Column
	Rows    []Row
}

type Column struct {
	Key   string
	Title string
}

type Row struct {
	Cells map[string]any
}
