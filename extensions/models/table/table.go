package table

import (
	"github.com/olekukonko/tablewriter"
	Errors "github.com/vortex14/gotyphoon/errors"
	"os"
)

const (
	NumberMarker = "№"
)

type H []string

type Table struct {
	stateRow int
	headers []string
	data [][]string
}

func (t *Table) GetCurrentRow()  int {
	return t.stateRow
}

func (t *Table) Render()  {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(t.headers)
	table.AppendBulk(t.data)
	table.Render()
}

func (t *Table) SetHeaders(headers H)  {
	t.headers = append([]string{NumberMarker}, headers...)
}

func (t *Table) GetHeaders() H {
	return t.headers
}

func (t *Table) Append(column string, data H) error {
	if len(t.headers) == 0 { return Errors.TableHeadersNotFound }
	if t.stateRow == 0 { t.stateRow ++ }

	return nil
}

func (t *Table) GetCountRow() int {
	return t.stateRow
}
