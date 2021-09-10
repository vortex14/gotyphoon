package table

import (
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	Errors "github.com/vortex14/gotyphoon/errors"
)

const (
	NumberMarker = "№"
)

type H []   string
type R []   string
type D [] [] string

type Table struct {
	data      D
	headers   H
	stateRow  int
	NumberRow bool
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
	if t.NumberRow {
		t.headers = append([]string{NumberMarker}, headers...)
	} else {
		t.headers = append(t.headers, headers...)
	}

}

func (t *Table) GetHeaders() H {
	return t.headers
}

func (t *Table) Append(row R) *Table {
	if len(t.headers) == 0 { color.Red(Errors.TableHeadersNotFound.Error()); return t}
	t.stateRow ++
	if t.data == nil { t.data = make([][]string, 0) }
	if t.NumberRow {
		t.data = append(t.data, append([]string{strconv.Itoa(t.stateRow)}, row...))
	} else {
		t.data = append(t.data, row)
	}
	return t
}

func (t *Table) AppendBulk(data D) *Table {
	if len(t.headers) == 0 { color.Red(Errors.TableHeadersNotFound.Error()); return t}
	if t.data == nil { t.data = make(D, 0) }
	if t.NumberRow {
		for ri := range data { t.stateRow ++; data[ri] = append([]string{strconv.Itoa(t.stateRow)}, data[ri]...) }
	}
	t.data = append(t.data, data...)

	return t
}

func (t *Table) GetCountRow() int {
	return t.stateRow
}
