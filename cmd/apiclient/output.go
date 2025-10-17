package apiclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// OutputJSON formats and prints data as JSON
func OutputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// OutputError prints an error message based on the format flag
func OutputError(err error) {
	if format == "json" {
		errData := map[string]string{"error": err.Error()}
		_ = OutputJSON(errData)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

// FormatRelation formats a relationship as "name(id)"
func FormatRelation(name string, id uint) string {
	return fmt.Sprintf("%s(%d)", name, id)
}

// TableWriter wraps go-pretty table for consistent table output
type TableWriter struct {
	table table.Writer
}

// NewTableWriter creates a new table writer
func NewTableWriter() *TableWriter {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleDefault)

	return &TableWriter{
		table: t,
	}
}

// WriteHeader writes the table header
func (tw *TableWriter) WriteHeader(columns ...string) {
	header := table.Row{}
	for _, col := range columns {
		header = append(header, col)
	}
	tw.table.AppendHeader(header)
}

// WriteRow writes a table row
func (tw *TableWriter) WriteRow(values ...string) {
	row := table.Row{}
	for _, val := range values {
		row = append(row, val)
	}
	tw.table.AppendRow(row)
}

// Flush flushes the table output
func (tw *TableWriter) Flush() error {
	tw.table.Render()
	return nil
}

// OutputTable is a convenience function for simple table output
func OutputTable(headers []string, rows [][]string) error {
	tw := NewTableWriter()
	tw.WriteHeader(headers...)
	for _, row := range rows {
		tw.WriteRow(row...)
	}
	return tw.Flush()
}
