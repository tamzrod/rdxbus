// cmd/rdxbus/easy_helpers.go
package main

import "github.com/tamzrod/rdxbus/internal/output"

func buildRows(start uint16, values []uint16) []output.Row {
	rows := make([]output.Row, 0, len(values))
	for i, v := range values {
		rows = append(rows, output.Row{
			Cells: map[string]any{
				"address": int(start) + i,
				"raw":     v,
			},
		})
	}
	return rows
}

func buildTable(rows []output.Row) *output.Table {
	return &output.Table{
		Columns: []output.Column{
			{Key: "address", Title: "Address"},
			{Key: "raw", Title: "Raw"},
		},
		Rows: rows,
	}
}
