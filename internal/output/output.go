package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

var globalFormat string

func SetFormat(f string) {
	globalFormat = f
}

func GetFormat() string {
	return globalFormat
}

func IsJSON() bool {
	return globalFormat == "json" || globalFormat == "jsonl"
}

func IsJSONL() bool {
	return globalFormat == "jsonl"
}

func PrintJSON(data any) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON 编码失败: %v\n", err)
		return
	}
	fmt.Println(string(out))
}

func PrintJSONL(items []map[string]any) {
	enc := json.NewEncoder(os.Stdout)
	for _, item := range items {
		enc.Encode(item)
	}
}

func PrintTable(headers []string, rows [][]string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	colConfigs := make([]table.ColumnConfig, len(headers))
	for i := range headers {
		colConfigs[i] = table.ColumnConfig{Number: i + 1, WidthMax: 60}
	}
	t.SetColumnConfigs(colConfigs)

	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	for _, row := range rows {
		r := make(table.Row, len(row))
		for i, v := range row {
			r[i] = v
		}
		t.AppendRow(r)
	}

	fmt.Println(t.Render())
}
