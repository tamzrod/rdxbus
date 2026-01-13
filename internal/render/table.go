package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/tamzrod/rdxbus/internal/output"
)

func Render(w io.Writer, out output.Output) {
	renderMeta(w, out.Meta)

	if out.Message != "" {
		fmt.Fprintln(w, out.Message)
	}

	if out.Error != "" {
		fmt.Fprintf(w, "ERROR: %s\n", out.Error)
		return
	}

	if out.Table != nil {
		renderTable(w, out.Table)
	}
}

func renderMeta(w io.Writer, m output.Meta) {
	if m.Mode == "" {
		return
	}

	fmt.Fprintf(w, "%s RESULT\n\n", strings.ToUpper(m.Mode))

	if m.Target != "" {
		fmt.Fprintf(w, "Target:   %s\n", m.Target)
	}
	if m.UnitID != 0 {
		fmt.Fprintf(w, "Unit ID:  %d\n", m.UnitID)
	}
	if m.Function != 0 {
		fmt.Fprintf(w, "Function: %d\n", m.Function)
	}
	if m.Latency > 0 {
		fmt.Fprintf(w, "Latency:  %s\n", m.Latency)
	}

	fmt.Fprintln(w)
}

func renderTable(w io.Writer, t *output.Table) {
	if len(t.Columns) == 0 {
		return
	}

	widths := make(map[string]int)

	for _, col := range t.Columns {
		widths[col.Key] = len(col.Title)
	}

	for _, row := range t.Rows {
		for _, col := range t.Columns {
			val := fmt.Sprintf("%v", row.Cells[col.Key])
			if len(val) > widths[col.Key] {
				widths[col.Key] = len(val)
			}
		}
	}

	for i, col := range t.Columns {
		if i > 0 {
			fmt.Fprint(w, "  ")
		}
		fmt.Fprintf(w, "%-*s", widths[col.Key], col.Title)
	}
	fmt.Fprintln(w)

	for i, col := range t.Columns {
		if i > 0 {
			fmt.Fprint(w, "  ")
		}
		fmt.Fprint(w, strings.Repeat("-", widths[col.Key]))
	}
	fmt.Fprintln(w)

	for _, row := range t.Rows {
		for i, col := range t.Columns {
			if i > 0 {
				fmt.Fprint(w, "  ")
			}
			val := fmt.Sprintf("%v", row.Cells[col.Key])
			fmt.Fprintf(w, "%-*s", widths[col.Key], val)
		}
		fmt.Fprintln(w)
	}
}
